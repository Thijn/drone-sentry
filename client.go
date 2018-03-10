package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type ClientConfig struct {
	Server       string
	Organization string
	Token        string
}

type Client interface {
	NewRelease(details *ReleaseDetails) (interface{}, error)
	NewDeploy(details *DeployDetails) (interface{}, error)
}

func NewClient(config *ClientConfig) Client {
	if config.Server == "" {
		config.Server = "https://app.getsentry.com"
	}

	return &client{
		Config: config,
	}
}

type client struct {
	Config *ClientConfig
}

func (c *client) request(method, url string, payload interface{}) (interface{}, error) {
	var body *bytes.Buffer
	if payload != nil {
		body = bytes.NewBuffer([]byte{})
		enc := json.NewEncoder(body)
		enc.SetIndent("", "  ")
		if err := enc.Encode(payload); err != nil {
			return nil, errors.Wrap(err, "failed to encode json request")
		}

		log.Println(body.String())
		log.Println()
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf8")
	}

	if c.Config.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Config.Token))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	var output interface{}
	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
			return nil, errors.Wrap(err, "failed to read json response")
		}
	} else {
		t, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response")
		}

		output = string(t)
	}

	if res.StatusCode >= 400 {
		return output, errors.Errorf("request failed with status %s", res.Status)
	}

	return output, nil
}

func (c *client) buildURL(paths ...string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(c.Config.Server, "/"), strings.TrimLeft(path.Join(paths...)+"/", "/"))
}
