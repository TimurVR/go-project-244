package code

import (
	"strings"
)

func FormatDiffOutputStylish(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	indentLevel := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.Contains(trimmed, "}") && !strings.Contains(trimmed, "{") {
			indentLevel--
		}
		indent := strings.Repeat("    ", indentLevel)
		formattedLine := indent + formatLineContent(trimmed)
		result = append(result, formattedLine)
		if strings.Contains(trimmed, "{") && !strings.Contains(trimmed, "}") {
			indentLevel++
		}
	}
	return strings.Join(result, "\n")
}
func formatLineContent(line string) string {
	formatted := line
	formatted = strings.ReplaceAll(formatted, "{ ", "{")
	formatted = strings.ReplaceAll(formatted, " }", "}")
	formatted = strings.ReplaceAll(formatted, ": {", ": {")
	formatted = strings.ReplaceAll(formatted, ": }", ": }")

	return formatted
}
