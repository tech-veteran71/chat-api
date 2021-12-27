package main

import (
	"context"
	"errors"
)

type userKey string

var (
	ctxUserKey     = userKey("username")
	ErrInvalidUser = errors.New("invalid user")
)

// User is an authenticated system user.
type User struct {
	ID    int64
	Name  string
	Label string // real name
	Token string
}

// ContextUser returns User from the context.
func ContextUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(ctxUserKey).(*User)
	return user, ok
}

// ContextWithUser returns a new context with the User.
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, ctxUserKey, user)
}
