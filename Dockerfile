# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.18 AS builder

WORKDIR /src

# Download dependencies first so they cache independently of source changes.
COPY go.mod go.sum ./
RUN go mod download

# Build a static binary (no libc dependency) so it runs on a distroless base.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/nootbot .

# --- Runtime stage ---
# distroless/static ships CA certificates, needed for the HTTPS calls to
# Discord and the weather API.
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Only the binary and non-secret config are copied in
COPY --from=builder /out/nootbot /app/nootbot
COPY notify.conf /app/notify.conf

# This image runs the bot in Discord mode
ENTRYPOINT ["/app/nootbot", "discord"]
