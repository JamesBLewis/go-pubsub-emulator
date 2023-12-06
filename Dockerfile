# syntax=docker/dockerfile:1

# steup 1 build the go binary
FROM golang:1.21 AS gobuild

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY src ./
RUN go mod download

# Build for arm
RUN CGO_ENABLED=0 go build -o ./dist/configure-pubsub ./cmd

FROM google/cloud-sdk:emulators

ENV PUBSUB_PROJECT testproject
ENV PUBSUB_TOPIC testtopic
ENV PUBSUB_SUBSCRIPTION testsubscription
ENV PUBSUB_PORT 8085
ENV PUBSUB_EMULATOR_HOST ${PUBSUB_PORT}

COPY --from=gobuild /app/dist /
COPY start.sh /
COPY wait-for-it.sh /

# Install glibc
RUN apt-get update && apt-get install -y libc6

# Set the library path
ENV LD_LIBRARY_PATH=/lib/x86_64-linux-gnu

# Create a volume for Pub/Sub data to reside
RUN mkdir -p /var/pubsub
VOLUME /var/pubsub

EXPOSE ${PUBSUB_PORT}

CMD ["./start.sh"]
