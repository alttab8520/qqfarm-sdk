package api

import (
	"context"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func (h *Hub) QQVipDailyStatus(ctx context.Context) (game.QQVipDailyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQVipDailyOut{}, err
	}
	return s.QQVipDailyStatus(ctx)
}

func (h *Hub) QQVipClaimDaily(ctx context.Context, in game.QQVipClaimDailyIn) (game.QQVipClaimDailyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQVipClaimDailyOut{}, err
	}
	return s.QQVipClaimDaily(ctx, in)
}

func (h *Hub) QQVipRefresh(ctx context.Context) (game.QQVipRefreshOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQVipRefreshOut{}, err
	}
	return s.QQVipRefresh(ctx)
}

func (h *Hub) QQVipClaimRewards(ctx context.Context, in game.QQVipClaimRewardsIn) (game.QQVipClaimRewardsOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQVipClaimRewardsOut{}, err
	}
	return s.QQVipClaimRewards(ctx, in)
}

func (h *Hub) QQVipRewardsStatus(ctx context.Context) (game.QQVipRewardsStatusOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQVipRewardsStatusOut{}, err
	}
	return s.QQVipRewardsStatus(ctx)
}

func (h *Hub) QQVipMarkRedpoint(ctx context.Context) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.QQVipMarkRedpoint(ctx)
}

func (h *Hub) Marquee(ctx context.Context) (game.MarqueeOut, error) {
	s, err := h.current()
	if err != nil {
		return game.MarqueeOut{}, err
	}
	return s.Marquee(ctx)
}

func (h *Hub) SystemUnlocked(ctx context.Context, in game.SystemOpenIn) (game.SystemOpenOut, error) {
	s, err := h.current()
	if err != nil {
		return game.SystemOpenOut{}, err
	}
	return s.SystemUnlocked(ctx, in)
}

func (h *Hub) MutantOpenInfo(ctx context.Context) (game.MutantOpenInfoOut, error) {
	s, err := h.current()
	if err != nil {
		return game.MutantOpenInfoOut{}, err
	}
	return s.MutantOpenInfo(ctx)
}

func (h *Hub) QQSubscribe(ctx context.Context) (game.QQSubscribeOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQSubscribeOut{}, err
	}
	return s.QQSubscribe(ctx)
}

func (h *Hub) WXSubscribe(ctx context.Context) (game.WXSubscribeOut, error) {
	s, err := h.current()
	if err != nil {
		return game.WXSubscribeOut{}, err
	}
	return s.WXSubscribe(ctx)
}

func (h *Hub) SetWXSubscribe(ctx context.Context, in game.WXSubscribeIn) (game.WXSubscribeOut, error) {
	s, err := h.current()
	if err != nil {
		return game.WXSubscribeOut{}, err
	}
	return s.SetWXSubscribe(ctx, in)
}

func (h *Hub) ModerateText(ctx context.Context, in game.ModerateTextIn) (game.ModerateTextOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ModerateTextOut{}, err
	}
	return s.ModerateText(ctx, in)
}

func (h *Hub) BatchModerateText(ctx context.Context, in game.ModerateTextBatchIn) (game.ModerateTextBatchOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ModerateTextBatchOut{}, err
	}
	return s.BatchModerateText(ctx, in)
}

func (h *Hub) ModeratePic(ctx context.Context, in game.ModeratePicIn) (game.ModeratePicOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ModeratePicOut{}, err
	}
	return s.ModeratePic(ctx, in)
}

func (h *Hub) BatchModeratePic(ctx context.Context, in game.ModeratePicBatchIn) (game.ModeratePicBatchOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ModeratePicBatchOut{}, err
	}
	return s.BatchModeratePic(ctx, in)
}
