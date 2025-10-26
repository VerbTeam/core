package avatarchecker

import (
	"os"

	checker "codeberg.org/VerbTeam/AvatarCheck"

	api "codeberg.org/VerbTeam/core/api/roproxy"
)

func Run(userid int) string {
	avatar, err := api.GetUserAvatar(userid)
	if err != nil {
		return err.Error()
	}

	if len(avatar.Data) == 0 {
		return "no avatar data, cannot check"
	}

	imageUrl := avatar.Data[0].ImageUrl
	if imageUrl == "" {
		return "avatar imageUrl is empty"
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "GEMINI_API_KEY not set"
	}

	res := checker.Check(apiKey, imageUrl)
	return res
}
