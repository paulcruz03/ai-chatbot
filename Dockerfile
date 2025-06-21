FROM golang:latest AS builder

WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main

# ----------

FROM scratch

WORKDIR /app

# Copy binary and env
COPY --from=builder /build/main ./main
COPY --from=builder /build/.env ./.env

# Expose HTTP port
EXPOSE 8080

ENTRYPOINT ["./main"]
