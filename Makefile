
install:
	cd pkg/order && go mod vendor
	cd configs && go mod vendor
	cd cmd/fdm && go mod vendor
	go mod vendor

lint:
	yamllint configs
	golangci-lint run --fix

test:


build:
	go build cmd/fdm/main.go

run:
	./main
	rm main

all: install lint test build run
