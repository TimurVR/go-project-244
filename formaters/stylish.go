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
		trimmed = strings.ReplaceAll(trimmed, "<nil>", "null")

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
			var prefix, content string

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
			content = strings.TrimSuffix(content, ":")
			result = append(result, baseIndent+prefix+content+": {")
			depth++
			continue
		}
		if strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			keyPart := strings.TrimSpace(parts[0])
			valuePart := strings.TrimSpace(parts[1])
			var prefix string
			var key string
			if strings.HasPrefix(keyPart, "- ") {
				prefix = "  - "
				key = strings.TrimPrefix(keyPart, "- ")
			} else if strings.HasPrefix(keyPart, "+ ") {
				prefix = "  + "
				key = strings.TrimPrefix(keyPart, "+ ")
			} else {
				prefix = "    "
				key = keyPart
			}

			result = append(result, baseIndent+prefix+key+": "+valuePart)
		} else {
			if strings.HasPrefix(trimmed, "- ") {
				result = append(result, baseIndent+"  - "+strings.TrimPrefix(trimmed, "- "))
			} else if strings.HasPrefix(trimmed, "+ ") {
				result = append(result, baseIndent+"  + "+strings.TrimPrefix(trimmed, "+ "))
			} else {
				result = append(result, baseIndent+"    "+trimmed)
			}
		}
	}
	return "{\n" + strings.Join(result, "\n") + "\n}"
}
