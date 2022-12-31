# notify

Wrapper for the AWS SNS Go SDK.

- [x] Send push notification
- [x] Register user device with SNS
- [x] Remove user device/endpoint from SNS

---

```shell
# copy demo env file 
cp ./examples/.env.example ./examples/.env

# replace values with your application ARN and device token
DEVICE_TOKEN=<my_device_token>
PLATFORM_APP_ARN=<my_app_arn>

# send yourself a notification
go run examples/examples.go
```

```go
// examples/examples.go

client := notify.NewClient(appARN, &session.Options{
	SharedConfigState: session.SharedConfigEnable,
})

alert := notify.Alert{
	Body:  aws.String("Alert body"),
	Title: aws.String("Alert title"),
}

notificationSound := "default"
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
```