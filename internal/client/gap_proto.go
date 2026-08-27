package client

import (
	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/pb"
)

func encodeBlockApps(block bool) []byte {
	req := pb.NewEncoder()
	req.Bool(1, block)
	return req.Bytes()
}

func encodeEncrypted(data string) []byte {
	if data == "" {
		return nil
	}
	req := pb.NewEncoder()
	req.String(1, data)
	return req.Bytes()
}

func encodeWXPage(offset, size int64) []byte {
	if size <= 0 {
		size = 20
	}
	req := pb.NewEncoder()
	req.Int(1, offset)
	req.Int(2, size)
	return req.Bytes()
}

func encodeSkinSet(ids []int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, ids)
	return req.Bytes()
}

func encodeSkinSetEffect(in game.SkinSetEffectIn) []byte {
	req := pb.NewEncoder()
	req.Int(1, in.SetID)
	req.Int(2, in.Type)
	req.Bool(3, in.Enabled)
	if in.TypeID != 0 {
		req.Int(4, in.TypeID)
	}
	if in.Param != 0 {
		inner := pb.NewEncoder()
		inner.Int(1, in.Param)
		req.Message(5, inner.Bytes())
	}
	return req.Bytes()
}

func encodeCookies(s string) []byte {
	if s == "" {
		return nil
	}
	req := pb.NewEncoder()
	req.String(1, s)
	return req.Bytes()
}

func encodeQQRecommend(in game.QQRecommendIn) []byte {
	req := pb.NewEncoder()
	if in.Class != "" {
		req.String(1, in.Class)
	}
	if in.Session != "" {
		req.String(2, in.Session)
	}
	if in.Scene != 0 {
		req.Int(3, in.Scene)
	}
	return req.Bytes()
}

func encodeCommunityID(id string) []byte {
	req := pb.NewEncoder()
	req.String(1, id)
	return req.Bytes()
}

func encodeQQRevoke(in game.QQRevokeIn) []byte {
	req := pb.NewEncoder()
	req.Int(1, in.GID)
	req.String(2, in.CommunityID)
	return req.Bytes()
}

func encodeUID(uid int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, uid)
	return req.Bytes()
}

func encodeGiftHistory(in game.GiftHistoryIn) []byte {
	req := pb.NewEncoder()
	if in.SourceType != 0 {
		req.Int(1, in.SourceType)
	}
	page := in.Page
	if page <= 0 {
		page = 1
	}
	size := in.PageSize
	if size <= 0 {
		size = 20
	}
	req.Int(2, page)
	req.Int(3, size)
	return req.Bytes()
}

func encodeOutBill(no string) []byte {
	req := pb.NewEncoder()
	req.String(1, no)
	return req.Bytes()
}

func encodeDisplay(in game.DisplayIn) []byte {
	req := pb.NewEncoder()
	if in.Name != "" {
		req.String(1, in.Name)
	}
	if in.Avatar != "" {
		req.String(2, in.Avatar)
	}
	if in.Signature != "" {
		req.String(3, in.Signature)
	}
	if in.Gender != 0 {
		req.Int(5, in.Gender)
	}
	if in.Remark != "" {
		req.String(6, in.Remark)
	}
	return req.Bytes()
}

func encodeUserSettings(in game.UserSettings) []byte {
	req := pb.NewEncoder()
	req.Bool(1, in.DisableNudge)
	req.Bool(2, in.DisableMonthCard)
	req.Bool(3, in.DisableQQSubscribe)
	req.Bool(4, in.DisableWXRecommend)
	req.Bool(5, in.DisableOfflineSummary)
	req.Bool(6, in.AllowArkVisit)
	return req.Bytes()
}

func encodeSettingsKeys(keys []int64) []byte {
	if len(keys) == 0 {
		return nil
	}
	req := pb.NewEncoder()
	req.RepeatedVarint(1, keys)
	return req.Bytes()
}

func encodeDeleteAccount(in game.DeleteAccountIn) []byte {
	req := pb.NewEncoder()
	req.String(1, in.Name)
	req.String(2, in.CertID)
	req.Int(3, in.CertType)
	return req.Bytes()
}

func encodeQQAuth(ok bool) []byte {
	req := pb.NewEncoder()
	if ok {
		req.Int(1, 1)
	} else {
		req.Int(1, 0)
	}
	return req.Bytes()
}

