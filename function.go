package code

import (
	format "code/formaters"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parsing(file string) (map[string]interface{}, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var data1 map[string]interface{}
	raz := filepath.Ext(file)
	switch raz {
	case ".json":
		err = json.Unmarshal([]byte(data), &data1)
		if err != nil {
			return nil, err
		}
	default:
		err = yaml.Unmarshal(data, &data1)
		if err != nil {
			return nil, err
		}
	}
	return data1, nil
}
func GenDiff(file1 string, file2 string, style string) (string, error) {
	map1, err := Parsing(file1)
	if err != nil {
		return "", err
	}
	map2, err := Parsing(file2)
	if err != nil {
		return "", err
	}
	var str string
	str, err = GenDifference(map1, map2)
	if err != nil {
		return "", err
	}
	switch style {
	case "json":
		temp, err := format.FormatDiffToJSON(str)
		if err != nil {
			return "", err
		}
		str = temp
	case "stylish":
		str = format.FormatDiffOutputStylish(str)
	case "plain":
		str = format.FormatDiffOutput(str)
	}
	return str, nil
}
func GenDifference(map1, map2 map[string]interface{}) (string, error) {
	return genDiffRecursive(map1, map2, 0, true)
}

func genDiffRecursive(map1, map2 map[string]interface{}, depth int, isRoot bool) (string, error) {
	indent := strings.Repeat("    ", depth)
	keys := getUniqueKeys(map1, map2)
	sort.Strings(keys)

	var result strings.Builder
	if !isRoot {
		result.WriteString("{\n")
	}

	for _, key := range keys {
		value1, exist1 := map1[key]
		value2, exist2 := map2[key]

		isValue1Object := isMap(value1)
		isValue2Object := isMap(value2)

		lineIndent := indent
		if depth > 0 {
			lineIndent = strings.Repeat("    ", depth)
		}

		switch {
		case isValue1Object && isValue2Object:
			nestedDiff, err := genDiffRecursive(
				value1.(map[string]interface{}),
				value2.(map[string]interface{}),
				depth+1,
				false,
			)
			if err != nil {
				return "", err
			}
			result.WriteString(fmt.Sprintf("%s    %s: %s\n", lineIndent, key, nestedDiff))

		case isValue1Object && !isValue2Object:
			nestedDiff, err := genDiffRecursive(
				value1.(map[string]interface{}),
				map[string]interface{}{},
				depth+1,
				false,
			)
			if err != nil {
				return "", err
			}
			result.WriteString(fmt.Sprintf("%s  - %s: %s\n", lineIndent, key, nestedDiff))
			if exist2 {
				result.WriteString(fmt.Sprintf("%s  + %s: %v\n", lineIndent, key, formatValue(value2)))
			}

		case !isValue1Object && isValue2Object:
			if exist1 {
				result.WriteString(fmt.Sprintf("%s  - %s: %v\n", lineIndent, key, formatValue(value1)))
			}
			nestedDiff, err := genDiffRecursive(
				map[string]interface{}{},
				value2.(map[string]interface{}),
				depth+1,
				false,
			)
			if err != nil {
				return "", err
			}
			result.WriteString(fmt.Sprintf("%s  + %s: %s\n", lineIndent, key, nestedDiff))

		default:
			if exist1 && exist2 {
				if equalValues(value1, value2) {
					result.WriteString(fmt.Sprintf("%s    %s: %v\n", lineIndent, key, formatValue(value1)))
				} else {
					result.WriteString(fmt.Sprintf("%s  - %s: %v\n", lineIndent, key, formatValue(value1)))
					result.WriteString(fmt.Sprintf("%s  + %s: %v\n", lineIndent, key, formatValue(value2)))
				}
			} else if exist1 {
				result.WriteString(fmt.Sprintf("%s  - %s: %v\n", lineIndent, key, formatValue(value1)))
			} else if exist2 {
				result.WriteString(fmt.Sprintf("%s  + %s: %v\n", lineIndent, key, formatValue(value2)))
			}
		}
	}

	if !isRoot {
		result.WriteString(fmt.Sprintf("%s}", indent))
	}
	return result.String(), nil
}
func getUniqueKeys(map1, map2 map[string]interface{}) []string {
	keysMap := make(map[string]bool)
	for key := range map1 {
		keysMap[key] = true
	}
	for key := range map2 {
		keysMap[key] = true
	}

	keys := make([]string, 0, len(keysMap))
	for key := range keysMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func isMap(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(map[string]interface{})
	return ok
}

func equalValues(v1, v2 interface{}) bool {
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}
	return fmt.Sprintf("%v", v1) == fmt.Sprintf("%v", v2)
}

func formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	if s, ok := value.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", value)
}
