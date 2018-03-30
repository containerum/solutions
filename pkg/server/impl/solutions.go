package impl

import (
	"context"

	"net/url"

	"strings"

	stypes "git.containerum.net/ch/json-types/solutions"
)

func (s *serverImpl) UpdateAvailableSolutionsList(ctx context.Context) error {
	//TODO: Don't download CSV every time
	solutions, err := s.svc.DownloadClient.DownloadCSV(ctx)
	if err != nil {
		return err
	}

	err = s.svc.DB.SaveAvailableSolutionsList(ctx, stypes.AvailableSolutionsList{solutions})
	if err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) GetAvailableSolutionsList(ctx context.Context) (*stypes.AvailableSolutionsList, error) {
	resp, err := s.svc.DB.GetAvailableSolutionsList(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *serverImpl) GetAvailableSolutionEnv(ctx context.Context, name string, branch string) (*stypes.SolutionEnv, error) {
	solution, err := s.svc.DB.GetAvailableSolution(ctx, name)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	sBranch := "master"
	if branch != "" {
		sBranch = strings.TrimSpace(branch)
	}
	sName := strings.TrimSpace(url.Path[1:])
	sFile := ".containerum.json"

	containerumJSONURL := "https://raw.githubusercontent.com/" + sName + "/" + sBranch + "/" + sFile

	solutionJSON, err := s.svc.DownloadClient.DownloadSolutionJSON(ctx, containerumJSONURL)
	if err != nil {
		return nil, err
	}

	resp := stypes.SolutionEnv{}

	for e, _ := range solutionJSON.Env {
		resp.Env = append(resp.Env, e)
	}

	return &resp, nil
}

func (s *serverImpl) GetAvailableSolutionResources(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error) {
	solution, err := s.svc.DB.GetAvailableSolution(ctx, name)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(solution.URL)
	if err != nil {
		return nil, err
	}

	sBranch := "master"
	if branch != "" {
		sBranch = strings.TrimSpace(branch)
	}
	sName := strings.TrimSpace(url.Path[1:])
	sFile := ".containerum.json"

	containerumJSONURL := "https://raw.githubusercontent.com/" + sName + "/" + sBranch + "/" + sFile

	solutionJSON, err := s.svc.DownloadClient.DownloadSolutionJSON(ctx, containerumJSONURL)
	if err != nil {
		return nil, err
	}

	resp := stypes.SolutionResources{map[string]int{}}

	for _, r := range solutionJSON.Run {
		resp.Resources[r.Type]++
	}

	return &resp, nil
}
