package local_model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type LocalModel struct {
	Content    string `json:"content"`
	Prediction string `json:"prediction"`
	Label      int    `json:"label"`
}

func Fetch(text string) string {
	main := log.New(os.Stdout, "[LOCAL MODEL]: ", log.Ldate|log.Ltime|log.Lshortfile)

	main.Println("fetching to local model....")
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:5000/?text=%s", text))

	if err != nil {
		main.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	main.Println("extracting the body mate")

	var PDIDDY LocalModel // i ran out of name to name this var soooooooooooo
	err = json.Unmarshal(body, &PDIDDY)

	if err != nil {
		main.Fatal(err)
	}

	main.Println("done job")
	return string(PDIDDY.Label) // 1 ands 0
}
