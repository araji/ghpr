package main

import (
	"log"
	"time"
)

func main() {

	ghc := GHClient{"octocat", "hello-world"}
	slacker := Slacker{"araji", "password"}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			over, under, err := ghc.GetPushRequests(365 * 24)
			if err != nil {
				log.Println("Fetch Error will be ignored ")
			}
			log.Println("over threshold", over, "under threshold := ", under)
			slacker.SendMessage(&SlackMessage{Channel: "mychannel", Text: "ALERT"})
		}
	}
}
