package utility

import (
	"os"
	"regexp"
)

func FlattenWhitespaces(input string) string {
	whitespacePattern := regexp.MustCompile(`\s+`)
	res := whitespacePattern.ReplaceAllString(input, " ")
	return res
}

func RemoveWhitespaces(input string) string {
	whitespacePattern := regexp.MustCompile(`\s+`)
	res := whitespacePattern.ReplaceAllString(input, "")
	return res
}

func CreateFileIfNotExists(filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			panic(err) // Handle error appropriately
		}
		file.Close()
	}
}
