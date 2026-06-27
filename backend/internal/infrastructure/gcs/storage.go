package gcs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	iamcredentials "cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	filerepo "github.com/jimi024rion/todo-go/backend/internal/domain/file/repository"
)

type gcsStorageRepository struct {
	client              *storage.Client
	bucketName          string
	emulatorHost        string
	serviceAccountEmail string
	iamClient           *iamcredentials.IamCredentialsClient
}

func NewStorageRepository(ctx context.Context) (filerepo.StorageRepository, error) {
	emulatorHost := os.Getenv("STORAGE_EMULATOR_HOST")

	var opts []option.ClientOption
	if emulatorHost != "" {
		opts = append(opts, option.WithoutAuthentication())
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("gcs client init: %w", err)
	}

	var saEmail string
	var iamClient *iamcredentials.IamCredentialsClient

	if emulatorHost == "" {
		saEmail, err = resolveServiceAccountEmail(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolve service account email: %w", err)
		}
		iamClient, err = iamcredentials.NewIamCredentialsClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("iam credentials client init: %w", err)
		}
	}

	return &gcsStorageRepository{
		client:              client,
		bucketName:          env.Cfg.GCSBucketName,
		emulatorHost:        emulatorHost,
		serviceAccountEmail: saEmail,
		iamClient:           iamClient,
	}, nil
}

// resolveServiceAccountEmail はADCからサービスアカウントのメールアドレスを取得する。
// SA キーファイルが指定されている場合はそこから取得し、
// Cloud Run 上ではメタデータサーバーにフォールバックする。
func resolveServiceAccountEmail(ctx context.Context) (string, error) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", fmt.Errorf("find default credentials: %w", err)
	}

	if len(creds.JSON) > 0 {
		var keyFile struct {
			ClientEmail string `json:"client_email"`
		}
		if jsonErr := json.Unmarshal(creds.JSON, &keyFile); jsonErr == nil && keyFile.ClientEmail != "" {
			return keyFile.ClientEmail, nil
		}
	}

	// Cloud Run のメタデータサーバーからサービスアカウントメールを取得
	email, err := metadata.EmailWithContext(ctx, "default")
	if err != nil {
		return "", fmt.Errorf("get service account email from metadata: %w", err)
	}
	return email, nil
}

func (r *gcsStorageRepository) GetDownloadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	if r.emulatorHost != "" {
		encoded := url.PathEscape(objectName)
		return fmt.Sprintf("http://%s/v0/b/%s/o/%s?alt=media", r.emulatorHost, r.bucketName, encoded), nil
	}

	signFunc := func(payload []byte) ([]byte, error) {
		resp, err := r.iamClient.SignBlob(ctx, &credentialspb.SignBlobRequest{
			Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", r.serviceAccountEmail),
			Payload: payload,
		})
		if err != nil {
			return nil, fmt.Errorf("sign blob: %w", err)
		}
		return resp.SignedBlob, nil
	}

	signedURL, err := storage.SignedURL(r.bucketName, objectName, &storage.SignedURLOptions{
		GoogleAccessID: r.serviceAccountEmail,
		SignBytes:      signFunc,
		Method:         "GET",
		Expires:        time.Now().Add(expiry),
		Scheme:         storage.SigningSchemeV4,
	})
	if err != nil {
		return "", fmt.Errorf("sign url: %w", err)
	}
	return signedURL, nil
}
