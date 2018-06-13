package clients

import (
	"context"

	"fmt"

	"time"

	"github.com/containerum/cherry"
	kube_types "github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/utils/httputil"
	utils "github.com/containerum/utils/httputil"
	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// ResourceClient is an interface to resource-service.
type ResourceClient interface {
	CreateDeployment(ctx context.Context, namespace string, deployment kube_types.Deployment) error
	CreateService(ctx context.Context, namespace string, service kube_types.Service) error
	DeleteDeployments(ctx context.Context, namespace, solutionName string) error
	DeleteServices(ctx context.Context, namespace, solutionName string) error
}

type httpResourceClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

// NewHTTPResourceClient returns client for resource-service working via restful api
func NewHTTPResourceClient(serverURL string, debug bool) ResourceClient {
	log := logrus.WithField("component", "resource_client")
	client := resty.New().
		SetHostURL(serverURL).
		SetLogger(log.WriterLevel(logrus.DebugLevel)).
		SetDebug(debug).
		SetTimeout(3*time.Second).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetError(cherry.Err{})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	return &httpResourceClient{
		rest: client,
		log:  log,
	}
}

func (c *httpResourceClient) CreateDeployment(ctx context.Context, namespace string, deployment kube_types.Deployment) error {
	c.log.Info("Creating deployment")

	resp, err := c.rest.R().SetContext(ctx).
		SetBody(deployment).
		SetHeaders(httputil.RequestXHeadersMap(ctx)).
		Post(fmt.Sprintf("/namespaces/%s/deployments", namespace))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpResourceClient) CreateService(ctx context.Context, namespace string, service kube_types.Service) error {
	c.log.Info("Creating service")
	resp, err := c.rest.R().SetContext(ctx).
		SetBody(service).
		SetHeaders(httputil.RequestXHeadersMap(ctx)).
		Post(fmt.Sprintf("/namespaces/%s/services", namespace))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpResourceClient) DeleteDeployments(ctx context.Context, namespace, solutionName string) error {
	c.log.Info("Deleting deployments")
	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(utils.RequestHeadersMap(ctx)).
		Delete(fmt.Sprintf("/namespaces/%s/solutions/%s/deployments", namespace, solutionName))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpResourceClient) DeleteServices(ctx context.Context, namespace, solutionName string) error {
	c.log.Info("Deleting services")
	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(utils.RequestHeadersMap(ctx)).
		Delete(fmt.Sprintf("/namespaces/%s/solutions/%s/deployments", namespace, solutionName))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}
