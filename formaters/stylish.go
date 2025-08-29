package code

import (
	"strings"
)

func FormatDiffOutputStylish(diffOutput string) string {

	trimmedOutput := strings.TrimSpace(diffOutput)
	if strings.HasPrefix(trimmedOutput, "{") && strings.HasSuffix(trimmedOutput, "}") {
		trimmedOutput = strings.TrimSpace(trimmedOutput[1 : len(trimmedOutput)-1])
	}

	lines := strings.Split(trimmedOutput, "\n")
	var result []string
	depth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		baseIndent := strings.Repeat("    ", depth)
		if trimmed == "}" {
			if depth > 0 {
				depth--
			}
			baseIndent = strings.Repeat("    ", depth)
			result = append(result, baseIndent+"}")
			continue
		}
		if strings.HasSuffix(trimmed, "{") {
			var prefix string
			var content string

			if strings.HasPrefix(trimmed, "- ") {
				prefix = "  - "
				content = strings.TrimSuffix(strings.TrimPrefix(trimmed, "- "), " {")
			} else if strings.HasPrefix(trimmed, "+ ") {
				prefix = "  + "
				content = strings.TrimSuffix(strings.TrimPrefix(trimmed, "+ "), " {")
			} else {
				prefix = "    "
				content = strings.TrimSuffix(trimmed, " {")
			}

			result = append(result, baseIndent+prefix+content+": {")
			depth++
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			result = append(result, baseIndent+"  - "+strings.TrimPrefix(trimmed, "- "))
		} else if strings.HasPrefix(trimmed, "+ ") {
			result = append(result, baseIndent+"  + "+strings.TrimPrefix(trimmed, "+ "))
		} else {
			result = append(result, baseIndent+"    "+trimmed)
		}
	}
	return "{\n" + strings.Join(result, "\n") + "\n}"
}
