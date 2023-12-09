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

COPY --from=gobuild /app/dist /

# make folder for config to be stored in
RUN mkdir -p /config
# Create a volume for Pub/Sub data to reside
RUN mkdir -p /var/pubsub
VOLUME /var/pubsub

EXPOSE ${PUBSUB_PORT}

# thanks chat gpt
CMD ["sh", "-c", "set -e; gcloud beta emulators pubsub start --project=${PUBSUB_PROJECT} --data-dir=/var/pubsub --host-port=0.0.0.0:${PUBSUB_PORT} --log-http --verbosity=debug --user-output-enabled & pid=$!; ./configure-pubsub; wait $pid"]
