.PHONY: bootstrap
bootstrap:
	go get -u golang.org/x/lint/golint
	go mod tidy

.PHONY: clean
clean:
	rm -rf ./bin
	mkdir ./bin

.PHONY: lint
lint: bootstrap
	golint ./...
	go vet ./...

.PHONY: test
test: bootstrap
	go test ./... -v -cover -race

.PHONY: run
run: bootstrap
	go run ./cmd/*.go

.PHONY: build
build:
	go build -o bin/framer-backend-interview-test ./cmd

.PHONY: image
image:
	docker build -t url_shortner .