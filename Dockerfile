# Build Stage
FROM golang:1.23.4-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux \
    GOARCH=arm64

RUN go build -ldflags "-s -w" -o app ./cmd

# Runtime Stage
FROM cgr.dev/chainguard/wolfi-base

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

RUN apk add --no-cache curl \
    && chmod +x ./app \
    && addgroup --system appgroup \
    && adduser --system --ingroup appgroup appuser

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl --fail http://localhost:8080/health || exit 1

ENTRYPOINT ["./app"]