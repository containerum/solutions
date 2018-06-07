package clients

import (
	"context"

	"fmt"

	kube_types "github.com/containerum/kube-client/pkg/model"

	"github.com/containerum/cherry"
	utils "github.com/containerum/utils/httputil"

	"time"

	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// KubeAPIClient is an interface to Kube-API.
type KubeAPIClient interface {
	GetUserDeployments(ctx context.Context, namespace, solutionName string) (*kube_types.DeploymentsList, error)
	GetUserServices(ctx context.Context, namespace, solutionName string) (*kube_types.ServicesList, error)
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

func (c *httpKubeAPIClient) GetUserDeployments(ctx context.Context, namespace, solutionName string) (*kube_types.DeploymentsList, error) {
	c.log.Info("Getting solution deployments")
	headersMap := utils.RequestHeadersMap(ctx)

	var dlist kube_types.DeploymentsList
	resp, err := c.rest.R().SetContext(ctx).
		SetResult(&dlist).
		SetHeaders(headersMap).
		Get(fmt.Sprintf("/namespaces/%s/solutions/%s/deployments", namespace, solutionName))
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error().(*cherry.Err)
	}

	return &dlist, nil
}

func (c *httpKubeAPIClient) GetUserServices(ctx context.Context, namespace, solutionName string) (*kube_types.ServicesList, error) {
	c.log.Info("Getting solution services")
	headersMap := utils.RequestHeadersMap(ctx)

	var slist kube_types.ServicesList

	resp, err := c.rest.R().SetContext(ctx).
		SetResult(&slist).
		SetHeaders(headersMap).
		Get(fmt.Sprintf("/namespaces/%s/solutions/%s/services", namespace, solutionName))
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error().(*cherry.Err)
	}

	return &slist, nil
}
