package code

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func FormatDiffOutput(diffOutput string) string {
	lines := strings.Split(diffOutput, "\n")
	changes := make(map[string]ChangeInfo)
	var currentPath []string
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
			fullPath := buildPath(append(currentPath, key))
			if value == "{" {
				if operation == "added" {
					changes[fullPath] = ChangeInfo{
						Operation: "added",
						Value:     "[complex value]",
						IsComplex: true,
					}
				} else if operation == "removed" {
					changes[fullPath] = ChangeInfo{
						Operation: "removed",
						Value:     "",
						IsComplex: true,
					}
				}
				currentPath = append(currentPath, key)
				stack = append(stack, make(map[string]bool))
				continue
			}
			if operation == "added" || operation == "removed" {
				changes[fullPath] = ChangeInfo{
					Operation: operation,
					Value:     value,
					IsComplex: false,
				}
			}
		}
	}
	return formatChanges(changes)
}

func formatChanges(changes map[string]ChangeInfo) string {
	var result []string
	processed := make(map[string]bool)
	for path, change := range changes {
		if processed[path] {
			continue
		}
		if change.Operation == "added" {
			if removedChange, exists := changes[path]; exists && removedChange.Operation == "removed" {
				fromValue := formatValue(removedChange.Value, removedChange.IsComplex)
				toValue := formatValue(change.Value, change.IsComplex)
				result = append(result, fmt.Sprintf("Property '%s' was updated. From %s to %s", path, fromValue, toValue))
				processed[path] = true
				continue
			}
		}
	}
	for path, change := range changes {
		if processed[path] {
			continue
		}
		switch change.Operation {
		case "added":
			value := formatValue(change.Value, change.IsComplex)
			result = append(result, fmt.Sprintf("Property '%s' was added with value: %s", path, value))
		case "removed":
			parentPath := getParentPath(path)
			if _, parentRemoved := changes[parentPath]; !parentRemoved || changes[parentPath].Operation != "removed" {
				result = append(result, fmt.Sprintf("Property '%s' was removed", path))
			}
		}
	}
	for path, change := range changes {
		if processed[path] {
			continue
		}

		if change.IsComplex && (change.Operation == "removed" || change.Operation == "added") {
			if change.Operation == "removed" {
				if !hasRemovedParent(path, changes) {
					result = append(result, fmt.Sprintf("Property '%s' was removed", path))
				}
			} else if change.Operation == "added" {
				if !hasAddedParent(path, changes) {
					value := formatValue(change.Value, change.IsComplex)
					result = append(result, fmt.Sprintf("Property '%s' was added with value: %s", path, value))
				}
			}
			processed[path] = true
		}
	}

	sort.Strings(result)
	return strings.Join(result, "\n")
}

func hasRemovedParent(path string, changes map[string]ChangeInfo) bool {
	parts := strings.Split(path, ".")
	for i := len(parts) - 1; i > 0; i-- {
		parentPath := strings.Join(parts[:i], ".")
		if change, exists := changes[parentPath]; exists && change.Operation == "removed" {
			return true
		}
	}
	return false
}

func hasAddedParent(path string, changes map[string]ChangeInfo) bool {
	parts := strings.Split(path, ".")
	for i := len(parts) - 1; i > 0; i-- {
		parentPath := strings.Join(parts[:i], ".")
		if change, exists := changes[parentPath]; exists && change.Operation == "added" {
			return true
		}
	}
	return false
}
func getParentPath(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) <= 1 {
		return ""
	}
	return strings.Join(parts[:len(parts)-1], ".")
}
func buildPath(path []string) string {
	return strings.Join(path, ".")
}
func formatValue(value string, isComplex bool) string {
	if isComplex {
		return "[complex value]"
	}
	if value == "true" || value == "false" || value == "null" || value == "<nil>" {
		if value == "<nil>" {
			return "null"
		}
		return value
	}
	if strings.TrimSpace(value) == "" {
		return "''"
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
	IsComplex bool
}
