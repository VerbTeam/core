package main

import (
	"fmt"

	server "codeberg.org/VerbTeam/core/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("error when importing env")
	}

	server.Start()
}
