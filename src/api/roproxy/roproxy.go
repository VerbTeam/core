package roproxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserInfo struct {
	Description string `json:"description"`
	Id          int    `json:"id"`
	IsBanned    bool   `json:"isBanned"`
}

func GetUserInfo(userid int) (UserInfo, error) {
	url := fmt.Sprintf("https://users.roproxy.com/v1/users/%d", userid)

	resp, err := http.Get(url)
	if err != nil {
		return UserInfo{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, fmt.Errorf("error reading response: %v", err)
	}

	var u UserInfo
	if err := json.Unmarshal(body, &u); err != nil {
		return UserInfo{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return u, nil
}
