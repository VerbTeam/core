package check

import (
	"strings"
)

func Check(wordtocheck string, wordList string) []string {
	finnalwordlist := strings.Split(wordList, "\n")

	slice := []string{}

	for _, word := range finnalwordlist {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		if strings.Contains(wordtocheck, word) {
			slice = append(slice, word)
		}
	}

	return slice
}
