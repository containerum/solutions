package impl

import (
	"context"
	"net/url"

	"strings"

	"fmt"

	"git.containerum.net/ch/solutions/pkg/server"
	stypes "github.com/containerum/kube-client/pkg/model"
	"github.com/json-iterator/go"
)

func (s *serverImpl) AddTemplate(ctx context.Context, solution stypes.AvailableSolution) error {
	err := s.svc.DB.CreateTemplate(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) UpdateTemplate(ctx context.Context, solution stypes.AvailableSolution) error {
	err := s.svc.DB.UpdateTemplate(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) DeleteTemplate(ctx context.Context, solution string) error {
	err := s.svc.DB.DeleteTemplate(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) ActivateTemplate(ctx context.Context, solution string) error {
	err := s.svc.DB.ActivateTemplate(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) DeactivateTemplate(ctx context.Context, solution string) error {
	err := s.svc.DB.DeactivateTemplate(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) GetTemplatesList(ctx context.Context, isAdmin bool) (*stypes.AvailableSolutionsList, error) {
	resp, err := s.svc.DB.GetTemplatesList(ctx, isAdmin)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *serverImpl) GetTemplatesEnvList(ctx context.Context, name string, branch string) (*stypes.SolutionEnv, error) {
	solution, err := s.svc.DB.GetTemplate(ctx, name)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	solurl, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	sName := strings.TrimSpace(solurl.Path[1:])
	sFile := ".containerum.json"

	containerumJSONURL := "https://raw.githubusercontent.com/" + sName + "/" + branch + "/" + sFile

	solutionJSON, err := s.svc.DownloadClient.DownloadFile(ctx, containerumJSONURL)
	if err != nil {
		return nil, err
	}

	var solutionStr server.Solution

	err = jsoniter.Unmarshal(solutionJSON, &solutionStr)
	if err != nil {
		return nil, err
	}

	resp := stypes.SolutionEnv{Env: solutionStr.Env}

	return &resp, nil
}

func (s *serverImpl) GetTemplatesResourcesList(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error) {
	solution, err := s.svc.DB.GetTemplate(ctx, name)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	urlcsv, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	sName := strings.TrimSpace(urlcsv.Path[1:])

	solutionJSON, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/.containerum.json", sName, branch))
	if err != nil {
		return nil, err
	}

	var solutionStr server.Solution

	err = jsoniter.Unmarshal(solutionJSON, &solutionStr)
	if err != nil {
		return nil, err
	}

	resp := stypes.SolutionResources{map[string]int{}}

	for _, r := range solutionStr.Run {
		resp.Resources[r.Type]++
	}

	return &resp, nil
}
