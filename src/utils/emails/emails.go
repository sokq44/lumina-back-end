package emails

import (
	"backend/config"
	"backend/utils/problems"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"log"
	"net/http"
)

type Client struct {
	From   string
	Client *ses.Client
}

var emails Client

func InitEmails() {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	emails.From = config.AwsSesFrom
	emails.Client = ses.NewFromConfig(cfg)

	log.Println("Initialized the AWS SES service.")
}

func GetEmails() *Client {
	return &emails
}

func (smtpClient *Client) SendEmail(recipient string, subject string, body string) *problems.Problem {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(body),
				},
			},
			Subject: &types.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(smtpClient.From),
	}

	_, err := smtpClient.Client.SendEmail(context.TODO(), input)
	if err != nil {
		return &problems.Problem{
			Type:          problems.EmailsProblem,
			ServerMessage: fmt.Sprintf("while trying to send an email -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (smtpClient *Client) SendVerificationEmail(receiver string, token string) *problems.Problem {
	front := config.FrontAddr
	emailBody := fmt.Sprintf("Verification Link: %s/email/%s", front, token)

	err := smtpClient.SendEmail(receiver, "Email Verification", emailBody)
	if err != nil {
		return err
	}

	return nil
}

func (smtpClient *Client) SendPasswordChangeEmail(receiver string, token string) *problems.Problem {
	front := config.FrontAddr
	emailBody := fmt.Sprintf("Change your password here: %s/user/password/%s", front, token)

	err := smtpClient.SendEmail(receiver, "Change Your Password", emailBody)
	if err != nil {
		return err
	}

	return nil
}
