package docker_backend

import (
  "strings"

  "gopkg.in/alecthomas/kingpin.v1"
  "github.com/samalba/dockerclient"
  log "github.com/Sirupsen/logrus"

  "github.com/nickschuch/close/backend"
)

var (
	cliDocker = kingpin.Flag("docker-endpoint", "The Docker backend connection.").Default("unix://var/run/docker.sock").OverrideDefaultFromEnvar("DOCKER_HOST").String()
)

type DockerBackend struct {}

func init() {
	backend.Register("docker", &DockerBackend{})
}

func (o *DockerBackend) List(id string) map[string]string {
  var list map[string]string

  // Loop over containers.
  client := getClient()
	containers, err := client.ListContainers(false, false, "")
  if err != nil {
		log.Fatalf("Cannot get a list of containers from the Docker backend.")
	}

	// Check if container has environment variable.
	for _, c := range containers {
		container, _ := client.InspectContainer(c.Id)

		// We try to find the domain environment variable. If we don't have one
		// then we have nothing left to do with this container.
		url := getContainerEnv(id, container.Config.Env)
		if len(url) <= 0 {
			continue
		}

    // Add it to the list.
    list[c.Id] = url
  }

  return list
}

func (o *DockerBackend) Close(id string) {
  client := getClient()
  client.RemoveContainer(id, true, true)
}

func getClient() *dockerclient.DockerClient {
  client, err := dockerclient.NewDockerClient(*cliDocker, nil)
  if err != nil {
		log.Fatalf("Cannot establish a Docker client")
	}
  return client
}

func getContainerEnv(key string, envs []string) string {
	for _, env := range envs {
		if strings.Contains(env, key) {
			envValue := strings.Split(env, "=")
			return envValue[1]
		}
	}
	return ""
}
