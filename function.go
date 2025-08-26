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
func GenDiff(map1, map2 map[string]interface{}, style string) (string, error) {
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
			temp, err := GenDiff(value1.(map[string]interface{}), value2.(map[string]interface{}), "")
			if err != nil {
				return "", err
			}
			str += fmt.Sprintf("    %s: %s\n", key, temp)
			continue
		}
		if isValue1Object && !isValue2Object {
			emptyMap := make(map[string]interface{})
			if exist2 {
				temp, err := GenDiff(value1.(map[string]interface{}), emptyMap, "")
				if err != nil {
					return "", err
				}
				str += fmt.Sprintf("  - %s: %s\n", key, temp)
				str += fmt.Sprintf("  + %s: %v\n", key, value2)
			} else {
				temp, err := GenDiff(value1.(map[string]interface{}), emptyMap, "")
				if err != nil {
					return "", err
				}
				str += fmt.Sprintf("  - %s: %s\n", key, temp)
			}
			continue
		}
		if !isValue1Object && isValue2Object {
			emptyMap := make(map[string]interface{})
			if exist1 {
				str += fmt.Sprintf("  - %s: %v\n", key, value1)
				temp, err := GenDiff(emptyMap, value2.(map[string]interface{}), "")
				if err != nil {
					return "", err
				}
				str += fmt.Sprintf("  + %s: %s\n", key, temp)
			} else {
				temp, err := GenDiff(emptyMap, value2.(map[string]interface{}), "")
				if err != nil {
					return "", nil
				}
				str += fmt.Sprintf("  + %s: %s\n", key, temp)
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
	switch style {
	case "json":
		str = format.FormatDiffToJSON(str)
	case "stylish":
		str = format.FormatDiffOutputStylish(str)
	case "plan":
		str = format.FormatDiffOutput(str)
	}
	return str, nil
}
