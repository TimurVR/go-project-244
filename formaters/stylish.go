package code

import (
	"strings"
)

func FormatDiffOutputStylish(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	result = append(result, "{")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		var lineType string
		var content string

		if strings.HasPrefix(trimmed, "- ") {
			lineType = "removed"
			content = strings.TrimPrefix(trimmed, "- ")
		} else if strings.HasPrefix(trimmed, "+ ") {
			lineType = "added"
			content = strings.TrimPrefix(trimmed, "+ ")
		} else {
			lineType = "unchanged"
			content = trimmed
		}
		var formattedLine string

		switch lineType {
		case "added":
			formattedLine = "  + " + content
		case "removed":
			formattedLine = "  - " + content
		case "unchanged":
			formattedLine = "    " + content
		}

		result = append(result, formattedLine)
	}
	result = append(result, "}")

	return strings.Join(result, "\n")
}
