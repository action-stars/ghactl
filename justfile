set windows-shell := ["pwsh", "-NoLogo", "-NoProfile", "-Command"]

default:
    just --list

setup: tidy
    {{ if os() == "windows" { 'New-Item -Name "dist" -ItemType "directory" -Force' } else { "mkdir -p ./dist" } }}

tidy:
    go mod tidy

fmt:
    golangci-lint fmt

lint:
    golangci-lint run --fix

test:
    go test -v -cover -timeout=120s -parallel=10 ./...

build:
    go build -o ./dist -v ./...

build-all:
    goreleaser build --clean --snapshot

mdfmt:
    rumdl fmt --fixable .

mdlint:
    rumdl check .
