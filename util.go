package main

import (
	"strings"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func check(e error) {
	if e != nil {
		log.Fatalf("ERROR: %v", e)
	}
}

func sliceUrl(u string) (string, string, int) {
  slice := strings.Split(u, "/")

  owner := slice[3]
  repo := slice[4]
  number, err := strconv.Atoi(slice[6])
  check(err)

  return owner, repo, number
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
