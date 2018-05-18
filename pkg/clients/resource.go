package clients

import (
	"context"

	"fmt"

	"time"

	"github.com/containerum/cherry"
	utils "github.com/containerum/utils/httputil"
	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// ResourceClient is an interface to resource-service.
type ResourceClient interface {
	CreateDeployment(ctx context.Context, namespace string, deployment string) error
	CreateService(ctx context.Context, namespace string, service string) error
	DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error
	DeleteService(ctx context.Context, namespace string, serviceName string) error
}

type httpResourceClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

// NewHTTPResourceClient returns client for resource-service working via restful api
func NewHTTPResourceClient(serverURL string) ResourceClient {
	log := logrus.WithField("component", "resource_client")
	client := resty.New().
		SetHostURL(serverURL).
		SetLogger(log.WriterLevel(logrus.DebugLevel)).
		SetDebug(true).
		SetTimeout(3 * time.Second).
		SetError(cherry.Err{})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	return &httpResourceClient{
		rest: client,
		log:  log,
	}
}

func (c *httpResourceClient) CreateDeployment(ctx context.Context, namespace string, deployment string) error {
	c.log.Info("Creating deployment")
	headersMap := utils.RequestHeadersMap(ctx)

	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(headersMap).
		SetBody(deployment).
		Post(fmt.Sprintf("/namespaces/%s/deployments", namespace))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpResourceClient) CreateService(ctx context.Context, namespace string, service string) error {
	c.log.Info("Creating service")
	headersMap := utils.RequestHeadersMap(ctx)
	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(headersMap).
		SetBody(service).
		Post(fmt.Sprintf("/namespaces/%s/services", namespace))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpResourceClient) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error {
	c.log.Info("Deleting deployment")
	headersMap := utils.RequestHeadersMap(ctx)
	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(headersMap).
		Delete(fmt.Sprintf("/namespaces/%s/deployments/%s", namespace, deploymentName))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpResourceClient) DeleteService(ctx context.Context, namespace string, serviceName string) error {
	c.log.Info("Deleting service")
	headersMap := utils.RequestHeadersMap(ctx)
	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(headersMap).
		Delete(fmt.Sprintf("/namespaces/%s/services/%s", namespace, serviceName))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}
