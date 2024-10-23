# .gitpod.Dockerfile
FROM golang:1.20

# Install essential packages
RUN apt-get update && apt-get install -y \
    vim \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install Hedera Go SDK
RUN go install github.com/hashgraph/hedera-sdk-go/v2@latest

# Set environment variables
ENV MIRROR_NODE_API_URL="https://mainnet.mirrornode.hedera.com/api/v1"
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

# Any additional setup can be added here
