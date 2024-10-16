package utils

import (
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

var Smtp SmtpClient

func (client *SmtpClient) OpenSmtpConnection(from string, user string, passwd string, host string, port string) (string, error) {
	auth := smtp.PlainAuth("", user, passwd, host)

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := smtp.Dial(addr)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the SMTP server: %v", err.Error())
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		return "", fmt.Errorf("failed to start TLS: %v", err)
	}

	if err = conn.Auth(auth); err != nil {
		return "", fmt.Errorf("failed to authenticate with the SMTP server: %v", err)
	}

	client.From = from
	client.User = user
	client.Passwd = passwd
	client.Host = host
	client.Port = port

	return fmt.Sprintf("connected to smtp server: %v:%v", host, port), nil
}

func (client *SmtpClient) SendEmail(receiver string, subject string, body string) error {
	auth := smtp.PlainAuth("", client.User, client.Passwd, client.Host)

	addr := fmt.Sprintf("%s:%s", client.Host, client.Port)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", client.From, receiver, subject, body))

	log.Printf("Sending email to %s with subject %s", receiver, subject)
	if err := smtp.SendMail(addr, auth, client.From, []string{receiver}, msg); err != nil {
		return err
	}
	log.Println("Email sent successfully")

	return nil
}
