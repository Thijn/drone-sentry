package main

import (
	"time"

	"github.com/pkg/errors"
)

type DeployDetails struct {
	Version string `json:"-"`

	Environment  string    `json:"environment"`
	Name         string    `json:"name,omitempty"`
	URL          string    `json:"url,omitempty"`
	DateStarted  time.Time `json:"dateStarted,omitempty"`
	DateFinished time.Time `json:"dateFinished,omitempty"`
}

func (c *client) NewDeploy(details *DeployDetails) (interface{}, error) {
	result, err := c.request("POST", c.buildURL("api/0/organizations", c.Config.Organization, "releases", details.Version, "deploys"), details)
	if err != nil {
		return result, errors.Wrap(err, "failed to create new deploy")
	}

	return result, nil
}
