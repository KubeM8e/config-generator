package utils

import "strings"

func ChangeTheNamesToCorrectFormat(name string) string {
	cleaned := strings.ReplaceAll(name, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ToLower(cleaned)
	return cleaned
}
