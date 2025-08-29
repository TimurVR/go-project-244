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
			// Определяем префикс на основе начала строки
			var prefix string
			var content string

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

			result = append(result, baseIndent+prefix+content+": {")
			depth++
			continue
		}

		// Для обычных строк просто добавляем отступы, сохраняя оригинальные префиксы
		// Определяем нужный отступ based on prefix
		if strings.HasPrefix(trimmed, "- ") {
			result = append(result, baseIndent+"  - "+strings.TrimPrefix(trimmed, "- "))
		} else if strings.HasPrefix(trimmed, "+ ") {
			result = append(result, baseIndent+"  + "+strings.TrimPrefix(trimmed, "+ "))
		} else {
			result = append(result, baseIndent+"    "+trimmed)
		}
	}

	// Возвращаем результат с обрамляющими скобками
	return "{\n" + strings.Join(result, "\n") + "\n}"
}
