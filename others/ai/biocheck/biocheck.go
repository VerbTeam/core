package biocheck

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/genai"

	"codeberg.org/VerbTeam/core/others/ai"
)

var FullPrompt = ai.IntroductionPrompt + "\n" + ai.BioCheckPrompt + "\n" + ai.Rating

func Check(geminiapikey string, bio string) string {
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

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-flash-lite-latest",
		genai.Text(fmt.Sprintf("Analyze this :%v", bio)),
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
