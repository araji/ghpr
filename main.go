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

//Environment Variables required to exist
var (
	SlackWebhook string
	GitAPIServer string
	GitToken     string
	GitOwner     string
	GitRepo      string
	Threshold    int
	PollPeriod   int
)

//CheckEnvironment check to make sure that all required env are present
func CheckEnvironment() error {
	//Chek Presences of Required Environment Variables.
	SlackWebhook, sexists := os.LookupEnv("SLACK_WEBHOOK")
	GitAPIServer, aexists := os.LookupEnv("GIT_API_SERVER")
	GitToken, tokexists := os.LookupEnv("GIT_TOKEN")
	GitOwner, gexists := os.LookupEnv("GIT_OWNER")
	GitRepo, rexists := os.LookupEnv("GIT_REPO")
	threshold, texists := os.LookupEnv("PR_THRESHOLD")
	pollPeriod, pexists := os.LookupEnv("POLL_PERIOD")

	if !aexists {
		GitAPIServer = "api.github.com"
	}
	if !sexists || !gexists || !rexists || !texists || !pexists || !tokexists {
		fmt.Printf("OWNER=%s, REPO= %s, Threshold= %s, POLL= %s, WEBHOOK=%s,Token = %s\n", GitOwner, GitRepo, threshold, pollPeriod, SlackWebhook, GitToken)
		return fmt.Errorf("one of the required environment variables is not set ")
	}
	_, ok := strconv.Atoi(threshold)
	if ok != nil {
		return fmt.Errorf("Threshold should be set to number of days , got :%s", threshold)
	}
	PollPeriod, ok := strconv.Atoi(pollPeriod)
	if ok != nil || PollPeriod == 0 {
		return fmt.Errorf("pollPeriod should be set to number of minutes between two polls and greater than 0 , got :%s", pollPeriod)
	}
	fmt.Printf("API = %s OWNER=%s, REPO= %s, Threshold= %s, POLL= %s ,WEBHOOK=%s\n", GitAPIServer, GitOwner, GitRepo, threshold, pollPeriod, SlackWebhook)

	return nil
}
func main() {

	err := CheckEnvironment()
	if err != nil {
		os.Exit(1)
	}
	// Initialize GithubClient and TimeTicker based on provided PollingPeriod
	// each tick pull pr list , send summary followed by details on the same slack channel
	//urlPath := fmt.Sprintf("%s/%s/%s", GitAPIServer, GitOwner, GitRepo)
	ghc := CreateGithubClient(GitAPIServer, GitOwner, GitRepo, GitToken)
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
