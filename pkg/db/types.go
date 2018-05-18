package db

import (
	"io"

	"context"

	"errors"

	stypes "github.com/containerum/kube-client/pkg/model"
)

// Errors which may occur in transactional operations
var (
	ErrTransactionBegin    = errors.New("transaction begin error")
	ErrTransactionRollback = errors.New("transaction rollback error")
	ErrTransactionCommit   = errors.New("transaction commit error")
)

// DB is an interface for persistent data storage (also sometimes called DAO).
type DB interface {
	CreateTemplate(ctx context.Context, solution stypes.AvailableSolution) error
	UpdateTemplate(ctx context.Context, solution stypes.AvailableSolution) error
	DeleteTemplate(ctx context.Context, solution string) error
	GetTemplatesList(ctx context.Context, isAdmin bool) (*stypes.AvailableSolutionsList, error)
	GetTemplate(ctx context.Context, name string) (*stypes.AvailableSolution, error)
	ActivateTemplate(ctx context.Context, solution string) error
	DeactivateTemplate(ctx context.Context, solution string) error

	GetSolutionsList(ctx context.Context, userID string) (*stypes.UserSolutionsList, error)
	AddSolution(ctx context.Context, solution stypes.UserSolution, userID, templateID, uuid, env string) error
	AddDeployment(ctx context.Context, name string, solutionID string) error
	AddService(ctx context.Context, name string, solutionID string) error
	DeleteSolution(ctx context.Context, name string, userID string) error
	CompletelyDeleteSolution(ctx context.Context, name string, userID string) error
	GetSolutionsDeployments(ctx context.Context, solutionName string, userID string) (deployments []string, ns *string, err error)
	GetSolutionsServices(ctx context.Context, solutionName string, userID string) (services []string, ns *string, err error)

	// Perform operations inside transaction
	// Transaction commits if `f` returns nil error, rollbacks and forwards error otherwise
	// May return ErrTransactionBegin if transaction start failed,
	// ErrTransactionCommit if commit failed, ErrTransactionRollback if rollback failed
	Transactional(ctx context.Context, f func(ctx context.Context, tx DB) error) error

	io.Closer
}
