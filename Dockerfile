FROM golang:1.16

LABEL maintainer="L Nascimento <lurianfn@gmail.com>"

WORKDIR /app/fdm

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY pkg pkg
COPY configs configs
COPY cmd cmd

# Build the Go app
RUN env GOOS=linux GOARCH=arm go build cmd/fdm/main.go

EXPOSE 8080

CMD ["./main"]