func encodeReportFlow(in game.ReportFlowIn) []byte {
	req := pb.NewEncoder()
	if in.OSType != 0 {
		req.Int(1, in.OSType)
	}
	if in.PlatType != 0 {
		req.Int(2, in.PlatType)
	}
	if in.OpenID != "" {
		req.String(3, in.OpenID)
	}
	if in.GID != 0 {
		req.Int(4, in.GID)
	}
	if in.Name != "" {
		req.String(5, in.Name)
	}
	if in.Now != 0 {
		req.Int(6, in.Now)
	}
	if in.Level != 0 {
		req.Int(7, in.Level)
	}
	if in.FlowType != 0 {
		req.Int(101, in.FlowType)
	}
	if in.FlowTypeStr != "" {
		req.String(102, in.FlowTypeStr)
	}
	if in.Int1 != 0 {
		req.Int(103, in.Int1)
	}
	if in.Int2 != 0 {
		req.Int(104, in.Int2)
	}
	if in.Int3 != 0 {
		req.Int(105, in.Int3)
	}
	if in.Int4 != 0 {
		req.Int(106, in.Int4)
	}
	if in.Int5 != 0 {
		req.Int(107, in.Int5)
	}
	if in.Str6 != "" {
		req.String(108, in.Str6)
	}
	if in.Str7 != "" {
		req.String(109, in.Str7)
	}
	if in.Str8 != "" {
		req.String(110, in.Str8)
	}
	if in.Str9 != "" {
		req.String(111, in.Str9)
	}
	if in.Str10 != "" {
		req.String(112, in.Str10)
	}
	return req.Bytes()
}

func encodeBatchReportFlow(flows []game.ReportFlowIn) []byte {
	req := pb.NewEncoder()
	for _, f := range flows {
		req.Message(1, encodeReportFlow(f))
	}
	return req.Bytes()
}

func encodeReportUser(in game.ReportUserIn) []byte {
	req := pb.NewEncoder()
	req.Int(1, in.GID)
	if in.Category != 0 {
		req.Int(2, in.Category)
	}
	req.RepeatedVarint(3, in.Reasons)
	if in.Scene != 0 {
		req.Int(4, in.Scene)
	}
	if in.Desc != "" {
		req.String(5, in.Desc)
	}
	if in.Content != "" {
		req.String(6, in.Content)
	}
	for _, s := range in.Pics {
		req.String(7, s)
	}
	for _, s := range in.Videos {
		req.String(8, s)
	}
	for _, s := range in.Voices {
		req.String(9, s)
	}
	if in.GroupID != "" {
		req.String(10, in.GroupID)
	}
	if in.GroupName != "" {
		req.String(11, in.GroupName)
	}
	if in.BattleID != "" {
		req.String(12, in.BattleID)
	}
	if in.BattleTime != 0 {
		req.Int(13, in.BattleTime)
	}
	if in.Entrance != 0 {
		req.Int(14, in.Entrance)
	}
	if in.MsgBoardID != "" {
		req.String(15, in.MsgBoardID)
	}
	return req.Bytes()
}

func decodeCleanFarmEvents(body []byte) (game.CleanFarmEventOut, error) {
	var out game.CleanFarmEventOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			e, err := decodeFarmSocial(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Events = append(out.Events, e)
		case 2:
			r, err := decodeFarmReward(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Rewards = append(out.Rewards, r)
		}
	}
	return out, nil
}

func decodeBlockApps(body []byte) (game.BlockAppsOut, error) {
	var out game.BlockAppsOut
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Block = pb.BoolField(m, 1)
	return out, nil
}

func decodeWXRecommend(body []byte) (game.WXRecommendOut, error) {
	var out game.WXRecommendOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			p, err := decodeWXPlayer(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Players = append(out.Players, p)
		case 2:
			out.Total = int64(f.Varint)
		case 3:
			out.HasMore = f.Varint != 0
		}
	}
	return out, nil
}

func decodeWXPlayer(raw []byte) (game.WXRecommendPlayer, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.WXRecommendPlayer{}, err
	}
	return game.WXRecommendPlayer{
		GID:     pb.IntField(m, 1),
		Name:    pb.StringField(m, 2),
		Avatar:  pb.StringField(m, 3),
		Level:   pb.IntField(m, 4),
		Applied: pb.BoolField(m, 7),
	}, nil
}

func decodeWXApply(body []byte) (game.WXApplyOut, error) {
	var out game.WXApplyOut
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
		out.Results = append(out.Results, game.WXApplyResult{
			GID:     pb.IntField(m, 1),
			Success: pb.BoolField(m, 2),
			Code:    pb.IntField(m, 3),
		})
	}
	return out, nil
}

func decodeSkinSets(body []byte) ([]game.SkinSetEffect, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.SkinSetEffect
	for _, f := range fields {
		if f.Num != 2 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.SkinSetEffect{
			SetID:  pb.IntField(m, 1),
			Type:   pb.IntField(m, 2),
			Active: pb.BoolField(m, 3),
		})
	}
	return out, nil
}

