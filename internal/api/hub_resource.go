package api

import (
	"context"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/resource"
)

type resourceAPI interface {
	Lookup(game.ResLookupIn) (game.ResListOut, error)
	List(game.ResListIn) (game.ResListOut, error)
	Tables() game.ResTablesOut
	Refresh(context.Context, game.ResRefreshIn) (game.ResRefreshOut, error)
}

func (h *Hub) WithResources(r resourceAPI) *Hub {
	h.mu.Lock()
	h.res = r
	h.mu.Unlock()
	return h
}

func (h *Hub) resSvc() resourceAPI {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.res == nil {
		h.res = resource.Open()
	}
	return h.res
}

func (h *Hub) ResLookup(in game.ResLookupIn) (game.ResListOut, error) {
	return h.resSvc().Lookup(in)
}

func (h *Hub) ResList(in game.ResListIn) (game.ResListOut, error) {
	return h.resSvc().List(in)
}

func (h *Hub) ResTables() game.ResTablesOut {
	return h.resSvc().Tables()
}

func (h *Hub) ResRefresh(ctx context.Context, in game.ResRefreshIn) (game.ResRefreshOut, error) {
	return h.resSvc().Refresh(ctx, in)
}
