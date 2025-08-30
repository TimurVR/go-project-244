package code

import (
	format "code/formaters"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	str= GenDifference(map1, map2)
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
func GenDifference(map1, map2 map[string]interface{}) string {
	keys := []string{}
	for key := range map1 {
		keys = append(keys, key)
	}
	for key := range map2 {
		_, exist := map1[key]
		if !exist {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	str := "{\n"
	for _, key := range keys {
		value1, exist1 := map1[key]
		value2, exist2 := map2[key]
		isValue1Object := false
		isValue2Object := false
		if exist1 && value1 != nil {
			if _, ok := value1.(map[string]interface{}); ok {
				isValue1Object = true
			}
		}
		if exist2 && value2 != nil {
			if _, ok := value2.(map[string]interface{}); ok {
				isValue2Object = true
			}
		}
		if isValue1Object && isValue2Object {
			str += fmt.Sprintf("    %s: %s\n", key, GenDifference(value1.(map[string]interface{}), value2.(map[string]interface{})))
			continue
		}
		if isValue1Object && !isValue2Object {
			emptyMap := make(map[string]interface{})
			if exist2 {
				str += fmt.Sprintf("  - %s: %s\n", key, GenDifference(value1.(map[string]interface{}), emptyMap))
				str += fmt.Sprintf("  + %s: %v\n", key, value2)
			} else {
				str += fmt.Sprintf("  - %s: %s\n", key, GenDifference(value1.(map[string]interface{}), emptyMap))
			}
			continue
		}
		if !isValue1Object && isValue2Object {
			emptyMap := make(map[string]interface{})
			if exist1 {
				str += fmt.Sprintf("  - %s: %v\n", key, value1)
				str += fmt.Sprintf("  + %s: %s\n", key, GenDifference(emptyMap, value2.(map[string]interface{})))
			} else {
				str += fmt.Sprintf("  + %s: %s\n", key, GenDifference(emptyMap, value2.(map[string]interface{})))
			}
			continue
		}
		if exist1 && exist2 {
			if value1 == value2 {
				str += fmt.Sprintf("    %s: %v\n", key, value1)
			} else {
				str += fmt.Sprintf("  - %s: %v\n", key, value1)
				str += fmt.Sprintf("  + %s: %v\n", key, value2)
			}
		} else if exist1 && !exist2 {
			str += fmt.Sprintf("  - %s: %v\n", key, value1)
		} else if !exist1 && exist2 {
			str += fmt.Sprintf("  + %s: %v\n", key, value2)
		}
	}
	str += "}"
	return str
}

