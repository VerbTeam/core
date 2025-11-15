package main

import (
	"log"
	"os"
	"slices"

	server "codeberg.org/VerbTeam/core/server"
	"github.com/joho/godotenv"
)

func main() {
	main := log.New(os.Stdout, "[MAIN]: ", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()

	if err != nil {
		main.Println("error when importing env, is the file even exist?")
	}

	enable := os.Getenv("ENABLE_LOCAL_MODEL")

	envs := []string{"LOCAL_MODEL_URL", "GEMINI_API_KEY", "SUPABASE_URL", "REDIS_PUBLIC_ENDPOINT", "REDIS_USERNAME", "REDIS_PASSWORDS"} // lazy way

	if enable == "false" {
		main.Println("local model is disabled, skipping envs that related to it")

		envs = slices.Delete(envs, 1, 2)
	} else {
		main.Println("local model is enabled")
	}

	for _, envkey := range envs {

		value, exists := os.LookupEnv(envkey)
		if !exists {
			main.Fatalf("%s is not set\n", envkey)
		} else {
			main.Printf("%s = %s\n", envkey, value)
		}
	}

	main.Println("check all passed, starting server...")

	server.Start()
}
