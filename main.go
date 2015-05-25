package main

import (
	"strings"
	"strconv"
	"net/http"

	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/google/go-github/github"
	log "github.com/Sirupsen/logrus"

	"github.com/nickschuch/close/backend"
	_ "github.com/nickschuch/close/backend/docker"
	_ "github.com/nickschuch/close/backend/ecs"
)

var (
	cliUsername = kingpin.Flag("user", "The Github account username.").Required().String()
	cliPassword = kingpin.Flag("pass", "The Github account password.").Required().String()
	cliEnv      = kingpin.Flag("env", "The Github account password.").Default("ISSUE_URL").String()
	cliBackend  = kingpin.Flag("backend", "The type of backend.").Default("docker").String()
)

type Transport struct {
	Username string
	Password string
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultTransport.RoundTrip(req)
}

func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func main() {
  kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Remove containers if a Pull Request status is closed."
	kingpin.Parse()

	// Get a list of all the environments from backend.
	b, err := backend.New(*cliBackend)
	if err != nil {
		log.Fatalf("Cannot find the backend: %v", *cliBackend)
	}

	// Loop over and check if they are still open.
	list := b.List(*cliEnv)
	for i, l := range list {
		closed := checkIssue(l)
		if ! closed {
			continue
		}

		// Close the environment with the backend indicator.
		b.Close(i)

		// Print it out to the screen for accountability.
		log.Println("Closed: ", l)
	}
}

func checkIssue(url string) bool {
	// Build a client with a token that we can use for authentication.
	t := &Transport{
		Username: *cliUsername,
		Password: *cliPassword,
	}
	clientGithub := github.NewClient(t.Client())

	owner, repo, number := sliceUrl(url)
	issue, _, err := clientGithub.Issues.Get(owner, repo, number)
	if err != nil {
		log.Fatalf("Failed to get the Github issue: %v", err)
	}

	// If the Pull Request is not closed then we have nothing left to do with
	// this container.
	if *issue.State != "closed" {
		return false
	}

	return true
}

func sliceUrl(u string) (string, string, int) {
  slice := strings.Split(u, "/")
  owner := slice[3]
  repo := slice[4]
  number, err := strconv.Atoi(slice[6])
  if err != nil {
		log.Fatalf("Cannot convert the issue to a number")
	}
  return owner, repo, number
}
