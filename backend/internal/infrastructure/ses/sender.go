package ses

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	emailuc "github.com/jimi024rion/todo-go/backend/internal/usecase/email"
)

type sesEmailSender struct {
	client      *sesv2.Client
	fromAddress string
}

func NewEmailSender(ctx context.Context) (emailuc.EmailSender, error) {
	endpoint := env.Cfg.SESEndpoint
	region := env.Cfg.AWSRegion
	fromAddress := env.Cfg.SESFromAddress

	var opts []func(*awsconfig.LoadOptions) error
	opts = append(opts, awsconfig.WithRegion(region))

	if endpoint != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		))
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("aws config load: %w", err)
	}

	var clientOpts []func(*sesv2.Options)
	if endpoint != "" {
		clientOpts = append(clientOpts, func(o *sesv2.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	client := sesv2.NewFromConfig(cfg, clientOpts...)

	return &sesEmailSender{
		client:      client,
		fromAddress: fromAddress,
	}, nil
}

func (s *sesEmailSender) Send(ctx context.Context, input *emailuc.SendEmailInput) error {
	rawMsg, err := buildRawMessage(s.fromAddress, input)
	if err != nil {
		return fmt.Errorf("build raw message: %w", err)
	}

	_, err = s.client.SendEmail(ctx, &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: rawMsg,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ses send email: %w", err)
	}

	return nil
}

func buildRawMessage(from string, input *emailuc.SendEmailInput) ([]byte, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	boundary := writer.Boundary()

	var raw bytes.Buffer

	raw.WriteString(fmt.Sprintf("From: %s\r\n", from))
	raw.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(input.To, ", ")))
	if len(input.CC) > 0 {
		raw.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(input.CC, ", ")))
	}
	raw.WriteString(fmt.Sprintf("Subject: %s\r\n", input.Subject))
	raw.WriteString("MIME-Version: 1.0\r\n")
	raw.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
	raw.WriteString("\r\n")

	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=UTF-8")
	textHeader.Set("Content-Transfer-Encoding", "7bit")
	textPart, err := writer.CreatePart(textHeader)
	if err != nil {
		return nil, fmt.Errorf("create text part: %w", err)
	}
	if _, err := textPart.Write([]byte(input.Body)); err != nil {
		return nil, fmt.Errorf("write text part: %w", err)
	}

	for _, att := range input.Attachments {
		attHeader := make(textproto.MIMEHeader)
		attHeader.Set("Content-Type", att.ContentType)
		attHeader.Set("Content-Transfer-Encoding", "base64")
		attHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", att.Filename))
		attPart, err := writer.CreatePart(attHeader)
		if err != nil {
			return nil, fmt.Errorf("create attachment part: %w", err)
		}
		encoded := base64.StdEncoding.EncodeToString(att.Data)
		if _, err := attPart.Write([]byte(encoded)); err != nil {
			return nil, fmt.Errorf("write attachment part: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	raw.Write(body.Bytes())

	return raw.Bytes(), nil
}
