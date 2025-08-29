package code

import (
	"strings"
)

func FormatDiffOutputStylish(diffOutput string) string {
	// Убираем обрамляющие фигурные скобки если они есть
	trimmedOutput := strings.TrimSpace(diffOutput)
	if strings.HasPrefix(trimmedOutput, "{") && strings.HasSuffix(trimmedOutput, "}") {
		trimmedOutput = strings.TrimSpace(trimmedOutput[1 : len(trimmedOutput)-1])
	}

	lines := strings.Split(trimmedOutput, "\n")
	var result []string
	depth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Базовый отступ - 4 пробела на уровень глубины
		baseIndent := strings.Repeat("    ", depth)

		// Обработка закрывающих скобок
		if trimmed == "}" {
			if depth > 0 {
				depth--
			}
			baseIndent = strings.Repeat("    ", depth)
			result = append(result, baseIndent+"}")
			continue
		}

		// Обработка открывающих скобок объектов
		if strings.HasSuffix(trimmed, "{") {
			var prefix, content string

			if strings.HasPrefix(trimmed, "- ") {
				prefix = "  - "
				content = strings.TrimSuffix(strings.TrimPrefix(trimmed, "- "), " {")
			} else if strings.HasPrefix(trimmed, "+ ") {
				prefix = "  + "
				content = strings.TrimSuffix(strings.TrimPrefix(trimmed, "+ "), " {")
			} else {
				prefix = "    "
				content = strings.TrimSuffix(trimmed, " {")
			}

			// Убираем двоеточие из content, если оно уже есть
			content = strings.TrimSuffix(content, ":")
			result = append(result, baseIndent+prefix+content+": {")
			depth++
			continue
		}

		// Обработка обычных свойств - проверяем, содержит ли строка уже двоеточие
		if strings.Contains(trimmed, ":") {
			// Если строка уже содержит двоеточие, просто добавляем отступы и префиксы
			if strings.HasPrefix(trimmed, "- ") {
				content := strings.TrimPrefix(trimmed, "- ")
				result = append(result, baseIndent+"  - "+content)
			} else if strings.HasPrefix(trimmed, "+ ") {
				content := strings.TrimPrefix(trimmed, "+ ")
				result = append(result, baseIndent+"  + "+content)
			} else {
				result = append(result, baseIndent+"    "+trimmed)
			}
		} else {
			// Если строка не содержит двоеточия, обрабатываем как раньше
			if strings.HasPrefix(trimmed, "- ") {
				content := strings.TrimPrefix(trimmed, "- ")
				result = append(result, baseIndent+"  - "+content)
			} else if strings.HasPrefix(trimmed, "+ ") {
				content := strings.TrimPrefix(trimmed, "+ ")
				result = append(result, baseIndent+"  + "+content)
			} else {
				result = append(result, baseIndent+"    "+trimmed)
			}
		}
	}

	// Возвращаем результат с обрамляющими скобками
	return "{\n" + strings.Join(result, "\n") + "\n}"
}
