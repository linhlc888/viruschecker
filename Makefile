all: build

GIT_COMMIT:=$(shell git rev-list -1 HEAD)
GIT_LAST_TAG:=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "unknown")
GIT_EXACT_TAG:=$(shell git name-rev --name-only --tags HEAD 2>/dev/null || echo "unknown")

VERSION_PATH:=github.com/linhlc888/viruschecker
LDFLAGS:=-X main.GitCommit=${GIT_COMMIT} \
	-X main.GitLastTag=${GIT_LAST_TAG} \
	-X main.GitExactTag=${GIT_EXACT_TAG}

build:
	go build -ldflags "$(LDFLAGS)" .

install:
	go install -ldflags "$(LDFLAGS)" .

releases:
	gox -ldflags "$(LDFLAGS)" -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

.PHONY: build install releases
