package code
import("strings"
"strconv"
"fmt"
"sort")

func FormatDiffOutput(diffOutput string) string {
	lines := strings.Split(diffOutput, "\n")
	var result []string
	var currentPath []string
	changes := make(map[string]ChangeInfo)
	stack := []map[string]bool{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if trimmed == "{" {
			stack = append(stack, make(map[string]bool))
			continue
		}

		if trimmed == "}" {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
			continue
		}
		if strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			operation := "unchanged"
			if strings.HasPrefix(key, "+ ") {
				operation = "added"
				key = strings.TrimPrefix(key, "+ ")
			} else if strings.HasPrefix(key, "- ") {
				operation = "removed"
				key = strings.TrimPrefix(key, "- ")
			}

			if value == "{" {
				fullPath := buildPath(append(currentPath, key))
				
				if operation == "added" {
					changes[fullPath] = ChangeInfo{
						Operation: "added",
						Value:     "[complex value]",
					}
				} else if operation == "removed" {
					changes[fullPath] = ChangeInfo{
						Operation: "removed",
						Value:     "",
					}
				}
				
				currentPath = append(currentPath, key)
				stack = append(stack, make(map[string]bool))
				continue
			}
			fullPath := buildPath(append(currentPath, key))
			
			if operation == "added" || operation == "removed" {
				changes[fullPath] = ChangeInfo{
					Operation: operation,
					Value:     value,
				}
			}
		}
	}
	processed := make(map[string]bool)
	for path, change := range changes {
		if processed[path] {
			continue
		}

		if change.Operation == "added" {
			if removedChange, exists := changes[path]; exists && removedChange.Operation == "removed" {
				fromValue := formatValue(removedChange.Value)
				toValue := formatValue(change.Value)
				result = append(result, fmt.Sprintf("Property '%s' was updated. From %s to %s", path, fromValue, toValue))
				processed[path] = true
			} else {
				value := formatValue(change.Value)
				result = append(result, fmt.Sprintf("Property '%s' was added with value: %s", path, value))
			}
		} else if change.Operation == "removed" {
			if addedChange, exists := changes[path]; exists && addedChange.Operation == "added" {
				continue
			} else {
				result = append(result, fmt.Sprintf("Property '%s' was removed", path))
			}
		}
	}
	sort.Strings(result)
	return strings.Join(result, "\n")
}

func buildPath(path []string) string {
	return strings.Join(path, ".")
}

func formatValue(value string) string {
	if value == "true" || value == "false" || value == "null" || value == "<nil>" {
		if value == "<nil>" {
			return "null"
		}
		return value
	}
	if strings.TrimSpace(value) == "" {
		return "''"
	}
	if value == "[complex value]" {
		return value
	}
	if _, err := strconv.Atoi(value); err == nil {
		return value
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return value
	}
	return fmt.Sprintf("'%s'", value)
}

type ChangeInfo struct {
	Operation string
	Value     string
}