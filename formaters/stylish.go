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
		if strings.HasPrefix(trimmed, "}") || strings.HasSuffix(trimmed, "}") {
			indentLevel--
		}
		indent := strings.Repeat("    ", max(0, indentLevel))
		formattedLine := formatStylishLine(trimmed, indent)
		result = append(result, formattedLine)
		if strings.HasSuffix(trimmed, "{") {
			indentLevel++
		}
	}

	return strings.Join(result, "\n")
}

func formatStylishLine(line string, indent string) string {
	cleanedLine := strings.TrimSpace(line)
	cleanedLine = strings.ReplaceAll(cleanedLine, "<nil>", "null")

	if cleanedLine == "{" || cleanedLine == "}" {
		return indent + cleanedLine
	} else if strings.HasPrefix(cleanedLine, "+ ") {
		content := strings.TrimPrefix(cleanedLine, "+ ")
		if strings.Contains(content, ": {") || strings.Contains(content, ": +") {
			content = strings.ReplaceAll(content, ": +", ": ")
			content = strings.ReplaceAll(content, "+ ", "")
		}
		return indent + "+ " + content
	} else if strings.HasPrefix(cleanedLine, "- ") {
		content := strings.TrimPrefix(cleanedLine, "- ")
		if strings.Contains(content, ": {") || strings.Contains(content, ": -") {
			content = strings.ReplaceAll(content, ": -", ": ")
			content = strings.ReplaceAll(content, "- ", "")
		}
		return indent + "- " + content
	} else {
		cleanedLine = strings.ReplaceAll(cleanedLine, "+ ", "")
		cleanedLine = strings.ReplaceAll(cleanedLine, "- ", "")
		return indent + cleanedLine
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
