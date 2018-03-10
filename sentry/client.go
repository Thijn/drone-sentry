package sentry

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
)

type Client interface {
	CreateRelease(*Release) error
	CreateDeploy(*Deploy) error
}

type client struct {
	sentryServer string
	org          string
	project      string
	apiKey       string
}

func NewClient(sentryServer, apiKey, org, project string) Client {
	if sentryServer == "" {
		sentryServer = "https://sentry.io"
	}

	return &client{
		sentryServer,
		org,
		project,
		apiKey,
	}
}

func (c *client) CreateRelease(msg *Release) error {

	body, _ := json.Marshal(msg)
	buf := bytes.NewReader(body)

	req, err := http.NewRequest("POST",
		path.Join(c.sentryServer, "api/0/projects/", c.org, c.project, "releases")+"/",
		buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	// 201 = Created, 208 = Already Exists
	if resp.StatusCode != 201 && resp.StatusCode != 208 {
		t, _ := ioutil.ReadAll(resp.Body)
		return &Error{resp.StatusCode, string(t)}
	}

	return nil
}

func (c *client) CreateDeploy(msg *Deploy) error {

	body, _ := json.Marshal(msg)
	buf := bytes.NewReader(body)

	req, err := http.NewRequest("POST",
		path.Join(c.sentryServer, "api/0/projects/", c.org, c.project, "releases", msg.Name, "deploys")+"/",
		buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	// 201 = Created, 208 = Already Exists
	if resp.StatusCode != 201 && resp.StatusCode != 208 {
		t, _ := ioutil.ReadAll(resp.Body)
		return &Error{resp.StatusCode, string(t)}
	}

	return nil
}
