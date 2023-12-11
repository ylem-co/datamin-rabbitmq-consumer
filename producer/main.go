package main

import (
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/streadway/amqp"
    "github.com/go-playground/validator/v10"
    log "github.com/sirupsen/logrus"
)


type SenderInput struct {
    Info        string  `json:"info" validate:"required"`
    CustomerID  int     `json:"customer_id" validate:"required,numeric"`
}

type IError struct {
    Field        string
    Tag          string
    Param        string
}

var Validator = validator.New()

func ValidatePostMessage(c *fiber.Ctx) error {
    var errors []*IError
    body := new(SenderInput)
    
    err := c.BodyParser(&body)
    if err != nil {
        log.Fatal(err)
    }

    err = Validator.Struct(body)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            var el IError
            el.Field = err.Field()
            el.Tag = err.Tag()
            el.Param = err.Param()
            errors = append(errors, &el)
        }
        return c.Status(fiber.StatusBadRequest).JSON(errors)
    }
    return c.Next()
}

func main() {
    amqpServerURL := os.Getenv("AMQP_SERVER_URL")
    queue := os.Getenv("QUEUE_NAME")

    connectRabbitMQ, err := amqp.Dial(amqpServerURL)
    if err != nil {
        panic(err)
    }
    defer connectRabbitMQ.Close()

    channelRabbitMQ, err := connectRabbitMQ.Channel()
    if err != nil {
        panic(err)
    }
    defer channelRabbitMQ.Close()

    app := fiber.New()
    app.Use(
        logger.New(), // add simple logger
    )

    app.Post("/message/", ValidatePostMessage, func(c *fiber.Ctx) error {
        body := new(SenderInput)
        err = c.BodyParser(&body)
        if err != nil {
            log.Fatal(err)
        }

        message := amqp.Publishing{
            ContentType: "application/json",
            Body:        []byte(c.Body()),
        }

        _, err = channelRabbitMQ.QueueDeclare(
            queue,           // queue name
            true,            // durable
            false,           // auto delete
            false,           // exclusive
            false,           // no wait
            nil,             // arguments
        )
        if err != nil {
            panic(err)
        }

        if err := channelRabbitMQ.Publish(
            "",              // exchange
            queue,           // queue name
            false,           // mandatory
            false,           // immediate
            message,         // message to publish
        ); err != nil {
            return err
        }

        return nil
    })

    log.Fatal(app.Listen(":3240"))
}
