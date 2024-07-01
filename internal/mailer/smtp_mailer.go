package mailer

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

type MailParams struct {
	To      string
	From    string
	Subject string
	Body    string
}

type Mailer interface {
	Send(*MailParams) error
}

type smtpMailer struct {
	host string
	port string
}

func NewSMTPMailer(host, port string) Mailer {
	return &smtpMailer{
		host: host,
		port: port,
	}
}

func (m *smtpMailer) Send(params *MailParams) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", params.From)
	msg.SetHeader("To", params.To)
	msg.SetHeader("Subject", params.Subject)
	msg.SetBody("text/html", params.Body)

	port, err := strconv.Atoi(m.port)
	if err != nil {
		return err
	}

	d := gomail.Dialer{Host: m.host, Port: port}
	return d.DialAndSend(msg)
}
