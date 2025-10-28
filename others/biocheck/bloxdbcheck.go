package biocheck

import (
	"fmt"
	"os"

	"codeberg.org/VerbTeam/core/others/biocheck/module/check"
	"codeberg.org/VerbTeam/core/others/biocheck/module/listdownloader"
)

func Check(bio string, wordslisturl string) []string {
	fmt.Println("Downloading word list")
	listdownloader.Download(wordslisturl)

	data, err := os.ReadFile("wordlist.txt")
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}

	matches := check.Check(bio, string(data))
	return matches
}
