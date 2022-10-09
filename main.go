package main

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

type Args struct {
	Server  string `envconfig:"server"`
	APIKey  string `envconfig:"auth_token"`
	Org     string `envconfig:"organization"`
	Project string `envconfig:"project"`
	Environ string `envconfig:"environment"`
	Version string `envconfig:"version"`
}

type DroneVars struct {
	BuildNumber   int    `envconfig:"build_number"`
	BuildFinished string `envconfig:"build_finished"`
	BuildStatus   string `envconfig:"build_status"`
	BuildLink     string `envconfig:"build_link"`
	CommitSha     string `envconfig:"commit_sha"`
	CommitBranch  string `envconfig:"commit_branch"`
	CommitAuthor  string `envconfig:"commit_author"`
	CommitLink    string `envconfig:"commit_link"`
	CommitMessage string `envconfig:"commit_message"`
	JobStarted    int64  `envconfig:"job_started"`
	Repo          string `envconfig:"build_link"`
	RepoLink      string `envconfig:"repo_link"`
	System        string
}

var version = "v1.0.0"

func main() {
	app := cli.NewApp()

	app.Name = "drone-sentry"
	app.Usage = "A Drone plugin which allows you to perform actions against your Sentry instance"

	app.Author = "Benjamin Pannell"
	app.Email = "admin@sierrasoftworks.com"
	app.Copyright = "Sierra Softworks © 2018"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo",
			Usage:  "git repo name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "prev.commit.sha",
			Usage:  "previous git commit sha",
			EnvVar: "DRONE_PREV_COMMIT_SHA",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "git commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git commit author",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git commit author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},

		cli.StringFlag{
			Name:   "sentry.server",
			Usage:  "sentry server URL",
			Value:  "https://app.getsentry.com",
			EnvVar: "PLUGIN_SERVER,SENTRY_SERVER",
		},
		cli.StringFlag{
			Name:   "sentry.token",
			Usage:  "sentry server access token",
			EnvVar: "PLUGIN_TOKEN,SENTRY_TOKEN",
		},
		cli.StringFlag{
			Name:   "sentry.organization",
			Usage:  "sentry organization short-name",
			EnvVar: "PLUGIN_ORGANIZATION",
		},

		cli.BoolFlag{
			Name:   "release",
			Usage:  "create a new Sentry release",
			EnvVar: "PLUGIN_RELEASE",
		},
		cli.BoolFlag{
			Name:   "deploy",
			Usage:  "create a new Sentry deployment",
			EnvVar: "PLUGIN_DEPLOY",
		},

		cli.StringFlag{
			Name:   "project",
			Usage:  "sentry project affected by this release",
			EnvVar: "PLUGIN_PROJECT",
		},
		cli.StringSliceFlag{
			Name:   "projects",
			Usage:  "sentry projects affected by this release",
			EnvVar: "PLUGIN_PROJECTS",
			Value:  &cli.StringSlice{},
		},

		cli.StringFlag{
			Name:   "release.version",
			Usage:  "the version of the release",
			EnvVar: "PLUGIN_RELEASE_VERSION,PLUGIN_DEPLOY_VERSION,PLUGIN_VERSION",
		},
		cli.StringFlag{
			Name:   "release.url",
			Usage:  "the url for viewing the release",
			EnvVar: "PLUGIN_RELEASE_URL",
		},

		cli.StringFlag{
			Name:   "deploy.environment",
			Usage:  "the environment that a release was deployed to",
			EnvVar: "PLUGIN_DEPLOY_ENVIRONMENT",
		},
		cli.StringFlag{
			Name:   "deploy.name",
			Usage:  "the name of a deployment",
			EnvVar: "PLUGIN_DEPLOY_NAME",
		},
		cli.StringFlag{
			Name:   "deploy.url",
			Usage:  "the url for viewing a deployment",
			EnvVar: "PLUGIN_DEPLOY_URL",
		},
	}

	var client Client

	app.Before = func(c *cli.Context) error {
		if c.String("sentry.token") == "" {
			return errors.New("must specify sentry.token")
		}

		if c.String("sentry.organization") == "" {
			return errors.New("must specify sentry.organization")
		}

		conf := &ClientConfig{
			Server:       c.String("sentry.server"),
			Organization: c.String("sentry.organization"),
			Token:        c.String("sentry.token"),
		}

		client = NewClient(conf)

		if !c.Bool("release") && !c.Bool("deploy") {
			return errors.New("must specify either release, deploy or both")
		}

		if len(StripEmptyStrings(append(c.StringSlice("projects"), c.String("project")))) == 0 {
			return errors.New("must specify at least one project")
		}

		if c.String("commit.sha") == "" {
			return errors.New("must specify commit sha")
		}

		if c.String("commit.ref") == "" {
			return errors.New("must specify commit ref")
		}

		if c.String("release.version") == "" {
			return errors.New("must specify version")
		}

		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "release",
			Usage: "Creates a new release",
			Action: func(c *cli.Context) error {
				result, err := client.NewRelease(&ReleaseDetails{
					Projects: []Ref{
						Ref{
							Name: c.GlobalString("project"),
						},
					},
					Version:  DefaultString(c.GlobalString("release.version"), c.GlobalString("commit.sha")),
					Ref:      c.GlobalString("commit.ref"),
					URL:      c.GlobalString("release.url"),
					Refs: []Ref{
						Ref{
							Repository:        c.GlobalString("repo"),
							CommitSHA:         c.GlobalString("commit.sha"),
							PreviousCommitSHA: c.GlobalString("prev.commit.sha"),
						},
					},
				})

				if result != nil {
					log.Println("Got response:")
					enc := json.NewEncoder(c.App.Writer)
					enc.SetIndent("", "  ")
					enc.Encode(result)
				}

				return err
			},
		},
		cli.Command{
			Name:  "deploy",
			Usage: "Creates a new deployment",
			Action: func(c *cli.Context) error {
				result, err := client.NewDeploy(&DeployDetails{
					Version:     c.GlobalString("release.version"),
					Environment: c.GlobalString("deploy.environment"),
					Name:        c.GlobalString("deploy.name"),
					URL:         c.GlobalString("deploy.url"),
				})

				if result != nil {
					log.Println("Got response:")
					enc := json.NewEncoder(c.App.Writer)
					enc.SetIndent("", "  ")
					enc.Encode(result)
				}

				return err
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("release") {
			if err := c.App.Command("release").Run(c); err != nil {
				return errors.Wrap(err, "task failed")
			}
		}

		if c.Bool("deploy") {
			if err := c.App.Command("deploy").Run(c); err != nil {
				return errors.Wrap(err, "task failed")
			}
		}

		return nil
	}

	app.RunAndExitOnError()
}
