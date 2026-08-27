package client

import (
	"context"
	"fmt"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func (c *Client) CleanFarmEvents(ctx context.Context) (game.CleanFarmEventOut, error) {
	user, err := c.Info()
	if err != nil {
		return game.CleanFarmEventOut{}, err
	}
	// 官方点青蛙走 Farming.clear_farm_social_item_ids=[5005]，空地。
	body, err := c.call(ctx, "Plant", "Farming", encodeFarming(nil, user.GID, []int64{frogPrankBottle}))
	if err != nil {
		return game.CleanFarmEventOut{}, err
	}
	op, err := decodeFarmingOut(body)
	if err != nil {
		return game.CleanFarmEventOut{}, err
	}
	return game.CleanFarmEventOut{Rewards: op.Rewards}, nil
}

func (c *Client) BlockApplications(ctx context.Context, in game.BlockAppsIn) (game.BlockAppsOut, error) {
	if _, err := c.Info(); err != nil {
		return game.BlockAppsOut{}, err
	}
	body, err := c.call(ctx, "Friend", "SetBlockApplications", encodeBlockApps(in.Block))
	if err != nil {
		return game.BlockAppsOut{}, err
	}
	out, err := decodeBlockApps(body)
	if err != nil {
		return game.BlockAppsOut{}, err
	}
	if !out.Block && in.Block {
		out.Block = in.Block
	}
	return out, nil
}

func (c *Client) WXRecommend(ctx context.Context, in game.WXRecommendIn) (game.WXRecommendOut, error) {
	if _, err := c.Info(); err != nil {
		return game.WXRecommendOut{}, err
	}
	body, err := c.call(ctx, "Friend", "GetWXRecommendations", encodeEncrypted(in.Encrypted))
	if err != nil {
		return game.WXRecommendOut{}, err
	}
	return decodeWXRecommend(body)
}

func (c *Client) WXRecommendPage(ctx context.Context, in game.WXRecommendPageIn) (game.WXRecommendOut, error) {
	if _, err := c.Info(); err != nil {
		return game.WXRecommendOut{}, err
	}
	body, err := c.call(ctx, "Friend", "GetWXRecommendationsPage", encodeWXPage(in.Offset, in.PageSize))
	if err != nil {
		return game.WXRecommendOut{}, err
	}
	return decodeWXRecommend(body)
}

func (c *Client) ApplyWXFriends(ctx context.Context, in game.GIDsIn) (game.WXApplyOut, error) {
	if _, err := c.Info(); err != nil {
		return game.WXApplyOut{}, err
	}
	if len(in.GIDs) == 0 {
		return game.WXApplyOut{}, fmt.Errorf("gids 不能为空")
	}
	body, err := c.call(ctx, "Friend", "ApplyWXFriends", encodeGIDs(in.GIDs))
	if err != nil {
		return game.WXApplyOut{}, err
	}
	return decodeWXApply(body)
}

func (c *Client) EquipSkinSet(ctx context.Context, in game.SkinSetIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if len(in.SkinIDs) == 0 {
		return fmt.Errorf("skin_ids 不能为空")
	}
	_, err := c.call(ctx, "Skin", "EquipSkinSet", encodeSkinSet(in.SkinIDs))
	return err
}

func (c *Client) SetSkinSetEffect(ctx context.Context, in game.SkinSetEffectIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.SetID <= 0 {
		return fmt.Errorf("set_id 不能为空")
	}
	_, err := c.call(ctx, "Skin", "SetSkinSetEffect", encodeSkinSetEffect(in))
	return err
}

func (c *Client) SkinSets(ctx context.Context) ([]game.SkinSetEffect, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Skin", "SkinsEquipped", nil)
	if err != nil {
		return nil, err
	}
	return decodeSkinSets(body)
}

func (c *Client) BuyPass(ctx context.Context) (game.BuyPassOut, error) {
	if _, err := c.Info(); err != nil {
		return game.BuyPassOut{}, err
	}
	body, err := c.call(ctx, "Season", "BuyPremiumBattlePass", nil)
	if err != nil {
		return game.BuyPassOut{}, err
	}
	out, err := decodeBuyPass(body)
	if err != nil {
		return game.BuyPassOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) MarkSeasonOpening(ctx context.Context) (bool, error) {
	if _, err := c.Info(); err != nil {
		return false, err
	}
	body, err := c.call(ctx, "Season", "MarkSeasonOpeningShown", nil)
	if err != nil {
		return false, err
	}
	return decodeSuccess(body)
}

func (c *Client) QQAuthGroups(ctx context.Context, in game.CookiesIn) (game.QQAuthGroupsOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQAuthGroupsOut{}, err
	}
	body, err := c.call(ctx, "QQGroup", "GetMyAuthGroups", encodeCookies(in.Cookies))
	if err != nil {
		return game.QQAuthGroupsOut{}, err
	}
	return decodeQQAuthGroups(body)
}

func (c *Client) QQRecommendGroups(ctx context.Context, in game.QQRecommendIn) (game.QQRecommendOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQRecommendOut{}, err
	}
	body, err := c.call(ctx, "QQGroup", "GetRecommendGroups", encodeQQRecommend(in))
	if err != nil {
		return game.QQRecommendOut{}, err
	}
	return decodeQQRecommend(body)
}

