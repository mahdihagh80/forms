package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SessionService interface {
	Create(ctx context.Context) (SessionData, error)
	Identify(ctx context.Context, session *SessionData) error
	Terminate(ctx context.Context, sessionId string) error
}

type SessionStore interface {
	Create(ctx context.Context, session SessionData) error
	Get(ctx context.Context, session *SessionData) error
	Delete(ctx context.Context, sessionId string) error
}

type SessionData struct {
	UserId    int
	SessionId string
	MaxAge    int
}

type Session struct {
	store SessionStore
}

func NewSessionService(store SessionStore) Session {
	return Session{
		store: store,
	}
}

func (s Session) Create(ctx context.Context) (SessionData, error) {
	userId, exists := ctx.Value("userId").(int)
	if !exists {
		return SessionData{}, errors.New("userId with type int doesn't exists in context's value")
	}

	sessionId, err := GenerateSessionID()
	if err != nil {
		return SessionData{}, fmt.Errorf("error while generating sessionId : %w", err)
	}

	session := SessionData{
		UserId:    userId,
		SessionId: sessionId,
		MaxAge:    3600,
	}

	if err := s.store.Create(ctx, session); err != nil {
		return SessionData{}, fmt.Errorf("error while creating session for user %d : %w", userId, err)
	}
	return session, nil
}

func (s Session) Identify(ctx context.Context, session *SessionData) error {
	err := s.store.Get(ctx, session)
	if err != nil {
		return fmt.Errorf("error while getting session from session store : %w", err)
	}
	return nil
}

func (s Session) Terminate(ctx context.Context, sessionId string) error {
	if err := s.store.Delete(ctx, sessionId); err != nil {
		return fmt.Errorf("error while Terminating session %s : %w", sessionId, err)
	}
	return nil
}

func GenerateSessionID() (string, error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuidObj.String(), nil
}
