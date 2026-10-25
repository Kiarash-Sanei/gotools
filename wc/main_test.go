package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failingReader struct {
	err error
}

func (reader failingReader) Read([]byte) (int, error) {
	return 0, reader.err
}

func TestCount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Counts
	}{
		{name: "empty", input: "", want: Counts{}},
		{name: "newline terminated", input: "one two\nthree\n", want: Counts{Lines: 2, Words: 3, Bytes: 14}},
		{name: "no trailing newline", input: "one two", want: Counts{Words: 2, Bytes: 7}},
		{name: "unicode bytes", input: "café\n", want: Counts{Lines: 1, Words: 1, Bytes: 6}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := count(strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("count returned error: %v", err)
			}
			if got != test.want {
				t.Errorf("count = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestCountReadError(t *testing.T) {
	wantErr := errors.New("read failed")
	_, err := count(failingReader{err: wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("count error = %v, want %v", err, wantErr)
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
			if exitCode := run(test.args, &stdout, &stderr, strings.NewReader("")); exitCode != 0 {
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

	if exitCode := run([]string{missingFile}, &stdout, &stderr, strings.NewReader("")); exitCode != 1 {
		t.Fatalf("run exit code = %d, want 1", exitCode)
	}
	if stdout.Len() != 0 {
		t.Errorf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "cannot read "+missingFile) {
		t.Errorf("stderr = %q, want missing file error", stderr.String())
	}
}

func TestRunStdin(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if exitCode := run([]string{"-l"}, &stdout, &stderr, strings.NewReader("a b\nc\n")); exitCode != 0 {
		t.Fatalf("run exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	if got, want := stdout.String(), "2\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Errorf("unexpected stderr: %q", stderr.String())
	}
}

func TestRunMultipleFiles(t *testing.T) {
	directory := t.TempDir()
	firstFile := filepath.Join(directory, "first.txt")
	secondFile := filepath.Join(directory, "second.txt")
	if err := os.WriteFile(firstFile, []byte("one\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(secondFile, []byte("two three\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	args := []string{"-l", "-w", firstFile, secondFile}
	if exitCode := run(args, &stdout, &stderr, strings.NewReader("")); exitCode != 0 {
		t.Fatalf("run exit code = %d, stderr = %q", exitCode, stderr.String())
	}

	want := "1 1 " + firstFile + "\n1 2 " + secondFile + "\n2 3 total\n"
	if got := stdout.String(); got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRunUnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if exitCode := run([]string{"-x"}, &stdout, &stderr, strings.NewReader("")); exitCode != 2 {
		t.Fatalf("run exit code = %d, want 2", exitCode)
	}
	if stdout.Len() != 0 {
		t.Errorf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "-x") {
		t.Errorf("stderr = %q, want it to mention -x", stderr.String())
	}
}
