#!/bin/bash

export PUBSUB_EMULATOR_HOST=localhost:${PUBSUB_PORT}

./wait-for-it.sh localhost:${PUBSUB_PORT} -- ./configure-pubsub ${PUBSUB_PROJECT} ${PUBSUB_TOPIC} ${PUBSUB_SUBSCRIPTION} &

gcloud beta emulators pubsub start --project=${PUBSUB_PROJECT} --data-dir=/var/pubsub --host-port=0.0.0.0:${PUBSUB_PORT} --log-http --verbosity=debug --user-output-enabled
