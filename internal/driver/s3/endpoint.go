package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

type endpointResolver struct {
	URL string
}

func (r *endpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL: r.URL,
	}, nil
}
