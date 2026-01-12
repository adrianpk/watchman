// Package hook provides the core hook evaluation logic.
package hook

import (
	"os"
	"strings"

	"github.com/adrianpk/watchman/internal/config"
	"github.com/adrianpk/watchman/internal/parser"
	"github.com/adrianpk/watchman/internal/policy"
)

// Input represents the hook input from Claude Code.
type Input struct {
	HookType  string
	ToolName  string
	ToolInput map[string]interface{}
}

// Result represents the evaluation result.
type Result struct {
	Allowed bool
	Reason  string
	Warning string
}

// Evaluator evaluates hook inputs against configured rules.
type Evaluator struct {
	cfg         *config.Config
	hookMatcher *HookMatcher
	hookExec    *HookExecutor
}

// NewEvaluator creates a new hook evaluator.
func NewEvaluator(cfg *config.Config) *Evaluator {
	return &Evaluator{
		cfg:         cfg,
		hookMatcher: NewHookMatcher(),
		hookExec:    NewHookExecutor(),
	}
}

// Evaluate processes the hook input and returns a result.
func (e *Evaluator) Evaluate(input Input) Result {
	// Check tool blocklist
	if e.isToolBlocked(input.ToolName) {
		return Result{Allowed: false, Reason: "tool is blocked by configuration: " + input.ToolName}
	}

	// Check tool allowlist
	if !e.isToolAllowed(input.ToolName) {
		return Result{Allowed: false, Reason: "tool is not in allowed list: " + input.ToolName}
	}

	// Non-filesystem tools are always allowed
	if !isFilesystemTool(input.ToolName) {
		return Result{Allowed: true}
	}

	// Check command blocklist for Bash
	if input.ToolName == "Bash" {
		if cmd, ok := input.ToolInput["command"].(string); ok {
			if blocked := e.isCommandBlocked(cmd); blocked != "" {
				return Result{Allowed: false, Reason: "command is blocked by configuration: " + blocked}
			}
		}
	}

	// Check protected paths
	paths := ExtractPaths(input.ToolName, input.ToolInput)
	for _, p := range paths {
		if policy.IsAlwaysProtected(p) {
			return Result{Allowed: false, Reason: "path is protected and cannot be accessed. User must perform this action manually."}
		}
	}

	// Apply workspace rule
	if e.cfg.Rules.Workspace {
		if result := e.evaluateWorkspace(input); !result.Allowed {
			return result
		}
	}

	// Apply scope rule
	if e.cfg.Rules.Scope {
		if result := e.evaluateScope(input); !result.Allowed {
			return result
		}
	}

	// Apply versioning rule
	if e.cfg.Rules.Versioning && input.ToolName == "Bash" {
		if result := e.evaluateVersioning(input); !result.Allowed {
			return result
		}
	}

	// Apply incremental rule
	if e.cfg.Rules.Incremental && isModificationTool(input.ToolName) {
		if result := e.evaluateIncremental(); !result.Allowed || result.Warning != "" {
			return result
		}
	}

	// Apply invariants rule
	if e.cfg.Rules.Invariants && isModificationTool(input.ToolName) {
		if result := e.evaluateInvariants(input); !result.Allowed {
			return result
		}
	}

	// Apply external hooks
	if len(e.cfg.Hooks) > 0 {
		if result := e.evaluateHooks(input); !result.Allowed || result.Warning != "" {
			return result
		}
	}

	return Result{Allowed: true}
}

func (e *Evaluator) evaluateWorkspace(input Input) Result {
	rule := policy.NewConfineToWorkspace(&e.cfg.Workspace)
	paths := ExtractPaths(input.ToolName, input.ToolInput)
	for _, p := range paths {
		parsed := parser.Command{Args: []string{p}}
		decision := rule.Evaluate(parsed)
		if !decision.Allowed {
			return Result{Allowed: false, Reason: decision.Reason}
		}
	}
	return Result{Allowed: true}
}

