# Dockerfile
# ---- Build stage ----
FROM golang:1.26.5 AS builder

WORKDIR /app

# cache deps separately so code changes don't invalidate the module cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# static binary, no CGO, smaller and portable across distros
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/api

# ---- Runtime stage ----
FROM alpine:3.20

# certs needed for outbound TLS (e.g. connecting to Neon over SSL)
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/server .

# if you run tern migrations as part of container startup, copy migrations too
COPY migrations ./migrations

EXPOSE 8080

CMD ["./server"]