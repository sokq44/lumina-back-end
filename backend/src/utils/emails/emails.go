package emails

import (
	"backend/config"
	"crypto/tls"
	"fmt"
	"log"
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
	from := config.Application.SMTP_FROM
	user := config.Application.SMTP_USER
	passwd := config.Application.SMTP_PASSWD
	host := config.Application.SMTP_HOST
	port := config.Application.SMTP_PORT

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

func (client *SmtpClient) SendEmail(receiver string, subject string, body string) error {
	auth := smtp.PlainAuth("", client.User, client.Passwd, client.Host)
	addr := fmt.Sprintf("%s:%s", client.Host, client.Port)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", client.From, receiver, subject, body))

	if err := smtp.SendMail(addr, auth, client.From, []string{receiver}, msg); err != nil {
		return err
	}

	return nil
}

func (client *SmtpClient) SendVerificationEmail(receiver string, token string) error {
	front := config.Application.FRONT_ADDR
	emailBody := fmt.Sprintf("Verification Link: %s/verify-email/%s", front, token)

	err := client.SendEmail(receiver, "Subject: Email Verification\r\n", emailBody)
	if err != nil {
		return fmt.Errorf("error while trying to send a verification email: %v", err.Error())
	}

	return nil
}
