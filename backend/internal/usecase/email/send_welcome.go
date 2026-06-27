package email

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
)

//go:embed testdata/dummy_qr.png
var dummyQRPNG []byte

type SendWelcomeUseCase struct {
	emailSender EmailSender
}

func NewSendWelcomeUseCase(emailSender EmailSender) *SendWelcomeUseCase {
	return &SendWelcomeUseCase{emailSender: emailSender}
}

type SendWelcomeInput struct {
	To string
}

type SendWelcomeOutput struct {
	Message string
}

func (uc *SendWelcomeUseCase) Execute(ctx context.Context, input *SendWelcomeInput) (*SendWelcomeOutput, error) {
	if input.To == "" {
		return nil, fmt.Errorf("recipient is required")
	}

	data := WelcomeTemplateData{
		UserName: "ユーザー",
	}

	var bodyBuf bytes.Buffer
	if err := welcomeBodyTmpl.Execute(&bodyBuf, data); err != nil {
		return nil, fmt.Errorf("execute welcome template: %w", err)
	}

	err := uc.emailSender.Send(ctx, &SendEmailInput{
		To:      input.To,
		Subject: welcomeSubject,
		Body:    bodyBuf.String(),
		Attachment: &Attachment{
			Filename:    "qr_code.png",
			ContentType: "image/png",
			Data:        dummyQRPNG,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send welcome email: %w", err)
	}

	return &SendWelcomeOutput{
		Message: "Welcome email sent successfully",
	}, nil
}
