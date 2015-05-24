package main


import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceUrl(t *testing.T) {
  owner, repo, number := sliceUrl("https://github.com/foo/bar/pull/1")
  assert.Equal(t, "foo", owner, "Found the owner.")
  assert.Equal(t, "bar", repo, "Found the repository.")
  assert.Equal(t, 1, number, "Found the issue number.")
}
