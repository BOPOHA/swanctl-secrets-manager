package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

const (
	charSet = "UTF-8"
)

type EmailConfig struct {
	emailTo    string
	emailBcc   string
	emailFrom  string
	body       string
	subject    string
	subjprefix string
}

func getDefaultEmailConfig() EmailConfig {
	return EmailConfig{
		emailBcc:   *emailBcc,
		emailFrom:  *emailFrom,
		subjprefix: "vpn credentials for ",
	}
}

func (c *EmailConfig) setEmailTo(to string) {
	c.emailTo = to
	c.subject = c.subjprefix + to
}

func getSESEmailInput(config EmailConfig) *sesv2.SendEmailInput {
	fmt.Printf("input: %+v\n", config)
	return &sesv2.SendEmailInput{
		Destination: &types.Destination{
			BccAddresses: []string{config.emailBcc},
			ToAddresses:  []string{config.emailTo},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String(charSet),
						Data:    aws.String(config.body),
					},
				},
				Subject: &types.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(config.subject),
				},
			},
		},
		ReplyToAddresses: []string{config.emailBcc},
		FromEmailAddress: aws.String(config.emailFrom),
	}
}
