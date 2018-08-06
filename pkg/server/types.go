package server

import (
	"context"

	"io"

	"git.containerum.net/ch/solutions/pkg/db"
	kube_types "github.com/containerum/kube-client/pkg/model"

	"git.containerum.net/ch/solutions/pkg/clients"
)

// SolutionsService is an interface for server "business logic"
type SolutionsService interface {
	AddTemplate(ctx context.Context, solution kube_types.SolutionTemplate) error
	UpdateTemplate(ctx context.Context, solution kube_types.SolutionTemplate) error
	GetTemplatesList(ctx context.Context, isAdmin bool) (*kube_types.SolutionsTemplatesList, error)
	GetTemplatesEnvList(ctx context.Context, name, branch string) (*kube_types.SolutionEnv, error)
	GetTemplatesResourcesList(ctx context.Context, name, branch string) (*kube_types.SolutionResources, error)
	ActivateTemplate(ctx context.Context, solution string) error
	DeactivateTemplate(ctx context.Context, solution string) error
	ValidateTemplate(ctx context.Context, solution kube_types.SolutionTemplate) error

	GetSolutionsList(ctx context.Context, isAdmin bool) (*kube_types.SolutionsList, error)
	GetNamespaceSolutionsList(ctx context.Context, namespace string, isAdmin bool) (*kube_types.SolutionsList, error)
	GetSolution(ctx context.Context, namespace, solutionName string, isAdmin bool) (*kube_types.Solution, error)
	GetSolutionDeployments(ctx context.Context, namespace, solutionName string) (*kube_types.DeploymentsList, error)
	GetSolutionServices(ctx context.Context, namespace, solutionName string) (*kube_types.ServicesList, error)
	RunSolution(ctx context.Context, solutionReq kube_types.Solution) (*kube_types.RunSolutionResponse, error)
	DeleteSolution(ctx context.Context, namespace, solution string) error
	DeleteSolutions(ctx context.Context) error
	DeleteNamespaceSolutions(ctx context.Context, namespace string) error
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
