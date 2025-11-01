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
		main.Println("error when importing env")
	}
	main.Println("starting server...")

	server.Start()
}
