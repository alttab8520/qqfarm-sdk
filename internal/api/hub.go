package api

import (
	"context"
	"errors"
	"sync"

	"github.com/alttab8520/qqfarm-sdk/internal/client"
	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

type Hub struct {
	newS game.Factory
	mu   sync.Mutex
	sess game.Session
}

func NewHub(newS game.Factory) *Hub {
	if newS == nil {
		newS = client.New
	}
	return &Hub{newS: newS}
}

func (h *Hub) Login(ctx context.Context, in game.LoginIn) (game.User, error) {
	s := h.newS()
	user, err := s.Login(ctx, in)
	if err != nil {
		_ = s.Close()
		return game.User{}, err
	}
	h.mu.Lock()
	if h.sess != nil {
		_ = h.sess.Close()
	}
	h.sess = s
	h.mu.Unlock()
	return user, nil
}

func (h *Hub) current() (game.Session, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.sess == nil {
		return nil, game.ErrNotLogin
	}
	return h.sess, nil
}

func (h *Hub) Info() (game.User, error) {
	s, err := h.current()
	if err != nil {
		return game.User{}, err
	}
	return s.Info()
}

func (h *Hub) Refresh(ctx context.Context) ([]game.Land, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Refresh(ctx)
}

func (h *Hub) Harvest(ctx context.Context, in game.HarvestIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Harvest(ctx, in)
}

func (h *Hub) Plant(ctx context.Context, in game.PlantIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Plant(ctx, in)
}

func (h *Hub) Friends(ctx context.Context) ([]game.Friend, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Friends(ctx)
}

func (h *Hub) Help(ctx context.Context, in game.HelpIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Help(ctx, in)
}

func failCode(err error) (int, string) {
	if err == nil {
		return 0, "ok"
	}
	if errors.Is(err, game.ErrNotLogin) {
		return 401, err.Error()
	}
	return 502, err.Error()
}
