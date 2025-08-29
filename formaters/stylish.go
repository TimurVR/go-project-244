package code

import (
	"strings"
)

func FormatDiffOutputStylish(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	indentLevel := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Уменьшаем уровень вложенности при закрывающих скобках
		if strings.Contains(trimmed, "}") && !strings.Contains(trimmed, "{") {
			indentLevel--
		}

		// Форматируем строку с правильными отступами
		formattedLine := formatStylishLine(trimmed, indentLevel)
		result = append(result, formattedLine)

		// Увеличиваем уровень вложенности при открывающих скобках
		if strings.Contains(trimmed, "{") && !strings.Contains(trimmed, "}") {
			indentLevel++
		}
	}

	return strings.Join(result, "\n")
}

func formatStylishLine(line string, indentLevel int) string {
	// Базовый отступ: 4 пробела на уровень
	baseIndent := strings.Repeat("    ", indentLevel)

	// Убираем лишние пробелы из исходной строки
	cleanedLine := strings.TrimSpace(line)

	if cleanedLine == "{" || cleanedLine == "}" {
		// Скобки: только базовый отступ
		return baseIndent + cleanedLine
	} else if strings.HasPrefix(cleanedLine, "+ ") {
		// Добавленные свойства: базовый отступ + содержание
		content := strings.TrimPrefix(cleanedLine, "+ ")
		return baseIndent + "+ " + content
	} else if strings.HasPrefix(cleanedLine, "- ") {
		// Удаленные свойства: базовый отступ + содержание
		content := strings.TrimPrefix(cleanedLine, "- ")
		return baseIndent + "- " + content
	} else {
		// Обычные свойства: базовый отступ + содержание
		return baseIndent + cleanedLine
	}
}
