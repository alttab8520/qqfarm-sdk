package client

import (
	"context"
	"fmt"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func (c *Client) QQVipDailyStatus(ctx context.Context) (game.QQVipDailyOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQVipDailyOut{}, err
	}
	body, err := c.call(ctx, "QQVip", "GetDailyGiftStatus", nil)
	if err != nil {
		return game.QQVipDailyOut{}, err
	}
	return decodeQQVipDaily(body)
}

func (c *Client) QQVipClaimDaily(ctx context.Context, in game.QQVipClaimDailyIn) (game.QQVipClaimDailyOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQVipClaimDailyOut{}, err
	}
	body, err := c.call(ctx, "QQVip", "ClaimDailyGift", encodeQQVipClaimDaily(in.ConfigID))
	if err != nil {
		return game.QQVipClaimDailyOut{}, err
	}
	return decodeQQVipClaimDaily(body)
}

func (c *Client) QQVipRefresh(ctx context.Context) (game.QQVipRefreshOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQVipRefreshOut{}, err
	}
	body, err := c.call(ctx, "QQVip", "RefreshVipInfo", nil)
	if err != nil {
		return game.QQVipRefreshOut{}, err
	}
	return decodeQQVipRefresh(body)
}

func (c *Client) QQVipClaimRewards(ctx context.Context, in game.QQVipClaimRewardsIn) (game.QQVipClaimRewardsOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQVipClaimRewardsOut{}, err
	}
	body, err := c.call(ctx, "QQVip", "ClaimQQVipRewards", encodeQQVipClaimRewards(in.ConfigIDs))
	if err != nil {
		return game.QQVipClaimRewardsOut{}, err
	}
	return decodeQQVipClaimRewards(body)
}

func (c *Client) QQVipRewardsStatus(ctx context.Context) (game.QQVipRewardsStatusOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQVipRewardsStatusOut{}, err
	}
	body, err := c.call(ctx, "QQVip", "GetQQVipRewardsStatus", nil)
	if err != nil {
		return game.QQVipRewardsStatusOut{}, err
	}
	return decodeQQVipRewardsStatus(body)
}

func (c *Client) QQVipMarkRedpoint(ctx context.Context) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	_, err := c.call(ctx, "QQVip", "MarkQQVipRedpointViewed", nil)
	return err
}

func (c *Client) Marquee(ctx context.Context) (game.MarqueeOut, error) {
	if _, err := c.Info(); err != nil {
		return game.MarqueeOut{}, err
	}
	body, err := c.call(ctx, "Marquee", "GetMarquee", nil)
	if err != nil {
		return game.MarqueeOut{}, err
	}
	return decodeMarquee(body)
}

func (c *Client) SystemUnlocked(ctx context.Context, in game.SystemOpenIn) (game.SystemOpenOut, error) {
	if _, err := c.Info(); err != nil {
		return game.SystemOpenOut{}, err
	}
	body, err := c.call(ctx, "SystemOpen", "SSIsSystemUnlocked", encodeSystemName(in.SystemName))
	if err != nil {
		return game.SystemOpenOut{}, err
	}
	return decodeSystemOpen(body)
}

func (c *Client) MutantOpenInfo(ctx context.Context) (game.MutantOpenInfoOut, error) {
	if _, err := c.Info(); err != nil {
		return game.MutantOpenInfoOut{}, err
	}
	body, err := c.call(ctx, "Mutant", "GetSystemOpenInfo", nil)
	if err != nil {
		return game.MutantOpenInfoOut{}, err
	}
	return decodeMutantOpenInfo(body)
}

func (c *Client) QQSubscribe(ctx context.Context) (game.QQSubscribeOut, error) {
	if _, err := c.Info(); err != nil {
		return game.QQSubscribeOut{}, err
	}
	body, err := c.call(ctx, "SubscribeQQ", "GetQQSubscribeStatus", nil)
	if err != nil {
		return game.QQSubscribeOut{}, err
	}
	return decodeQQSubscribe(body)
}

func (c *Client) WXSubscribe(ctx context.Context) (game.WXSubscribeOut, error) {
	if _, err := c.Info(); err != nil {
		return game.WXSubscribeOut{}, err
	}
	body, err := c.call(ctx, "SubscribeWX", "GetSubscribeMessageStatus", nil)
	if err != nil {
		return game.WXSubscribeOut{}, err
	}
	return decodeWXSubscribe(body)
}

func (c *Client) SetWXSubscribe(ctx context.Context, in game.WXSubscribeIn) (game.WXSubscribeOut, error) {
	if _, err := c.Info(); err != nil {
		return game.WXSubscribeOut{}, err
	}
	if len(in.Templates) == 0 {
		return game.WXSubscribeOut{}, fmt.Errorf("templates 不能为空")
	}
	body, err := c.call(ctx, "SubscribeWX", "SetSubscribeMessageStatus", encodeWXSubscribe(in.Templates))
	if err != nil {
		return game.WXSubscribeOut{}, err
	}
	out, err := decodeWXSubscribe(body)
	if err != nil {
		return game.WXSubscribeOut{}, err
	}
	if len(out.Templates) == 0 {
		out.Templates = in.Templates
	}
	return out, nil
}

func (c *Client) ModerateText(ctx context.Context, in game.ModerateTextIn) (game.ModerateTextOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ModerateTextOut{}, err
	}
	if in.Text == "" {
		return game.ModerateTextOut{}, fmt.Errorf("text 不能为空")
	}
	body, err := c.call(ctx, "Uicproxy", "ModerateText", encodeModerateText(in))
	if err != nil {
		return game.ModerateTextOut{}, err
	}
	return decodeModerateText(body)
}

func (c *Client) BatchModerateText(ctx context.Context, in game.ModerateTextBatchIn) (game.ModerateTextBatchOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ModerateTextBatchOut{}, err
	}
	if len(in.Items) == 0 {
		return game.ModerateTextBatchOut{}, fmt.Errorf("text_items 不能为空")
	}
	body, err := c.call(ctx, "Uicproxy", "BatchModerateText", encodeBatchModerateText(in.Items))
	if err != nil {
		return game.ModerateTextBatchOut{}, err
	}
	return decodeBatchModerateText(body)
}

func (c *Client) ModeratePic(ctx context.Context, in game.ModeratePicIn) (game.ModeratePicOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ModeratePicOut{}, err
	}
	if in.URL == "" {
		return game.ModeratePicOut{}, fmt.Errorf("pic_url 不能为空")
	}
	body, err := c.call(ctx, "Uicproxy", "ModeratePic", encodeModeratePic(in))
	if err != nil {
		return game.ModeratePicOut{}, err
	}
	return decodeModeratePic(body)
}

func (c *Client) BatchModeratePic(ctx context.Context, in game.ModeratePicBatchIn) (game.ModeratePicBatchOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ModeratePicBatchOut{}, err
	}
	if len(in.Items) == 0 {
		return game.ModeratePicBatchOut{}, fmt.Errorf("pic_items 不能为空")
	}
	body, err := c.call(ctx, "Uicproxy", "BatchModeratePic", encodeBatchModeratePic(in.Items))
	if err != nil {
		return game.ModeratePicBatchOut{}, err
	}
	return decodeBatchModeratePic(body)
}
