package ses

import (
	"bytes"
	"context"
	"fmt"

	mail "github.com/wneessen/go-mail"

	emailuc "github.com/jimi024rion/todo-go/backend/internal/usecase/email"
)

type emailSender struct {
	client *Client
}

func NewEmailSender(client *Client) emailuc.EmailSender {
	return &emailSender{client: client}
}

func (s *emailSender) Send(ctx context.Context, input *emailuc.SendEmailInput) error {
	m := mail.NewMsg()

	if err := m.From(s.client.FromAddress()); err != nil {
		return fmt.Errorf("set from: %w", err)
	}
	if err := m.To(input.To...); err != nil {
		return fmt.Errorf("set to: %w", err)
	}
	if len(input.CC) > 0 {
		if err := m.Cc(input.CC...); err != nil {
			return fmt.Errorf("set cc: %w", err)
		}
	}

	m.Subject(input.Subject)
	m.SetBodyString(mail.TypeTextPlain, input.Body)

	for _, att := range input.Attachments {
		m.AttachReader(att.Filename, bytes.NewReader(att.Data),
			mail.WithFileContentType(mail.ContentType(att.ContentType)),
		)
	}

	var buf bytes.Buffer
	if _, err := m.WriteTo(&buf); err != nil {
		return fmt.Errorf("build raw message: %w", err)
	}

	if err := s.client.SendRawEmail(ctx, buf.Bytes()); err != nil {
		return err
	}

	return nil
}
