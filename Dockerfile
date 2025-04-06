FROM golang:latest AS builder

WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main

# ----------

FROM debian:bookworm-slim

WORKDIR /app

# Copy binary and env
COPY --from=builder /build/ws_tester.html ./ws_tester.html
COPY --from=builder /build/main ./main
COPY --from=builder /build/.env ./.env

# Install ca-certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

# Expose HTTP port
EXPOSE 8080

ENTRYPOINT ["./main"]
