FROM golang:1.25 AS build

WORKDIR /build

COPY go.mod go.sum .
RUN go mod download

COPY . .

# TODO: refactor
RUN <<EOF
    go install github.com/swaggo/swag/cmd/swag@v1.16.6
    swag init
    go mod tidy

    go build -tags=swagger -o app
EOF
    
FROM ubuntu:24.04

WORKDIR /app

COPY --from=build /build/app .
COPY --from=build /build/platform/migrations ./migrations

CMD ["./app"]
