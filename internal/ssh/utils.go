package ssh

import (
	"bufio"
	"strings"
)

// indentOutput indents the output for better readability
func indentOutput(output string) string {
	if output == "" {
		return ""
	}

	var indented strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		indented.WriteString("    " + scanner.Text() + "\n")
	}
	return indented.String()
}
