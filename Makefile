VERSION ?= 1.0.12
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/elsni/lagerator/args.appVersion=$(VERSION) \
	-X github.com/elsni/lagerator/args.buildCommit=$(COMMIT) \
	-X github.com/elsni/lagerator/args.buildDate=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/lgrt main.go

release:
	go build -ldflags "-s -w $(LDFLAGS)" -o bin/lgrt main.go

run:
	go run main.go

test:
	go test ./...

install: release
	sudo cp bin/lgrt /usr/local/bin/

getdata:
	cp ~/.lgrt/lgrtdata.json .

putdata:
	-mkdir ~/.lgrt
	-mv ~/.lgrt/lgrtdata.json ~/.lgrt/lgrtdata.json.bak
	cp ./lgrtdata.json ~/.lgrt/lgrtdata.json
