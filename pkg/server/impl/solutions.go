package impl

import (
	"context"
	"net/url"

	"strings"

	"fmt"

	stypes "git.containerum.net/ch/json-types/solutions"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/solutions"
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

func (s *serverImpl) GetAvailableSolutionsList(ctx context.Context) (*stypes.AvailableSolutionsList, error) {
	resp, err := s.svc.DB.GetAvailableSolutionsList(ctx)
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
		return nil, cherry.ErrSolutionNotExist()
	}

	solurl, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	sBranch := "master"
	if branch != "" {
		sBranch = strings.TrimSpace(branch)
	}
	sName := strings.TrimSpace(solurl.Path[1:])
	sFile := ".containerum.json"

	containerumJSONURL := "https://raw.githubusercontent.com/" + sName + "/" + sBranch + "/" + sFile

	solutionJSON, err := s.svc.DownloadClient.DownloadFile(ctx, containerumJSONURL)
	if err != nil {
		return nil, err
	}

	var solutionStr server.Solution

	err = jsoniter.Unmarshal(solutionJSON, &solutionStr)
	if err != nil {
		return nil, err
	}

	resp := stypes.SolutionEnv{}

	for e := range solutionStr.Env {
		resp.Env = append(resp.Env, e)
	}

	return &resp, nil
}

func (s *serverImpl) GetAvailableSolutionResourcesList(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error) {
	solution, err := s.svc.DB.GetAvailableSolution(ctx, name)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}
	if solution == nil {
		return nil, cherry.ErrSolutionNotExist()
	}

	urlcsv, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	if branch != "" {
		branch = strings.TrimSpace(branch)
	} else {
		branch = "master"
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
