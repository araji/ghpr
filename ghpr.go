package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	githubAPI = "api.github.com"
	ppage     = 50
)

/*
GHClient github client for a specific repo
*/
type GHClient struct {
	owner, repo string
}

//GHPR subset of github push request info , we focus on what we need
type GHPR struct {
	Number    int       `json:"number"`
	State     string    `json:"state"`
	Title     string    `json:"title"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
}

//GetPushRequests returns list of all push requests in two maps data structures
//we check the link header for the absence of "next page" signaling that we need to stop iterating
//if any of the requests generates an error or returns a non 200 , we return error to the caller.
func (ghc *GHClient) GetPushRequests(threshold int) (over, under map[int]GHPR, err error) {
	var ghpr []GHPR
	re := regexp.MustCompile("rel=\"next\"")
	var pageNumber int = 1
	overThresholdPR := make(map[int]GHPR)
	underThresholdPR := make(map[int]GHPR)

	var lastPage bool = false
	now := time.Now()
	log.Printf("checking repo: github.com/%s/%s\n ", ghc.owner, ghc.repo)

	for lastPage == false {
		url := fmt.Sprintf("https://%s/repos/%s/%s/pulls?per_page=%d&page=%d", githubAPI, ghc.owner, ghc.repo, ppage, pageNumber)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println("unable to get request")
			err = fmt.Errorf("Failed to fetch data from url %s ", url)
			return nil, nil, err
		}
		defer resp.Body.Close()
		linkHeader := resp.Header.Get("Link")
		if !re.MatchString(linkHeader) {
			lastPage = true
		}
		pageNumber++
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &ghpr)
		for _, pr := range ghpr {
			if now.Sub(pr.CreatedAt).Hours() >= float64(threshold) {
				overThresholdPR[pr.Number] = pr
			} else {
				underThresholdPR[pr.Number] = pr
			}
		}
	}

	return overThresholdPR, underThresholdPR, nil
}
