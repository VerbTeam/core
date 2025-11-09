package main

import (
	"log"
	"os"

	server "codeberg.org/VerbTeam/core/server"
	"github.com/joho/godotenv"
)

func main() {
	main := log.New(os.Stdout, "[MAIN]: ", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()

	if err != nil {
		main.Println("error when importing env, is the file even exist?")
	}

	envs := []string{"GEMINI_API_KEY", "SUPABASE_URL", "REDIS_PUBLIC_ENDPOINT", "REDIS_USERNAME", "REDIS_PASSWORDS"} // lazy way

	for _, envkey := range envs {

		value, exists := os.LookupEnv(envkey)
		if !exists {
			main.Fatalf("%s is not set\n", envkey)
		} else {
			main.Printf("%s = %s\n", envkey, value)
		}
	}

	main.Println("starting server...")

	server.Start()
}
