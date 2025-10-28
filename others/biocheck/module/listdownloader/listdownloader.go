package listdownloader

import (
	"io"
	"net/http"
	"os"
)

func Download(url string) {
	if url == "" {
		url = "https://raw.githubusercontent.com/bloxDB/wordlist/main/data/all.txt"
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("wordlist.txt")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, _ = io.Copy(out, resp.Body)
}
