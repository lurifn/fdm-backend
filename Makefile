
install:
	cd pkg/order && go mod vendor
	cd configs && go mod vendor
	cd cmd/fdm && go mod vendor
	go mod vendor

lint:
	yamllint configs .github/workflows
	golangci-lint run --fix
	gofmt -l -w -s .
	golint ./...

test:


build:
	go build cmd/fdm/main.go

run:
	./main
	rm main

docker:
	docker-compose down
	docker-compose up --build

all: install lint test docker
