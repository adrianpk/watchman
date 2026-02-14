// Package git provides utilities for interacting with git repositories.
package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// StagedDiff represents the staged changes in a git repository.
type StagedDiff struct {
	Files   []string
	Content string
	Size    int
}

// GetStagedDiff retrieves the staged diff from the current git repository.
func GetStagedDiff(workDir string) (*StagedDiff, error) {
	files, err := getStagedFiles(workDir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return &StagedDiff{}, nil
	}

	content, err := getStagedContent(workDir)
	if err != nil {
		return nil, err
	}

	return &StagedDiff{
		Files:   files,
		Content: content,
		Size:    len(content),
	}, nil
}

// getStagedFiles returns the list of staged file paths.
func getStagedFiles(workDir string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = workDir

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return nil, nil
	}

	return strings.Split(output, "\n"), nil
}

// getStagedContent returns the full staged diff content.
func getStagedContent(workDir string) (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	cmd.Dir = workDir

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return stdout.String(), nil
}

// IsCommitCommand checks if a bash command is a git commit.
func IsCommitCommand(command string) bool {
	return strings.Contains(command, "git") && strings.Contains(command, "commit")
}

// IsAddCommand checks if a bash command is a git add.
func IsAddCommand(command string) bool {
	return strings.Contains(command, "git") && strings.Contains(command, "add")
}

// ExtractAddFiles extracts file paths from a git add command.
func ExtractAddFiles(command string) []string {
	var files []string
	parts := strings.Fields(command)

	skipNext := false
	inGitAdd := false

	for i, part := range parts {
		if skipNext {
			skipNext = false
			continue
		}

		// Skip flags and their values
		if strings.HasPrefix(part, "-") {
			if part == "-C" || part == "--git-dir" || part == "--work-tree" {
				skipNext = true
			}
			continue
		}

		// Detect git add
		if part == "git" && i+1 < len(parts) {
			continue
		}
		if part == "add" {
			inGitAdd = true
			continue
		}

		// Collect files after "add"
		if inGitAdd && part != "&&" && part != ";" && part != "||" {
			// Stop at command separators
			if part == "&&" || part == ";" || part == "||" {
				break
			}
			files = append(files, part)
		}
	}

	return files
}

// ReadFiles reads the content of multiple files and concatenates them.
func ReadFiles(workDir string, files []string) (string, error) {
	var content strings.Builder

	for _, file := range files {
		path := file
		if !strings.HasPrefix(file, "/") {
			path = workDir + "/" + file
		}

		data, err := os.ReadFile(path)
		if err != nil {
			continue // Skip files that can't be read
		}

		content.WriteString("=== " + file + " ===\n")
		content.Write(data)
		content.WriteString("\n")
	}

	return content.String(), nil
}
