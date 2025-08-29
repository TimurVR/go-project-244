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
			depth--
		}
		if trimmed == "" {
			continue
		}
		indent := strings.Repeat("    ", depth)
		if strings.HasPrefix(trimmed, "- ") {
			content := strings.TrimPrefix(trimmed, "- ")
			result = append(result, fmt.Sprintf("%s  - %s", indent, content))
		} else if strings.HasPrefix(trimmed, "+ ") {
			content := strings.TrimPrefix(trimmed, "+ ")
			result = append(result, fmt.Sprintf("%s  + %s", indent, content))
		} else if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "- ") && !strings.HasPrefix(trimmed, "+ ") {
			result = append(result, fmt.Sprintf("%s    %s", indent, trimmed))
		} else if trimmed == "{" {
			result = append(result, fmt.Sprintf("%s{", indent))
		} else if trimmed == "}" {
			result = append(result, fmt.Sprintf("%s}", indent))
		}
		if trimmed == "{" {
			depth++
		}
	}
	return strings.Join(result, "\n")
}
