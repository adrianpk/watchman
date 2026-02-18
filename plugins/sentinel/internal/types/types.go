package types

type HookInput struct {
	ToolName   string         `json:"tool_name"`
	ToolInput  map[string]any `json:"tool_input"`
	Paths      []string       `json:"paths"`
	WorkingDir string         `json:"working_dir"`
}

type HookOutput struct {
	Decision   string      `json:"decision"`
	Reason     string      `json:"reason,omitempty"`
	Warning    string      `json:"warning,omitempty"`
	Violations []Violation `json:"violations,omitempty"`
}

type EvalRequest struct {
	ToolName  string
	FilePath  string
	Content   string
	Standards string
}

type Violation struct {
	Rule       string  `json:"rule"`
	File       string  `json:"file"`
	Lines      string  `json:"lines"`
	Detail     string  `json:"detail"`
	Confidence float64 `json:"confidence"`
}

type EvalResult struct {
	Decision   string
	Reason     string
	Warning    string
	Violations []Violation
	Confidence float64
}
