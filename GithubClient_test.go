package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	testAPIServer string = "localhost"
	testGitOwner  string = "testowner"
	testGitRepo   string = "testrepo"
	testGitToken  string = "ABCDEFGHIJK123LMNO"
	testThreshold int    = 365
)

// How to test main : take stuff out of it ?
// how to use subtests to iterate over a set of test data
// goal here is to use subtests through go test -run
func TestMissingEnv(t *testing.T) {
	tt := []struct {
		name         string
		env          string
		defaultValue string
	}{
		{"webhook present", "SLACK_WEBHOOK", "WEBHOOK"},
		//{"api_server present", "GIT_API_SERVER", "SERVER"},
		{"git token present ", "GIT_TOKEN", "TOKEN"},
		{"git owner present ", "GIT_OWNER", "OWNER"},
		{"git repo present ", "GIT_REPO", "REPO"},
		{"threshold present ", "PR_THRESHOLD", "5"},
		{"poll_period present ", "POLL_PERIOD", "365"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, envlist := range tt {
				os.Setenv(envlist.env, envlist.defaultValue)
			}
			os.Unsetenv(tc.env)
			err := CheckEnvironment()
			if err == nil {
				t.Fatalf("test %s should have generated an error ", tc.name)
			}
		})
	}
}
func TestBadAuth(t *testing.T) {
	fmt.Println("testing something")
	t.Errorf("failed miserably....")

}

//origial implementation in  http://www.inanzzz.com/index.php/post/fb0m/mocking-and-testing-http-clients-in-golang
func TestValidPullRequest(t *testing.T) {
	srv := fakeServer()
	defer srv.Close()

	ghc := CreateGithubClient(testAPIServer, testGitOwner, testGitRepo, testGitToken)
	over, under, err := ghc.GetPushRequests(testThreshold * 24)
	if err != nil {
		t.Fatalf("testing valid pullrewquest parsing failed : %v", err)
	}
	fmt.Printf("length %d, %d ", len(over), len(under))
}

//TODO: multiple page
//TODO: add ssl ?
func fakeServer(owner, repo string) *httptest.Server {
	url := fmt.Sprintf("http://%s/repos/%s/%s/pulls", "localhost", owner, repo)
	handler := http.NewServeMux()
	handler.HandleFunc(url, pullRequests)
	srv := httptest.NewServer(handler)

	return srv
}

//TODO: make the request tell us about what the response should be
func pullRequests(w http.ResponseWriter, r *http.Request) {
	//load the test response and send it as a json payload

	_, _ = w.Write([]byte(r.Header.Get("FakeResponse")))
}
