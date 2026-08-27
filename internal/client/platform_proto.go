package client

import (
	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/pb"
	"google.golang.org/protobuf/encoding/protowire"
)

func encodeQQVipClaimDaily(configID int64) []byte {
	if configID == 0 {
		return nil
	}
	req := pb.NewEncoder()
	req.Int(1, configID)
	return req.Bytes()
}

func encodeQQVipClaimRewards(ids []int64) []byte {
	if len(ids) == 0 {
		return nil
	}
	req := pb.NewEncoder()
	req.PackedVarints(1, ids)
	return req.Bytes()
}

func encodeSystemName(name int64) []byte {
	if name == 0 {
		return nil
	}
	req := pb.NewEncoder()
	req.Int(1, name)
	return req.Bytes()
}

func encodeWXSubscribe(items []game.WXTemplateStatus) []byte {
	req := pb.NewEncoder()
	for _, it := range items {
		inner := pb.NewEncoder()
		inner.String(1, it.TemplateID)
		inner.Bool(2, it.Subscribed)
		req.Message(1, inner.Bytes())
	}
	return req.Bytes()
}

func encodeModerateText(in game.ModerateTextIn) []byte {
	req := pb.NewEncoder()
	req.String(1, in.Text)
	req.String(2, in.Reason)
	return req.Bytes()
}

func encodeBatchModerateText(items []game.ModerateTextIn) []byte {
	req := pb.NewEncoder()
	for _, it := range items {
		req.Message(1, encodeModerateText(it))
	}
	return req.Bytes()
}

func encodeModeratePic(in game.ModeratePicIn) []byte {
	req := pb.NewEncoder()
	req.String(1, in.URL)
	req.String(2, in.Reason)
	return req.Bytes()
}

func encodeBatchModeratePic(items []game.ModeratePicIn) []byte {
	req := pb.NewEncoder()
	for _, it := range items {
		req.Message(1, encodeModeratePic(it))
	}
	return req.Bytes()
}

func decodeQQVipDaily(body []byte) (game.QQVipDailyOut, error) {
	var out game.QQVipDailyOut
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.IsQQVip = pb.BoolField(m, 1)
	out.CanClaim = pb.BoolField(m, 2)
	out.ClaimedToday = pb.BoolField(m, 3)
	rewards, err := decodeItemsAt(body, 4)
	if err != nil {
		return out, err
	}
	out.Rewards = rewards
	return out, nil
}

func decodeQQVipClaimDaily(body []byte) (game.QQVipClaimDailyOut, error) {
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return game.QQVipClaimDailyOut{}, err
	}
	return game.QQVipClaimDailyOut{Rewards: items}, nil
}

func decodeQQVipRefresh(body []byte) (game.QQVipRefreshOut, error) {
	var out game.QQVipRefreshOut
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.IsQQVip = pb.BoolField(m, 1)
	out.VIPLevel = pb.IntField(m, 2)
	return out, nil
}

func decodeQQVipConfig(raw []byte) (game.QQVipConfig, error) {
	var out game.QQVipConfig
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.Type = pb.IntField(m, 1)
	out.Enable = pb.BoolField(m, 3)
	out.SeasonID = pb.IntField(m, 4)
	out.ID = pb.IntField(m, 5)
	out.Multiplier = pb.IntField(m, 6)
	out.Start = pb.IntField(m, 7)
	out.End = pb.IntField(m, 8)
	rewards, err := decodeItemsAt(raw, 2)
	if err != nil {
		return out, err
	}
	out.Rewards = rewards
	return out, nil
}

func decodeQQVipRewardsStatus(body []byte) (game.QQVipRewardsStatusOut, error) {
	var out game.QQVipRewardsStatusOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.IsQQVip = f.Varint != 0
		case 2:
			out.CanClaim = f.Varint != 0
		case 3:
			out.ClaimedToday = f.Varint != 0
		case 4:
			out.RewardsCanClaim = f.Varint != 0
		case 5:
			cfg, err := decodeQQVipConfig(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Configs = append(out.Configs, cfg)
		case 6:
			out.HasRedpoint = f.Varint != 0
		}
	}
	return out, nil
}

