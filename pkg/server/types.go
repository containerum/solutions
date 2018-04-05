package server

import (
	"context"

	"io"

	stypes "git.containerum.net/ch/json-types/solutions"

	"git.containerum.net/ch/solutions/pkg/clients"
	"git.containerum.net/ch/solutions/pkg/models"
)

// SolutionsService is an interface for server "business logic"
type SolutionsService interface {
	UpdateAvailableSolutionsList(ctx context.Context) error
	GetAvailableSolutionsList(ctx context.Context) (*stypes.AvailableSolutionsList, error)
	GetAvailableSolutionEnvList(ctx context.Context, name string, branch string) (*stypes.SolutionEnv, error)
	GetAvailableSolutionResourcesList(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error)
	GetUserSolutionsList(ctx context.Context) (*stypes.UserSolutionsList, error)
	RunSolution(ctx context.Context, solution stypes.UserSolution) error
	DeleteSolution(ctx context.Context, solution string) error
	GetUserSolutionDeployments(ctx context.Context, solutionName string) (*stypes.DeploymentsList, error)
	GetUserSolutionServices(ctx context.Context, solutionName string) (*stypes.ServicesList, error)
	io.Closer
}

// Services is a collection of resources needed for server functionality.
type Services struct {
	DB             models.DB
	DownloadClient clients.DownloadClient
	KubeAPI        clients.KubeAPIClient
}

type Solution struct {
	Env map[string]string `json:"env"`
	Run []ConfigFile      `json:"run,omitempty"`
}

type ConfigFile struct {
	Name string `json:"config_file"`
	Type string `json:"type"`
}

type ResName struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}
