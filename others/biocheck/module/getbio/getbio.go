package getbio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserInfo struct {
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	IsBanned    bool      `json:"isBanned"`
	ID          int       `json:"id"`
	Name        string    `json:"name"`
}

func Getbio(userid int) string {
	url := fmt.Sprintf("https://users.roproxy.com/v1/users/%v", userid)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func GetbioWithCheck(userid int) string {
	datastr := Getbio(userid)
	var u UserInfo

	json.Unmarshal([]byte(datastr), &u)

	if u.IsBanned {
		return "banned"
	}
	return u.Description
}
