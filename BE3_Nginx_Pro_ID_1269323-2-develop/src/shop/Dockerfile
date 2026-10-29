FROM golang:1.25.6-alpine
WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN swag init \
          -d ./cmd,./internal \
          -g main.go \
          --parseInternal \
          --outputTypes go,json,yaml \
          -o ./docs/openapi && \
    go build -o main ./cmd/main.go

ENTRYPOINT ["./main"]
