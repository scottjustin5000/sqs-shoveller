package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqssvc *sqs.SQS

func getApproximateQueueCount(sourceQueue string) int {
	params := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(sourceQueue),
		AttributeNames: []*string{
			aws.String("ApproximateNumberOfMessages"),
			aws.String("ApproximateNumberOfMessagesDelayed"),
			aws.String("ApproximateNumberOfMessagesNotVisible"),
		},
	}
	resp, err := sqssvc.GetQueueAttributes(params)
	if err != nil {
		log.Fatal(err)
	}
	var total = 0
	for attrib := range resp.Attributes {
		prop := resp.Attributes[attrib]
		i, _ := strconv.Atoi(*prop)
		fmt.Println(attrib, i)
		total += i
	}
	return total
}

func moveMessage(destinationQueue string, message sqs.Message) error {
	sendParams := &sqs.SendMessageInput{
		MessageBody: message.Body,
		QueueUrl:    aws.String(destinationQueue),
	}
	_, err := sqssvc.SendMessage(sendParams)
	return err

}

func removeMessage(sourceQueue string, message sqs.Message) error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sourceQueue),
		ReceiptHandle: message.ReceiptHandle,
	}
	_, err := sqssvc.DeleteMessage(deleteParams)
	return err
}

func shovel(sourceQueue string, destinationQueue string) {
	queueCount := getApproximateQueueCount(sourceQueue)
	if queueCount == 0 {
		fmt.Println("No more messages to shovel")
		return
	}
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sourceQueue),
		MaxNumberOfMessages: aws.Int64(3),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(5),
	}
	receiveResp, err := sqssvc.ReceiveMessage(receiveParams)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for _, message := range receiveResp.Messages {

		fmt.Println("moving message...")
		moveError := moveMessage(destinationQueue, *message)
		if moveError != nil {
			log.Println(moveError)
			os.Exit(1)
		}
		deleteError := removeMessage(sourceQueue, *message)
		if deleteError != nil {
			log.Println(deleteError)
			os.Exit(1)
		}
	}
	shovel(sourceQueue, destinationQueue)

}

func main() {
	s := flag.String("s", "", "source queue")
	d := flag.String("d", "", "destination queue")
	r := flag.String("r", "", "aws region")

	flag.Parse()

	if *s == "" || *d == "" || *r == "" {
		log.Fatal("Source queue, destination queue, and region required")
	}
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(*r),
	})

	sqssvc = sqs.New(sess)
	shovel(*s, *d)
}
