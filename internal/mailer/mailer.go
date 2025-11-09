package mailer

import (
	"context"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
	SendHTML(ctx context.Context, to, subject, htmlBody string) error
}
