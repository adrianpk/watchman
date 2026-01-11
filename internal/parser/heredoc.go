package parser

import "regexp"

// heredocPattern matches heredoc start: << or <<- followed by delimiter
var heredocPattern = regexp.MustCompile(`<<-?\s*['"]?(\w+)['"]?`)

// stripHeredocs removes heredoc content from a command string.
// This prevents heredoc content from being parsed as command arguments.
func stripHeredocs(cmd string) string {
	matches := heredocPattern.FindAllStringSubmatchIndex(cmd, -1)
	if len(matches) == 0 {
		return cmd
	}

	result := cmd
	// Process matches in reverse to preserve indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		fullMatchEnd := match[1]
		delimiter := cmd[match[2]:match[3]]

		// Find the closing delimiter (must be on its own line)
		closingPattern := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(delimiter) + `$`)
		remaining := cmd[fullMatchEnd:]
		closingMatch := closingPattern.FindStringIndex(remaining)

		if closingMatch != nil {
			// Remove from after << DELIM to end of closing delimiter
			contentEnd := fullMatchEnd + closingMatch[1]
			result = result[:fullMatchEnd] + result[contentEnd:]
		}
	}

	return result
}
