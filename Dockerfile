FROM golang:alpine

RUN apk add --no-cache git

WORKDIR /app/dns-docker-helper

COPY go/go.mod .
COPY go/go.sum .

RUN go mod download
COPY go/main.go .

RUN go build -o dns-docker-helper .

CMD ["dns-docker-helper"]