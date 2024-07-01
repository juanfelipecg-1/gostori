package notification

import (
	"bytes"
	"github.com/juanfcgarcia/gostori/internal/mailer"
	"html/template"
)

const (
	fromEmail = "noreply@stori.com"
	subject   = "The summary of your transactions is here"
)

type MonthlySummary struct {
	TransactionCount int
}

type SendSummaryEmailParams struct {
	Email               string
	TotalBalance        float64
	AverageCredit       float64
	AverageDebit        float64
	TransactionsByMonth map[string]MonthlySummary
}

type Notifier interface {
	SendSummaryEmail(params SendSummaryEmailParams) error
}

type notifications struct {
	mailSvc mailer.Mailer
}

func NewNotifier(mailSvc mailer.Mailer) Notifier {
	return &notifications{mailSvc: mailSvc}
}

func (n *notifications) SendSummaryEmail(params SendSummaryEmailParams) error {
	tmpl, err := template.ParseFiles("templates/summary_email.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, params)
	if err != nil {
		return err
	}

	mailParams := &mailer.MailParams{
		To:      params.Email,
		Subject: subject,
		From:    fromEmail,
		Body:    body.String(),
	}

	return n.mailSvc.Send(mailParams)
}
