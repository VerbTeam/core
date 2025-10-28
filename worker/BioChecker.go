package worker

import (
	wordcheck "codeberg.org/VerbTeam/WordsCheck"

	"encoding/json"

	api "codeberg.org/VerbTeam/core/api/roproxy"
)

func BioRun(userid int) string {
	bio, err := api.GetUserInfo(userid)
	if err != nil {
		return ""
	}

	res := wordcheck.Check(bio.Description, "")

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
