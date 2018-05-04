package impl

import (
	"context"
	"net/url"

	"strings"

	"fmt"

	stypes "git.containerum.net/ch/solutions/pkg/models"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/json-iterator/go"
)

func (s *serverImpl) UpdateAvailableSolutionsList(ctx context.Context) error {
	solutions, err := s.svc.DownloadClient.DownloadSolutionsCSV(ctx)
	if err != nil {
		return err
	}

	err = s.svc.DB.SaveAvailableSolutionsList(ctx, stypes.AvailableSolutionsList{solutions})
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) AddAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error {
	err := s.svc.DB.CreateAvailableSolution(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) UpdateAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error {
	err := s.svc.DB.UpdateAvailableSolution(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) DeleteAvailableSolution(ctx context.Context, solution string) error {
	err := s.svc.DB.DeleteAvailableSolution(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) ActivateAvailableSolution(ctx context.Context, solution string) error {
	err := s.svc.DB.ActivateAvailableSolution(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) DeactivateAvailableSolution(ctx context.Context, solution string) error {
	err := s.svc.DB.DeactivateAvailableSolution(ctx, solution)
	if err := s.handleDBError(err); err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) GetAvailableSolutionsList(ctx context.Context, isAdmin bool) (*stypes.AvailableSolutionsList, error) {
	resp, err := s.svc.DB.GetAvailableSolutionsList(ctx, isAdmin)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *serverImpl) GetAvailableSolutionEnvList(ctx context.Context, name string, branch string) (*stypes.SolutionEnv, error) {
	solution, err := s.svc.DB.GetAvailableSolution(ctx, name)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}
	if solution == nil {
		return nil, sErrors.ErrSolutionNotExist()
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

func (s *serverImpl) GetAvailableSolutionResourcesList(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error) {
	solution, err := s.svc.DB.GetAvailableSolution(ctx, name)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}
	if solution == nil {
		return nil, sErrors.ErrSolutionNotExist()
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
