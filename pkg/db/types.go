package db

import (
	"io"

	"context"

	"errors"

	stypes "git.containerum.net/ch/solutions/pkg/models"
)

// Errors which may occur in transactional operations
var (
	ErrTransactionBegin    = errors.New("transaction begin error")
	ErrTransactionRollback = errors.New("transaction rollback error")
	ErrTransactionCommit   = errors.New("transaction commit error")
)

// DB is an interface for persistent data storage (also sometimes called DAO).
type DB interface {
	SaveAvailableSolutionsList(ctx context.Context, solution stypes.AvailableSolutionsList) error
	CreateAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error
	UpdateAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error
	DeleteAvailableSolution(ctx context.Context, solution string) error
	GetAvailableSolutionsList(ctx context.Context, isAdmin bool) (*stypes.AvailableSolutionsList, error)
	GetAvailableSolution(ctx context.Context, name string) (*stypes.AvailableSolution, error)
	ActivateAvailableSolution(ctx context.Context, solution string) error
	DeactivateAvailableSolution(ctx context.Context, solution string) error

	GetUserSolutionsList(ctx context.Context, userID string) (*stypes.UserSolutionsList, error)
	GetUserSolution(ctx context.Context, solutionName string) (*stypes.UserSolution, error)
	AddSolution(ctx context.Context, solution stypes.UserSolution, userID string, uuid string, env string) error
	AddDeployment(ctx context.Context, name string, solutionID string) error
	AddService(ctx context.Context, name string, solutionID string) error
	DeleteSolution(ctx context.Context, name string) error
	GetUserSolutionsDeployments(ctx context.Context, solutionName string) (deployments []string, ns *string, err error)
	GetUserSolutionsServices(ctx context.Context, solutionName string) (services []string, ns *string, err error)

	// Perform operations inside transaction
	// Transaction commits if `f` returns nil error, rollbacks and forwards error otherwise
	// May return ErrTransactionBegin if transaction start failed,
	// ErrTransactionCommit if commit failed, ErrTransactionRollback if rollback failed
	Transactional(ctx context.Context, f func(ctx context.Context, tx DB) error) error

	io.Closer
}
