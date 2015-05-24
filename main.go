package main

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/google/go-github/github"
	"github.com/samalba/dockerclient"
)

var (
	username = kingpin.Flag("user", "The Github account username.").Required().String()
	password = kingpin.Flag("pass", "The Github account password.").Required().String()
	env      = kingpin.Flag("env", "The Github account password.").Default("ISSUE_URL").String()
	docker   = kingpin.Flag("docker", "The Docker endpoint.").Default("unix:///var/run/docker.sock").String()
)

func main() {
    kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Remove Docker containers if Pull Request status is closed."
	kingpin.Parse() 

	// Build a client with a token that we can use for
	// authentication.
	t := &Transport{
		Username: *username,
		Password: *password,
	}
	clientGithub := github.NewClient(t.Client())

	// Get a connection to the Docker instance.
	clientDocker, err := dockerclient.NewDockerClient(*docker, nil)
	check(err)

	// Loop over containers.
	containers, err := clientDocker.ListContainers(false, false, "")
	check(err)

	// Check if container has environment variable.
	for _, c := range containers {
		container, _ := clientDocker.InspectContainer(c.Id)

		// We try to find the domain environment variable. If we don't have one
		// then we have nothing left to do with this container.
		url := getContainerEnv(*env, container.Config.Env)
		if len(url) <= 0 {
			continue
		}

		owner, repo, number := sliceUrl(url)
	    issue, _, err := clientGithub.Issues.Get(owner, repo, number)
	    check(err)

		// If the Pull Request is no closed then we have nothing left to do with
		// this container.
		if *issue.State != "closed" {
			continue
		}

		clientDocker.RemoveContainer(container.Id, true, true)
	}
}
