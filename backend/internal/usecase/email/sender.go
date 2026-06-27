package email

import "context"

type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

type SendEmailInput struct {
	To          string
	Subject     string
	Body        string
	Attachment  *Attachment
}

type EmailSender interface {
	Send(ctx context.Context, input *SendEmailInput) error
}
