package biocheck

import (
	"fmt"
	"os"

	"codeberg.org/VerbTeam/core/others/biocheck/module/check"
	"codeberg.org/VerbTeam/core/others/biocheck/module/listdownloader"
)

func Check(bio string, wordslisturl string) []string {
	if _, err := os.Stat("wordlist.txt"); os.IsNotExist(err) {
		fmt.Println("wordlist.txt not found, downloading word list")
		listdownloader.Download(wordslisturl)

	} else if err != nil {
		fmt.Println("stat error:", err)
		return nil
	}

	data, err := os.ReadFile("wordlist.txt")
	if err != nil {
		fmt.Println("read error:", err)
		return nil
	}

	matches := check.Check(bio, string(data))
	return matches
}
