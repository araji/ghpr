package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

//Colors Used for slack messages
const (
	RED   string = "#ff0000"
	GREEN string = "#008000"
	BLACK string = "#000000"
)

func main() {

	//Chek Presences of Required Environment Variables.
	SlackWebhook, sexists := os.LookupEnv("SLACK_WEBHOOK")
	GitOwner, gexists := os.LookupEnv("GIT_OWNER")
	GitRepo, rexists := os.LookupEnv("GIT_REPO")
	threshold, texists := os.LookupEnv("PR_THRESHOLD")
	pollPeriod, pexists := os.LookupEnv("POLL_PERIOD")

	if !sexists || !gexists || !rexists || !texists || !pexists {
		fmt.Printf("OWNER=%s, REPO= %s, Threshold= %s, POLL= %s \n\tWEBHOOK=%s\n", GitOwner, GitRepo, threshold, pollPeriod, SlackWebhook)
		log.Fatalf("one of the required environment variables is not set ")
	}
	Threshold, ok := strconv.Atoi(threshold)
	if ok != nil {
		log.Fatalf("Threshold should be set to number of days , got :%s", threshold)
	}
	PollPeriod, ok := strconv.Atoi(pollPeriod)
	if ok != nil || PollPeriod == 0 {
		log.Fatalf("pollPeriod should be set to number of minutes between two polls and greater than 0 , got :%s", pollPeriod)
	}

	// Initialize GithubClient and TimeTicker based on provided PollingPeriod
	// each tick pull pr list , send summary followed by details on the same slack channel
	ghc := GHClient{GitOwner, GitRepo}
	ticker := time.NewTicker(time.Duration(PollPeriod) * time.Minute)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		over, under, err := ghc.GetPushRequests(Threshold * 24)
		if err != nil {
			log.Println("Fetch Error, skipping ... ")
			continue
		}
		log.Println("over threshold", len(over), "under threshold := ", len(under))
		SendSlackMessage(SlackWebhook, fmt.Sprintf("over threshold= %d , under threshold = %d ", len(over), len(under)), BLACK)
		for _, pr := range over {
			SendSlackMessage(SlackWebhook, pr.HTMLURL, RED)
		}
		for _, pr := range under {
			SendSlackMessage(SlackWebhook, pr.HTMLURL, GREEN)
		}

	}
}
