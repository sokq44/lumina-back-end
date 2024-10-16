package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SmtpClient struct {
	From       string
	SMTPUser   string
	SMTPPasswd string
	SMTPHost   string
	SMTPPort   string
}

var Smtp SmtpClient

func (client *SmtpClient) OpenSmtpConnection(from string, smtpUser string, smtpPasswd string, smtpHost string, smtpPort string) (string, error) {
	auth := smtp.PlainAuth("", smtpUser, smtpPasswd, smtpHost)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	conn, err := smtp.Dial(addr)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the SMTP server: %v", err.Error())
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		return "", fmt.Errorf("failed to start TLS: %v", err)
	}

	if err = conn.Auth(auth); err != nil {
		return "", fmt.Errorf("failed to authenticate with the SMTP server: %v", err)
	}

	client.From = from
	client.SMTPUser = smtpUser
	client.SMTPPasswd = smtpPasswd
	client.SMTPHost = smtpHost
	client.SMTPPort = smtpPort

	return fmt.Sprintf("connected to smtp server: %v:%v", smtpHost, smtpPort), nil
}
