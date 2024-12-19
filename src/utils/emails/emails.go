package emails

import (
	"backend/config"
	"backend/utils/problems"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
)

type SmtpClient struct {
	From   string
	User   string
	Passwd string
	Host   string
	Port   string
}

var emails SmtpClient

func InitEmails() {
	from := config.SmtpFrom
	user := config.SmtpUser
	passwd := config.SmtpPass
	host := config.SmtpHost
	port := config.SmtpPort

	auth := smtp.PlainAuth("", user, passwd, host)

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := smtp.Dial(addr)
	if err != nil {
		log.Fatalf("failed to connect to the SMTP server: %v", err.Error())
	}

	defer func(conn *smtp.Client) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close the SMTP connection: %v", err)
		}
	}(conn)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		log.Fatalf("failed to start TLS: %v", err)
	}

	if err = conn.Auth(auth); err != nil {
		log.Fatalf("failed to authenticate with the SMTP server: %v", err)
	}

	emails.From = from
	emails.User = user
	emails.Passwd = passwd
	emails.Host = host
	emails.Port = port

	log.Printf("initialized the smtp service (%v:%v)", host, port)
}

func GetEmails() *SmtpClient {
	return &emails
}

func (client *SmtpClient) SendEmail(receiver string, subject string, body string) *problems.Problem {
	auth := smtp.PlainAuth("", client.User, client.Passwd, client.Host)
	addr := fmt.Sprintf("%s:%s", client.Host, client.Port)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", client.From, receiver, subject, body))

	if err := smtp.SendMail(addr, auth, client.From, []string{receiver}, msg); err != nil {
		return &problems.Problem{
			Type:          problems.EmailsProblem,
			ServerMessage: fmt.Sprintf("while trying to send an email -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (client *SmtpClient) SendVerificationEmail(receiver string, token string) *problems.Problem {
	front := config.FrontAddr
	emailBody := fmt.Sprintf("Verification Link: %s/email/%s", front, token)

	err := client.SendEmail(receiver, "Subject: Email Verification\r\n", emailBody)
	if err != nil {
		return err
	}

	return nil
}

func (client *SmtpClient) SendPasswordChangeEmail(receiver string, token string) *problems.Problem {
	front := config.FrontAddr
	emailBody := fmt.Sprintf("Change your password here: %s/user/password/%s", front, token)

	err := client.SendEmail(receiver, "Subject: Change Your Password\r\n", emailBody)
	if err != nil {
		return err
	}

	return nil
}
