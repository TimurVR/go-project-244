package code

import (
	"fmt"
	"strconv"
	"strings"
	"sort"
)

func GenDifferencePlain(map1, map2 map[string]interface{}) string {
    return buildPlainDiff(map1, map2, []string{})
}

func buildPlainDiff(map1, map2 map[string]interface{}, path []string) string {
    var result []string
    allKeys := make(map[string]bool)
    for k := range map1 {
        allKeys[k] = true
    }
    for k := range map2 {
        allKeys[k] = true
    }
    keys := make([]string, 0, len(allKeys))
    for k := range allKeys {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, key := range keys {
        currentPath := append([]string{}, path...)
        currentPath = append(currentPath, key)
        fullPath := strings.Join(currentPath, ".")

        value1, exists1 := map1[key]
        value2, exists2 := map2[key]

        if exists1 && exists2 {
            if isNestedMap(value1) && isNestedMap(value2) {
                nestedResult := buildPlainDiff(
                    value1.(map[string]interface{}),
                    value2.(map[string]interface{}),
                    currentPath,
                )
                if nestedResult != "" {
                    result = append(result, nestedResult)
                }
            } else if !isEqual(value1, value2) {
                fromValue := formatPlainValue(value1)
                toValue := formatPlainValue(value2)
                result = append(result, fmt.Sprintf("Property '%s' was updated. From %s to %s", fullPath, fromValue, toValue))
            }
        } else if exists1 && !exists2 {
            result = append(result, fmt.Sprintf("Property '%s' was removed", fullPath))
        } else if !exists1 && exists2 {
            value := formatPlainValue(value2)
            result = append(result, fmt.Sprintf("Property '%s' was added with value: %s", fullPath, value))
        }
    }

    return strings.Join(result, "\n")
}

func isNestedMap(value interface{}) bool {
    if value == nil {
        return false
    }
    _, ok := value.(map[string]interface{})
    return ok
}

func formatPlainValue(value interface{}) string {
    if value == nil {
        return "null"
    }
    if isNestedMap(value) {
        return "[complex value]"
    }
    switch v := value.(type) {
    case string:
        if v == "true" || v == "false" || v == "null" {
            return v
        }
        if _, err := strconv.Atoi(v); err == nil {
            return v
        }
        if _, err := strconv.ParseFloat(v, 64); err == nil {
            return v
        }
        return "'" + v + "'"
    default:
        return fmt.Sprintf("%v", v)
    }
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