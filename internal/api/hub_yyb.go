package api

import (
	"context"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/yyb"
)

type yybAPI interface {
	Accounts(context.Context) ([]game.YYBAccount, error)
	CreateQR(context.Context) (game.YYBQROut, error)
	QRImage(string) ([]byte, error)
	Poll(context.Context, string) (game.YYBPollOut, error)
	Confirm(context.Context, string) (game.YYBAccount, error)
	Refresh(context.Context, string) (game.YYBRefreshOut, error)
	GetCode(context.Context, string, string) (game.YYBCodeOut, error)
	Delete(context.Context, string) (game.YYBDeleteOut, error)
	Profile(context.Context, string) (game.YYBAccount, error)
	Phone(context.Context, string, string) (game.YYBRawOut, error)
	WXData(context.Context, game.YYBWXDataIn) (game.YYBRawOut, error)
}

func (h *Hub) WithYYB(y yybAPI) *Hub {
	h.mu.Lock()
	h.yyb = y
	h.mu.Unlock()
	return h
}

func (h *Hub) yybSvc() (yybAPI, error) {
	h.mu.Lock()
	if h.yyb != nil {
		y := h.yyb
		h.mu.Unlock()
		return y, nil
	}
	h.mu.Unlock()

	svc, err := yyb.Open(yyb.DirFromEnv())
	if err != nil {
		return nil, err
	}
	h.mu.Lock()
	if h.yyb == nil {
		h.yyb = svc
	} else {
		_ = svc.Close()
	}
	y := h.yyb
	h.mu.Unlock()
	return y, nil
}

func (h *Hub) YYBAccounts(ctx context.Context) ([]game.YYBAccount, error) {
	y, err := h.yybSvc()
	if err != nil {
		return nil, err
	}
	return y.Accounts(ctx)
}

func (h *Hub) YYBCreateQR(ctx context.Context) (game.YYBQROut, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBQROut{}, err
	}
	return y.CreateQR(ctx)
}

func (h *Hub) YYBImage(sessionID string) ([]byte, error) {
	y, err := h.yybSvc()
	if err != nil {
		return nil, err
	}
	return y.QRImage(sessionID)
}

func (h *Hub) YYBPoll(ctx context.Context, sessionID string) (game.YYBPollOut, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBPollOut{}, err
	}
	return y.Poll(ctx, sessionID)
}

func (h *Hub) YYBConfirm(ctx context.Context, sessionID string) (game.YYBAccount, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBAccount{}, err
	}
	return y.Confirm(ctx, sessionID)
}

func (h *Hub) YYBRefresh(ctx context.Context, ref string) (game.YYBRefreshOut, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBRefreshOut{}, err
	}
	return y.Refresh(ctx, ref)
}

func (h *Hub) YYBCode(ctx context.Context, ref, appID string) (game.YYBCodeOut, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBCodeOut{}, err
	}
	return y.GetCode(ctx, ref, appID)
}

func (h *Hub) YYBLogin(ctx context.Context, in game.YYBRefIn) (game.User, error) {
	code, err := h.YYBCode(ctx, in.Ref, in.AppID)
	if err != nil {
		return game.User{}, err
	}
	return h.Login(ctx, game.LoginIn{Code: code.Code, OpenID: code.OpenID})
}

func (h *Hub) YYBDelete(ctx context.Context, ref string) (game.YYBDeleteOut, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBDeleteOut{}, err
	}
	return y.Delete(ctx, ref)
}

func (h *Hub) YYBProfile(ctx context.Context, ref string) (game.YYBAccount, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBAccount{}, err
	}
	return y.Profile(ctx, ref)
}

func (h *Hub) YYBPhone(ctx context.Context, ref, appID string) (game.YYBRawOut, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBRawOut{}, err
	}
	return y.Phone(ctx, ref, appID)
}

func (h *Hub) YYBWXData(ctx context.Context, in game.YYBWXDataIn) (game.YYBRawOut, error) {
	y, err := h.yybSvc()
	if err != nil {
		return game.YYBRawOut{}, err
	}
	return y.WXData(ctx, in)
}
