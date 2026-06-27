package email

import "context"

type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

type SendEmailInput struct {
	To          []string
	CC          []string
	Subject     string
	Body        string
	Attachments []Attachment
}

type EmailSender interface {
	Send(ctx context.Context, input *SendEmailInput) error
}
