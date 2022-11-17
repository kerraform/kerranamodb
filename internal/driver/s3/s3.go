package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kerraform/kerranamodb/internal/driver"
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
