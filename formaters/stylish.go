package code

import (
	"strings"
)

func FormatDiffOutputStylish(diffOutput string) string {
	trimmedOutput := strings.TrimSpace(diffOutput)
	if strings.HasPrefix(trimmedOutput, "{\n") && strings.HasSuffix(trimmedOutput, "\n}") {
		return trimmedOutput
	}
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
		if trimmed == "}" || strings.HasSuffix(trimmed, "}") && !strings.Contains(trimmed, "{") {
			if depth > 0 {
				depth--
			}
			baseIndent = strings.Repeat("    ", depth)
			result = append(result, baseIndent+"}")
			continue
		}
		if strings.Contains(trimmed, "{") {
			parts := strings.SplitN(trimmed, "{", 2)
			keyPart := strings.TrimSpace(parts[0])

			var prefix string
			if strings.HasPrefix(keyPart, "- ") {
				prefix = "  - "
				keyPart = strings.TrimPrefix(keyPart, "- ")
			} else if strings.HasPrefix(keyPart, "+ ") {
				prefix = "  + "
				keyPart = strings.TrimPrefix(keyPart, "+ ")
			} else {
				prefix = "    "
			}

			result = append(result, baseIndent+prefix+keyPart+": {")
			depth++
			continue
		}
		var prefix string
		var content string

		if strings.HasPrefix(trimmed, "- ") {
			prefix = "  - "
			content = strings.TrimPrefix(trimmed, "- ")
		} else if strings.HasPrefix(trimmed, "+ ") {
			prefix = "  + "
			content = strings.TrimPrefix(trimmed, "+ ")
		} else {
			prefix = "    "
			content = trimmed
		}

		result = append(result, baseIndent+prefix+content)
	}
	res := "{\n" + strings.Join(result, "\n") + "\n}"
	expected := strings.ReplaceAll(res, "  +", "      +")
	expected = strings.ReplaceAll(expected, "  -", "      -")
	expected = strings.ReplaceAll(expected, "    {", "        {")
	expected = strings.ReplaceAll(expected, "    }", "        }")
	lines1 := strings.Split(expected, "\n")
	var otvet strings.Builder
	for _, line := range lines1 {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		spaceCount := 0
		for _, char := range line {
			if char != ' ' {
				break
			}
			spaceCount++
		}
		newIndent := strings.Repeat("    ", spaceCount/2)
		otvet.WriteString(newIndent + trimmed + "\n")
	}
	return strings.TrimSpace(otvet.String())
}
