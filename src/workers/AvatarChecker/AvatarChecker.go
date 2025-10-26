package avatarchecker

import (
	"os"

	checker "codeberg.org/VerbTeam/AvatarCheck"
)

func Run(userid int) {
	checker.Check(os.Getenv("GEMINI_API_KEY"), "https://tr.rbxcdn.com/30DAY-Avatar-88AF4A3282EF465B02D7B2424B210B39-Png/352/352/Avatar/Webp/noFilter")
}
