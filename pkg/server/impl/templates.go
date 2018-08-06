package impl

import (
	"context"
	"fmt"
	"net/url"

	"git.containerum.net/ch/solutions/pkg/db"
	"git.containerum.net/ch/solutions/pkg/server"
	kube_types "github.com/containerum/kube-client/pkg/model"
	"github.com/json-iterator/go"
)

func (s *serverImpl) GetTemplatesList(ctx context.Context, isAdmin bool) (*kube_types.AvailableSolutionsList, error) {
	resp, err := s.svc.DB.GetTemplatesList(ctx, isAdmin)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	if !isAdmin {
		for i := range resp.Solutions {
			resp.Solutions[i].ID = ""
			resp.Solutions[i].Active = false
		}
	}

	return resp, nil
}

func (s *serverImpl) GetTemplatesEnvList(ctx context.Context, name string, branch string) (*kube_types.SolutionEnv, error) {
	solution, err := s.svc.DB.GetTemplate(ctx, name)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	solurl, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	solutionJSON, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/.containerum.json", solurl.Path[1:], branch))
	if err != nil {
		return nil, err
	}

	var solutionStr server.Solution

	if err = jsoniter.Unmarshal(solutionJSON, &solutionStr); err != nil {
		return nil, err
	}

	resp := kube_types.SolutionEnv{Env: solutionStr.Env}

	return &resp, nil
}

func (s *serverImpl) GetTemplatesResourcesList(ctx context.Context, name string, branch string) (*kube_types.SolutionResources, error) {
	solution, err := s.svc.DB.GetTemplate(ctx, name)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	solurl, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	solutionJSON, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/.containerum.json", solurl.Path[1:], branch))
	if err != nil {
		return nil, err
	}

	var solutionStr server.Solution
	err = jsoniter.Unmarshal(solutionJSON, &solutionStr)
	if err != nil {
		return nil, err
	}

	resp := kube_types.SolutionResources{Resources: map[string]int{}}

	for _, r := range solutionStr.Run {
		resp.Resources[r.Type]++
	}

	return &resp, nil
}

func (s *serverImpl) AddTemplate(ctx context.Context, solution kube_types.AvailableSolution) error {
	err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.CreateTemplate(ctx, solution)
	})
	return s.handleDBError(err)
}

func (s *serverImpl) UpdateTemplate(ctx context.Context, solution kube_types.AvailableSolution) error {
	err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.UpdateTemplate(ctx, solution)
	})
	return s.handleDBError(err)
}

func (s *serverImpl) ActivateTemplate(ctx context.Context, solution string) error {
	err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.ActivateTemplate(ctx, solution)
	})
	return s.handleDBError(err)
}

func (s *serverImpl) DeactivateTemplate(ctx context.Context, solution string) error {
	err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.DeactivateTemplate(ctx, solution)
	})
	return s.handleDBError(err)
}

func (s *serverImpl) ValidateTemplate(ctx context.Context, solution kube_types.AvailableSolution) error {
	solurl, err := url.Parse(solution.URL)
	if err != nil {
		return err
	}

	solutionJSON, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/master/.containerum.json", solurl.Path[1:]))
	if err != nil {
		return err
	}

	var solutionStr server.Solution
	err = jsoniter.Unmarshal(solutionJSON, &solutionStr)
	if err != nil {
		return err
	}

	for _, r := range solutionStr.Run {
		if _, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/master/%s", solurl.Path[1:], r.Name)); err != nil {
			return err
		}
	}
	return nil
}
