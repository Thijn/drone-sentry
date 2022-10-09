package main

import (
	"log"
	"time"

	"github.com/pkg/errors"
)

type Commit struct {
	SHA         string    `json:"id"`
	Repository  string    `json:"repository,omitempty"`
	Message     string    `json:"message,omitempty"`
	AuthorName  string    `json:"author_name,omitempty"`
	AuthorEmail string    `json:"author_email,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
}

type Project struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Ref struct {
	CommitSHA         string `json:"commit"`
	Repository        string `json:"repository"`
	PreviousCommitSHA string `json:"previousCommit,omitempty"`
}

type ReleaseDetails struct {
	Version      string    `json:"version"`
	Ref          string    `json:"ref,omitempty"`
	URL          string    `json:"url,omitempty"`
	Projects     []Project `json:"projects"`
	DateReleased time.Time `json:"dateReleased,omitempty"`
	Commits      []Commit  `json:"commits,omitempty"`
	Refs         []Ref     `json:"refs,omitempty"`
}

func (c *client) NewRelease(details *ReleaseDetails) (interface{}, error) {
	if details.DateReleased.IsZero() {
		details.DateReleased = time.Now().UTC()
	}

	log.Println("Creating new release")
	result, err := c.request("POST", c.buildURL("api/0/organizations", c.Config.Organization, "releases"), details)
	if err != nil {
		return result, errors.Wrap(err, "failed to create new release")
	}

	return result, nil
}
