package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
)

type smtpMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

var _ Mailer = (*smtpMailer)(nil)

func NewSMTPMailer(host string, port int, email, password, from string) *smtpMailer {
	return &smtpMailer{
		host:     host,
		port:     port,
		username: email,
		password: password,
		from:     from,
	}
}

func (m *smtpMailer) Send(ctx context.Context, to, subject, body string) error {
	message := "From: " + m.from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	return m.send(ctx, to, message)
}

func (m *smtpMailer) SendHTML(ctx context.Context, to, subject, htmlBody string) error {
	message := "From: " + m.from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		htmlBody

	return m.send(ctx, to, message)
}

func (m *smtpMailer) send(ctx context.Context, to, message string) error {
	err := smtp.SendMail(
		m.smtpAddr(),
		smtp.PlainAuth("", m.username, m.password, m.host),
		m.from, []string{to}, []byte(message),
	)

	if err != nil {
		slog.Error("failed to send email", "err", err)
		return err
	}

	slog.Info("email sent", "to", to, "message", message)
	return nil
}

func (m *smtpMailer) smtpAddr() string {
	return fmt.Sprintf("%s:%d", m.host, m.port)
}
