package session

import (
	"context"
	"errors"
)

type Session struct {
	UserID   string
	Username string
}

var (
	ErrNoAuth = errors.New("no session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)

	if !ok || sess == nil {
		return nil, ErrNoAuth
	}

	return sess, nil

}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
