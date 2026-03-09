FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o deckhand ./cmd/deckhand

FROM alpine:latest
RUN apk add --no-cache rsync openssh-client
COPY --from=builder /app/deckhand /usr/local/bin/deckhand
ENTRYPOINT ["/usr/local/bin/deckhand"]