func (c *Client) QQBind(ctx context.Context, in game.QQBindIn) (game.QQBindOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQBindOut{}, err
	}
	if in.CommunityID == "" {
		return game.QQBindOut{}, fmt.Errorf("community_id 不能为空")
	}
	body, err := c.call(ctx, "QQGroup", "BindCommunity", encodeCommunityID(in.CommunityID))
	if err != nil {
		return game.QQBindOut{}, err
	}
	return decodeQQBind(body)
}

func (c *Client) QQLeave(ctx context.Context) (game.QQLeaveOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQLeaveOut{}, err
	}
	body, err := c.call(ctx, "QQGroup", "LeaveCommunity", nil)
	if err != nil {
		return game.QQLeaveOut{}, err
	}
	return decodeQQLeave(body)
}

func (c *Client) QQCommunity(ctx context.Context, in game.PageIn) (game.QQCommunityOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQCommunityOut{}, err
	}
	page := in.From
	if page <= 0 {
		page = 1
	}
	req := encodeUID(page)
	body, err := c.call(ctx, "QQGroup", "GetCommunityInfo", req)
	if err != nil {
		return game.QQCommunityOut{}, err
	}
	return decodeQQCommunityOut(body)
}

func (c *Client) QQBindInfo(ctx context.Context) (game.QQBindInfoOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQBindInfoOut{}, err
	}
	body, err := c.call(ctx, "QQGroup", "GetBindInfo", nil)
	if err != nil {
		return game.QQBindInfoOut{}, err
	}
	return decodeQQBindInfo(body)
}

func (c *Client) QQClaimReward(ctx context.Context) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "QQGroup", "ClaimCommunityReward", nil)
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) QQRevokeAuth(ctx context.Context, in game.QQRevokeIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.CommunityID == "" {
		return fmt.Errorf("community_id 不能为空")
	}
	_, err := c.call(ctx, "QQGroup", "RevokeAuth", encodeQQRevoke(in))
	return err
}

func (c *Client) UseGiftToken(ctx context.Context, in game.UIDIn) (game.GiftTokenOut, error) {
	if _, err := c.Info(); err != nil {
		return game.GiftTokenOut{}, err
	}
	if in.UID <= 0 {
		return game.GiftTokenOut{}, fmt.Errorf("uid 不能为空")
	}
	body, err := c.call(ctx, "Gift", "UseGiftToken", encodeUID(in.UID))
	if err != nil {
		return game.GiftTokenOut{}, err
	}
	return decodeGiftToken(body)
}

func (c *Client) GiftHistory(ctx context.Context, in game.GiftHistoryIn) (game.GiftHistoryOut, error) {
	if _, err := c.Info(); err != nil {
		return game.GiftHistoryOut{}, err
	}
	body, err := c.call(ctx, "Gift", "GetGiftClaimHistory", encodeGiftHistory(in))
	if err != nil {
		return game.GiftHistoryOut{}, err
	}
	return decodeGiftHistory(body)
}

func (c *Client) TransferStatus(ctx context.Context, in game.TransferIn) (game.TransferOut, error) {
	if _, err := c.Info(); err != nil {
		return game.TransferOut{}, err
	}
	if in.OutBillNo == "" {
		return game.TransferOut{}, fmt.Errorf("out_bill_no 不能为空")
	}
	body, err := c.call(ctx, "Gift", "QueryTransferStatus", encodeOutBill(in.OutBillNo))
	if err != nil {
		return game.TransferOut{}, err
	}
	return decodeTransfer(body)
}

func (c *Client) CancelTransfer(ctx context.Context, in game.UIDIn) (game.TransferOut, error) {
	if _, err := c.Info(); err != nil {
		return game.TransferOut{}, err
	}
	if in.UID <= 0 {
		return game.TransferOut{}, fmt.Errorf("uid 不能为空")
	}
	body, err := c.call(ctx, "Gift", "CancelRedpacketTransfer", encodeUID(in.UID))
	if err != nil {
		return game.TransferOut{}, err
	}
	return decodeTransfer(body)
}

func (c *Client) FollowGiftStatus(ctx context.Context) (game.FollowGiftOut, error) {
	if _, err := c.Info(); err != nil {
		return game.FollowGiftOut{}, err
	}
	body, err := c.call(ctx, "Misc", "GetFollowGiftStatus", nil)
	if err != nil {
		return game.FollowGiftOut{}, err
	}
	return decodeFollowGift(body)
}

