package main

import (
	"fmt"
	"log"
	"time"
)

func main() {

	ghc := GHClient{"uber", "makisu"}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			over, under, err := ghc.GetPushRequests(365 * 24)
			if err != nil {
				log.Println("Fetch Error will be ignored ")
			}
			log.Println("over threshold", len(over), "under threshold := ", len(under))
			SendSlackMessage(fmt.Sprintf("over threshold= %d , under threshold = %d ", len(over), len(under)), "#000000")
			for _, pr := range over {
				SendSlackMessage(pr.HTMLURL, "#ff0000")
			}
			for _, pr := range under {
				SendSlackMessage(pr.HTMLURL, "#008000")
			}

		}
	}
}
