# Alpine image to build
FROM alpine:latest AS builder

# Install Go for building
RUN apk add --no-cache go git

# Set working directory to /app
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the application
COPY . .
RUN go build -o noah-mqtt cmd/noah-mqtt/main.go

# scratch image to run
FROM scratch

# Copy built binaries
COPY --from=builder /app/noah-mqtt /noah-mqtt
COPY LICENSE /
COPY passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set permissions and entry point
USER gouser
ENTRYPOINT ["/noah-mqtt"]
