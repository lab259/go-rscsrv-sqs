package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/lab259/go-rscsrv"
	"github.com/lab259/go-rscsrv-prometheus/promhermes"
	"github.com/lab259/go-rscsrv-sqs/examples/01_trafficking_messages/services"
	h "github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
)

func main() {
	serviceStarter := rscsrv.DefaultServiceStarter(
		&services.DefaultPromService,
		&services.DefaultSQSService,
	)
	if err := serviceStarter.Start(); err != nil {
		panic(err)
	}

	router := h.DefaultRouter()
	router.Use(middlewares.RecoverableMiddleware, middlewares.LoggingMiddleware)
	router.Get("/metrics", promhermes.Handler(&services.DefaultPromService))

	app := h.NewApplication(h.ApplicationConfig{
		ServiceStarter: serviceStarter,
		HTTP: h.FasthttpServiceConfiguration{
			Bind: ":3000",
		},
	}, router)

	log.Println("Go to http://localhost:3000/metrics")

	// Updating metrics
	go func() {
		var opt int = 1
		var lastMessage string

		for {
			fmt.Println()
			switch opt {
			case 1:
				msg := strconv.Itoa(rand.Int())
				log.Println("Sending new message ", msg, "...")
				sOut, err := services.DefaultSQSService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String(string(msg)),
				})
				if err != nil {
					log.Println("with error: ", err)
				} else {
					log.Println("with success! {ID: ", *sOut.MessageId, "}")
				}
				time.Sleep(2 * time.Second)
				opt++
			case 2:
				log.Print("Receiving messages ... ")
				rOut, err := services.DefaultSQSService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				if err != nil {
					log.Println("with error: ", err)
				} else {
					log.Println("with success! amount received", len(rOut.Messages), "}")
					for i, msg := range rOut.Messages {
						log.Println("Message: #", i, " ", *msg.Body, " Len: ", len(*msg.Body))
					}
					lastMessage = *rOut.Messages[0].ReceiptHandle
				}
				opt++
			case 3:
				log.Print("Deleting last message ... ")
				_, err := services.DefaultSQSService.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: aws.String(lastMessage),
				})
				if err != nil {
					log.Println("with error: ", err)
				} else {
					log.Println("with success!")
				}
				opt = 1
			default:
				panic("invalid option")
			}
		}
	}()
	app.Start()
}
