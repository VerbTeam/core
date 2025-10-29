package worker

import (
	"os"

	checker "codeberg.org/VerbTeam/core/others/ai/biocheck"

	"encoding/json"

	api "codeberg.org/VerbTeam/core/api/roproxy"
)

func BioRunAI(userid int) string {
	bio, err := api.GetUserInfo(userid)
	if err != nil {
		return ""
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "GEMINI_API_KEY not set"
	}

	res := checker.Check(apiKey, bio.Description)

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
