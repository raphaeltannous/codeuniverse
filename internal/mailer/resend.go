package mailer

import (
	"context"

	"github.com/resend/resend-go/v2"
)

type resendMailer struct {
	from string

	client *resend.Client
}

func (m *resendMailer) SendHTML(ctx context.Context, to string, subject string, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := m.client.Emails.Send(params)
	return err
}

func (m *resendMailer) Send(ctx context.Context, to, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{to},
		Subject: subject,
		Text:    body,
	}

	_, err := m.client.Emails.Send(params)
	return err
}

func NewResendMailer(apiKey string, from string) Mailer {
	return &resendMailer{
		from:   from,
		client: resend.NewClient(apiKey),
	}
}
