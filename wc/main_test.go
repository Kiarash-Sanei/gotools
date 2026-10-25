package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		lines     int
		words     int
		byteCount int
	}{
		{name: "empty", input: "", lines: 0, words: 0, byteCount: 0},
		{name: "newline terminated", input: "one two\nthree\n", lines: 2, words: 3, byteCount: 14},
		{name: "no trailing newline", input: "one two", lines: 0, words: 2, byteCount: 7},
		{name: "unicode bytes", input: "café\n", lines: 1, words: 1, byteCount: 6},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lines, words, byteCount, err := count(strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("count returned error: %v", err)
			}
			if lines != test.lines || words != test.words || byteCount != test.byteCount {
				t.Errorf("count = (%d, %d, %d), want (%d, %d, %d)", lines, words, byteCount, test.lines, test.words, test.byteCount)
			}
		})
	}
}

func TestRun(t *testing.T) {
	fileName := filepath.Join(t.TempDir(), "input file.txt")
	if err := os.WriteFile(fileName, []byte("one two\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{name: "default counts", args: []string{fileName}, wantOutput: "1 2 8 " + fileName + "\n"},
		{name: "selected count", args: []string{"-w", fileName}, wantOutput: "2 " + fileName + "\n"},
		{name: "combined counts", args: []string{"-c", "-l", fileName}, wantOutput: "1 8 " + fileName + "\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if exitCode := run(test.args, &stdout, &stderr); exitCode != 0 {
				t.Fatalf("run exit code = %d, stderr = %q", exitCode, stderr.String())
			}
			if stdout.String() != test.wantOutput {
				t.Errorf("stdout = %q, want %q", stdout.String(), test.wantOutput)
			}
			if stderr.Len() != 0 {
				t.Errorf("unexpected stderr: %q", stderr.String())
			}
		})
	}
}

func TestRunMissingFile(t *testing.T) {
	var stdout, stderr bytes.Buffer
	missingFile := filepath.Join(t.TempDir(), "missing.txt")

	if exitCode := run([]string{missingFile}, &stdout, &stderr); exitCode != 1 {
		t.Fatalf("run exit code = %d, want 1", exitCode)
	}
	if stdout.Len() != 0 {
		t.Errorf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "cannot read "+missingFile) {
		t.Errorf("stderr = %q, want missing file error", stderr.String())
	}
}
