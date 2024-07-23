# datamin-rabbitmq-consumer
Golang-based RabbitMQ consumer to consume messages from your queues, forward messages to Datamin pipelines, and trigger them in real-time.

![GitHub branch check runs](https://img.shields.io/github/check-runs/datamin-io/datamin-rabbitmq-consumer/main?color=green)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/datamin-io/datamin-rabbitmq-consumer?color=blue)
<a href="https://github.com/datamin-io/ylem?tab=Apache-2.0-1-ov-file">![Static Badge](https://img.shields.io/badge/license-Apache%202.0-blue)</a>
<a href="https://datamin.io" target="_blank">![Static Badge](https://img.shields.io/badge/website-datamin.io-blue)</a>
<a href="https://docs.datamin.io" target="_blank">![Static Badge](https://img.shields.io/badge/documentation-docs.datamin.io-blue)</a>
<a href="https://join.slack.com/t/datamincommunity/shared_invite/zt-2nawzl6h0-qqJ0j7Vx_AEHfnB45xJg2Q" target="_blank">![Static Badge](https://img.shields.io/badge/community-join%20Slack-blue)</a>

# Usage

## ENV variables

Create the new `.env` and copy the content from `.env.dist` to it and then add your values to the list of parameters:

* AMQP_SERVER_URL - the URL or your Rabbit MQ server in the format of: `amqp://guest:guest@datamin-message-broker:5672/`
* DATAMIN_API_URL - If you use the cloud version of Datamin, it is `https://api.datamin.io`, othervise add URL of your custom Datamin API instance here
* DATAMIN_API_CLIENT_ID - API Client ID. Here is how to create it: https://docs.datamin.io/datamin-api/oauth-clients
* DATAMIN_API_CLIENT_SECRET - API Client Secret. Here is how to create it: https://docs.datamin.io/datamin-api/oauth-clients
* QUEUE_NAME - Queue name, the consumer should listen to
* WORKFLOW_UUID - UUID of the pipeline to trigger in case of the new messages in the queue

For example:
```
AMQP_SERVER_URL=amqp://guest:guest@datamin-message-broker:5672/ 
DATAMIN_API_URL=https://api.datamin.io
DATAMIN_API_CLIENT_ID=be48d317-924c-4b1c-809a-026638e447b7
DATAMIN_API_CLIENT_SECRET=cfa9dd5e33cd8b2f0f604d94b5gggipgyggggguje2cefb
QUEUE_NAME=datamin_queue
WORKFLOW_UUID=2c334d54-c807-4629-99eb-a4def2455557
```

## Run the service

`make run`

or in the debug mode:

`docker-compose up --build`

By default it launches three docker containers (see `docker-compose.yml`):

* Message broker
* Producer
* Consumer

The first two are optional and are not needed in case you already have your own RabbitMQ service and data producers that populate queues. 

In this case you can run the consumer container only:

```
docker run -itd --rm \
  -e AMQP_SERVER_URL='%YOUR_AMQP_SERVER_URL%' \
  -e DATAMIN_API_URL='%DATAMIN_API_URL%' \
  -e DATAMIN_API_CLIENT_ID='%DATAMIN_API_CLIENT_ID%' \
  -e DATAMIN_API_CLIENT_SECRET='%DATAMIN_API_CLIENT_SECRET%' \
  -e QUEUE_NAME='%QUEUE_NAME%' \
  -e WORKFLOW_UUID='%WORKFLOW_UUID%' \
	--network=datamin-rabbitmq-consumer_datamin-rmq-network   \
	datamin-rabbitmq-consumer-datamin-consumer
```

Using the same command you can run more containers with consumers listening to the same or other queues.

# GUI

RabbitMQ administrative interface is available here: http://localhost:15672/

# Test flow 

## 1. Create the pipeline

As an example, we will use a simple pipeline that receives a message from RabbitMQ and sends an Email notification:

<img width="593" alt="Screenshot 2023-12-11 at 21 53 35" src="https://github.com/datamin-io/datamin-rabbitmq-consumer/assets/346908/7b6682b8-cfaa-45e4-8a35-da9e5ffbdeda">

Notification task:

<img width="783" alt="Screenshot 2023-12-11 at 21 58 51" src="https://github.com/datamin-io/datamin-rabbitmq-consumer/assets/346908/52472292-879c-483a-956d-d7db875babcd">

## 2. Send test message to the producer API:

API Endpoint: http://localhost:3240/message

API Endpoint expects the POST request with the JSON body

```
{
    "info": "some information", // just an example
    "customer_id": 1 // just an example
}
```

**Valid input**

<img width="857" alt="Screenshot 2023-12-11 at 21 55 13" src="https://github.com/datamin-io/datamin-rabbitmq-consumer/assets/346908/518fb0dd-9fcc-49ef-8a5c-a413a35f60e0">

**Invalid input. Validation error is returned**

<img width="863" alt="Screenshot 2023-12-11 at 21 55 29" src="https://github.com/datamin-io/datamin-rabbitmq-consumer/assets/346908/5ba1719a-a7a9-4021-97a0-944650c988da">

## 3. Check CLI output to make sure it was received by the producer and consumed

<img width="1207" alt="Screenshot 2023-12-11 at 21 56 43" src="https://github.com/datamin-io/datamin-rabbitmq-consumer/assets/346908/a0169217-ffed-4ef8-85c5-b79d9e7b4983">

## 4. Check you Email inbox and confirm that the message is received

<img width="330" alt="Screenshot 2023-12-11 at 21 59 20" src="https://github.com/datamin-io/datamin-rabbitmq-consumer/assets/346908/5e2aea30-eebc-4098-ae75-7681ade03eee">

# Full documentation

https://docs.datamin.io/ 

# Development

## Linter

### Install Golang linter on MacOS

``` bash
$ brew install golangci-lint
$ brew upgrade golangci-lint
```

### Check the code with it

``` bash
$ golangci-lint run
```
More information is in the official documentation: https://golangci-lint.run/

