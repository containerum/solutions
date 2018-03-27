package server

import (
	"context"

	"io"

	"git.containerum.net/ch/solutions/pkg/models"
)

// SolutionsService is an interface for server "business logic"
type SolutionsService interface {
	Test1(ctx context.Context) error
	Test2(ctx context.Context) error

	io.Closer
}

// Services is a collection of resources needed for server functionality.
type Services struct {
	DB models.DB
}