func decodeBuyPass(body []byte) (game.BuyPassOut, error) {
	var out game.BuyPassOut
	items, err := decodeItemsAt(body, 2)
	if err != nil {
		return out, err
	}
	out.Items = items
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Success = pb.BoolField(m, 1)
	if raw := m[3].Bytes; len(raw) > 0 {
		out.Pass, err = decodeBattlePass(raw)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func decodeSuccess(body []byte) (bool, error) {
	if len(body) == 0 {
		return true, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return false, err
	}
	return pb.BoolField(m, 1), nil
}

func decodeQQCommunity(raw []byte) (game.QQCommunity, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.QQCommunity{}, err
	}
	return game.QQCommunity{
		OpenID: pb.StringField(m, 1),
		Name:   pb.StringField(m, 2),
		Avatar: pb.StringField(m, 3),
	}, nil
}

func decodeQQAuthGroups(body []byte) (game.QQAuthGroupsOut, error) {
	var out game.QQAuthGroupsOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Groups = append(out.Groups, game.QQGroup{
				OpenID: pb.StringField(m, 1),
				Name:   pb.StringField(m, 2),
				Bound:  pb.IntField(m, 3),
				Avatar: pb.StringField(m, 4),
			})
		case 2:
			out.Cookies = string(f.Bytes)
		}
	}
	return out, nil
}

func decodeQQRecommend(body []byte) (game.QQRecommendOut, error) {
	var out game.QQRecommendOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Groups = append(out.Groups, game.QQRecommendGroup{
				OpenID:  pb.StringField(m, 1),
				Name:    pb.StringField(m, 2),
				Avatar:  pb.StringField(m, 3),
				Auth:    pb.StringField(m, 4),
				Jump:    pb.StringField(m, 5),
				Bound:   pb.IntField(m, 6),
				Members: pb.IntField(m, 7),
			})
		case 2:
			out.Ended = f.Varint != 0
		case 3:
			out.Pos = int64(f.Varint)
		case 4:
			out.Session = string(f.Bytes)
		}
	}
	return out, nil
}

func decodeQQBind(body []byte) (game.QQBindOut, error) {
	var out game.QQBindOut
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	if raw := m[1].Bytes; len(raw) > 0 {
		out.Community, err = decodeQQCommunity(raw)
		if err != nil {
			return out, err
		}
	}
	out.BoundAt = pb.IntField(m, 2)
	out.RewardClaimed = pb.BoolField(m, 3)
	return out, nil
}

func decodeQQLeave(body []byte) (game.QQLeaveOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.QQLeaveOut{}, err
	}
	return game.QQLeaveOut{
		QuitLeft: pb.IntField(m, 1),
		Cooldown: pb.IntField(m, 2),
	}, nil
}

func decodeQQCommunityOut(body []byte) (game.QQCommunityOut, error) {
	var out game.QQCommunityOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.Community, err = decodeQQCommunity(f.Bytes)
			if err != nil {
				return out, err
			}
		case 2:
			out.BoundAt = int64(f.Varint)
		case 3:
			out.Members = int64(f.Varint)
		case 4:
			fr, err := decodeFriend(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Friends = append(out.Friends, fr)
		case 5:
			out.HasMore = f.Varint != 0
		}
	}
	return out, nil
}

func decodeQQBindInfo(body []byte) (game.QQBindInfoOut, error) {
	var out game.QQBindInfoOut
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	if raw := m[1].Bytes; len(raw) > 0 {
		out.Community, err = decodeQQCommunity(raw)
		if err != nil {
			return out, err
		}
	}
	out.BoundAt = pb.IntField(m, 2)
	out.Cooldown = pb.IntField(m, 3)
	out.RewardClaimed = pb.BoolField(m, 4)
	out.QuitLeft = pb.IntField(m, 5)
	out.Rewards, err = decodeItemsAt(body, 6)
	if err != nil {
		return out, err
	}
	out.MaxQuit = pb.IntField(m, 7)
	out.CooldownDays = pb.IntField(m, 8)
	return out, nil
}

func decodeGiftToken(body []byte) (game.GiftTokenOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.GiftTokenOut{}, err
	}
	return game.GiftTokenOut{
		SourceType:     pb.IntField(m, 1),
		PresentOrder:   pb.StringField(m, 2),
		RedeemCode:     pb.StringField(m, 3),
		PlatformURL:    pb.StringField(m, 4),
		DisplayName:    pb.StringField(m, 5),
		PackageInfo:    pb.StringField(m, 6),
		OutBillNo:      pb.StringField(m, 7),
		TransferAmount: pb.IntField(m, 8),
		MchID:          pb.StringField(m, 9),
	}, nil
}

func decodeGiftHistory(body []byte) (game.GiftHistoryOut, error) {
	var out game.GiftHistoryOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			r, err := decodeGiftRecord(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Records = append(out.Records, r)
		case 2:
			out.Total = int64(f.Varint)
		}
	}
	return out, nil
}

