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
		if trimmed == "}" || strings.HasSuffix(trimmed, "}") {
			indentLevel = max(0, indentLevel-1)
		}
		indent := strings.Repeat("    ", indentLevel)
		formattedLine := formatStylishLine(trimmed, indent)
		result = append(result, formattedLine)
		if trimmed == "{" || strings.HasSuffix(trimmed, "{") {
			indentLevel++
		}
	}

	return strings.Join(result, "\n")
}

func formatStylishLine(line string, indent string) string {
	cleanedLine := strings.TrimSpace(line)
	cleanedLine = strings.ReplaceAll(cleanedLine, "<nil>", "null")
	if strings.HasSuffix(cleanedLine, ": ") {
		cleanedLine = strings.TrimSuffix(cleanedLine, " ")
	}
	switch {
	case cleanedLine == "{" || cleanedLine == "}":
		return indent + cleanedLine

	case strings.HasPrefix(cleanedLine, "+ "):
		content := strings.TrimPrefix(cleanedLine, "+ ")
		if strings.Contains(content, ": {") {
			parts := strings.SplitN(content, ": {", 2)
			return indent + "+ " + parts[0] + ": {" + parts[1]
		}
		return indent + "+ " + content

	case strings.HasPrefix(cleanedLine, "- "):
		content := strings.TrimPrefix(cleanedLine, "- ")
		if strings.Contains(content, ": {") {
			parts := strings.SplitN(content, ": {", 2)
			return indent + "- " + parts[0] + ": {" + parts[1]
		}
		return indent + "- " + content

	default:
		content := strings.TrimPrefix(cleanedLine, "+ ")
		content = strings.TrimPrefix(content, "- ")
		return indent + content
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
