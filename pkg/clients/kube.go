package clients

import (
	"context"

	"fmt"

	kube_types "git.containerum.net/ch/kube-api/pkg/model"

	"github.com/containerum/cherry"
	utils "github.com/containerum/utils/httputil"

	"time"

	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// KubeAPIClient is an interface to Kube-API.
type KubeAPIClient interface {
	GetUserDeployments(ctx context.Context, namespace string, depl []string) (*kube_types.DeploymentsList, error)
	GetUserServices(ctx context.Context, namespace string, svc []string) (*kube_types.ServicesList, error)
}

type httpKubeAPIClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

// NewHTTPKubeAPIClient returns client for resource-service working via restful api
func NewHTTPKubeAPIClient(serverURL string, debug bool) KubeAPIClient {
	log := logrus.WithField("component", "kube_api_client")
	client := resty.New().
		SetHostURL(serverURL).
		SetLogger(log.WriterLevel(logrus.DebugLevel)).
		SetDebug(debug).
		SetTimeout(3 * time.Second).
		SetError(cherry.Err{})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	return &httpKubeAPIClient{
		rest: client,
		log:  log,
	}
}

func (c *httpKubeAPIClient) GetUserDeployments(ctx context.Context, namespace string, depl []string) (*kube_types.DeploymentsList, error) {
	c.log.Info("Getting user deployments")
	headersMap := utils.RequestHeadersMap(ctx)

	var dlist kube_types.DeploymentsList

	dlist.Deployments = make([]kube_types.DeploymentWithOwner, 0)

	for _, d := range depl {
		var depl kube_types.DeploymentWithOwner

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

		dlist.Deployments = append(dlist.Deployments, depl)
	}

	return &dlist, nil
}

func (c *httpKubeAPIClient) GetUserServices(ctx context.Context, namespace string, svc []string) (*kube_types.ServicesList, error) {
	c.log.Info("Getting user services")
	headersMap := utils.RequestHeadersMap(ctx)

	var dlist kube_types.ServicesList

	dlist.Services = make([]kube_types.ServiceWithOwner, 0)

	for _, r := range svc {
		var service kube_types.ServiceWithOwner

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

		dlist.Services = append(dlist.Services, service)
	}

	return &dlist, nil
}
