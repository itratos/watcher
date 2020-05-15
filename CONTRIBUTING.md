# Contributing Guide

## Develop

1. Install go
1. Install dependencies with `go get -u github.com/dideler/watcher`
1. Run tests with `go test`
1. Build a dev executable with `cd cmd/watcher/ && go build .`

## Update

Since this is a fork, it can benefit from certain upstream changes.

If there's a commit or merge request that's relevant, append ".patch" to the URL
and apply it to the fork as `curl -sL $PATCH_URL | git am`

## Release

1. Install [goreleaser](https://github.com/goreleaser/goreleaser)
1. Create a GitHub Personal Access Token with repo scope and export `GITHUB_TOKEN` 
1. Update watcher's semantic version and push a new git tag
1. Test the release from the project's root directory: `goreleaser --snapshot --skip-publish --rm-dist`
1. If the dry run was succesful, publish the release: `goreleaser --rm-dist`
