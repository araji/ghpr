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

//GHPR github push request info we are interested in
// date (time.Time) 2020-07-20 02:17:13 +0000 UTC
// date (string)    2020-07-20 02:17:13Z
type GHPR struct {
	Number    int       `json:"number"`
	State     string    `json:"state"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

//GetPushRequests return list of all push requests
//pull requests older that threshold (hours ) will need to be reviewed
func (ghc *GHClient) GetPushRequests(threshold int) (over, under int, err error) {
	var ghpr []GHPR
	re := regexp.MustCompile("rel=\"next\"")
	var pageNumber int = 1
	var overThreshold, underThreshold int
	var lastPage bool = false
	now := time.Now()
	log.Printf("checking repo: github.com/%s/%s\n ", ghc.owner, ghc.repo)
	for lastPage == false {
		url := fmt.Sprintf("https://%s/repos/%s/%s/pulls?per_page=%d&page=%d", githubAPI, ghc.owner, ghc.repo, ppage, pageNumber)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println("unable to get request")
			err = fmt.Errorf("Failed to fetch data from url %s ", url)
			return 0, 0, err
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
				overThreshold++
			} else {
				underThreshold++
			}
		}
	}

	return overThreshold, underThreshold, nil
}
