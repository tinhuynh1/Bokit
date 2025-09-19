package domain

import (
	"context"
)

type SessionCollection interface {
	CreateSession(ctx context.Context, session *SessionCreate) error
	GetSessionByCode(ctx context.Context, code string) (*Session, error)
	// GetSessions(ctx context.Context) ([]Session, error)
	// UpdateSession(ctx context.Context, session *Session) error
}
