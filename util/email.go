package util

import (
	"fmt"

	"github.com/srikanta0427/rest_api_design/config"
	"github.com/wneessen/go-mail"
)

func SendMail(to, subject, body string) error {
	username := config.GetString("SMTP_USERNAME", "")
	password := config.GetString("SMTP_PASSWORD", "")

	if username == "" || password == "" {
		return fmt.Errorf("SMTP credentials are missing")
	}

	msg := mail.NewMsg()
	if err := msg.From(username); err != nil {
		return err
	}
	if err := msg.To(to); err != nil {
		return err
	}
	msg.Subject(subject)
	msg.SetBodyString("text/html", body)

	client, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		return err
	}
	fmt.Println("email sent")
	return client.DialAndSend(msg)
}
