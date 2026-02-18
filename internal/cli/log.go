package cli

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"
)

const logFile = "/tmp/watchman.log"

type logEntry struct {
	Timestamp string `json:"ts"`
	ID        string `json:"id"`
	Decision  string `json:"decision"`
	Tool      string `json:"tool"`
	CWD       string `json:"cwd"`
	Reason    string `json:"reason"`
	Command   string `json:"cmd,omitempty"`
	Path      string `json:"path,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}

func RunLog(args []string) error {
	n := 1
	if len(args) > 0 {
		var err error
		n, err = strconv.Atoi(args[0])
		if err != nil || n < 1 {
			n = 1
		}
	}

	entries, err := readLastEntries(n)
	if err != nil {
		return err
	}

	return writeEntries(os.Stdout, entries)
}

func readLastEntries(n int) ([]logEntry, error) {
	f, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return nil, nil
	}

	start := len(lines) - n
	if start < 0 {
		start = 0
	}

	var entries []logEntry
	for _, line := range lines[start:] {
		var entry logEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func writeEntries(w io.Writer, entries []logEntry) error {
	for i, e := range entries {
		if i > 0 {
			io.WriteString(w, "\n")
		}
		io.WriteString(w, "┌─ "+e.Timestamp+" ["+e.ID+"] "+e.Decision+"\n")
		io.WriteString(w, "│  Tool:   "+e.Tool+"\n")
		io.WriteString(w, "│  CWD:    "+e.CWD+"\n")
		io.WriteString(w, "│  Reason: "+e.Reason+"\n")
		if e.Command != "" {
			cmd := e.Command
			if len(cmd) > 80 {
				cmd = cmd[:80] + "..."
			}
			io.WriteString(w, "│  Cmd:    "+cmd+"\n")
		}
		if e.Path != "" {
			io.WriteString(w, "│  Path:   "+e.Path+"\n")
		}
		if e.Pattern != "" {
			io.WriteString(w, "│  Pattern: "+e.Pattern+"\n")
		}
		io.WriteString(w, "└─\n")
	}
	return nil
}
