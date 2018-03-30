package models

import (
	"io"

	"context"

	"errors"

	stypes "git.containerum.net/ch/json-types/solutions"
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
	GetAvailableSolutionsList(ctx context.Context) (*stypes.AvailableSolutionsList, error)
	GetAvailableSolution(ctx context.Context, name string) (*stypes.AvailableSolution, error)

	// Perform operations inside transaction
	// Transaction commits if `f` returns nil error, rollbacks and forwards error otherwise
	// May return ErrTransactionBegin if transaction start failed,
	// ErrTransactionCommit if commit failed, ErrTransactionRollback if rollback failed
	Transactional(ctx context.Context, f func(ctx context.Context, tx DB) error) error

	io.Closer
}
