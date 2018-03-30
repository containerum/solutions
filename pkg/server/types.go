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
	GetAvailableSolutionEnv(ctx context.Context, name string, branch string) (*stypes.SolutionEnv, error)
	GetAvailableSolutionResources(ctx context.Context, name string, branch string) (*stypes.SolutionResources, error)
	io.Closer
}

// Services is a collection of resources needed for server functionality.
type Services struct {
	DB             models.DB
	DownloadClient clients.DownloadClient
}
