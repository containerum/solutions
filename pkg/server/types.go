package server

import (
	"context"

	"io"

	kube_types "git.containerum.net/ch/kube-api/pkg/model"
	"git.containerum.net/ch/solutions/pkg/db"
	stypes "github.com/containerum/kube-client/pkg/model"

	"git.containerum.net/ch/solutions/pkg/clients"
)

// SolutionsService is an interface for server "business logic"
type SolutionsService interface {
	AddTemplate(ctx context.Context, solution stypes.AvailableSolution) error
	UpdateTemplate(ctx context.Context, solution stypes.AvailableSolution) error
	DeleteTemplate(ctx context.Context, solution string) error
	GetTemplatesList(ctx context.Context, isAdmin bool) (*stypes.AvailableSolutionsList, error)
	GetTemplatesEnvList(ctx context.Context, name string, branch string) (*stypes.SolutionEnv, error)
	GetTemplatesResourcesList(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error)
	ActivateTemplate(ctx context.Context, solution string) error
	DeactivateTemplate(ctx context.Context, solution string) error

	GetSolutionsList(ctx context.Context, isAdmin bool) (*stypes.UserSolutionsList, error)
	RunSolution(ctx context.Context, solutionReq stypes.UserSolution) (*stypes.RunSolutionResponce, error)
	DeleteSolution(ctx context.Context, solution string) error
	GetSolutionDeployments(ctx context.Context, solutionName string) (*kube_types.DeploymentsList, error)
	GetSolutionServices(ctx context.Context, solutionName string) (*kube_types.ServicesList, error)
	io.Closer
}

// Services is a collection of resources needed for server functionality.
type Services struct {
	DB             db.DB
	DownloadClient clients.DownloadClient
	ResourceClient clients.ResourceClient
	KubeAPIClient  clients.KubeAPIClient
}

type Solution struct {
	Env map[string]string `json:"env"`
	Run []ConfigFile      `json:"run,omitempty"`
}

type ConfigFile struct {
	Name string `json:"config_file"`
	Type string `json:"type"`
}
