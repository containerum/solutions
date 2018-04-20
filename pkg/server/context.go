package server

import (
	"context"
)

// Keys to inject data to context
const (
	UserIDContextKey = iota
	UserRoleContextKey
)

// MustGetUserID attempts to extract user ID using SessionIDContextKey from context.
// It panics if value was not found in context.
func MustGetUserID(ctx context.Context) string {
	uid, ok := ctx.Value(UserIDContextKey).(string)
	if !ok {
		panic("user id not found in context")
	}
	return uid
}

// MustGetUserRole attempts to extract user ID using SessionIDContextKey from context.
// It panics if value was not found in context.
func MustGetUserRole(ctx context.Context) string {
	uid, ok := ctx.Value(UserRoleContextKey).(string)
	if !ok {
		panic("user role not found in context")
	}
	return uid
}
