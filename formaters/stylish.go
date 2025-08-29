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
			if indentLevel < 0 {
				indentLevel = 0
			}
		}
		formattedLine := formatStylishLine(trimmed, indentLevel)
		result = append(result, formattedLine)
		if strings.Contains(trimmed, "{") && !strings.Contains(trimmed, "}") {
			indentLevel++
			if indentLevel == 7 {
				indentLevel = 6
			}
		}
	}
	return strings.Join(result, "\n")
}

func formatStylishLine(line string, indentLevel int) string {
	baseIndent := strings.Repeat("    ", indentLevel)
	cleanedLine := strings.TrimSpace(line)

	if cleanedLine == "{" || cleanedLine == "}" {
		return baseIndent + cleanedLine
	} else if strings.HasPrefix(cleanedLine, "+ ") {
		content := strings.TrimPrefix(cleanedLine, "+ ")
		return baseIndent + "+ " + content
	} else if strings.HasPrefix(cleanedLine, "- ") {
		content := strings.TrimPrefix(cleanedLine, "- ")
		return baseIndent + "- " + content
	} else {
		return baseIndent + cleanedLine
	}
}
