package main

import (
	"fmt"

	bioWorker "codeburg.org/VerbTeam/core/workers/BioCheck"
)

func main() {
	var id int

	fmt.Scanln(&id)

	resChan := make(chan string) // make a channel

	// run this in a diffrent thread (or that what i know)
	go func() {
		resChan <- bioWorker.Run(id) // catch it
	}()

	resBio := <-resChan // goes to this

	fmt.Println(resBio)
}