func decodeGiftRecord(raw []byte) (game.GiftClaimRecord, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.GiftClaimRecord{}, err
	}
	return game.GiftClaimRecord{
		ClaimID:        pb.StringField(m, 1),
		ItemID:         pb.IntField(m, 2),
		SourceType:     pb.IntField(m, 3),
		Status:         pb.IntField(m, 4),
		Assigned:       pb.IntField(m, 5),
		Expire:         pb.IntField(m, 6),
		Title:          pb.StringField(m, 10),
		Subtitle:       pb.StringField(m, 11),
		DetailURL:      pb.StringField(m, 12),
		Code:           pb.StringField(m, 20),
		PresentOrderID: pb.StringField(m, 30),
		WXActivityID:   pb.StringField(m, 31),
	}, nil
}

func decodeTransfer(body []byte) (game.TransferOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.TransferOut{}, err
	}
	return game.TransferOut{
		State:      pb.IntField(m, 1),
		FailReason: pb.StringField(m, 2),
	}, nil
}

func decodeFollowGift(body []byte) (game.FollowGiftOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.FollowGiftOut{}, err
	}
	return game.FollowGiftOut{
		Followed: pb.BoolField(m, 1),
		Claimed:  pb.BoolField(m, 2),
		RedDot:   pb.BoolField(m, 3),
	}, nil
}

func decodeRechargeBonus(body []byte) (game.RechargeBonusOut, error) {
	var out game.RechargeBonusOut
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Active = pb.BoolField(m, 1)
	cfg := m[2].Bytes
	if len(cfg) == 0 {
		return out, nil
	}
	cm, err := pb.FieldMap(cfg)
	if err != nil {
		return out, err
	}
	out.Start = pb.IntField(cm, 2)
	out.End = pb.IntField(cm, 3)
	out.Unlock = pb.IntField(cm, 5)
	fields, err := pb.Walk(cfg)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 4 {
			continue
		}
		rm, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Ranges = append(out.Ranges, game.RechargeRange{
			Min:   pb.IntField(rm, 1),
			Max:   pb.IntField(rm, 2),
			Ratio: pb.IntField(rm, 3),
		})
	}
	return out, nil
}

func decodeRechargeData(body []byte) (game.RechargeDataOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.RechargeDataOut{}, err
	}
	return game.RechargeDataOut{
		Recharged: pb.IntField(m, 1),
		Returned:  pb.IntField(m, 2),
	}, nil
}

func decodeDisplay(body []byte) (game.DisplayOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.DisplayOut{}, err
	}
	return game.DisplayOut{
		Name:       pb.StringField(m, 1),
		Avatar:     pb.StringField(m, 2),
		Signature:  pb.StringField(m, 3),
		Gender:     pb.IntField(m, 5),
		Remark:     pb.StringField(m, 6),
		Authorized: pb.IntField(m, 8),
	}, nil
}

func decodeUserSettings(body []byte) (game.UserSettings, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.UserSettings{}, err
	}
	raw := m[1].Bytes
	if len(raw) == 0 {
		return decodeUserSettingsMsg(body)
	}
	return decodeUserSettingsMsg(raw)
}

func decodeUserSettingsMsg(raw []byte) (game.UserSettings, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.UserSettings{}, err
	}
	return game.UserSettings{
		DisableNudge:          pb.BoolField(m, 1),
		DisableMonthCard:      pb.BoolField(m, 2),
		DisableQQSubscribe:    pb.BoolField(m, 3),
		DisableWXRecommend:    pb.BoolField(m, 4),
		DisableOfflineSummary: pb.BoolField(m, 5),
		AllowArkVisit:         pb.BoolField(m, 6),
	}, nil
}

func decodeDeleteAccount(body []byte) (game.DeleteAccountOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.DeleteAccountOut{}, err
	}
	return game.DeleteAccountOut{
		Success:   pb.BoolField(m, 1),
		Msg:       pb.StringField(m, 2),
		RequestAt: pb.IntField(m, 3),
		DeleteAt:  pb.IntField(m, 4),
	}, nil
}

func decodeDecrypt(body []byte) (string, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return "", err
	}
	return pb.StringField(m, 1), nil
}

func decodeQQAuth(body []byte) (game.QQAuthOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.QQAuthOut{}, err
	}
	return game.QQAuthOut{Authorized: pb.IntField(m, 1)}, nil
}

func decodeBatchReport(body []byte) (game.BatchReportFlowOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.BatchReportFlowOut{}, err
	}
	return game.BatchReportFlowOut{
		Success: pb.IntField(m, 1),
		Fail:    pb.IntField(m, 2),
	}, nil
}

func decodeReportUser(body []byte) (game.ReportUserOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.ReportUserOut{}, err
	}
	return game.ReportUserOut{
		Success: pb.BoolField(m, 1),
		Msg:     pb.StringField(m, 2),
		TraceID: pb.StringField(m, 3),
	}, nil
}
