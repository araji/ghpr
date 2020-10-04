package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
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
	//fmt.Println("not implemented")
}

func TestValidMultipleRequests(t *testing.T) {
	testGitRepo = "all"
	ts := fakeServer(testGitOwner, testGitRepo)
	defer ts.Close()
	testAPIServer := ts.URL
	os.Setenv("GIT_API_SERVER", testAPIServer)
	ghc := CreateGithubClient(testAPIServer, testGitOwner, testGitRepo, testGitToken)
	over, under, err := ghc.GetPushRequests(testThreshold * 24)
	if err != nil {
		t.Fatalf("testing full pull request parsing failed : %v", err)
	}
	if len(over)+len(under) != 30 {
		t.Errorf("expected %d PR's but got %d", 30, len(over)+len(under))
	}
}
func TestValidOnePullRequest(t *testing.T) {
	testGitRepo = "one"
	ts := fakeServer(testGitOwner, testGitRepo)
	defer ts.Close()
	testAPIServer := ts.URL
	os.Setenv("GIT_API_SERVER", testAPIServer)
	ghc := CreateGithubClient(testAPIServer, testGitOwner, testGitRepo, testGitToken)
	over, under, err := ghc.GetPushRequests(testThreshold * 24)
	if err != nil {
		t.Fatalf("testing valid pull request parsing failed : %v", err)
	}
	if len(over)+len(under) != 1 {
		t.Errorf("expected %d PR's but got %d", 1, len(over)+len(under))
	}
	//fmt.Printf("length %d, %d ", len(over), len(under))
}

func TestBadRepo(t *testing.T) {
	testGitRepo = "norepo"
	ts := fakeServer(testGitOwner, testGitRepo)
	defer ts.Close()
	testAPIServer := ts.URL
	os.Setenv("GIT_API_SERVER", testAPIServer)
	ghc := CreateGithubClient(testAPIServer, testGitOwner, testGitRepo, testGitToken)
	_, _, err := ghc.GetPushRequests(testThreshold * 24)
	if err == nil {
		t.Errorf("expected an error for repositoruy not found")
	}

}
func TestNoPullRequest(t *testing.T) {
	testGitRepo = "none"
	ts := fakeServer(testGitOwner, testGitRepo)
	defer ts.Close()
	testAPIServer := ts.URL
	os.Setenv("GIT_API_SERVER", testAPIServer)
	ghc := CreateGithubClient(testAPIServer, testGitOwner, testGitRepo, testGitToken)
	over, under, err := ghc.GetPushRequests(testThreshold * 24)
	if err != nil {
		t.Fatalf("testing valid pull request parsing failed : %v", err)
	}
	if len(over)+len(under) != 0 {
		t.Errorf("expected %d PR's but got %d", 0, len(over)+len(under))
	}
	fmt.Printf("length %d, %d ", len(over), len(under))
}

func fakeServer(owner, repo string) *httptest.Server {
	url := fmt.Sprintf("/repos/%s/%s/pulls/", owner, repo)
	mux := http.NewServeMux()
	mux.HandleFunc(url, pullRequest)
	srv := httptest.NewServer(mux)
	return srv
}

//TODO: test multiple page
//TODO: add tls support
func pullRequest(w http.ResponseWriter, r *http.Request) {
	var fid *os.File
	var err error
	fragments := strings.Split(r.URL.Path, "/")
	// -3 to get the folder before last
	base := fragments[len(fragments)-3]
	switch base {
	case "none":
		{
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Link", "rel=\"next\"")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("[]"))
			return
		}
	case "one":
		{
			fid, err = os.Open("testdata/pr1.json")
			if err != nil {
				log.Fatalf("unable to open file , error = %s", err)
			}
		}
	case "all":
		{
			fid, err = os.Open("testdata/pr.json")
			if err != nil {
				log.Fatalf("unable to open file , error = %s", err)
			}
		}
	default:
		{
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	data, err := ioutil.ReadAll(fid)
	if err != nil {
		log.Fatalf("unable to open file , error = %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Link", "rel=\"next\"")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
