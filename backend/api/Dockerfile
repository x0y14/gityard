FROM golang:1.24-bookworm AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o apiserver

FROM debian:bookworm-slim
# RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    # ca-certificates && \
    # rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/apiserver /app/apiserver

# Run the web service on container startup.
CMD ["/app/apiserver"]