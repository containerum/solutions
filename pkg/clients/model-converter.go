package clients

import (
	"context"
	"time"

	"fmt"

	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// ConverterClient is an interface to resource-service.
type ConverterClient interface {
	ConvertDeployment(ctx context.Context, deployment string) (convertedDeploy *string, err error)
	ConvertService(ctx context.Context, service string) (convertedDeploy *string, err error)
}

type httpConverterClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

// NewHTTPConverterClient returns client for resource-service working via restful api
func NewHTTPConverterClient(serverURL string) ConverterClient {
	log := logrus.WithField("component", "model_converter_client")
	client := resty.New().
		SetHostURL(serverURL).
		SetLogger(log.WriterLevel(logrus.DebugLevel)).
		SetDebug(true).
		SetTimeout(3 * time.Second).
		SetError(solutionsErrorsErr{})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	return &httpConverterClient{
		rest: client,
		log:  log,
	}
}

func (c *httpConverterClient) ConvertDeployment(ctx context.Context, deployment string) (convertedDeploy *string, err error) {
	c.log.Info("Converting deployment")

	resp, err := c.rest.R().SetContext(ctx).
		SetBody(deployment).
		Post("/convert/fromkube/deployment")
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error().(*solutionsErrorsErr)
	}
	result := fmt.Sprintf("%v", string(resp.Body()))
	return &result, nil
}

func (c *httpConverterClient) ConvertService(ctx context.Context, service string) (convertedDeploy *string, err error) {
	c.log.Info("Converting service")

	resp, err := c.rest.R().SetContext(ctx).
		SetBody(service).
		Post("/convert/fromkube/service")
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error().(*solutionsErrorsErr)
	}
	result := fmt.Sprintf("%v", string(resp.Body()))
	return &result, nil
}
