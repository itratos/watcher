# Contributing Guide

## Develop

1. Install go
2. Run tests with `go test`
3. Build a dev executable with `cd cmd/watcher/ && go build .`

## Release

1. Install [goreleaser](https://github.com/goreleaser/goreleaser)
2. Update watcher's semantic version and push a new git tag
3. Test the release; from the project's root directory, run: `goreleaser --snapshot --skip-publish --rm-dist`
4. If the dry run was succesful, publish the release: `goreleaser`