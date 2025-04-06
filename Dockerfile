FROM golang:latest as builder

WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main

# ----------

FROM debian:bullseye-slim

WORKDIR /app

# Copy binary and env
COPY --from=builder /build/main ./main
COPY --from=builder /build/.env ./.env

# Copy TLS cert and key
COPY --from=builder /build/key.pem ./key.pem
COPY --from=builder /build/cert.pem ./cert.pem

# Install ca-certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

# Expose HTTPS port
EXPOSE 443

ENTRYPOINT ["./main"]
