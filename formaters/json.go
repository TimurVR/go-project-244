package code

import (
	"encoding/json"
	"strconv"
	"strings"
)

func FormatDiffToJSON(input string) (string) {
	lines := strings.Split(input, "\n")
	result := map[string]interface{}{
		"differences": []map[string]interface{}{},
		"common":      []map[string]interface{}{},
	}
	differences := parseDifferences(lines)
	result["differences"] = differences
	buildCommonFromDifferences(differences, &result)

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
func parseDifferences(lines []string) []map[string]interface{} {
	differences := []map[string]interface{}{}
	currentPath := []string{}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		if strings.Contains(trimmedLine, ":") && !strings.HasPrefix(trimmedLine, "+ ") && !strings.HasPrefix(trimmedLine, "- ") {
			parts := strings.SplitN(trimmedLine, ":", 2)
			key := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])
			
			if valueStr == "" || valueStr == "{" {
				currentPath = append(currentPath, key)
			}
		} else if trimmedLine == "}" {
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
		}
		if strings.HasPrefix(trimmedLine, "+ ") || strings.HasPrefix(trimmedLine, "- ") {
			parts := strings.SplitN(trimmedLine[2:], ":", 2)
			if len(parts) < 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])

			if valueStr == "{" {
				continue
			}
			fullPath := make([]string, len(currentPath))
			copy(fullPath, currentPath)
			fullPath = append(fullPath, key)
			pathStr := strings.Join(fullPath, ".")
			pathStr = strings.TrimPrefix(pathStr, "common.")

			diff := map[string]interface{}{
				"key": pathStr,
			}

			value := parseValue(valueStr)
			
			if strings.HasPrefix(trimmedLine, "+ ") {
				diff["type"] = "added"
				diff["newValue"] = value
				diff["oldValue"] = nil
			} else {
				diff["type"] = "removed"
				diff["oldValue"] = value
				diff["newValue"] = nil
			}

			differences = append(differences, diff)
		}
	}

	return differences
}
func buildCommonFromDifferences(differences []map[string]interface{}, result *map[string]interface{}) {
	commonArray := []map[string]interface{}{}
	commonMap := make(map[string]interface{})
	finalValues := make(map[string]interface{})
	for _, diff := range differences {
		key := diff["key"].(string)
		
		if diff["type"] == "added" {
			finalValues[key] = diff["newValue"]
		} else if diff["type"] == "removed" {
			delete(finalValues, key)
		}
	}

	for keyPath, value := range finalValues {
		setNestedValue(commonMap, keyPath, value)
	}
	for key, value := range commonMap {
		commonArray = append(commonArray, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}

	(*result)["common"] = commonArray
}
func setNestedValue(obj map[string]interface{}, keyPath string, value interface{}) {
	keys := strings.Split(keyPath, ".")
	current := obj

	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = value
		} else {
			if current[key] == nil {
				current[key] = make(map[string]interface{})
			}
			if nested, ok := current[key].(map[string]interface{}); ok {
				current = nested
			}
		}
	}
}
func parseValue(valueStr string) interface{} {
	valueStr = strings.TrimSpace(valueStr)
	
	if valueStr == "{" {
		return nil
	}
	if valueStr == "\u003cnil\u003e" {
		return nil
	}
	if valueStr == `"<nil>"` {
		return nil
	}
	if valueStr == "null" || valueStr == "undefined" || valueStr == "" {
		return nil
	}
	switch valueStr {
	case "true":
		return true
	case "false":
		return false
	}

	if num, err := strconv.Atoi(valueStr); err == nil {
		return num
	}
	if num, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return num
	}
	if strings.HasPrefix(valueStr, `"`) && strings.HasSuffix(valueStr, `"`) {
		unquoted := valueStr[1 : len(valueStr)-1]
		unquoted = strings.ReplaceAll(unquoted, ``, "<")
		unquoted = strings.ReplaceAll(unquoted, ``, ">")
		return unquoted
	}
	return valueStr
}
