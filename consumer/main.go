package main

import (
    "os"

    "github.com/streadway/amqp"
    log "github.com/sirupsen/logrus"
)

func main() {
    amqpServerURL := os.Getenv("AMQP_SERVER_URL")
    queueName := os.Getenv("QUEUE_NAME")
    workflowUuid := os.Getenv("WORKFLOW_UUID")
    apiUrl := os.Getenv("DATAMIN_API_URL")
    apiClientID := os.Getenv("DATAMIN_API_CLIENT_ID")
    apiClientSecret := os.Getenv("DATAMIN_API_CLIENT_SECRET")

    // Create a new RabbitMQ connection.
    connectRabbitMQ, err := amqp.Dial(amqpServerURL)
    if err != nil {
        panic(err)
    }
    defer connectRabbitMQ.Close()

    // Opening a channel to the RabbitMQ instance over
    // the connection we have already established.
    channelRabbitMQ, err := connectRabbitMQ.Channel()
    if err != nil {
        panic(err)
    }
    defer channelRabbitMQ.Close()

    log.Println("Successfully connected to RabbitMQ")
    log.Println("Waiting for messages from " + queueName)

    _, err = channelRabbitMQ.QueueDeclare(
        queueName,       // queue name
        true,            // durable
        false,           // auto delete
        false,           // exclusive
        false,           // no wait
        nil,             // arguments
    )
    if err != nil {
        panic(err)
    }

    messages, err := channelRabbitMQ.Consume(
        queueName,       // queue name
        "",              // consumer
        true,            // auto-ack
        false,           // exclusive
        false,           // no local
        false,           // no wait
        nil,             // arguments
    )
    if err != nil {
        log.Println(err)
    }

    apiClient := NewWorkflowClient(apiUrl, apiClientID, apiClientSecret, "dtmntst", "dtmntst500")

    // Make a channel to receive messages into infinite loop.
    forever := make(chan bool)

    go func() {
        for message := range messages {
            runUuid, err := apiClient.RunWorkflow(workflowUuid, message.Body)
            if err != nil {
                log.Error(err)
            } else {
                log.Printf("Successfully triggered a workflow, WORKFLOW_UUID: %s, RUN_UUID: %s", workflowUuid, runUuid)
            }
        }
    }()

    <-forever
}
