package server

import (
	"context"

	"io"

	kube_types "git.containerum.net/ch/kube-api/pkg/model"
	"git.containerum.net/ch/solutions/pkg/db"
	"git.containerum.net/ch/solutions/pkg/models"
	stypes "git.containerum.net/ch/solutions/pkg/models"

	"git.containerum.net/ch/solutions/pkg/clients"
)

// SolutionsService is an interface for server "business logic"
type SolutionsService interface {
	UpdateAvailableSolutionsList(ctx context.Context) error
	AddAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error
	UpdateAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error
	DeleteAvailableSolution(ctx context.Context, solution string) error
	GetAvailableSolutionsList(ctx context.Context, isAdmin bool) (*stypes.AvailableSolutionsList, error)
	GetAvailableSolutionEnvList(ctx context.Context, name string, branch string) (*stypes.SolutionEnv, error)
	GetAvailableSolutionResourcesList(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error)
	GetUserSolutionsList(ctx context.Context) (*stypes.UserSolutionsList, error)
	ActivateAvailableSolution(ctx context.Context, solution string) error
	DeactivateAvailableSolution(ctx context.Context, solution string) error

	DownloadSolutionConfig(ctx context.Context, solutionReq stypes.UserSolution) (solutionFile []byte, solutionName *string, err error)
	ParseSolutionConfig(ctx context.Context, solutionBody []byte, solutionReq stypes.UserSolution) (solutionConfig *Solution, solutionUUID *string, err error)
	CreateSolutionResources(ctx context.Context, solutionConfig Solution, solutionReq stypes.UserSolution, solutionName string, solutionUUID string) (*models.RunSolutionResponce, error)
	DeleteSolution(ctx context.Context, solution string) error
	GetUserSolutionDeployments(ctx context.Context, solutionName string) (*kube_types.DeploymentsList, error)
	GetUserSolutionServices(ctx context.Context, solutionName string) (*kube_types.ServicesList, error)
	io.Closer
}

// Services is a collection of resources needed for server functionality.
type Services struct {
	DB              db.DB
	DownloadClient  clients.DownloadClient
	ResourceClient  clients.ResourceClient
	KubeAPIClient   clients.KubeAPIClient
	ConverterClient clients.ConverterClient
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
