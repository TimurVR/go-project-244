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
		"common":      map[string]interface{}{},
	}
	currentPath := []string{}
	inCommon := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "+ ") {
			processDiffLine(line[2:], "added", result, currentPath)
		} else if strings.HasPrefix(line, "- ") {
			processDiffLine(line[2:], "removed", result, currentPath)
		} else if strings.Contains(line, ":") {
			processCommonLine(line, result, currentPath, &inCommon)
		} else if line == "}" {
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
		}
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

func processDiffLine(line, diffType string, result map[string]interface{}, currentPath []string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])
	fullPath := append([]string{}, currentPath...)
	fullPath = append(fullPath, key)
	pathStr := strings.Join(fullPath, ".")
	diff := map[string]interface{}{
		"key":      pathStr,
		"type":     diffType,
		"oldValue": nil,
		"newValue": nil,
	}
	if diffType == "added" {
		diff["newValue"] = parseValue(valueStr)
	} else if diffType == "removed" {
		diff["oldValue"] = parseValue(valueStr)
	}
	differences := result["differences"].([]map[string]interface{})
	result["differences"] = append(differences, diff)
}

func processCommonLine(line string, result map[string]interface{}, currentPath []string, inCommon *bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])
	if valueStr == "" || valueStr == "{" {
		currentPath = append(currentPath, key)
	} else {
		common := result["common"].(map[string]interface{})
		fullPath := append([]string{}, currentPath...)
		fullPath = append(fullPath, key)
		
		// Создаем вложенную структуру
		current := common
		for _, pathPart := range fullPath[:len(fullPath)-1] {
			if _, exists := current[pathPart]; !exists {
				current[pathPart] = make(map[string]interface{})
			}
			current = current[pathPart].(map[string]interface{})
		}
		
		// Устанавливаем значение
		current[fullPath[len(fullPath)-1]] = parseValue(valueStr)
	}
}

func parseValue(valueStr string) interface{} {
	valueStr = strings.TrimSpace(valueStr)
	
	switch valueStr {
	case "true":
		return true
	case "false":
		return false
	case "null", "undefined", "":
		return nil
	}

	// Пробуем распарсить как число
	if num, err := strconv.Atoi(valueStr); err == nil {
		return num
	}
	if num, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return num
	}
	if strings.HasPrefix(valueStr, `"`) && strings.HasSuffix(valueStr, `"`) {
		return valueStr[1 : len(valueStr)-1]
	}

	return valueStr
}