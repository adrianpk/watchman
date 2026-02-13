package types

type HookInput struct {
	ToolName   string         `json:"tool_name"`
	ToolInput  map[string]any `json:"tool_input"`
	Paths      []string       `json:"paths"`
	WorkingDir string         `json:"working_dir"`
}

type HookOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
	Warning  string `json:"warning,omitempty"`
}

type EvalRequest struct {
	ToolName  string
	FilePath  string
	Content   string
	Standards string
}

type EvalResult struct {
	Decision   string
	Reason     string
	Warning    string
	Violations []string
}
