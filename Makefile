
all: run

run:
	go run ./cmd/... -config=dev.yml

setup:
	go install github.com/goreleaser/goreleaser@latest

build:
	$$GOPATH/bin/goreleaser build --config=.github/goreleaser.yml --clean --snapshot

clean:
	rm -r dist/ aikido || true
