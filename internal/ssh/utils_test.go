package ssh

import (
	"testing"
)

func TestIndentOutput_Empty(t *testing.T) {
	result := indentOutput("")
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}

func TestIndentOutput_SingleLine(t *testing.T) {
	input := "hello world"
	expected := "    hello world\n"
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}

func TestIndentOutput_MultipleLines(t *testing.T) {
	input := "line1\nline2\nline3"
	expected := "    line1\n    line2\n    line3\n"
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}

func TestIndentOutput_WithEmptyLines(t *testing.T) {
	input := "line1\n\nline3"
	expected := "    line1\n    \n    line3\n"
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}

func TestIndentOutput_WithTrailingNewline(t *testing.T) {
	input := "line1\nline2\n"
	expected := "    line1\n    line2\n" // The function doesn't add an extra newline
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}
