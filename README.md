# Google Cloud Pub/Sub Emulator

This repository contains configuration for a Docker image that contains the Google Cloud Pub/Sub emulator. It is Dockerised version of https://cloud.google.com/pubsub/docs/emulator that provides additional options for configuration.

## Usage

### Command Line
This shell statement shows the basic usage, which will create a Pub/Sub emulator running on port 8085, with a project name of "testproject", a default topic of "testtopic" and a subscription of "testsubscription".

#### Basic

```shell script
docker run --rm --publish 8085:8085 jamesblewis/go-pubsub-emulator:latest
```

Much like jesse's original pkg you can also configure the Pub/Sub project name, topic name, subscription name and exposure port:

```shell script
docker run --rm --env PUBSUB_PROJECT=myproject --env PUBSUB_TOPIC=mytopic --env PUBSUB_SUBSCRIPTION=mysubscription --env PUBSUB_PORT=10101 --publish 10101:10101 jamesblewis/go-pubsub-emulator:latest
```

#### Advanced

This image allows you to specify more complex config via yaml for when you need multiple topics and subscriptions.

```shell script
docker run --rm --env --mount source=./path/to/your/config,target=config/ --publish 10101:10101 jamesblewis/go-pubsub-emulator:latest
```

The yaml config should look something like this:

```yaml
projectID: myproject
topics:
  topic1:
    - subscription1
    - subscription2
  topic2:
    - subscription1
    - subscription2

```

_Note if you choose to specify a CONFIG file it will override any other environment variables you pass in._

### Docker Compose

This is the equivalent configuration for `docker-compose`, with custom environment variables for the project name, topic name, subscription name and Pub/Sub port.

```docker-compose
version: '3'
  services:
    pubsub:
      image: jamesblewis/go-pubsub-emulator:latest
      platform: linux/amd64
      ports:
        - "10101:10101"
      environment:
        - PUBSUB_PORT=10101
        - PUBSUB_PROJECT=myproject
        - PUBSUB_TOPIC=mytopic
        - PUBSUB_SUBSCRIPTION=mysubscription
```

## Attribution
This image is based on [jacksonjesse/pubsub-emulator](https://github.com/jacksonjesse/pubsub-emulator).
