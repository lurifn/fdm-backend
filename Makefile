install:
	go mod tidy
	go get -u ./...

lint:
	yamllint configs .github/workflows
	golangci-lint run --fix
	gofmt -l -w -s .
	golint ./...

test:
	go test -v ./...

build:
	go build cmd/fdm/main.go

run:
	./main
	rm main

docker:
	docker-compose down
	docker-compose up --build

all: install lint test docker
