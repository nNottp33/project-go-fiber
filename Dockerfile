# Stage 1: Build
FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download && go mod verify

ENV PATH="/app:${PATH}"
RUN apk update && apk add --no-cache curl
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s && mv ./bin/air /app

COPY . /app
RUN go build -o server server.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY .env .

ENTRYPOINT ["/app/server", "/app/air"]
CMD ["/app/air"]