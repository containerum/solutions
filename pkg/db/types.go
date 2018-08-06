package db

import (
	"io"

	"context"

	"errors"

	kube_types "github.com/containerum/kube-client/pkg/model"
)

// Errors which may occur in transactional operations
var (
	ErrTransactionBegin    = errors.New("transaction begin error")
	ErrTransactionRollback = errors.New("transaction rollback error")
	ErrTransactionCommit   = errors.New("transaction commit error")
)

// DB is an interface for persistent data storage (also sometimes called DAO).
type DB interface {
	CreateTemplate(ctx context.Context, solution kube_types.SolutionTemplate) error
	UpdateTemplate(ctx context.Context, solution kube_types.SolutionTemplate) error
	DeleteTemplate(ctx context.Context, solution string) error
	GetTemplatesList(ctx context.Context, isAdmin bool) (*kube_types.SolutionsTemplatesList, error)
	GetTemplate(ctx context.Context, name string) (*kube_types.SolutionTemplate, error)
	ActivateTemplate(ctx context.Context, solution string) error
	DeactivateTemplate(ctx context.Context, solution string) error

	GetSolutionsList(ctx context.Context, userID string) (*kube_types.SolutionsList, error)
	GetNamespaceSolutionsList(ctx context.Context, namespace string) (*kube_types.SolutionsList, error)
	GetSolution(ctx context.Context, namespace, solutionName string) (*kube_types.Solution, error)
	AddSolution(ctx context.Context, solution kube_types.Solution, userID, templateID, uuid, env string) error
	DeleteSolution(ctx context.Context, namespace, solutionName string) error
	CompletelyDeleteSolution(ctx context.Context, namespace, solutionName string) error
	CompletelyDeleteSolutions(ctx context.Context, userID string) error
	CompletelyDeleteNamespaceSolutions(ctx context.Context, namespace string) error

	// Perform operations inside transaction
	// Transaction commits if `f` returns nil error, rollbacks and forwards error otherwise
	// May return ErrTransactionBegin if transaction start failed,
	// ErrTransactionCommit if commit failed, ErrTransactionRollback if rollback failed
	Transactional(ctx context.Context, f func(ctx context.Context, tx DB) error) error

	io.Closer
}
