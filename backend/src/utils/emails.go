package utils

import (
	"backend/config"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// TODO:
// Better looking email verification template.

type SmtpClient struct {
	From   string
	User   string
	Passwd string
	Host   string
	Port   string
}

var Smtp SmtpClient

func init() {
	from := config.AppContext["SMTP_FROM"]
	user := config.AppContext["SMTP_USER"]
	passwd := config.AppContext["SMTP_PASSWD"]
	host := config.AppContext["SMTP_HOST"]
	port := config.AppContext["SMTP_PORT"]

	auth := smtp.PlainAuth("", user, passwd, host)

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := smtp.Dial(addr)
	if err != nil {
		log.Fatalf("failed to connect to the SMTP server: %v", err.Error())
	}
	defer conn.Close()

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

	Smtp.From = from
	Smtp.User = user
	Smtp.Passwd = passwd
	Smtp.Host = host
	Smtp.Port = port

	log.Printf("connected to smtp server: %v:%v", host, port)
}

func (client *SmtpClient) SendEmail(receiver string, subject string, body string) error {
	auth := smtp.PlainAuth("", client.User, client.Passwd, client.Host)

	addr := fmt.Sprintf("%s:%s", client.Host, client.Port)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", client.From, receiver, subject, body))

	if err := smtp.SendMail(addr, auth, client.From, []string{receiver}, msg); err != nil {
		return err
	}

	log.Println("Email sent to:", receiver)

	return nil
}

func (client *SmtpClient) SendVerificationEmail(receiver string, token string) error {
	emailBody := fmt.Sprintf("Verification Link: %s/verify-email/%s", config.AppContext["FRONT_ADDR"], token)

	err := client.SendEmail(receiver, "Subject: Email Verification\r\n", emailBody)
	if err != nil {
		return fmt.Errorf("error while trying to send a verification email: %v", err.Error())
	}

	return nil
}
