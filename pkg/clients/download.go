package clients

import (
	"context"

	"encoding/csv"
	"io"
	"strings"

	"errors"

	stypes "git.containerum.net/ch/json-types/solutions"
	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var csvURL string

type Solution struct {
	Env map[string]string `json:"env"`
	Run []ConfigFile      `json:"run,omitempty"`
}

type ConfigFile struct {
	Name string `json:"config_file"`
	Type string `json:"type"`
}

// ResourceServiceClient is an interface to resource-service.
type DownloadClient interface {
	DownloadCSV(ctx context.Context) ([]stypes.AvailableSolution, error)
	DownloadSolutionJSON(ctx context.Context, url string) (*Solution, error)
}

type httpDownloadClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

// NewHTTPResourceServiceClient returns client for resource-service working via restful api
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

func (c *httpDownloadClient) DownloadCSV(ctx context.Context) ([]stypes.AvailableSolution, error) {
	c.log.Info("Downloading CSV")

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
				Limits: &stypes.Limits{
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

func (c *httpDownloadClient) DownloadSolutionJSON(ctx context.Context, url string) (*Solution, error) {
	c.log.Info("Downloading solution config")

	solution := Solution{}

	resp, err := c.rest.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() < 399 {
		err = jsoniter.Unmarshal(resp.Body(), &solution)
		if err != nil {
			return nil, err
		}

		return &solution, nil
	} else {
		return nil, errors.New("Unable to download file")
	}
}
