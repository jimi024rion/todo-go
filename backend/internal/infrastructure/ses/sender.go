package ses

import (
	"context"

	emailuc "github.com/jimi024rion/todo-go/backend/internal/usecase/email"
)

type emailSender struct {
	client *Client
}

func NewEmailSender(client *Client) emailuc.EmailSender {
	return &emailSender{client: client}
}

func (s *emailSender) Send(ctx context.Context, input *emailuc.SendEmailInput) error {
	sesInput := &SendEmailInput{
		To:      []string{input.To},
		Subject: input.Subject,
		Body:    input.Body,
	}

	if input.Attachment != nil {
		sesInput.Attachments = []Attachment{
			{
				Filename:    input.Attachment.Filename,
				ContentType: input.Attachment.ContentType,
				Data:        input.Attachment.Data,
			},
		}
	}

	return s.client.SendEmail(ctx, sesInput)
}
