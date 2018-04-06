package clients

import (
	"context"

	"fmt"

	stypes "git.containerum.net/ch/json-types/solutions"

	"time"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/utils"
	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// KubeAPIClient is an interface to resource-service.
type KubeAPIClient interface {
	GetUserDeployments(ctx context.Context, namespace string, depl []string) (*stypes.DeploymentsList, error)
	GetUserServices(ctx context.Context, namespace string, svc []string) (*stypes.ServicesList, error)
	CreateDeployment(ctx context.Context, namespace string, deployment string) error
	CreateService(ctx context.Context, namespace string, service string) error
	DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error
	DeleteService(ctx context.Context, namespace string, serviceName string) error
}

type httpKubeAPIClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

// NewHTTPResourceServiceClient returns client for resource-service working via restful api
func NewHTTPKubeAPIClient(serverURL string) KubeAPIClient {
	log := logrus.WithField("component", "kube_api_client")
	client := resty.New().
		SetHostURL(serverURL).
		SetLogger(log.WriterLevel(logrus.DebugLevel)).
		SetDebug(true).
		SetTimeout(3 * time.Second).
		SetError(cherry.Err{})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	return &httpKubeAPIClient{
		rest: client,
		log:  log,
	}
}

func (c *httpKubeAPIClient) CreateDeployment(ctx context.Context, namespace string, deployment string) error {
	c.log.Info("Creating deployment")
	headersMap := utils.RequestHeadersMap(ctx)
	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(headersMap).
		SetBody(deployment).
		Post(fmt.Sprintf("/namespaces/%s/deployments/file", namespace))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpKubeAPIClient) CreateService(ctx context.Context, namespace string, service string) error {
	c.log.Info("Creating service")
	headersMap := utils.RequestHeadersMap(ctx)
	resp, err := c.rest.R().SetContext(ctx).
		SetHeaders(headersMap).
		SetBody(service).
		Post(fmt.Sprintf("/namespaces/%s/services/file", namespace))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error().(*cherry.Err)
	}
	return nil
}

func (c *httpKubeAPIClient) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error {
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

func (c *httpKubeAPIClient) DeleteService(ctx context.Context, namespace string, serviceName string) error {
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

func (c *httpKubeAPIClient) GetUserDeployments(ctx context.Context, namespace string, depl []string) (*stypes.DeploymentsList, error) {
	c.log.Info("Getting user deployments")
	headersMap := utils.RequestHeadersMap(ctx)

	var dlist stypes.DeploymentsList

	dlist.Deployments = make([]*interface{}, 0)

	for _, d := range depl {
		var depl interface{}

		resp, err := c.rest.R().SetContext(ctx).
			SetResult(&depl).
			SetHeaders(headersMap).
			Get(fmt.Sprintf("/namespaces/%s/deployments/%s", namespace, d))
		if err != nil {
			return nil, err
		}
		if resp.Error() != nil {
			return nil, resp.Error().(*cherry.Err)
		}

		dlist.Deployments = append(dlist.Deployments, &depl)
	}

	return &dlist, nil
}

func (c *httpKubeAPIClient) GetUserServices(ctx context.Context, namespace string, svc []string) (*stypes.ServicesList, error) {
	c.log.Info("Getting user services")
	headersMap := utils.RequestHeadersMap(ctx)

	var dlist stypes.ServicesList

	dlist.Services = make([]*interface{}, 0)

	for _, r := range svc {
		var service interface{}

		resp, err := c.rest.R().SetContext(ctx).
			SetResult(&service).
			SetHeaders(headersMap).
			Get(fmt.Sprintf("/namespaces/%s/services/%s", namespace, r))
		if err != nil {
			return nil, err
		}
		if resp.Error() != nil {
			return nil, resp.Error().(*cherry.Err)
		}

		dlist.Services = append(dlist.Services, &service)
	}

	return &dlist, nil
}
