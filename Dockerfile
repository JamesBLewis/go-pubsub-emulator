# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY src ./
RUN go mod download

# Build
RUN GOOS=linux go build -o ./dist/configure-pubsub ./cmd

FROM google/cloud-sdk:456.0.0-slim

ENV PUBSUB_PROJECT testproject
ENV PUBSUB_TOPIC testtopic
ENV PUBSUB_SUBSCRIPTION testsubscription
ENV PUBSUB_PORT 8085
ENV PUBSUB_EMULATOR_HOST ${PUBSUB_PORT}

COPY --from=builder /app/dist /bin

# Create a volume for Pub/Sub data to reside
RUN mkdir -p /var/pubsub
VOLUME /var/pubsub

# Install Java for the Pub/Sub emulator, and the emulator
RUN apt-get -yq install openjdk-8-jdk google-cloud-sdk-pubsub-emulator

COPY start.sh /
COPY wait-for-it.sh /

EXPOSE ${PUBSUB_PORT}

CMD ["./start.sh"]
