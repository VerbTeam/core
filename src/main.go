package main

import (
	"fmt"

	"github.com/joho/godotenv"

	avatarWorker "codeberg.org/VerbTeam/core/workers/AvatarChecker"
	bioWorker "codeberg.org/VerbTeam/core/workers/BioCheck"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	var id int

	fmt.Scanln(&id)

	resChan := make(chan string) // make a channel

	// run this in a diffrent thread (or that what i know)
	go func() {
		resChan <- bioWorker.Run(id) // catch it
	}()

	go func() {
		resChan <- avatarWorker.Run(id) // catch it
	}()

	resBio := <-resChan // goes to this

	fmt.Println(resBio)
}
