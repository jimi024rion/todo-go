package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
)

type Client struct {
	ses         *sesv2.Client
	fromAddress string
}

func NewClient(ctx context.Context) (*Client, error) {
	endpoint := env.Cfg.SESEndpoint
	region := env.Cfg.AWSRegion

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

	return &Client{
		ses:         sesv2.NewFromConfig(cfg, clientOpts...),
		fromAddress: env.Cfg.SESFromAddress,
	}, nil
}

func (c *Client) SendRawEmail(ctx context.Context, rawMsg []byte) error {
	_, err := c.ses.SendEmail(ctx, &sesv2.SendEmailInput{
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

func (c *Client) FromAddress() string {
	return c.fromAddress
}
