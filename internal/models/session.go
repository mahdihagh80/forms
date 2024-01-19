package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mahdihagh80/forms/internal/services"
	"github.com/redis/go-redis/v9"
)

type sessionStore struct {
	store *redis.Client
}

func NewSessionStore(store *redis.Client) sessionStore {
	return sessionStore{
		store: store,
	}
}

func (ss sessionStore) Create(ctx context.Context, session services.SessionData) error {
	fmt.Println("helpppppppppppppppppppp\nppppppppppppppppppppppp")
	ttl := time.Duration(session.MaxAge) * time.Second
	fmt.Println("ttl : ", ttl)
	err := ss.store.Set(ctx, session.SessionId, session.UserId, time.Duration(session.MaxAge)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error while setting session in session store : %w", err)
	}
	return nil
}

func (ss sessionStore) Get(ctx context.Context, session *services.SessionData) error {
	userIdStr, err := ss.store.Get(ctx, session.SessionId).Result()
	if err != nil {
		return fmt.Errorf("error while getting session from session store : %w", err)
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return fmt.Errorf("error while converting session value : %w", err)
	}
	session.UserId = userId
	return nil
}

func (ss sessionStore) Delete(ctx context.Context, sessionId string) error {
	err := ss.store.Del(ctx, sessionId).Err()
	if err != nil {
		return fmt.Errorf("error while deleting session from session store : %w", err)
	}
	return nil
}
