FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags='-w -s' \
    -a \
    -o web-analyzer ./cmd/server/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup -g '' appuser

COPY --from=builder /app/web-analyzer /web-analyzer
COPY --from=builder /app/web /web

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

ENV GIN_MODE=release

USER appuser:appgroup

ENTRYPOINT ["/web-analyzer"]