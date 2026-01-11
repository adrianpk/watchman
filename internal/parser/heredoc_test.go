package parser

import "testing"

func TestStripHeredocs(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want string
	}{
		{
			name: "no heredoc",
			cmd:  "echo hello",
			want: "echo hello",
		},
		{
			name: "basic heredoc",
			cmd: `cat <<EOF
hello world
EOF`,
			want: "cat <<EOF",
		},
		{
			name: "heredoc with dash",
			cmd: `cat <<-EOF
	indented content
EOF`,
			want: "cat <<-EOF",
		},
		{
			name: "heredoc with quoted delimiter",
			cmd: `cat <<'EOF'
$variable stays literal
EOF`,
			want: "cat <<'EOF'",
		},
		{
			name: "heredoc with double quoted delimiter",
			cmd: `cat <<"EOF"
content here
EOF`,
			want: `cat <<"EOF"`,
		},
		{
			name: "heredoc with absolute path inside",
			cmd: `cat <<EOF
/etc/passwd
/home/user/.bashrc
EOF`,
			want: "cat <<EOF",
		},
		{
			name: "multiple heredocs",
			cmd: `cat <<EOF1
first
EOF1
cat <<EOF2
second
EOF2`,
			want: `cat <<EOF1
cat <<EOF2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripHeredocs(tt.cmd)
			if got != tt.want {
				t.Errorf("stripHeredocs() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseWithHeredoc(t *testing.T) {
	// Test that heredoc content does not appear in parsed Args
	// This is the main purpose: prevent paths inside heredocs from being detected as violations

	cmd := `cat <<EOF
/etc/passwd
/home/user/.bashrc
EOF`

	parsed := Parse(cmd)

	if parsed.Program != "cat" {
		t.Errorf("Program = %q, want %q", parsed.Program, "cat")
	}

	// The key test: absolute paths from heredoc content should NOT appear in Args
	for _, arg := range parsed.Args {
		if arg == "/etc/passwd" || arg == "/home/user/.bashrc" {
			t.Errorf("Args contains heredoc content %q which should have been stripped", arg)
		}
	}
}

func TestParseWithHeredocPreservesRealArgs(t *testing.T) {
	// Ensure that real arguments before the heredoc are preserved
	cmd := `cat /tmp/real-file.txt <<EOF
/etc/passwd
EOF`

	parsed := Parse(cmd)

	if parsed.Program != "cat" {
		t.Errorf("Program = %q, want %q", parsed.Program, "cat")
	}

	// Real arg should be present
	found := false
	for _, arg := range parsed.Args {
		if arg == "/tmp/real-file.txt" {
			found = true
		}
		if arg == "/etc/passwd" {
			t.Errorf("Args contains heredoc content /etc/passwd which should have been stripped")
		}
	}

	if !found {
		t.Errorf("Args should contain /tmp/real-file.txt, got %v", parsed.Args)
	}
}
