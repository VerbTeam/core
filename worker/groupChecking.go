package worker

import (
	"encoding/json"
	"fmt"

	api "codeberg.org/VerbTeam/core/api/roproxy"
)

type Response struct {
	Data []struct {
		Group struct {
			ID int64 `json:"id"`
		} `json:"group"`
	} `json:"data"`
}

func RunGroupCheck(userid int) string {
	ress, err := api.GetUserGroups(userid)

	if err != nil {
		fmt.Println(err)
	}

	b, _ := json.Marshal(ress)

	var res Response
	if err := json.Unmarshal([]byte(string(b)), &res); err != nil {
		panic(err)
	}

	// step 1: make a slice of all ids
	var allIDs []int64
	for _, item := range res.Data {
		allIDs = append(allIDs, item.Group.ID)
	}

	// step 2: make a "keep" list
	keepList := []int64{35396105, 35065141, 34664468, 693308, 3331780, 143355, 33904411, 36058773, 17284863, 6046690, 1018746, 35991486, 36048528} // more group ids can be put here but as of right now we just use a fixed list

	// step 3: use a map for quick lookup
	keepMap := make(map[int64]bool)
	for _, id := range keepList {
		keepMap[id] = true
	}

	// step 4: filter out ids not in keepList
	var filtered []int64
	for _, id := range allIDs {
		if keepMap[id] {
			filtered = append(filtered, id)
		}
	}

	str, _ := json.Marshal(filtered)
	return string(str)
}
