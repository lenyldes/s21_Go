FROM golang:tip-alpine3.24
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy && go build -o server ./cmd/server

ENTRYPOINT ["./server"]
