FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/stream-file-locally ./cmd/api

FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /out/stream-file-locally /app/stream-file-locally
COPY docs/openapi.yaml /app/docs/openapi.yaml
RUN mkdir -p /data/stream-file-locally && chown -R app:app /data/stream-file-locally /app

USER app
EXPOSE 8080

ENTRYPOINT ["/app/stream-file-locally"]
