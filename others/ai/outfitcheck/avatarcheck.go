package avatarcheck

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"google.golang.org/genai"

	"codeberg.org/VerbTeam/core/others/ai"
)

var FullPrompt = ai.IntroductionPrompt + "\n" + ai.AvatarPrompt + "\n" + ai.Rating

func Check(geminiapikey string, imageurl string) string {
	ctx := context.Background()

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(FullPrompt, genai.RoleUser),
		ResponseMIMEType:  "application/json",

		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"status": {
					Type: genai.TypeString,
				},
				"reason": {
					Type: genai.TypeString,
				},
				"rating": {
					Type: genai.TypeNumber,
				},
			},
			Required:         []string{"status", "reason", "rating"},
			PropertyOrdering: []string{"status", "reason", "rating"},
		},
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: geminiapikey,
	})
	if err != nil {
		log.Fatal(err)
	}

	imageResp, _ := http.Get(imageurl)

	imageBytes, _ := io.ReadAll(imageResp.Body)

	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/jpeg"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-flash-lite-latest",
		contents,
		config,
	)
	if err != nil {
		log.Fatal(err)
	}
	clean := strings.ReplaceAll(result.Text(), "\n", "") // bloat killer v2
	clean = strings.ReplaceAll(clean, "  ", "")
	clean = strings.TrimSpace(clean)
	return clean
}
