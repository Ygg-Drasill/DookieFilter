# Use official Go image
FROM golang:1.23 AS builder

RUN rm -rf /var/lib/apt/lists/* && \
    apt-get update --fix-missing && \
    apt-get install -y --no-install-recommends libzmq3-dev && \
    apt-get clean && rm -rf /var/lib/apt/lists/

# Set working directory
WORKDIR /app

# Mark GOPRIVATE repository
ENV GOPRIVATE=github.com/Ygg-Drasill/DookieFilter

# Configure Git authentication using secret
RUN --mount=type=secret,id=github_token \
    GITHUB_TOKEN=$(cat /run/secrets/github_token) && \
    git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=1 go build -o /mock-transmitter

# Minimal runtime image
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update --fix-missing && \
    apt-get install -y --no-install-recommends libzmq5 && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /mock-transmitter /mock-transmitter

COPY raw.jsonl ./
# Run the application
CMD ["/mock-transmitter"]
