package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " the Winner is here",
			expected: []string{"the", "winner", "is", "here"},
		},
		{
			input:    "hello World ",
			expected: []string{"hello", "world"},
		},
	}
	for _, c := range cases {
		actual := CleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Fatalf(
				"expected %d words, got %d (%v)",
				len(c.expected),
				len(actual),
				actual,
			)
		}

		for _, c := range cases {
			actual := CleanInput(c.input)
			// Check the length of the actual slice against the expected slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			for i := range actual {
				word := actual[i]
				expectedWord := c.expected[i]
				// Check each word in the slice
				// if they don't match, use t.Errorf to print an error message
				// and fail the test
				if word != expectedWord {
					t.Errorf(word, expectedWord)
				}
			}
		}
	}
}

func TestCommandExit(t *testing.T) {
	// If this env var is set, we're in the child process
	if os.Getenv("TEST_EXIT") == "1" {
		cfg := &config{}
		CommandExit(cfg)
		return
	}

	// Re-run this test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestCommandExit")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")

	err := cmd.Run()

	// We EXPECT an error because os.Exit was called
	if err == nil {
		t.Fatal("expected process to exit, but it did not")
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}

	if exitErr.ExitCode() != 0 {
		t.Fatalf("expected exit code 0, got %d", exitErr.ExitCode())
	}
}

func TestHelpCommand(t *testing.T) {
	cfg := &config{}
	commands := buildCommands(cfg)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commands["help"].callback(cfg)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "help") {
		t.Errorf("expected output to contain help command")
	}

	if !strings.Contains(output, "exit") {
		t.Errorf("expected output to contain exit command")
	}
}
