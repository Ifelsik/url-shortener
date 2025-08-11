FROM golang:1.23-alpine3.22 AS build

WORKDIR /usr/src/url-shortener

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build ./cmd/main.go

FROM alpine:3.22

WORKDIR /url-shortener

COPY --from=build /usr/src/url-shortener/main .

CMD ["./main"]
