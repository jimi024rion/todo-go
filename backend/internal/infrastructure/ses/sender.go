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
	return s.client.SendEmail(ctx, input)
}
