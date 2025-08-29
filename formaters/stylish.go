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
		if strings.HasPrefix(trimmed, "}") || (strings.Contains(trimmed, "}") && !strings.Contains(trimmed, "{")) {
			indentLevel = max(0, indentLevel-1)
		}
		indent := strings.Repeat("    ", indentLevel)
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
	if strings.Contains(cleanedLine, ": {") {
		if strings.HasPrefix(cleanedLine, "+ ") {
			content := strings.TrimPrefix(cleanedLine, "+ ")
			return indent + "+ " + removeInnerPrefixes(content)
		} else if strings.HasPrefix(cleanedLine, "- ") {
			content := strings.TrimPrefix(cleanedLine, "- ")
			return indent + "- " + removeInnerPrefixes(content)
		} else {
			return indent + removeInnerPrefixes(cleanedLine)
		}
	}
	if cleanedLine == "{" || cleanedLine == "}" {
		return indent + cleanedLine
	} else if strings.HasPrefix(cleanedLine, "+ ") {
		content := strings.TrimPrefix(cleanedLine, "+ ")
		return indent + "+ " + content
	} else if strings.HasPrefix(cleanedLine, "- ") {
		content := strings.TrimPrefix(cleanedLine, "- ")
		return indent + "- " + content
	} else {
		return indent + cleanedLine
	}
}
func removeInnerPrefixes(content string) string {
	content = strings.ReplaceAll(content, ": + ", ": ")
	content = strings.ReplaceAll(content, ": - ", ": ")
	content = strings.ReplaceAll(content, "+ ", "")
	content = strings.ReplaceAll(content, "- ", "")
	return content
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
