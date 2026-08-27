package api

import (
	"context"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func (h *Hub) CleanFarmEvents(ctx context.Context) (game.CleanFarmEventOut, error) {
	s, err := h.current()
	if err != nil {
		return game.CleanFarmEventOut{}, err
	}
	return s.CleanFarmEvents(ctx)
}

func (h *Hub) BlockApplications(ctx context.Context, in game.BlockAppsIn) (game.BlockAppsOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BlockAppsOut{}, err
	}
	return s.BlockApplications(ctx, in)
}

func (h *Hub) WXRecommend(ctx context.Context, in game.WXRecommendIn) (game.WXRecommendOut, error) {
	s, err := h.current()
	if err != nil {
		return game.WXRecommendOut{}, err
	}
	return s.WXRecommend(ctx, in)
}

func (h *Hub) WXRecommendPage(ctx context.Context, in game.WXRecommendPageIn) (game.WXRecommendOut, error) {
	s, err := h.current()
	if err != nil {
		return game.WXRecommendOut{}, err
	}
	return s.WXRecommendPage(ctx, in)
}

func (h *Hub) ApplyWXFriends(ctx context.Context, in game.GIDsIn) (game.WXApplyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.WXApplyOut{}, err
	}
	return s.ApplyWXFriends(ctx, in)
}

func (h *Hub) EquipSkinSet(ctx context.Context, in game.SkinSetIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.EquipSkinSet(ctx, in)
}

func (h *Hub) SetSkinSetEffect(ctx context.Context, in game.SkinSetEffectIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.SetSkinSetEffect(ctx, in)
}

func (h *Hub) SkinSets(ctx context.Context) ([]game.SkinSetEffect, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.SkinSets(ctx)
}

func (h *Hub) BuyPass(ctx context.Context) (game.BuyPassOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BuyPassOut{}, err
	}
	return s.BuyPass(ctx)
}

func (h *Hub) MarkSeasonOpening(ctx context.Context) (bool, error) {
	s, err := h.current()
	if err != nil {
		return false, err
	}
	return s.MarkSeasonOpening(ctx)
}

func (h *Hub) QQAuthGroups(ctx context.Context, in game.CookiesIn) (game.QQAuthGroupsOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQAuthGroupsOut{}, err
	}
	return s.QQAuthGroups(ctx, in)
}

func (h *Hub) QQRecommendGroups(ctx context.Context, in game.QQRecommendIn) (game.QQRecommendOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQRecommendOut{}, err
	}
	return s.QQRecommendGroups(ctx, in)
}

func (h *Hub) QQBind(ctx context.Context, in game.QQBindIn) (game.QQBindOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQBindOut{}, err
	}
	return s.QQBind(ctx, in)
}

func (h *Hub) QQLeave(ctx context.Context) (game.QQLeaveOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQLeaveOut{}, err
	}
	return s.QQLeave(ctx)
}

func (h *Hub) QQCommunity(ctx context.Context, in game.PageIn) (game.QQCommunityOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQCommunityOut{}, err
	}
	return s.QQCommunity(ctx, in)
}

func (h *Hub) QQBindInfo(ctx context.Context) (game.QQBindInfoOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQBindInfoOut{}, err
	}
	return s.QQBindInfo(ctx)
}

func (h *Hub) QQClaimReward(ctx context.Context) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.QQClaimReward(ctx)
}

func (h *Hub) QQRevokeAuth(ctx context.Context, in game.QQRevokeIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.QQRevokeAuth(ctx, in)
}

func (h *Hub) UseGiftToken(ctx context.Context, in game.UIDIn) (game.GiftTokenOut, error) {
	s, err := h.current()
	if err != nil {
		return game.GiftTokenOut{}, err
	}
	return s.UseGiftToken(ctx, in)
}

func (h *Hub) GiftHistory(ctx context.Context, in game.GiftHistoryIn) (game.GiftHistoryOut, error) {
	s, err := h.current()
	if err != nil {
		return game.GiftHistoryOut{}, err
	}
	return s.GiftHistory(ctx, in)
}

func (h *Hub) TransferStatus(ctx context.Context, in game.TransferIn) (game.TransferOut, error) {
	s, err := h.current()
	if err != nil {
		return game.TransferOut{}, err
	}
	return s.TransferStatus(ctx, in)
}

func (h *Hub) CancelTransfer(ctx context.Context, in game.UIDIn) (game.TransferOut, error) {
	s, err := h.current()
	if err != nil {
		return game.TransferOut{}, err
	}
	return s.CancelTransfer(ctx, in)
}

func (h *Hub) FollowGiftStatus(ctx context.Context) (game.FollowGiftOut, error) {
	s, err := h.current()
	if err != nil {
		return game.FollowGiftOut{}, err
	}
	return s.FollowGiftStatus(ctx)
}

func (h *Hub) SetFollowGift(ctx context.Context, in game.FollowGiftIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.SetFollowGift(ctx, in)
}

func (h *Hub) ClaimFollowGift(ctx context.Context) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimFollowGift(ctx)
}

func (h *Hub) RechargeBonus(ctx context.Context) (game.RechargeBonusOut, error) {
	s, err := h.current()
	if err != nil {
		return game.RechargeBonusOut{}, err
	}
	return s.RechargeBonus(ctx)
}

func (h *Hub) RechargeBonusData(ctx context.Context) (game.RechargeDataOut, error) {
	s, err := h.current()
	if err != nil {
		return game.RechargeDataOut{}, err
	}
	return s.RechargeBonusData(ctx)
}

func (h *Hub) SetDisplay(ctx context.Context, in game.DisplayIn) (game.DisplayOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DisplayOut{}, err
	}
	return s.SetDisplay(ctx, in)
}

func (h *Hub) GetSettings(ctx context.Context, in game.SettingsKeysIn) (game.UserSettings, error) {
	s, err := h.current()
	if err != nil {
		return game.UserSettings{}, err
	}
	return s.GetSettings(ctx, in)
}

func (h *Hub) SetSettings(ctx context.Context, in game.UserSettings) (game.UserSettings, error) {
	s, err := h.current()
	if err != nil {
		return game.UserSettings{}, err
	}
	return s.SetSettings(ctx, in)
}

func (h *Hub) DeleteAccount(ctx context.Context, in game.DeleteAccountIn) (game.DeleteAccountOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DeleteAccountOut{}, err
	}
	return s.DeleteAccount(ctx, in)
}

func (h *Hub) DecryptOpenData(ctx context.Context, in game.DecryptIn) (string, error) {
	s, err := h.current()
	if err != nil {
		return "", err
	}
	return s.DecryptOpenData(ctx, in)
}

func (h *Hub) SetQQRecommendAuth(ctx context.Context, in game.QQAuthIn) (game.QQAuthOut, error) {
	s, err := h.current()
	if err != nil {
		return game.QQAuthOut{}, err
	}
	return s.SetQQRecommendAuth(ctx, in)
}

func (h *Hub) ReportFlow(ctx context.Context, in game.ReportFlowIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.ReportFlow(ctx, in)
}

func (h *Hub) BatchReportFlow(ctx context.Context, in game.BatchReportFlowIn) (game.BatchReportFlowOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BatchReportFlowOut{}, err
	}
	return s.BatchReportFlow(ctx, in)
}

func (h *Hub) ReportUser(ctx context.Context, in game.ReportUserIn) (game.ReportUserOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ReportUserOut{}, err
	}
	return s.ReportUser(ctx, in)
}