func (c *Client) SetFollowGift(ctx context.Context, in game.FollowGiftIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	_, err := c.call(ctx, "Misc", "SetFollowGiftStatus", encodeBlockApps(in.Followed))
	return err
}

func (c *Client) ClaimFollowGift(ctx context.Context) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Misc", "ClaimFollowGift", nil)
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) RechargeBonus(ctx context.Context) (game.RechargeBonusOut, error) {
	if _, err := c.Info(); err != nil {
		return game.RechargeBonusOut{}, err
	}
	body, err := c.call(ctx, "RechargeBonus", "GetConfig", nil)
	if err != nil {
		return game.RechargeBonusOut{}, err
	}
	return decodeRechargeBonus(body)
}

func (c *Client) RechargeBonusData(ctx context.Context) (game.RechargeDataOut, error) {
	if _, err := c.Info(); err != nil {
		return game.RechargeDataOut{}, err
	}
	body, err := c.call(ctx, "RechargeBonus", "GetPlayerData", nil)
	if err != nil {
		return game.RechargeDataOut{}, err
	}
	return decodeRechargeData(body)
}

func (c *Client) SetDisplay(ctx context.Context, in game.DisplayIn) (game.DisplayOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DisplayOut{}, err
	}
	body, err := c.call(ctx, "User", "SetDisplayInfo", encodeDisplay(in))
	if err != nil {
		return game.DisplayOut{}, err
	}
	return decodeDisplay(body)
}

func (c *Client) GetSettings(ctx context.Context, in game.SettingsKeysIn) (game.UserSettings, error) {
	if _, err := c.Info(); err != nil {
		return game.UserSettings{}, err
	}
	body, err := c.call(ctx, "User", "GetUserSettings", encodeSettingsKeys(in.Keys))
	if err != nil {
		return game.UserSettings{}, err
	}
	return decodeUserSettings(body)
}

func (c *Client) SetSettings(ctx context.Context, in game.UserSettings) (game.UserSettings, error) {
	if _, err := c.Info(); err != nil {
		return game.UserSettings{}, err
	}
	body, err := c.call(ctx, "User", "SetUserSettings", encodeUserSettings(in))
	if err != nil {
		return game.UserSettings{}, err
	}
	return decodeUserSettings(body)
}

func (c *Client) DeleteAccount(ctx context.Context, in game.DeleteAccountIn) (game.DeleteAccountOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DeleteAccountOut{}, err
	}
	if in.Name == "" || in.CertID == "" {
		return game.DeleteAccountOut{}, fmt.Errorf("name 和 cert_id 不能为空")
	}
	body, err := c.call(ctx, "User", "DeleteAccount", encodeDeleteAccount(in))
	if err != nil {
		return game.DeleteAccountOut{}, err
	}
	return decodeDeleteAccount(body)
}

func (c *Client) DecryptOpenData(ctx context.Context, in game.DecryptIn) (string, error) {
	if _, err := c.Info(); err != nil {
		return "", err
	}
	if in.Encrypted == "" {
		return "", fmt.Errorf("encrypted_data 不能为空")
	}
	body, err := c.call(ctx, "User", "DecryptOpenData", encodeEncrypted(in.Encrypted))
	if err != nil {
		return "", err
	}
	return decodeDecrypt(body)
}

func (c *Client) SetQQRecommendAuth(ctx context.Context, in game.QQAuthIn) (game.QQAuthOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQAuthOut{}, err
	}
	body, err := c.call(ctx, "User", "SetQQFriendRecommendAuthorized", encodeQQAuth(in.Authorized))
	if err != nil {
		return game.QQAuthOut{}, err
	}
	return decodeQQAuth(body)
}

func (c *Client) ReportFlow(ctx context.Context, in game.ReportFlowIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	_, err := c.call(ctx, "User", "ClientReportFlow", encodeReportFlow(in))
	return err
}

func (c *Client) BatchReportFlow(ctx context.Context, in game.BatchReportFlowIn) (game.BatchReportFlowOut, error) {
	if _, err := c.Info(); err != nil {
		return game.BatchReportFlowOut{}, err
	}
	if len(in.Flows) == 0 {
		return game.BatchReportFlowOut{}, fmt.Errorf("flows 不能为空")
	}
	body, err := c.call(ctx, "User", "BatchClientReportFlow", encodeBatchReportFlow(in.Flows))
	if err != nil {
		return game.BatchReportFlowOut{}, err
	}
	return decodeBatchReport(body)
}

func (c *Client) ReportUser(ctx context.Context, in game.ReportUserIn) (game.ReportUserOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ReportUserOut{}, err
	}
	if in.GID <= 0 {
		return game.ReportUserOut{}, fmt.Errorf("gid 不能为空")
	}
	body, err := c.call(ctx, "User", "Report", encodeReportUser(in))
	if err != nil {
		return game.ReportUserOut{}, err
	}
	return decodeReportUser(body)
}
