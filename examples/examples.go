package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/drewalth/notify"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	check(err)

	// the platform application arn
	appARN := os.Getenv("PLATFORM_APP_ARN")
	// the user's device token from a database or another source
	deviceToken := os.Getenv("DEVICE_TOKEN")

	client := notify.NewClient(appARN, &session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	alert := notify.Alert{
		Body:  aws.String("Alert body"),
		Title: aws.String("Alert title"),
	}

	notificationSound := "default" // TODO: fix sound option
	notificationBadge := 0

	pushData := notify.Push{
		Alert: &alert,
		Sound: &notificationSound,
		Badge: &notificationBadge,
	}

	endpointArn, err := client.GetTokenArn(deviceToken)

	check(err)

	result, err := client.Send(endpointArn, &pushData)

	check(err)
	fmt.Println(*result.MessageId)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
