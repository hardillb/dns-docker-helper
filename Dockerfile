FROM golang:alpine AS build

RUN apk add --no-cache git

WORKDIR /app/dns-docker-helper

COPY go/go.mod .
COPY go/go.sum .

RUN go mod download
COPY go/main.go .

RUN go build -o dns-docker-helper .

FROM alpine
WORKDIR /
COPY --from=build /app/dns-docker-helper/dns-docker-helper /dns-docker-helper

CMD ["dns-docker-helper"]