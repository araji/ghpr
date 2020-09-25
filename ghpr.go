package main

import "fmt"

/*
GHClient github client for a specific repo
*/
type GHClient struct {
	owner, repo string
}

//GetPushRequests return list of all push requests
func (ghc *GHClient) GetPushRequests() {
	fmt.Printf("\nchecking repo: github.com/%s/%s\n ", ghc.owner, ghc.repo)
}