func (e *Evaluator) evaluateScope(input Input) Result {
	rule := policy.NewScopeToFiles(&e.cfg.Scope)
	paths := ExtractPaths(input.ToolName, input.ToolInput)
	for _, p := range paths {
		parsed := parser.Command{Args: []string{p}}
		decision := rule.Evaluate(input.ToolName, parsed)
		if !decision.Allowed {
			return Result{Allowed: false, Reason: decision.Reason}
		}
	}
	return Result{Allowed: true}
}

func (e *Evaluator) evaluateVersioning(input Input) Result {
	cmd, ok := input.ToolInput["command"].(string)
	if !ok {
		return Result{Allowed: true}
	}
	rule := policy.NewVersioningRule(&e.cfg.Versioning)
	decision := rule.Evaluate(cmd)
	return Result{Allowed: decision.Allowed, Reason: decision.Reason}
}

func (e *Evaluator) evaluateIncremental() Result {
	rule := policy.NewIncrementalRule(&e.cfg.Incremental)
	decision := rule.Evaluate()
	return Result{Allowed: decision.Allowed, Reason: decision.Reason, Warning: decision.Warning}
}

func (e *Evaluator) evaluateInvariants(input Input) Result {
	rule := policy.NewInvariantsRule(&e.cfg.Invariants)
	paths := ExtractPaths(input.ToolName, input.ToolInput)

	// Get content for content-based checks
	content := ""
	if c, ok := input.ToolInput["content"].(string); ok {
		content = c
	}

	for _, p := range paths {
		decision := rule.Evaluate(input.ToolName, p, content)
		if !decision.Allowed {
			return Result{Allowed: false, Reason: decision.Reason}
		}
	}
	return Result{Allowed: true}
}

func (e *Evaluator) evaluateHooks(input Input) Result {
	paths := ExtractPaths(input.ToolName, input.ToolInput)

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	hookInput := HookInput{
		ToolName:   input.ToolName,
		ToolInput:  input.ToolInput,
		Paths:      paths,
		WorkingDir: cwd,
	}

	var warnings []string

	for i := range e.cfg.Hooks {
		hookCfg := &e.cfg.Hooks[i]

		if !e.hookMatcher.Matches(hookCfg, input.ToolName, paths) {
			continue
		}

		result := e.hookExec.Execute(hookCfg, hookInput)

		if !result.Allowed {
			return Result{
				Allowed: false,
				Reason:  hookCfg.Name + ": " + result.Reason,
			}
		}

		if result.Warning != "" {
			warnings = append(warnings, hookCfg.Name+": "+result.Warning)
		}
	}

	if len(warnings) > 0 {
		return Result{Allowed: true, Warning: strings.Join(warnings, "; ")}
	}

	return Result{Allowed: true}
}

func (e *Evaluator) isToolBlocked(tool string) bool {
	for _, t := range e.cfg.Tools.Block {
		if strings.EqualFold(t, tool) {
			return true
		}
	}
	return false
}

func (e *Evaluator) isToolAllowed(tool string) bool {
	if len(e.cfg.Tools.Allow) == 0 {
		return true
	}
	for _, t := range e.cfg.Tools.Allow {
		if strings.EqualFold(t, tool) {
			return true
		}
	}
	return false
}

func (e *Evaluator) isCommandBlocked(cmd string) string {
	for _, pattern := range e.cfg.Commands.Block {
		if strings.Contains(cmd, pattern) {
			return pattern
		}
	}
	return ""
}

var filesystemTools = map[string]bool{
	"Bash":  true,
	"Read":  true,
	"Write": true,
	"Edit":  true,
	"Glob":  true,
	"Grep":  true,
}

func isFilesystemTool(tool string) bool {
	return filesystemTools[tool]
}

func isModificationTool(tool string) bool {
	switch tool {
	case "Write", "Edit", "NotebookEdit":
		return true
	}
	return false
}
