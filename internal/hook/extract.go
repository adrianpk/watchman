package hook

import "github.com/adrianpk/watchman/internal/parser"

// ExtractPaths extracts filesystem paths from tool input.
func ExtractPaths(toolName string, toolInput map[string]interface{}) []string {
	switch toolName {
	case "Bash":
		return extractBashPaths(toolInput)
	case "Read", "Write", "Edit":
		return extractFilePath(toolInput)
	case "Glob":
		return extractGlobPaths(toolInput)
	case "Grep":
		return extractGrepPaths(toolInput)
	}
	return nil
}

func extractBashPaths(toolInput map[string]interface{}) []string {
	cmdStr, ok := toolInput["command"].(string)
	if !ok {
		return nil
	}
	cmd := parser.Parse(cmdStr)
	var paths []string
	paths = append(paths, cmd.Args...)
	for _, v := range cmd.Flags {
		if v != "" {
			paths = append(paths, v)
		}
	}
	for _, v := range cmd.Env {
		paths = append(paths, v)
	}
	return paths
}

func extractFilePath(toolInput map[string]interface{}) []string {
	if fp, ok := toolInput["file_path"].(string); ok {
		return []string{fp}
	}
	return nil
}

func extractGlobPaths(toolInput map[string]interface{}) []string {
	var paths []string
	if p, ok := toolInput["path"].(string); ok {
		paths = append(paths, p)
	}
	if pattern, ok := toolInput["pattern"].(string); ok {
		paths = append(paths, pattern)
	}
	return paths
}

func extractGrepPaths(toolInput map[string]interface{}) []string {
	if p, ok := toolInput["path"].(string); ok {
		return []string{p}
	}
	return nil
}
