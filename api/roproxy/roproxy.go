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

type AvatarItem struct {
	TargetId int    `json:"targetId"`
	State    string `json:"state"`
	ImageUrl string `json:"imageUrl"`
	Version  string `json:"version"`
}

type AvatarResponse struct {
	Data []AvatarItem `json:"data"`
}

type Groups struct {
	Data []GroupsItem `json:"data"`
}

type GroupsItem struct {
	Group GroupInfo `json:"group"`
}

type GroupInfo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

func GetUserAvatar(userid int) (AvatarResponse, error) {
	url := fmt.Sprintf("https://thumbnails.roblox.com/v1/users/avatar?userIds=%d&size=420x420&format=Png&isCircular=false", userid)

	resp, err := http.Get(url)
	if err != nil {
		return AvatarResponse{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AvatarResponse{}, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AvatarResponse{}, fmt.Errorf("error reading response: %v", err)
	}

	var res AvatarResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return AvatarResponse{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return res, nil
}

func GetUserGroups(userid int) (Groups, error) {
	url := fmt.Sprintf("https://groups.roblox.com/v1/users/%d/groups/roles?includeLocked=true", userid)

	resp, err := http.Get(url)
	if err != nil {
		return Groups{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Groups{}, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Groups{}, fmt.Errorf("error reading response: %v", err)
	}

	var gp Groups
	if err := json.Unmarshal(body, &gp); err != nil {
		return Groups{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return gp, nil
}
