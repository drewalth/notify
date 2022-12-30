// Package notify provides helpers for setting up, maintaining, and
// sending Push notifications with the AWS SNS SDK
//
// wip
package notify

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Client struct {
	sns    *sns.SNS
	appArn *string
}

type Alert struct {
	Title        *string        `json:"title"`
	Body         *string        `json:"body"`
	LocKey       *string        `json:"loc-key,omitempty"`
	LocArgs      *[]interface{} `json:"loc-args,omitempty"`
	ActionLocKey *string        `json:"action-loc-key,omitempty"`
}

type Push struct {
	Alert *Alert      `json:"alert,omitempty"`
	Sound *string     `json:"sound,omitempty"`
	Data  interface{} `json:"custom_data,omitempty"`
	Badge *int        `json:"badge,omitempty"`
}

type iosPush struct {
	APS Push `json:"aps"`
}

type wrapper struct {
	APNS        string `json:"APNS"`
	APNSSandbox string `json:"APNS_SANDBOX"`
	Default     string `json:"default"`
}

// Create a new SNS client
func NewClient(applicationArn string, sessionOptions *session.Options) *Client {
	client := new(Client)
	client.appArn = &applicationArn

	sess := session.Must(session.NewSessionWithOptions(*sessionOptions))
	client.sns = sns.New(sess)
	return client
}

// Sends a message to a specific Endpoint ARN
func (client *Client) Send(arn string, data *Push) (*sns.PublishOutput, error) {
	msg := wrapper{}
	ios := iosPush{
		APS: *data,
	}
	b, err := json.Marshal(ios)
	if err != nil {
		return &sns.PublishOutput{}, err
	}
	msg.APNS = string(b[:])
	msg.APNSSandbox = msg.APNS

	pushData, err := json.Marshal(msg)
	if err != nil {
		return &sns.PublishOutput{}, err
	}
	m := string(pushData[:])
	params := &sns.PublishInput{
		Message:          aws.String(m),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(arn),
	}
	result, err := client.sns.Publish(params)

	return result, err
}

// Get the Endpoint ARN for the provided device token
func (client *Client) GetTokenArn(deviceToken string) (string, error) {

	var input = sns.ListEndpointsByPlatformApplicationInput{
		PlatformApplicationArn: client.appArn,
	}

	result, err := client.sns.ListEndpointsByPlatformApplication(&input)

	if err != nil {
		return "", err
	}

	// this will take forever if there are a lot of
	// app installs/subscribers
	for _, t := range result.Endpoints {
		var tAtt = t.Attributes["Token"]
		if *tAtt == deviceToken {
			return *t.EndpointArn, nil
		}
	}

	return "", errors.New("no valid match")
}

// Unregister SNS Endpoint
func (client *Client) Unregister(arn string) error {
	params := &sns.DeleteEndpointInput{
		EndpointArn: aws.String(arn),
	}
	_, err := client.sns.DeleteEndpoint(params)
	return err
}

// Register the user device token into SNS
func (client *Client) Register(deviceToken string) (string, error) {

	params := &sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: client.appArn,
		Token:                  aws.String(deviceToken),
		Attributes: map[string]*string{
			"Token":   aws.String(deviceToken),
			"Enabled": aws.String("true"),
		},
	}

	result, err := client.sns.CreatePlatformEndpoint(params)

	return *result.EndpointArn, err
}
