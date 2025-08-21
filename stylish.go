package code
import("strings"
"fmt")
func FormatDiffOutputStylish(diffOutput string) string {
	lines := strings.Split(diffOutput, "\n")
	var result []string
	depth := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "}") {
			depth--
		}
		if trimmed == "" {
			continue
		}
		dotIndent := strings.Repeat(".", depth*2)	
		if strings.HasPrefix(trimmed, "- ") {
			content := strings.TrimPrefix(trimmed, "- ")
			result = append(result, fmt.Sprintf("%s- %s", dotIndent, content))
		} else if strings.HasPrefix(trimmed, "+ ") {
			content := strings.TrimPrefix(trimmed, "+ ")
			result = append(result, fmt.Sprintf("%s+ %s", dotIndent, content))
		} else if strings.Contains(trimmed, ":") {
			result = append(result, fmt.Sprintf("%s%s", dotIndent, trimmed))
		} else if trimmed == "{" {
			result = append(result, fmt.Sprintf("%s{", dotIndent))
		} else if trimmed == "}" {
			result = append(result, fmt.Sprintf("%s}", dotIndent))
		}
		if strings.Contains(trimmed, "{") {
			depth++
		}
	}
	return strings.Join(result, "\n")
}