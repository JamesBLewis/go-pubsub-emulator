# syntax=docker/dockerfile:1

# steup 1 build the go binary
FROM golang:1.21 AS gobuild

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY src ./
RUN go mod download

# Build
RUN GOOS=linux go build -o ./dist/configure-pubsub ./cmd

# step 2 install the pubsub emulator
FROM google/cloud-sdk:alpine AS cloud-sdk
RUN gcloud components install pubsub-emulator

# step 3 only copy what we need onto the runtime image
FROM --platform=linux/amd64 openjdk:jre-alpine

ENV PUBSUB_PROJECT testproject
ENV PUBSUB_TOPIC testtopic
ENV PUBSUB_SUBSCRIPTION testsubscription
ENV PUBSUB_PORT 8085
ENV PUBSUB_EMULATOR_HOST ${PUBSUB_PORT}

# Create a volume for Pub/Sub data to reside
RUN mkdir -p /var/pubsub
VOLUME /var/pubsub

COPY --from=gobuild /app/dist /
COPY --from=cloud-sdk /google-cloud-sdk/platform/pubsub-emulator /pubsub-emulator

COPY start.sh /
COPY wait-for-it.sh /

RUN apk --update --no-cache add tini bash

ENTRYPOINT ["/sbin/tini", "--"]

EXPOSE ${PUBSUB_PORT}

CMD ["./start.sh"]
