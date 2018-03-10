# drone-sentry [![Build Status](https://travis-ci.org/SierraSoftworks/drone-sentry.svg?branch=master)](https://travis-ci.org/SierraSoftworks/drone-sentry)
**Sentry Plugin for Drone**

This is a small plugin for Drone which allows you to manage the creation of Sentry
releases and deployments from within your Drone pipelines. It provides the ability
to customize most aspects of its behaviour and aims to be a comprehensive user of
the Sentry API while keeping its interface as simple as possible.

## Usage

```yaml
pipeline:
  sentry:
    image: sierrasoftworks/drone-sentry
    # Specify your custom sentry server, if you want to
    server: https://sentry.example.org

    # Specify the organization to create releases in
    organization: my-org-short-name

    # Specify the project that you're deploying or creating a release for
    project: my-project-name

    # If your release affects more than one project, you can list them here as well
    projects:
      - my-other-project

    # Tell the plugin that we want to create a new release
    release: true
    release_version: "${DRONE_COMMIT_SHA}"
    release_url: "${DRONE_COMMIT_LINK}"

    # Tell the plugin that we want to create a new deployment
    deploy: true
    deploy_environment: "${DRONE_DEPLOY_TO}"
    deploy_name: "Deploying ${DRONE_COMMIT_SHA} to ${DRONE_DEPLOY_TO}"
    deploy_url: "api.${DRONE_DEPLOY_TO}.myservice.io"

    # Specify that we're using the SENTRY_TOKEN secret
    secrets: [ sentry_token ]
```