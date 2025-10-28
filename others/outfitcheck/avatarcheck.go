package avatarcheck

import (
	"context"
	"io"
	"log"
	"net/http"

	"google.golang.org/genai"
)

func Check(geminiapikey string, imageurl string) string {
	ctx := context.Background()

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(`
you’re an ai avatar moderator for roblox,a game mostly for kids, run by "VerbTeam" (3rd party tool). your job is to check each outfit image a user is wearing and see if it’s safe for kids. focus on stuff that’s sexually suggestive, explicit, or just straight-up not kid-friendly. each outfit will come with a name and an image. your response should clearly say if it’s appropriate or inappropriate, with a short explanation. keep it kid-safe and to the point.

your job:

* check if the avatar is **appropriate** or **inappropriate**.
* only reply with **one of these 2 answers**:

  1. **Appropriate** – explain briefly why it's fine.
  2. **Inappropriate** – explain briefly why it's not suitable for kids.

keep it short, clear, and kid-safe.
happy moderating!
`, genai.RoleUser), // pls dont copy this i spent my life on this  prompt
		ResponseMIMEType: "application/json",

		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"status": {
					Type: genai.TypeString,
				},
				"reason": {
					Type: genai.TypeString,
				},
			},
			PropertyOrdering: []string{"status", "reason"},
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

	return result.Text()
}
