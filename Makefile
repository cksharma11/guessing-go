PKGS := $(shell go list ./... | grep -v /vendor)

guessing:
	CGO_ENABLED=0 go build -o bin/guessing ./pkg/main/