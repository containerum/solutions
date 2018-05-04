package clients

import (
	"context"

	"encoding/csv"
	"io"
	"strings"

	"errors"

	stypes "github.com/containerum/kube-client/pkg/model"
	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var csvURL string

// DownloadClient is an interface to resource-service.
type DownloadClient interface {
	DownloadSolutionsCSV(ctx context.Context) ([]stypes.AvailableSolution, error)
	DownloadFile(ctx context.Context, url string) ([]byte, error)
}

type httpDownloadClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

// NewHTTPDownloadClient returns client for resource-service working via restful api
func NewHTTPDownloadClient(serverURL string) DownloadClient {
	log := logrus.WithField("component", "download_client")
	client := resty.New().
		SetLogger(log.WriterLevel(logrus.DebugLevel))
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal

	csvURL = serverURL

	return &httpDownloadClient{
		rest: client,
		log:  log,
	}
}

func (c *httpDownloadClient) DownloadSolutionsCSV(ctx context.Context) ([]stypes.AvailableSolution, error) {
	c.log.Infoln("Downloading CSV")

	resp, err := c.rest.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		Get(csvURL)

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(resp.RawBody())
	defer resp.RawBody().Close()
	var solutions []stypes.AvailableSolution
	header := true
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !header {
			solutions = append(solutions, stypes.AvailableSolution{
				Name: line[0],
				Limits: &stypes.SolutionLimits{
					CPU: line[1],
					RAM: line[2],
				},
				Images: strings.Split(line[3], ";"),
				URL:    line[4],
			})
		}
		header = false
	}
	return solutions, nil
}

func (c *httpDownloadClient) DownloadFile(ctx context.Context, url string) ([]byte, error) {
	c.log.WithField("URL", url).Infoln("Downloading file")

	resp, err := c.rest.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() > 399 {
		return nil, errors.New("Unable to download file")
	}
	return resp.Body(), nil
}
