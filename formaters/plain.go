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
	stack := []map[string]interface{}{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if trimmed == "{" {
			stack = append(stack, make(map[string]interface{}))
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
			cleanKey := key

			if strings.HasPrefix(key, "+ ") {
				operation = "added"
				cleanKey = strings.TrimPrefix(key, "+ ")
			} else if strings.HasPrefix(key, "- ") {
				operation = "removed"
				cleanKey = strings.TrimPrefix(key, "- ")
			}

			fullPath := buildPath(append(currentPath, cleanKey))

			// Обработка вложенных объектов
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
				} else {
					// Для неизмененных сложных объектов просто продолжаем путь
					currentPath = append(currentPath, cleanKey)
					stack = append(stack, make(map[string]interface{}))
					continue
				}
				currentPath = append(currentPath, cleanKey)
				stack = append(stack, make(map[string]interface{}))
				continue
			}

			// Обработка простых значений
			if operation != "unchanged" {
				// Проверяем, есть ли противоположная операция для этого пути
				if existing, exists := changes[fullPath]; exists {
					if (existing.Operation == "added" && operation == "removed") ||
						(existing.Operation == "removed" && operation == "added") {
						// Это обновление значения
						fromValue := existing.Value
						if existing.Operation == "removed" {
							fromValue = existing.Value
						}
						changes[fullPath] = ChangeInfo{
							Operation: "updated",
							OldValue:  fromValue,
							Value:     value,
							IsComplex: false,
						}
					}
				} else {
					// Новая операция
					changes[fullPath] = ChangeInfo{
						Operation: operation,
						Value:     value,
						IsComplex: false,
					}
				}
			}
		}
	}

	return formatChanges(changes)
}

func formatChanges(changes map[string]ChangeInfo) string {
	var result []string
	processed := make(map[string]bool)

	// Сначала обрабатываем обновления
	for path, change := range changes {
		if change.Operation == "updated" {
			fromValue := formatValue(change.OldValue, false)
			toValue := formatValue(change.Value, false)
			result = append(result, fmt.Sprintf("Property '%s' was updated. From %s to %s", path, fromValue, toValue))
			processed[path] = true
		}
	}

	// Затем обрабатываем добавления и удаления
	for path, change := range changes {
		if processed[path] {
			continue
		}

		// Пропускаем если родительский объект был удален/добавлен
		if hasParentOperation(path, changes, "removed") && change.Operation == "removed" {
			continue
		}
		if hasParentOperation(path, changes, "added") && change.Operation == "added" {
			continue
		}

		switch change.Operation {
		case "added":
			value := formatValue(change.Value, change.IsComplex)
			result = append(result, fmt.Sprintf("Property '%s' was added with value: %s", path, value))
		case "removed":
			result = append(result, fmt.Sprintf("Property '%s' was removed", path))
		}
	}

	sort.Strings(result)
	return strings.Join(result, "\n")
}

func hasParentOperation(path string, changes map[string]ChangeInfo, operation string) bool {
	parts := strings.Split(path, ".")
	for i := len(parts) - 1; i > 0; i-- {
		parentPath := strings.Join(parts[:i], ".")
		if change, exists := changes[parentPath]; exists && change.Operation == operation {
			return true
		}
	}
	return false
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
	OldValue  string
	IsComplex bool
}
