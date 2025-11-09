package biocheck

import (
	"log"
	"os"

	"codeberg.org/VerbTeam/core/others/biocheck/module/check"
	"codeberg.org/VerbTeam/core/others/biocheck/module/listdownloader"
)

func Check(bio string, wordslisturl string) []string {

	main := log.New(os.Stdout, "[BLOXDB]: ", log.Ldate|log.Ltime|log.Lshortfile)

	if _, err := os.Stat("wordlist.txt"); os.IsNotExist(err) {
		main.Println("wordlist.txt not found, downloading word list")
		listdownloader.Download("https://raw.githubusercontent.com/whoschip/wordlist/main/data/all.txt")

	} else if err != nil {
		main.Println("stat error:", err)
		return nil
	}

	data, err := os.ReadFile("wordlist.txt")
	if err != nil {
		main.Println("read error:", err)
		return nil
	}

	matches := check.Check(bio, string(data))
	return matches
}