func decodeQQVipClaimRewards(body []byte) (game.QQVipClaimRewardsOut, error) {
	var out game.QQVipClaimRewardsOut
	if len(body) == 0 {
		return out, nil
	}
	out.SkinIDs = pb.PackedInts(body, 1)
	out.FrameIDs = pb.PackedInts(body, 2)
	rewards, err := decodeItemsAt(body, 3)
	if err != nil {
		return out, err
	}
	out.Rewards = rewards
	return out, nil
}

func decodeMarquee(body []byte) (game.MarqueeOut, error) {
	var out game.MarqueeOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Msgs = append(out.Msgs, game.MarqueeMsg{
			UUID:         pb.IntField(m, 1),
			ConfigID:     pb.IntField(m, 2),
			Expire:       pb.IntField(m, 3),
			Type:         pb.IntField(m, 4),
			Content:      pb.StringField(m, 5),
			Priority:     pb.IntField(m, 6),
			DisplayCount: pb.IntField(m, 7),
		})
	}
	return out, nil
}

func decodeSystemOpen(body []byte) (game.SystemOpenOut, error) {
	var out game.SystemOpenOut
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Unlocked = pb.BoolField(m, 1)
	return out, nil
}

func decodeMutantOpenInfo(body []byte) (game.MutantOpenInfoOut, error) {
	var out game.MutantOpenInfoOut
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Tips = pb.StringField(m, 1)
	rewards, err := decodeItemsAt(body, 2)
	if err != nil {
		return out, err
	}
	out.Rewards = rewards
	return out, nil
}

func decodeQQSubscribe(body []byte) (game.QQSubscribeOut, error) {
	var out game.QQSubscribeOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		if f.Kind == protowire.VarintType {
			out.Status = int64(f.Varint)
			out.Subscribed = f.Varint != 0
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		item := game.QQSubscribeItem{ID: pb.IntField(m, 1), Subscribed: pb.BoolField(m, 2)}
		out.Items = append(out.Items, item)
		if item.Subscribed {
			out.Subscribed = true
		}
	}
	return out, nil
}

func decodeWXSubscribe(body []byte) (game.WXSubscribeOut, error) {
	var out game.WXSubscribeOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Templates = append(out.Templates, game.WXTemplateStatus{
			TemplateID: pb.StringField(m, 1),
			Subscribed: pb.BoolField(m, 2),
		})
	}
	return out, nil
}

func decodeModerateTextResult(raw []byte) (game.ModerateTextOut, error) {
	var out game.ModerateTextOut
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.Text = pb.StringField(m, 2)
	out.Dirty = pb.BoolField(m, 3)
	out.Reason = pb.StringField(m, 4)
	return out, nil
}

func decodeModerateText(body []byte) (game.ModerateTextOut, error) {
	var out game.ModerateTextOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	var replyReason string
	for _, f := range fields {
		switch f.Num {
		case 1:
			got, err := decodeModerateTextResult(f.Bytes)
			if err != nil {
				return out, err
			}
			out = got
		case 2:
			replyReason = string(f.Bytes)
		}
	}
	if out.Reason == "" {
		out.Reason = replyReason
	}
	return out, nil
}

func decodeBatchModerateText(body []byte) (game.ModerateTextBatchOut, error) {
	var out game.ModerateTextBatchOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		item, err := decodeModerateTextResult(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Items = append(out.Items, item)
	}
	return out, nil
}

func decodeModeratePicResult(raw []byte) (game.ModeratePicOut, error) {
	var out game.ModeratePicOut
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.URL = pb.StringField(m, 1)
	out.Dirty = pb.BoolField(m, 2)
	out.DirtyType = pb.IntField(m, 3)
	out.Reason = pb.StringField(m, 4)
	return out, nil
}

func decodeModeratePic(body []byte) (game.ModeratePicOut, error) {
	var out game.ModeratePicOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	var replyReason string
	for _, f := range fields {
		switch f.Num {
		case 1:
			got, err := decodeModeratePicResult(f.Bytes)
			if err != nil {
				return out, err
			}
			out = got
		case 2:
			replyReason = string(f.Bytes)
		}
	}
	if out.Reason == "" {
		out.Reason = replyReason
	}
	return out, nil
}

func decodeBatchModeratePic(body []byte) (game.ModeratePicBatchOut, error) {
	var out game.ModeratePicBatchOut
	if len(body) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		item, err := decodeModeratePicResult(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Items = append(out.Items, item)
	}
	return out, nil
}
