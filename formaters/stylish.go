package code

import (
	"fmt"
	"strings"
)

func FormatDiffOutputStylish(diffOutput string) string {
	lines := strings.Split(diffOutput, "\n")
	var result []string
	depth := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "}" {
			if depth > 0 {
				depth--
			}
			indent := strings.Repeat("    ", depth)
			result = append(result, fmt.Sprintf("%s}", indent))
			continue
		}
		if trimmed == "" {
			continue
		}
		indent := strings.Repeat("    ", depth)
		if trimmed == "{" {
			result = append(result, fmt.Sprintf("%s{", indent))
			depth++
			continue
		}
		contentIndent := strings.Repeat("    ", depth)
		switch {
		case strings.HasPrefix(trimmed, "- "):
			content := strings.TrimPrefix(trimmed, "- ")
			result = append(result, fmt.Sprintf("%s  - %s", contentIndent, content))
		case strings.HasPrefix(trimmed, "+ "):
			content := strings.TrimPrefix(trimmed, "+ ")
			result = append(result, fmt.Sprintf("%s  + %s", contentIndent, content))
		case strings.Contains(trimmed, ":"):
			result = append(result, fmt.Sprintf("%s    %s", contentIndent, trimmed))
		default:
			result = append(result, fmt.Sprintf("%s%s", contentIndent, trimmed))
		}
	}
	formattedResult := strings.Join(result, "\n")
	formattedResult = strings.ReplaceAll(formattedResult, "<nil>", "null")

	return formattedResult
}
