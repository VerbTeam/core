package local_model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type LocalModel struct {
	Content    string `json:"content"`
	Prediction string `json:"prediction"`
	Label      int    `json:"label"`
}

func Fetch(text string) string {
	main := log.New(os.Stdout, "[LOCAL MODEL]: ", log.Ldate|log.Ltime|log.Lshortfile)

	main.Println("fetching to local model....")
	encoded := url.QueryEscape(text)
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:5000/?text=%s", encoded))

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

	return strconv.Itoa(PDIDDY.Label) // 1 ands 0
}
