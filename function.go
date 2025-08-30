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
	case "plain":
		str = format.FormatDiffOutput(str)
	}
	return str, nil
}

func GenDifference(map1, map2 map[string]interface{}) (string, error) {
	return genDifference(map1, map2, 0)
}

func genDifference(map1, map2 map[string]interface{}, indentLevel int) (string, error) {
	indent := strings.Repeat("    ", indentLevel)

	keys := getUniqueKeys(map1, map2)
	sort.Strings(keys)

	var result strings.Builder
	result.WriteString("{\n")

	for _, key := range keys {
		value1, exist1 := map1[key]
		value2, exist2 := map2[key]
		isNested1 := isNestedObject(value1)
		isNested2 := isNestedObject(value2)
		if !exist1 && exist2 {
			if isNested2 {
				nestedDiff, err := genDifference(map[string]interface{}{}, value2.(map[string]interface{}), indentLevel+1)
				if err != nil {
					return "", err
				}
				result.WriteString(fmt.Sprintf("%s  + %s: %s\n", indent, key, nestedDiff))
			} else {
				result.WriteString(fmt.Sprintf("%s  + %s: %s\n", indent, key, formatValue(value2)))
			}
		} else if exist1 && !exist2 {
			if isNested1 {
				nestedDiff, err := genDifference(value1.(map[string]interface{}), map[string]interface{}{}, indentLevel+1)
				if err != nil {
					return "", err
				}
				result.WriteString(fmt.Sprintf("%s  - %s: %s\n", indent, key, nestedDiff))
			} else {
				result.WriteString(fmt.Sprintf("%s  - %s: %s\n", indent, key, formatValue(value1)))
			}
		} else {
			if isNested1 && isNested2 {
				nestedDiff, err := genDifference(value1.(map[string]interface{}), value2.(map[string]interface{}), indentLevel+1)
				if err != nil {
					return "", err
				}
				result.WriteString(fmt.Sprintf("%s    %s: %s\n", indent, key, nestedDiff))
			} else if isNested1 && !isNested2 {
				nestedDiff, err := genDifference(value1.(map[string]interface{}), map[string]interface{}{}, indentLevel+1)
				if err != nil {
					return "", err
				}
				result.WriteString(fmt.Sprintf("%s  - %s: %s\n", indent, key, nestedDiff))
				result.WriteString(fmt.Sprintf("%s  + %s: %s\n", indent, key, formatValue(value2)))
			} else if !isNested1 && isNested2 {
				result.WriteString(fmt.Sprintf("%s  - %s: %s\n", indent, key, formatValue(value1)))
				nestedDiff, err := genDifference(map[string]interface{}{}, value2.(map[string]interface{}), indentLevel+1)
				if err != nil {
					return "", err
				}
				result.WriteString(fmt.Sprintf("%s  + %s: %s\n", indent, key, nestedDiff))
			} else {
				if isEqual(value1, value2) {
					result.WriteString(fmt.Sprintf("%s    %s: %s\n", indent, key, formatValue(value1)))
				} else {
					result.WriteString(fmt.Sprintf("%s  - %s: %s\n", indent, key, formatValue(value1)))
					result.WriteString(fmt.Sprintf("%s  + %s: %s\n", indent, key, formatValue(value2)))
				}
			}
		}
	}

	result.WriteString(fmt.Sprintf("%s}", indent))
	return result.String(), nil
}
func getUniqueKeys(map1, map2 map[string]interface{}) []string {
	keys := make(map[string]bool)
	for k := range map1 {
		keys[k] = true
	}
	for k := range map2 {
		keys[k] = true
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func isNestedObject(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(map[string]interface{})
	return ok
}

func formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", value)
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
