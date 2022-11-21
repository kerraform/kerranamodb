package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/id"
	modelv1 "github.com/kerraform/kerranamodb/internal/model/v1"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	ErrS3NotAllowed = errors.New("uploads to s3 are done by presigned url")
)

type d struct {
	bucket string
	logger *zap.Logger
	tracer trace.Tracer
	s3     *s3.Client
}

type DriverOpts struct {
	AccessKey    string
	Bucket       string
	Endpoint     string
	SecretKey    string
	Tracer       trace.Tracer
	UsePathStyle bool
}

func NewDriver(logger *zap.Logger, opts *DriverOpts) (driver.Driver, error) {
	if opts == nil {
		return nil, fmt.Errorf("invalid s3 credentials")
	}

	cred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, ""))
	loadOpts := []func(*config.LoadOptions) error{
		config.WithCredentialsProvider(cred),
	}

	if opts.Endpoint != "" {
		endpointResolver := &endpointResolver{
			URL: opts.Endpoint,
		}
		loadOpts = append(loadOpts, config.WithEndpointResolverWithOptions(endpointResolver))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = opts.UsePathStyle
	})

	return &d{
		bucket: opts.Bucket,
		logger: logger,
		tracer: opts.Tracer,
		s3:     s3Client,
	}, nil
}

func (d *d) DeleteLock(ctx context.Context, table string, lid id.LockID) error {
	keyPath := fmt.Sprintf("%s/%s", table, lid)
	_, err := d.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(keyPath),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *d) GetLock(ctx context.Context, table string, lid id.LockID) (driver.Info, error) {
	keyPath := fmt.Sprintf("%s/%s", table, lid)
	b := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(d.s3)
	_, err := downloader.Download(ctx, b, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(keyPath),
	})

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			return "", driver.ErrLockNotFound
		}

		return "", err
	}

	return driver.Info(string(b.Bytes())), nil
}

func (d *d) HasLock(ctx context.Context, table string, lid id.LockID) (bool, error) {
	_, err := d.GetLock(ctx, table, lid)
	if err != nil {
		if errors.Is(err, driver.ErrLockNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *d) SaveLock(ctx context.Context, table string, lid id.LockID, info driver.Info) error {
	keyPath := fmt.Sprintf("%s/%s", table, lid)

	b := bytes.NewBuffer([]byte(info))
	uploader := manager.NewUploader(d.s3)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(keyPath),
		Body:   b,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *d) CreateTenant(ctx context.Context, table string, token string) error {
	uploader := manager.NewUploader(d.s3)

	b := bytes.NewBuffer([]byte(token))
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", table, driver.TokenFile)),
		Body:   b,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *d) GetTenant(ctx context.Context, table string) (*modelv1.Tenant, error) {
	b := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(d.s3)
	_, err := downloader.Download(ctx, b, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", table, driver.TokenFile)),
	})

	if err != nil {
		var bne *types.NoSuchKey
		if errors.As(err, &bne) {
			return nil, driver.ErrTenantNotFound
		}
		return nil, err
	}

	return &modelv1.Tenant{
		Table: table,
	}, nil
}
