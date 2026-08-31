package client

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/pb"
	"google.golang.org/protobuf/encoding/protowire"
)

const (
	fallbackGameVersion = "1.13.3.16_20260826"
	channelID           = "1256"
)

func gameVersion() string {
	if v := os.Getenv("FARM_GAME_VER"); v != "" {
		return v
	}
	return fallbackGameVersion
}

func encodeLogin() []byte {
	dev := pb.NewEncoder()
	dev.String(1, gameVersion())
	dev.String(2, "Windows Unknown x64")
	dev.String(5, "wifi")
	dev.Int(10, 98089)
	dev.String(13, "microsoft")
	report := pb.NewEncoder()
	report.String(5, "other")
	report.Int(6, 2)
	req := pb.NewEncoder()
	req.Message(5, dev.Bytes())
	req.String(7, channelID)
	req.Message(8, report.Bytes())
	return req.Bytes()
}

func encodeAllLands(hostGID int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	return req.Bytes()
}

func visitReason(selfGID, hostGID int64) int64 {
	if hostGID > 0 && hostGID != selfGID {
		return 2
	}
	return 0
}

func encodeHarvest(landIDs []int64, hostGID int64, isAll bool, reason int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, hostGID)
	req.Bool(3, isAll)
	_ = reason
	return req.Bytes()
}

func encodePlant(seedID int64, landIDs []int64) []byte {
	item := pb.NewEncoder()
	item.Int(1, seedID)
	item.RepeatedVarint(2, landIDs)
	req := pb.NewEncoder()
	req.Message(2, item.Bytes())
	return req.Bytes()
}

func encodeEnter(hostGID, reason int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	req.Int(2, reason)
	return req.Bytes()
}

func encodeGIDs(gids []int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, gids)
	return req.Bytes()
}

func encodeShareFlag() []byte {
	req := pb.NewEncoder()
	req.Bool(1, true)
	return req.Bytes()
}

func decodeCanShare(body []byte) (bool, error) {
	if len(body) == 0 {
		return false, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return false, err
	}
	return pb.BoolField(m, 1), nil
}

func encodeWater(landIDs []int64, hostGID int64, reason int64) []byte {
	return encodeLandOp(landIDs, hostGID, reason)
}

func encodeLandOp(landIDs []int64, hostGID int64, reason int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, hostGID)
	_ = reason
	return req.Bytes()
}

func encodeFertilize(landIDs []int64, fertilizerID int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, fertilizerID)
	return req.Bytes()
}

func encodeAnti(data []byte) []byte {
	req := pb.NewEncoder()
	req.BytesField(1, data)
	return req.Bytes()
}

func encodeShopInfo(shopID int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, shopID)
	return req.Bytes()
}

func encodeBuy(goodsID, num, price int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, goodsID)
	req.Int(2, num)
	req.Int(3, price)
	return req.Bytes()
}

func encodeCoreItem(it game.BagItem) []byte {
	item := pb.NewEncoder()
	item.Int(1, it.ID)
	item.Int(2, it.Count)
	item.Int(6, it.UID)
	return item.Bytes()
}

func encodeSell(items []game.BagItem) []byte {
	req := pb.NewEncoder()
	for _, it := range items {
		req.Message(1, encodeCoreItem(it))
	}
	return req.Bytes()
}

func encodeUse(in game.UseIn) []byte {
	req := pb.NewEncoder()
	req.Message(1, encodeCoreItem(game.BagItem{ID: in.ID, Count: in.Count, UID: in.UID}))
	if in.HostGID > 0 || len(in.LandIDs) > 0 {
		req.Message(2, encodeUseTarget(in.HostGID, in.LandIDs))
	}
	return req.Bytes()
}

func encodeUseTarget(hostGID int64, landIDs []int64) []byte {
	// v59 UseTarget: 1 host_gid, 2 land_ids, 3 use_config_id。官方青蛙写 0，省略即可。
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	req.RepeatedVarint(2, landIDs)
	return req.Bytes()
}

func encodeLeave(hostGID int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	return req.Bytes()
}

func encodeRemove(landIDs []int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	return req.Bytes()
}

func encodeLandID(landID int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, landID)
	return req.Bytes()
}

func encodePage(from, count, fallback int64) []byte {
	if count <= 0 {
		count = fallback
	}
	req := pb.NewEncoder()
	req.Int(1, from)
	req.Int(2, count)
	return req.Bytes()
}

func encodeFeed(foodID, count int64) []byte {
	if count <= 0 {
		count = 1
	}
	req := pb.NewEncoder()
	req.Int(1, foodID)
	req.Int(2, count)
	return req.Bytes()
}

func encodeRank(typ, page int64) []byte {
	if typ <= 0 {
		typ = 1
	}
	if page <= 0 {
		page = 1
	}
	req := pb.NewEncoder()
	req.Int(1, typ)
	req.Int(2, page)
	return req.Bytes()
}

func encodeAvatarOwned(typ int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, typ)
	return req.Bytes()
}

func encodeAvatarEquip(id int64, off bool) []byte {
	req := pb.NewEncoder()
	req.Int(1, id)
	if off {
		req.Varint(2, 0)
	} else {
		req.Bool(2, true)
	}
	return req.Bytes()
}

func encodeSkinEquip(current, next int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, current)
	req.Varint(2, uint64(next))
	return req.Bytes()
}

func encodeAchieve(kind, id int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, kind)
	req.Int(2, id)
	return req.Bytes()
}

func encodeAchieveGoal(kind, id, goalID int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, goalID)
	req.Int(2, kind)
	req.Int(3, id)
	return req.Bytes()
}

func encodeSetTags() []byte {
	return nil
}

func encodeFarming(landIDs []int64, hostGID int64, itemIDs []int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, hostGID)
	req.RepeatedVarint(5, itemIDs)
	return req.Bytes()
}

func encodeRefreshLands(landIDs []int64, hostGID int64, reason int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	req.RepeatedVarint(2, landIDs)
	req.Int(4, reason)
	return req.Bytes()
}

func encodeCleanSocial(landIDs, itemIDs []int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.RepeatedVarint(2, itemIDs)
	return req.Bytes()
}

func encodeBatchUse(items []game.UseIn) []byte {
	req := pb.NewEncoder()
	for _, it := range items {
		req.Message(1, encodeCoreItem(game.BagItem{ID: it.ID, Count: it.Count, UID: it.UID}))
	}
	return req.Bytes()
}

func encodeHeartbeat(hostGID int64, ver string) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	req.String(2, ver)
	return req.Bytes()
}

func encodeGroupID(id int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, id)
	return req.Bytes()
}

func encodePutSocial(hostGID int64, landIDs []int64, reason int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	if len(landIDs) == 1 {
		req.Int(2, landIDs[0])
	} else {
		req.RepeatedVarint(2, landIDs)
	}
	_ = reason
	return req.Bytes()
}

func encodeReportInvite(openID, key string) []byte {
	req := pb.NewEncoder()
	req.String(1, openID)
	req.String(2, key)
	return req.Bytes()
}

func encodeDogBuy(id, price int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, id)
	req.Int(2, price)
	return req.Bytes()
}

func encodeOpenIDs(field protowire.Number, ids []string) []byte {
	req := pb.NewEncoder()
	for _, id := range ids {
		req.String(field, id)
	}
	return req.Bytes()
}

func encodeGIDsPacked(gids []int64) []byte {
	req := pb.NewEncoder()
	req.PackedVarints(1, gids)
	return req.Bytes()
}

func encodeArkClick(in game.ArkIn) []byte {
	req := pb.NewEncoder()
	req.Int(1, in.GID)
	req.String(2, in.OpenID)
	scene := in.Scene
	if scene == "" {
		scene = channelID
	}
	req.String(3, scene)
	req.Int(4, in.ShareID)
	req.String(5, in.Key)
	return req.Bytes()
}

func encodeLogout(reason string) []byte {
	req := pb.NewEncoder()
	req.String(1, reason)
	return req.Bytes()
}

func encodeUIDs(uids []int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, uids)
	return req.Bytes()
}

func encodeIllustratedLevels(typ int64) []byte {
	if typ <= 0 {
		typ = 1
	}
	req := pb.NewEncoder()
	req.Int(1, typ)
	return req.Bytes()
}

func encodeCheckCanOperate(hostGID, operationID int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
	req.Int(2, operationID)
	req.Varint(3, 0)
	return req.Bytes()
}

func encodeClaimTask(id int64) []byte {
	return encodeClaimTaskShared(id, false)
}

func encodeClaimTaskShared(id int64, shared bool) []byte {
	req := pb.NewEncoder()
	req.Int(1, id)
	req.Bool(2, shared)
	return req.Bytes()
}

func encodeBatchClaim(ids []int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, ids)
	return req.Bytes()
}

func encodeClaimDaily(typ int64, ids []int64) []byte {
	_ = typ
	req := pb.NewEncoder()
	req.RepeatedVarint(1, ids)
	return req.Bytes()
}

func encodeBox(box int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, box)
	return req.Bytes()
}

func encodeEmailRef(box int64, id string) []byte {
	req := pb.NewEncoder()
	req.Int(1, box)
	req.String(2, id)
	return req.Bytes()
}

func encodeSeasonInfo() []byte {
	req := pb.NewEncoder()
	req.Bool(1, true)
	return req.Bytes()
}

func encodeMallList(slot int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, slot)
	return req.Bytes()
}

func encodeOperate(id, cmd int64, field int, inner []byte) []byte {
	req := pb.NewEncoder()
	req.Int(1, id)
	req.Int(2, cmd)
	if field > 0 {
		req.MessageAlways(protowire.Number(field), inner)
	}
	return req.Bytes()
}

func encodeShopBuy(id, goodsID, count int64) []byte {
	if count <= 0 {
		count = 1
	}
	inner := pb.NewEncoder()
	inner.Int(1, goodsID)
	inner.Int(2, count)
	return encodeOperate(id, 1, 101, inner.Bytes())
}

func encodeRandBuy(id, goodsID, count int64) []byte {
	if count <= 0 {
		count = 1
	}
	inner := pb.NewEncoder()
	inner.Int(1, goodsID)
	inner.Int(2, count)
	return encodeOperate(id, 2, 102, inner.Bytes())
}

func encodeRandRefresh(id int64) []byte {
	return encodeOperate(id, 3, 0, nil)
}

func encodeShopBatch(id int64, items []game.ShopBuyItem) []byte {
	inner := pb.NewEncoder()
	for _, it := range items {
		sel := pb.NewEncoder()
		sel.Int(1, it.GoodsID)
		count := it.Count
		if count <= 0 {
			count = 1
		}
		sel.Int(2, count)
		inner.Message(1, sel.Bytes())
	}
	return encodeOperate(id, 34, 133, inner.Bytes())
}

func encodeMegaClaim(id int64) []byte {
	return encodeOperate(id, 21, 119, nil)
}

func encodeTechSubmit(id, nodeID int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, nodeID)
	return encodeOperate(id, 40, 140, inner.Bytes())
}

func encodeLottery(id, hostGID, free, paid int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, free)
	inner.Int(2, paid)
	inner.Int(3, hostGID)
	return encodeOperate(id, 9, 107, inner.Bytes())
}

func encodeBrewStart(id int64, items []game.BrewItem) []byte {
	inner := pb.NewEncoder()
	for _, it := range items {
		sel := pb.NewEncoder()
		sel.Int(1, it.UID)
		sel.Int(2, it.Count)
		inner.Message(1, sel.Bytes())
	}
	return encodeOperate(id, 14, 112, inner.Bytes())
}

func encodeBrewStep(id int64) []byte {
	return encodeOperate(id, 15, 113, nil)
}

func encodeBrewClaim(id, claimType int64) []byte {
	if claimType <= 0 {
		claimType = 1
	}
	inner := pb.NewEncoder()
	inner.Int(1, claimType)
	return encodeOperate(id, 16, 114, inner.Bytes())
}

func encodeRecallClaim(id int64) []byte {
	return encodeOperate(id, 19, 117, nil)
}

func encodeReturnGift(id int64) []byte {
	return encodeOperate(id, 20, 118, nil)
}

func encodeInviteClaim(id, rewardType int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, rewardType)
	return encodeOperate(id, 22, 120, inner.Bytes())
}

func encodeNewcomerClaim(id int64) []byte {
	return encodeOperate(id, 23, 121, nil)
}

func encodeInvitees(id int64) []byte {
	return encodeOperate(id, 24, 122, nil)
}

func encodeSendGift(id, gid, msgTextID int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, gid)
	inner.Int(2, msgTextID)
	return encodeOperate(id, 26, 124, inner.Bytes())
}

func encodeSignin(id, rewardID int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, rewardID)
	return encodeOperate(id, 4, 103, inner.Bytes())
}

func encodeProgress(id, step int64) []byte {
	inner := pb.NewEncoder()
	inner.Varint(1, uint64(step))
	return encodeOperate(id, 25, 125, inner.Bytes())
}

func encodeDraw(id, count int64) []byte {
	if count <= 0 {
		count = 1
	}
	inner := pb.NewEncoder()
	inner.Int(1, count)
	return encodeOperate(id, 5, 104, inner.Bytes())
}

func encodeDrawHistory(id int64) []byte {
	return encodeOperate(id, 6, 105, nil)
}

func encodeMarkViewed(id int64) []byte {
	return encodeOperate(id, 7, 0, nil)
}

func encodeRandBatch(id int64, items []game.ShopBuyItem) []byte {
	inner := pb.NewEncoder()
	for _, it := range items {
		sel := pb.NewEncoder()
		sel.Int(1, it.GoodsID)
		count := it.Count
		if count <= 0 {
			count = 1
		}
		sel.Int(2, count)
		inner.Message(1, sel.Bytes())
	}
	return encodeOperate(id, 8, 106, inner.Bytes())
}

func encodeLotteryHistory(id int64) []byte {
	return encodeOperate(id, 10, 108, nil)
}

func encodeCheerJoin(id, campID int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, campID)
	return encodeOperate(id, 11, 109, inner.Bytes())
}

func encodeCheerSubmit(id, count int64) []byte {
	if count <= 0 {
		count = 1
	}
	inner := pb.NewEncoder()
	inner.Int(1, count)
	return encodeOperate(id, 12, 110, inner.Bytes())
}

func encodeCheerClaim(id, tier int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, tier)
	return encodeOperate(id, 13, 111, inner.Bytes())
}

func encodeRecallable(id int64) []byte {
	return encodeOperate(id, 17, 115, nil)
}

func encodeRecalled(id int64) []byte {
	return encodeOperate(id, 18, 116, nil)
}

func encodeCharityShare(id int64) []byte {
	return encodeOperate(id, 35, 134, nil)
}

func encodeCharityDonate(id int64) []byte {
	return encodeOperate(id, 36, 135, nil)
}

func encodeCharityClaim(id, score int64) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, score)
	return encodeOperate(id, 37, 136, inner.Bytes())
}

func encodeCharityXhh(id int64) []byte {
	return encodeOperate(id, 38, 137, nil)
}

func encodeCharityAgree(id int64, agreed bool) []byte {
	inner := pb.NewEncoder()
	if agreed {
		inner.Varint(1, 1)
	} else {
		inner.Varint(1, 0)
	}
	return encodeOperate(id, 39, 138, inner.Bytes())
}

// 宠物寻宝。活动号 2026090101，仅活动号，开停看 List 带 hunt 那条的 start / end。中间没有 912。
const (
	cmdHuntFinishCG     = 901
	cmdHuntGuide        = 902
	cmdHuntFeed         = 903
	cmdHuntDraw         = 904
	cmdHuntLog          = 905
	cmdHuntClaimStory   = 906
	cmdHuntClaimSeed    = 907
	cmdHuntRefreshCharm = 908
	cmdHuntEquip        = 909
	cmdHuntBattle       = 910
	cmdHuntPlunderedLog = 911
	cmdHuntOpen         = 913
	cmdHuntEscort       = 914
	cmdHuntCompensate   = 915
	cmdHuntFriendInfo   = 916
)

func encodeHunt(id, cmd int64) []byte {
	return encodeOperate(id, cmd, 0, nil)
}

func encodeHuntEquip(id int64, charmIDs []int64) []byte {
	inner := pb.NewEncoder()
	for _, charmID := range charmIDs {
		inner.Int(1, charmID)
	}
	return encodeOperate(id, cmdHuntEquip, cmdHuntEquip, inner.Bytes())
}

func encodeHuntBattle(id, defenderGID int64, treasureID string) []byte {
	inner := pb.NewEncoder()
	inner.Int(1, defenderGID)
	if treasureID != "" {
		inner.String(2, treasureID)
	}
	return encodeOperate(id, cmdHuntBattle, cmdHuntBattle, inner.Bytes())
}

func encodeIllustratedList(rarity, typ int64) []byte {
	req := pb.NewEncoder()
	req.Varint(1, uint64(rarity))
	if typ <= 0 {
		typ = 1
	}
	req.Int(2, typ)
	return req.Bytes()
}

func encodeMallBuy(id, num int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, id)
	req.Int(2, num)
	return req.Bytes()
}

func encodeEmailIDs(box int64, ids []string) []byte {
	req := pb.NewEncoder()
	req.Int(1, box)
	for _, id := range ids {
		req.String(2, id)
	}
	return req.Bytes()
}

func encodeTaskReport(id, progress int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, id)
	req.Int(2, progress)
	return req.Bytes()
}

func encodeCareerGID() []byte {
	return nil
}

func encodeVisitPage(page int64) []byte {
	if page <= 0 {
		page = 1
	}
	req := pb.NewEncoder()
	req.Int(1, page)
	return req.Bytes()
}

func decodeAnti(body []byte) ([]byte, error) {
	if len(body) == 0 {
		return nil, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return nil, err
	}
	return append([]byte{}, m[1].Bytes...), nil
}

func decodeUser(body []byte) (game.User, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return game.User{}, err
	}
	if len(top[1].Bytes) == 0 {
		return game.User{}, fmt.Errorf("登录回包没有用户资料")
	}
	return decodeBasic(top[1].Bytes)
}

func decodeBasic(raw []byte) (game.User, error) {
	var u game.User
	basic, err := pb.FieldMap(raw)
	if err != nil {
		return u, err
	}
	u.GID = pb.IntField(basic, 1)
	u.Name = pb.StringField(basic, 2)
	u.Level = pb.IntField(basic, 3)
	u.Exp = pb.IntField(basic, 4)
	u.Gold = pb.IntField(basic, 5)
	u.OpenID = pb.StringField(basic, 6)
	u.Avatar = pb.StringField(basic, 7)
	u.Remark = pb.StringField(basic, 8)
	u.Signature = pb.StringField(basic, 9)
	u.Gender = pb.IntField(basic, 10)
	u.Authorized = pb.IntField(basic, 13)
	u.LastOnline = pb.IntField(basic, 17)
	return u, nil
}

func decodeLoginExtra(body []byte) (bool, game.Weather) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return false, game.Weather{}
	}
	var weather game.Weather
	if raw := m[18].Bytes; len(raw) > 0 {
		weather, _ = decodeWeather(raw)
	}
	return false, weather
}

func decodeLands(body []byte) ([]game.Land, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var lands []game.Land
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		land, err := decodeLand(f.Bytes)
		if err != nil {
			return nil, err
		}
		lands = append(lands, land)
	}
	return lands, nil
}

type landOpSpec struct {
	drops   int32
	cost    int32
	social  int32
	events  int32
	rewards int32
}

func decodeLandOp(body []byte) (game.LandOpOut, error) {
	return walkLandOp(body, landOpSpec{drops: 3})
}

func decodeFarmingOut(body []byte) (game.LandOpOut, error) {
	return walkLandOp(body, landOpSpec{drops: 3, rewards: 4})
}

func decodeFertilizeOut(body []byte) (game.LandOpOut, error) {
	return walkLandOp(body, landOpSpec{drops: 4, cost: 3})
}

func decodeSocialOp(body []byte) (game.LandOpOut, error) {
	return walkLandOp(body, landOpSpec{drops: 3, social: 4})
}

func decodeUseLandOp(body []byte) (game.LandOpOut, error) {
	var out game.LandOpOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch int32(f.Num) {
		case 4:
			land, err := decodeLand(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Lands = append(out.Lands, land)
		case 5:
			d, err := decodeLandDrop(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Drops = append(out.Drops, d)
		case 6:
			r, err := decodeFarmReward(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Rewards = append(out.Rewards, r)
		}
	}
	return out, nil
}

func decodeRefreshOut(body []byte) (game.LandOpOut, error) {
	return walkLandOp(body, landOpSpec{events: 3})
}

func walkLandOp(body []byte, spec landOpSpec) (game.LandOpOut, error) {
	var out game.LandOpOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		n := int32(f.Num)
		switch {
		case n == 1:
			land, err := decodeLand(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Lands = append(out.Lands, land)
		case n == 2:
			lim, err := decodeOpLimit(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Limits = append(out.Limits, lim)
		case spec.drops != 0 && n == spec.drops:
			d, err := decodeLandDrop(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Drops = append(out.Drops, d)
		case spec.cost != 0 && n == spec.cost:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Costs = append(out.Costs, it)
		case spec.social != 0 && n == spec.social:
			s, err := decodeLandSocial(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Social = append(out.Social, s)
		case spec.events != 0 && n == spec.events:
			e, err := decodeFarmSocial(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Events = append(out.Events, e)
		case spec.rewards != 0 && n == spec.rewards:
			r, err := decodeFarmReward(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Rewards = append(out.Rewards, r)
		}
	}
	return out, nil
}

func decodeOpLimit(raw []byte) (game.OpLimit, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.OpLimit{}, err
	}
	return game.OpLimit{
		ID:         pb.IntField(m, 1),
		Used:       pb.IntField(m, 2),
		Limit:      pb.IntField(m, 3),
		ShareID:    pb.IntField(m, 4),
		ExpUsed:    pb.IntField(m, 5),
		ExpLimit:   pb.IntField(m, 6),
		ExpShareID: pb.IntField(m, 7),
	}, nil
}

func decodeLandDrop(raw []byte) (game.LandDrop, error) {
	var d game.LandDrop
	fields, err := pb.Walk(raw)
	if err != nil {
		return d, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			d.LandID = int64(f.Varint)
		case 2:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return d, err
			}
			d.Rewards = append(d.Rewards, it)
		case 3:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return d, err
			}
			d.Costs = append(d.Costs, it)
		}
	}
	d.Skills = pb.PackedInts(raw, 4)
	return d, nil
}

func decodeFarmSocial(raw []byte) (game.FarmSocial, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.FarmSocial{}, err
	}
	return game.FarmSocial{
		ItemID:  pb.IntField(m, 1),
		OwnerID: pb.IntField(m, 2),
		PutTime: pb.IntField(m, 3),
	}, nil
}

func decodeFarmReward(raw []byte) (game.FarmReward, error) {
	var r game.FarmReward
	fields, err := pb.Walk(raw)
	if err != nil {
		return r, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			r.SourceItemID = int64(f.Varint)
		case 2:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return r, err
			}
			r.Items = append(r.Items, it)
		}
	}
	return r, nil
}

func decodeLandSocial(raw []byte) (game.LandSocial, error) {
	var s game.LandSocial
	fields, err := pb.Walk(raw)
	if err != nil {
		return s, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			s.LandID = int64(f.Varint)
		case 2:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return s, err
			}
			s.Items = append(s.Items, it)
		}
	}
	return s, nil
}

func decodeItem(raw []byte) (game.Item, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.Item{}, err
	}
	return game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)}, nil
}

func decodeItemAt(body []byte, field int32) (game.Item, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return game.Item{}, err
	}
	for _, f := range fields {
		if int32(f.Num) != field {
			continue
		}
		return decodeItem(f.Bytes)
	}
	return game.Item{}, nil
}

func decodeItemOp(body []byte, usedField, gotField, compField int32, usedSingle bool) (game.ItemOpOut, error) {
	var out game.ItemOpOut
	if usedField != 0 {
		if usedSingle {
			it, err := decodeItemAt(body, usedField)
			if err != nil {
				return out, err
			}
			if it.ID != 0 || it.Count != 0 {
				out.Used = []game.Item{it}
			}
		} else {
			used, err := decodeItemsAt(body, usedField)
			if err != nil {
				return out, err
			}
			out.Used = used
		}
	}
	if gotField != 0 {
		got, err := decodeItemsAt(body, gotField)
		if err != nil {
			return out, err
		}
		out.Items = got
	}
	if compField != 0 {
		comp, err := decodeItemsAt(body, compField)
		if err != nil {
			return out, err
		}
		out.Compensated = comp
	}
	return out, nil
}

func decodeLand(raw []byte) (game.Land, error) {
	var land game.Land
	m, err := pb.FieldMap(raw)
	if err != nil {
		return land, err
	}
	land.ID = pb.IntField(m, 1)
	land.Unlocked = pb.BoolField(m, 2)
	land.Level = pb.IntField(m, 3)
	land.MaxLevel = pb.IntField(m, 4)
	land.CouldUnlock = pb.BoolField(m, 5)
	land.Shared = pb.BoolField(m, 11)
	land.CanShare = pb.BoolField(m, 12)
	land.MasterLandID = pb.IntField(m, 13)
	land.SlaveLandIDs = pb.PackedInts(raw, 14)
	land.LandSize = pb.IntField(m, 15)
	land.LandsLevel = pb.IntField(m, 16)
	if need := m[6].Bytes; len(need) > 0 {
		n, err := decodeLandNeed(need, true)
		if err != nil {
			return land, err
		}
		land.Unlock = &n
	}
	if need := m[8].Bytes; len(need) > 0 {
		n, err := decodeLandNeed(need, false)
		if err != nil {
			return land, err
		}
		land.Upgrade = &n
	}
	if plantRaw := m[10].Bytes; len(plantRaw) > 0 {
		p, err := pb.FieldMap(plantRaw)
		if err != nil {
			return land, err
		}
		land.PlantID = pb.IntField(p, 1)
		land.PlantName = pb.StringField(p, 2)
		land.Season = pb.IntField(p, 5)
		land.DryNum = pb.IntField(p, 6)
		land.StoleNum = pb.IntField(p, 9)
		land.FruitID = pb.IntField(p, 10)
		land.FruitNum = pb.IntField(p, 11)
		weeds := pb.PackedInts(plantRaw, 12)
		insects := pb.PackedInts(plantRaw, 13)
		land.HasWeed = len(weeds) > 0
		land.HasInsect = len(insects) > 0
		land.Stealers = pb.PackedInts(plantRaw, 14)
		land.GrowSec = pb.IntField(p, 15)
		land.Stealable = pb.BoolField(p, 16)
		land.FertLeft = pb.IntField(p, 17)
		land.LeftFruit = pb.IntField(p, 18)
		land.Intimacy = pb.IntField(p, 19)
		land.MutantIDs = pb.PackedInts(plantRaw, 20)
		land.Mutant = len(land.MutantIDs) > 0
		land.PlantExp = pb.IntField(p, 22)
		land.LastFruit = pb.IntField(p, 24)
		land.LastFruitNum = pb.IntField(p, 25)
		land.SourceFruit = pb.IntField(p, 26)
		land.SourceFruitNum = pb.IntField(p, 27)
		land.SourcePriceID = pb.IntField(p, 28)
		land.SourcePrice = pb.IntField(p, 29)
		land.MutantPriceID = pb.IntField(p, 30)
		land.MutantPrice = pb.IntField(p, 31)
		land.Phase = currentPhase(plantRaw)
		items, err := decodeCropItems(plantRaw, 35, true)
		if err != nil {
			return land, err
		}
		crop, err := decodeCropItems(plantRaw, 38, false)
		if err != nil {
			return land, err
		}
		land.SocialItems = append(items, crop...)
	}
	return land, nil
}

func decodeLandNeed(raw []byte, unlock bool) (game.LandNeed, error) {
	var n game.LandNeed
	m, err := pb.FieldMap(raw)
	if err != nil {
		return n, err
	}
	if unlock {
		n.LandID = pb.IntField(m, 1)
		n.Level = pb.IntField(m, 2)
		n.Items, err = decodeItemsAt(raw, 3)
	} else {
		n.Lands = pb.IntField(m, 1)
		n.Items, err = decodeItemsAt(raw, 2)
	}
	return n, err
}

func decodeLandBuff(raw []byte) (game.LandBuff, error) {
	var b game.LandBuff
	m, err := pb.FieldMap(raw)
	if err != nil {
		return b, err
	}
	b.Yield = pb.IntField(m, 1)
	b.Time = pb.IntField(m, 2)
	b.Exp = pb.IntField(m, 3)
	b.Mutant = pb.IntField(m, 4)
	b.Pass = pb.IntField(m, 5)
	return b, nil
}

func decodeCropItems(raw []byte, num protowire.Number, withType bool) ([]game.CropItem, error) {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var out []game.CropItem
	for _, f := range fields {
		if f.Num != num {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		it := game.CropItem{
			ItemID: pb.IntField(m, 1),
		}
		if withType {
			it.Count = pb.IntField(m, 2)
			it.Type = pb.IntField(m, 3)
			it.PutTime = pb.IntField(m, 4)
			it.End = pb.IntField(m, 5)
		} else {
			it.Owner = pb.IntField(m, 2)
			it.PutTime = pb.IntField(m, 3)
			it.LandID = pb.IntField(m, 4)
		}
		out = append(out, it)
	}
	return out, nil
}

func currentPhase(plantRaw []byte) int64 {
	fields, err := pb.Walk(plantRaw)
	if err != nil {
		return 0
	}
	var phase, latest int64
	for _, f := range fields {
		if f.Num != 4 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			continue
		}
		begin := pb.IntField(m, 2)
		if begin >= latest {
			latest = begin
			phase = pb.IntField(m, 1)
		}
	}
	return phase
}

func decodeItems(body []byte) ([]game.Item, error) {
	return decodeItemsAt(body, 2)
}

func decodeShareClaim(body []byte) ([]game.Item, error) {
	out, err := decodeShareClaimOut(body)
	if err != nil {
		return nil, err
	}
	return out.Items, nil
}

func decodeShareClaimOut(body []byte) (game.ShareClaimOut, error) {
	var out game.ShareClaimOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.Success = f.Varint != 0
		case 2:
			out.HasReward = f.Varint != 0
		case 3:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Items = append(out.Items, it)
		}
	}
	return out, nil
}

func decodeItemsAt(body []byte, field int32) ([]game.Item, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var items []game.Item
	for _, f := range fields {
		if int32(f.Num) != field {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		items = append(items, game.Item{
			ID:    pb.IntField(m, 1),
			Count: pb.IntField(m, 2),
		})
	}
	return items, nil
}

func decodeBagReply(body []byte) ([]game.BagItem, error) {
	out, err := decodeBagOut(body)
	if err != nil {
		return nil, err
	}
	return out.Items, nil
}

func decodeBagOut(body []byte) (game.BagOut, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return game.BagOut{}, err
	}
	return decodeItemBagOut(top[1].Bytes)
}

func decodeLoginBag(body []byte) ([]game.BagItem, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return nil, err
	}
	return decodeItemBag(top[2].Bytes)
}

func decodeItemBag(raw []byte) ([]game.BagItem, error) {
	out, err := decodeItemBagOut(raw)
	if err != nil {
		return nil, err
	}
	return out.Items, nil
}

func decodeItemBagOut(raw []byte) (game.BagOut, error) {
	var out game.BagOut
	if len(raw) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			item, err := decodeBagItem(f.Bytes)
			if err != nil {
				return out, err
			}
			if item.ID != 0 {
				out.Items = append(out.Items, item)
			}
		}
	}
	return out, nil
}

func decodeBagItem(raw []byte) (game.BagItem, error) {
	var it game.BagItem
	m, err := pb.FieldMap(raw)
	if err != nil {
		return it, err
	}
	it.ID = pb.IntField(m, 1)
	it.Count = pb.IntField(m, 2)
	it.UID = pb.IntField(m, 6)
	it.Mutants = append(pb.PackedInts(raw, 7), pb.PackedInts(raw, 8)...)
	if show := m[100].Bytes; len(show) > 0 {
		it.SellID, it.SellPrice = parseSellPrice(show)
	}
	return it, nil
}

func parseSellPrice(show []byte) (int64, int64) {
	m, err := pb.FieldMap(show)
	if err != nil {
		return 0, 0
	}
	f1 := pb.IntField(m, 1)
	if f1 > 0 && f1 != 1001 {
		return 0, f1
	}
	sells, err := decodeItemsAt(show, 4)
	if err != nil || len(sells) == 0 {
		return 0, 0
	}
	if sells[0].ID == 1001 && sells[0].Count > 0 {
		return 1001, sells[0].Count
	}
	if sells[0].ID > 0 && sells[0].ID != 1001 {
		return sells[0].ID, sells[0].Count
	}
	return 0, 0
}

func decodeShops(body []byte) ([]game.Shop, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var shops []game.Shop
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		shops = append(shops, game.Shop{
			ID:   pb.IntField(m, 1),
			Name: pb.StringField(m, 2),
			Type: pb.IntField(m, 3),
		})
	}
	return shops, nil
}

func decodeGoods(body []byte) ([]game.Goods, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var goods []game.Goods
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		item, err := decodeGoodsItem(f.Bytes)
		if err != nil {
			return nil, err
		}
		goods = append(goods, item)
	}
	return goods, nil
}

func decodeGoodsItem(raw []byte) (game.Goods, error) {
	var g game.Goods
	m, err := pb.FieldMap(raw)
	if err != nil {
		return g, err
	}
	g.ID = pb.IntField(m, 1)
	g.Bought = pb.IntField(m, 2)
	g.Price = pb.IntField(m, 3)
	g.Limit = pb.IntField(m, 4)
	g.Unlocked = pb.BoolField(m, 5)
	g.ItemID = pb.IntField(m, 6)
	g.ItemCount = pb.IntField(m, 7)
	g.Countdown = pb.BoolField(m, 9)
	g.EndTime = pb.IntField(m, 10)
	fields, err := pb.Walk(raw)
	if err != nil {
		return g, err
	}
	for _, f := range fields {
		if f.Num != 8 {
			continue
		}
		cm, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return g, err
		}
		g.Conds = append(g.Conds, game.ShopCond{Type: pb.IntField(cm, 1), Param: pb.IntField(cm, 2)})
	}
	return g, nil
}

func decodeShopBuy(body []byte) (game.BuyOut, error) {
	var out game.BuyOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			g, err := decodeGoodsItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Goods = &g
		case 2:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Items = append(out.Items, it)
		case 3:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Cost = append(out.Cost, it)
		}
	}
	return out, nil
}

func decodeMallBuy(body []byte) (game.BuyOut, error) {
	var out game.BuyOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 2:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Items = append(out.Items, it)
		case 3:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Cost = append(out.Cost, it)
		case 4, 5:
			lim, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Bought = pb.IntField(lim, 2)
			out.Limit = pb.IntField(lim, 3)
		}
	}
	return out, nil
}

func decodeVisit(body []byte) (game.Visit, error) {
	var out game.Visit
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			host, err := decodeBasic(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Host = host
		case 2:
			land, err := decodeLand(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Lands = append(out.Lands, land)
		case 3:
			dog, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.DogID = pb.IntField(dog, 1)
			out.DogFoodSec = pb.IntField(dog, 2)
		case 6:
			out.AtHome = f.Varint != 0
		case 7:
			rel, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Friend = pb.BoolField(rel, 2)
			out.Community = pb.BoolField(rel, 3)
		case 8:
			out.ServerMs = int64(f.Varint)
		case 9:
			out.Tourist = f.Varint != 0
		case 10:
			out.ApplyToken = string(f.Bytes)
		case 11:
			lim, err := decodeOpLimit(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Limits = append(out.Limits, lim)
		case 13:
			w, err := decodeWeather(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Weather = w
		}
	}
	return out, nil
}

func decodeFriends(body []byte) ([]game.Friend, error) {
	out, err := decodeFriendsOut(body)
	if err != nil {
		return nil, err
	}
	return out.Friends, nil
}

func decodeFriendsOut(body []byte) (game.FriendsOut, error) {
	var out game.FriendsOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			friend, err := decodeFriend(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Friends = append(out.Friends, friend)
		case 3:
			out.ApplicationCount = int64(f.Varint)
		case 4:
			b, err := decodeBlocked(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Blocked = append(out.Blocked, b)
		case 5:
			b, err := decodeBlocked(f.Bytes)
			if err != nil {
				return out, err
			}
			out.BlockedBy = append(out.BlockedBy, b)
		}
	}
	return out, nil
}

func decodeFriend(raw []byte) (game.Friend, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.Friend{}, err
	}
	f := game.Friend{
		GID:        pb.IntField(m, 1),
		OpenID:     pb.StringField(m, 2),
		Name:       pb.StringField(m, 3),
		Avatar:     pb.StringField(m, 4),
		Remark:     pb.StringField(m, 5),
		Level:      pb.IntField(m, 6),
		Gold:       pb.IntField(m, 7),
		Authorized: pb.IntField(m, 10),
		LastLogin:  pb.IntField(m, 19),
		Exp:        pb.IntField(m, 14),
		Banned:     pb.BoolField(m, 15),
		Returner:   pb.BoolField(m, 17),
	}
	if tagsRaw := m[8].Bytes; len(tagsRaw) > 0 {
		t, err := pb.FieldMap(tagsRaw)
		if err != nil {
			return f, err
		}
		f.New = pb.BoolField(t, 1)
		f.Follow = pb.BoolField(t, 2)
	}
	if plantRaw := m[9].Bytes; len(plantRaw) > 0 {
		p, err := pb.FieldMap(plantRaw)
		if err != nil {
			return f, err
		}
		f.DryNum = pb.IntField(p, 7)
		f.WeedNum = pb.IntField(p, 8)
		f.InsectNum = pb.IntField(p, 9)
		f.StealNum = pb.IntField(p, 6)
	}
	if album := m[12].Bytes; len(album) > 0 {
		a, err := pb.FieldMap(album)
		if err != nil {
			return f, err
		}
		f.AlbumNormal = pb.IntField(a, 1)
		f.AlbumPremium = pb.IntField(a, 2)
	}
	if pass := m[18].Bytes; len(pass) > 0 {
		p, err := pb.FieldMap(pass)
		if err != nil {
			return f, err
		}
		f.PassSeason = pb.IntField(p, 1)
		f.PassLevel = pb.IntField(p, 2)
	}
	if raw := m[20].Bytes; len(raw) > 0 {
		w, err := decodeWeather(raw)
		if err != nil {
			return f, err
		}
		f.Weather = &w
	}
	return f, nil
}

func decodeApplications(body []byte) (game.ApplicationsOut, error) {
	var out game.ApplicationsOut
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
			out.Applications = append(out.Applications, game.Application{
				GID:    pb.IntField(m, 1),
				Time:   pb.IntField(m, 2),
				OpenID: pb.StringField(m, 3),
				Name:   pb.StringField(m, 4),
				Avatar: pb.StringField(m, 5),
			})
		case 2:
			out.Blocked = f.Varint != 0
		}
	}
	return out, nil
}

func decodeAccept(body []byte) (game.AcceptOut, error) {
	var out game.AcceptOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			friend, err := decodeFriend(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Friends = append(out.Friends, friend)
		case 2:
			out.Success = int64(f.Varint)
		case 3:
			out.Failed = int64(f.Varint)
		case 4:
			out.Full = int64(f.Varint)
		}
	}
	return out, nil
}

func decodeReject(body []byte) (game.RejectOut, error) {
	var out game.RejectOut
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Count = pb.IntField(m, 1)
	return out, nil
}

func decodeSetTags(body []byte) (game.Friend, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.Friend{}, err
	}
	if len(m[1].Bytes) == 0 {
		return game.Friend{}, nil
	}
	return decodeFriend(m[1].Bytes)
}

func decodeOneLand(body []byte) (game.Land, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return game.Land{}, err
	}
	if len(top[1].Bytes) == 0 {
		return game.Land{}, nil
	}
	return decodeLand(top[1].Bytes)
}

func decodeCanOperate(body []byte) (game.CanOperateOut, error) {
	var out game.CanOperateOut
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.OK = pb.BoolField(m, 1)
	out.StealNum = pb.IntField(m, 2)
	return out, nil
}

func decodeHarvest(body []byte) (game.HarvestOut, error) {
	var out game.HarvestOut
	items, err := decodeItemsAt(body, 2)
	if err != nil {
		return out, err
	}
	lost, err := decodeItemsAt(body, 3)
	if err != nil {
		return out, err
	}
	extra, err := decodeItemsAt(body, 6)
	if err != nil {
		return out, err
	}
	lands, err := decodeLands(body)
	if err != nil {
		return out, err
	}
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 4:
			lim, err := decodeOpLimit(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Limits = append(out.Limits, lim)
		case 5:
			em, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Warnings = append(out.Warnings, game.HarvestWarn{
				LandID: pb.IntField(em, 1),
				Text:   pb.StringField(em, 2),
			})
		case 7:
			d, err := decodeLandDrop(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Drops = append(out.Drops, d)
		case 8:
			b, err := decodeBuff(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Buffs = append(out.Buffs, b)
		}
	}
	out.Items = items
	out.Lost = lost
	out.Extra = extra
	out.Lands = lands
	return out, nil
}

func decodeWeatherStatus(body []byte) (game.Weather, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return game.Weather{}, err
	}
	raw := top[1].Bytes
	if len(raw) == 0 {
		raw = body
	}
	return decodeWeather(raw)
}

func decodeWeather(raw []byte) (game.Weather, error) {
	var w game.Weather
	m, err := pb.FieldMap(raw)
	if err != nil {
		return w, err
	}
	w.Type = pb.IntField(m, 1)
	w.Source = pb.IntField(m, 2)
	w.Begin = pb.IntField(m, 3)
	w.End = pb.IntField(m, 4)
	w.Active = pb.BoolField(m, 5)
	w.Afterglow = pb.IntField(m, 6)
	w.AfterglowTyp = pb.IntField(m, 7)
	w.CanCollect = pb.BoolField(m, 8)
	w.BlockReason = pb.IntField(m, 9)
	return w, nil
}

func decodeCurrentWeather(body []byte) (game.Weather, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return game.Weather{}, err
	}
	raw := top[1].Bytes
	if len(raw) == 0 {
		return game.Weather{}, nil
	}
	return decodeWeatherInfo(raw)
}

func decodeTodayWeather(body []byte) ([]game.Weather, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Weather
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		w, err := decodeWeatherInfo(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, nil
}

func decodeWeatherInfo(raw []byte) (game.Weather, error) {
	var w game.Weather
	m, err := pb.FieldMap(raw)
	if err != nil {
		return w, err
	}
	w.Type = pb.IntField(m, 1)
	w.Begin = pb.IntField(m, 2)
	w.End = pb.IntField(m, 3)
	w.Name = pb.StringField(m, 4)
	return w, nil
}

func decodeTaskBoard(body []byte) (game.TaskBoard, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return game.TaskBoard{}, err
	}
	raw := top[1].Bytes
	if len(raw) == 0 {
		return game.TaskBoard{}, nil
	}
	return decodeTaskInfo(raw)
}

func decodeTaskInfo(raw []byte) (game.TaskBoard, error) {
	var out game.TaskBoard
	if len(raw) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			task, err := decodeTask(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Growth = append(out.Growth, task)
		case 2:
			task, err := decodeTask(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Daily = append(out.Daily, task)
		case 4:
			active, err := decodeActive(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Actives = append(out.Actives, active)
		}
	}
	return out, nil
}

func decodeTaskClaim(body []byte) (game.TaskClaimOut, error) {
	var out game.TaskClaimOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Items = append(out.Items, it)
		case 2:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Compensated = append(out.Compensated, it)
		}
	}
	return out, nil
}

func decodeTask(raw []byte) (game.Task, error) {
	var t game.Task
	fields, err := pb.Walk(raw)
	if err != nil {
		return t, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			t.ID = int64(f.Varint)
		case 2:
			t.Progress = int64(f.Varint)
		case 3:
			t.Claimed = f.Varint != 0
		case 4:
			t.Unlocked = f.Varint != 0
		case 5:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return t, err
			}
			t.Rewards = append(t.Rewards, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		case 6:
			t.Total = int64(f.Varint)
		case 7:
			t.Share = int64(f.Varint)
		case 9:
			t.Desc = string(f.Bytes)
		case 10:
			t.Type = int64(f.Varint)
		case 11:
			t.Group = int64(f.Varint)
		case 12:
			t.Cond = int64(f.Varint)
		case 14:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return t, err
			}
			t.Extra = append(t.Extra, it)
		}
	}
	return t, nil
}

func decodeActive(raw []byte) (game.Active, error) {
	var a game.Active
	fields, err := pb.Walk(raw)
	if err != nil {
		return a, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			a.Type = int64(f.Varint)
		case 2:
			a.Progress = int64(f.Varint)
		case 3:
			box, err := decodeActiveBox(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Boxes = append(a.Boxes, box)
		}
	}
	return a, nil
}

func decodeActiveBox(raw []byte) (game.ActiveBox, error) {
	var b game.ActiveBox
	fields, err := pb.Walk(raw)
	if err != nil {
		return b, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			b.ID = int64(f.Varint)
		case 2:
			b.Need = int64(f.Varint)
		case 3:
			b.Status = int64(f.Varint)
			b.CanClaim = b.Status == 1
			b.Claimed = b.Status == 2
		case 4:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return b, err
			}
			b.Rewards = append(b.Rewards, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		}
	}
	return b, nil
}

func decodeEmails(body []byte) ([]game.Email, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Email
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		mail, err := decodeEmailItem(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, mail)
	}
	return out, nil
}

func decodeEmailClaim(body []byte) (game.EmailClaimOut, error) {
	var out game.EmailClaimOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Items = append(out.Items, it)
		case 2:
			out.Unclaimed = append(out.Unclaimed, string(f.Bytes))
		}
	}
	return out, nil
}

func decodeEmailItem(raw []byte) (game.Email, error) {
	var e game.Email
	m, err := pb.FieldMap(raw)
	if err != nil {
		return e, err
	}
	e.ID = pb.StringField(m, 1)
	e.Type = pb.IntField(m, 2)
	e.Title = pb.StringField(m, 3)
	e.Read = pb.BoolField(m, 4)
	e.HasReward = pb.BoolField(m, 5)
	e.SendTime = pb.IntField(m, 6)
	e.Subtitle = pb.StringField(m, 7)
	e.Expire = pb.IntField(m, 8)
	e.Tips = pb.StringField(m, 9)
	return e, nil
}

func decodeEmailDetail(body []byte) (game.EmailDetail, error) {
	var out game.EmailDetail
	top, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	raw := top[1].Bytes
	if len(raw) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.ID = string(f.Bytes)
		case 2:
			out.Type = int64(f.Varint)
		case 3:
			out.Title = string(f.Bytes)
		case 4:
			out.Content = string(f.Bytes)
		case 5:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Rewards = append(out.Rewards, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		case 6:
			out.SendTime = int64(f.Varint)
		case 7:
			out.Expire = int64(f.Varint)
		case 8:
			out.Read = f.Varint != 0
		case 9:
			out.Reason = string(f.Bytes)
		case 10:
			out.Claimed = f.Varint != 0
		case 11:
			out.Tips = string(f.Bytes)
		}
	}
	return out, nil
}

func decodeHeartbeat(body []byte) (game.HeartbeatOut, error) {
	var out game.HeartbeatOut
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.ServerMs = pb.IntField(m, 1)
	return out, nil
}

func decodeGetGroup(body []byte) (game.ActivityGroup, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return game.ActivityGroup{}, err
	}
	for _, f := range fields {
		if f.Num == 1 {
			return decodeActivityGroupFull(f.Bytes)
		}
	}
	return game.ActivityGroup{}, nil
}

func decodeActivities(body []byte) ([]game.Activity, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Activity
	seen := map[int64]bool{}
	for _, f := range fields {
		switch f.Num {
		case 1:
			acts, err := decodeActivityGroup(f.Bytes)
			if err != nil {
				return nil, err
			}
			for _, a := range acts {
				if a.ID != 0 && seen[a.ID] {
					continue
				}
				if a.ID != 0 {
					seen[a.ID] = true
				}
				out = append(out, a)
			}
		case 2:
			a, err := decodeActivitySummary(f.Bytes)
			if err != nil {
				return nil, err
			}
			if a.ID != 0 && seen[a.ID] {
				continue
			}
			if a.ID != 0 {
				seen[a.ID] = true
			}
			out = append(out, a)
		}
	}
	return out, nil
}

func decodeActivityGroup(raw []byte) ([]game.Activity, error) {
	g, err := decodeActivityGroupFull(raw)
	if err != nil {
		return nil, err
	}
	return g.Activities, nil
}

func decodeActivityGroupFull(raw []byte) (game.ActivityGroup, error) {
	var g game.ActivityGroup
	fields, err := pb.Walk(raw)
	if err != nil {
		return g, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			head, err := decodeActivityHead(f.Bytes)
			if err != nil {
				return g, err
			}
			g.ID = head.ID
			g.Name = head.Name
			g.Type = head.Type
			g.Start = head.Start
			g.End = head.End
			g.Status = head.Status
		case 2:
			act, err := decodeActivityData(f.Bytes)
			if err != nil {
				return g, err
			}
			if act.ID == 0 {
				continue
			}
			if act.GroupID == 0 {
				act.GroupID = g.ID
			}
			g.Activities = append(g.Activities, act)
		}
	}
	return g, nil
}

func decodeActivityHead(raw []byte) (game.Activity, error) {
	var a game.Activity
	m, err := pb.FieldMap(raw)
	if err != nil {
		return a, err
	}
	a.ID = pb.IntField(m, 1)
	a.GroupID = pb.IntField(m, 2)
	a.Type = pb.IntField(m, 3)
	a.Name = pb.StringField(m, 4)
	a.Desc = pb.StringField(m, 5)
	a.Start = pb.IntField(m, 6) // Unix 秒。活动号不是开日，以这里为准
	a.End = pb.IntField(m, 7)   // Unix 秒。本地配置没有结束日，以这里为准
	a.ClientID = pb.IntField(m, 8)
	a.Status = pb.IntField(m, 20)
	a.SplashOrder = pb.IntField(m, 21)
	a.Splashed = pb.BoolField(m, 22)
	a.RedDot = pb.BoolField(m, 23)
	return a, nil
}

func decodeActivityData(raw []byte) (game.Activity, error) {
	var a game.Activity
	fields, err := pb.Walk(raw)
	if err != nil {
		return a, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			head, err := decodeActivityHead(f.Bytes)
			if err != nil {
				return a, err
			}
			head.SigninClaimed = a.SigninClaimed
			head.SigninRewardID = a.SigninRewardID
			head.Signin = a.Signin
			head.Shop = a.Shop
			head.Nodes = a.Nodes
			head.Lottery = a.Lottery
			head.Brew = a.Brew
			head.Recall = a.Recall
			head.Invite = a.Invite
			head.Gift = a.Gift
			head.RandShop = a.RandShop
			head.Mega = a.Mega
			head.Progress = a.Progress
			head.Draw = a.Draw
			head.Cheer = a.Cheer
			head.Charity = a.Charity
			head.Drops = a.Drops
			head.Hunt = a.Hunt
			a = head
		case 101:
			shop, err := decodeRandShop(f.Bytes)
			if err != nil {
				return a, err
			}
			a.RandShop = &shop
		case 102:
			goods, err := decodeActivityShop(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Shop = goods
		case 103:
			signin, err := decodeSigninBody(f.Bytes)
			if err != nil {
				return a, err
			}
			a.SigninClaimed = signin.Claimed
			a.Signin = signin.Rewards
			if len(signin.Rewards) > 0 {
				a.SigninRewardID = signin.Rewards[0].ID
			}
		case 104:
			draw, err := decodeDrawInfo(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Draw = &draw
		case 105:
			lot, err := decodeLotteryInfo(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Lottery = &lot
		case 107:
			cheer, err := decodeCampCheer(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Cheer = &cheer
		case 108:
			brew, err := decodeBrewInfo(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Brew = &brew
		case 109:
			rec, err := decodeRecallInfo(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Recall = &rec
		case 111:
			inv, err := decodeInviteInfo(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Invite = &inv
		case 110:
			mega, err := decodeMegaEvent(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Mega = &mega
		case 112:
			prog, err := decodeProgressReward(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Progress = &prog
		case 113:
			gift, err := decodeGiftState(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Gift = &gift
		case 115:
			hunt, err := decodeHuntBody(f.Bytes)
			if err != nil {
				return a, err
			}
			if hunt != nil {
				a.Hunt = hunt
			}
		case 116:
			ch, err := decodeCharityFlower(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Charity = &ch
		case 117:
			drops, err := decodeDropPreviews(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Drops = drops
		case 118:
			nodes, err := decodeTechNodes(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Nodes = nodes
		}
	}
	return a, nil
}

type signinBody struct {
	Claimed bool
	Rewards []game.SigninReward
}

func decodeSigninBody(raw []byte) (signinBody, error) {
	var out signinBody
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.Claimed = f.Varint != 0
		case 2:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			items, err := decodeItemsAt(f.Bytes, 3)
			if err != nil {
				return out, err
			}
			out.Rewards = append(out.Rewards, game.SigninReward{
				ID:    pb.IntField(m, 1),
				Desc:  pb.StringField(m, 2),
				Items: items,
			})
		}
	}
	return out, nil
}

func decodeMegaEvent(raw []byte) (game.MegaEvent, error) {
	var out game.MegaEvent
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.Day = int64(f.Varint)
		case 2:
			out.Total = int64(f.Varint)
		case 3:
			out.Lookback = int64(f.Varint)
		case 4:
			rew, err := decodeMegaReward(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Rewards = append(out.Rewards, rew)
		case 5:
			seg, err := decodeMegaSegment(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Events = append(out.Events, seg)
		}
	}
	return out, nil
}

func decodeMegaReward(raw []byte) (game.MegaReward, error) {
	var r game.MegaReward
	m, err := pb.FieldMap(raw)
	if err != nil {
		return r, err
	}
	r.Day = pb.IntField(m, 1)
	r.Unlocked = pb.BoolField(m, 2)
	r.Claimed = pb.BoolField(m, 3)
	r.Claimable = pb.BoolField(m, 4)
	items, err := decodeItemsAt(raw, 5)
	if err != nil {
		return r, err
	}
	r.Rewards = items
	return r, nil
}

func decodeMegaSegment(raw []byte) (game.MegaSegment, error) {
	var s game.MegaSegment
	m, err := pb.FieldMap(raw)
	if err != nil {
		return s, err
	}
	s.Level = pb.IntField(m, 1)
	s.Unlocked = pb.BoolField(m, 2)
	s.Name = pb.StringField(m, 3)
	s.Art = pb.StringField(m, 4)
	s.Desc = pb.StringField(m, 5)
	return s, nil
}

func decodeProgressReward(raw []byte) (game.ProgressReward, error) {
	var out game.ProgressReward
	fields, err := pb.Walk(raw)
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
			out.Preview = append(out.Preview, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		case 2:
			step, err := decodeProgressStep(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Steps = append(out.Steps, step)
		case 3:
			out.Unlocked = int64(f.Varint)
		case 4:
			out.Completed = f.Varint != 0
		}
	}
	return out, nil
}

func decodeProgressStep(raw []byte) (game.ProgressStep, error) {
	var s game.ProgressStep
	m, err := pb.FieldMap(raw)
	if err != nil {
		return s, err
	}
	s.Step = pb.IntField(m, 1)
	s.Status = pb.IntField(m, 4)
	cost, err := decodeItemsAt(raw, 2)
	if err != nil {
		return s, err
	}
	s.Cost = cost
	rew, err := decodeItemsAt(raw, 3)
	if err != nil {
		return s, err
	}
	s.Rewards = rew
	return s, nil
}

func decodeLotteryInfo(raw []byte) (game.LotteryInfo, error) {
	var l game.LotteryInfo
	m, err := pb.FieldMap(raw)
	if err != nil {
		return l, err
	}
	l.FreeLeft = pb.IntField(m, 1)
	l.FreeLimit = pb.IntField(m, 2)
	l.PaidLeft = pb.IntField(m, 3)
	l.PaidLimit = pb.IntField(m, 4)
	l.CostID = pb.IntField(m, 5)
	l.CostCount = pb.IntField(m, 6)
	l.Diamond = pb.IntField(m, 7)
	l.BizType = pb.IntField(m, 9)
	fields, err := pb.Walk(raw)
	if err != nil {
		return l, err
	}
	for _, f := range fields {
		if f.Num != 8 {
			continue
		}
		pm, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return l, err
		}
		if id := pb.IntField(pm, 1); id > 0 {
			l.Preview = append(l.Preview, id)
		}
	}
	return l, nil
}

func decodeBrewInfo(raw []byte) (game.BrewInfo, error) {
	var b game.BrewInfo
	m, err := pb.FieldMap(raw)
	if err != nil {
		return b, err
	}
	b.Value = pb.IntField(m, 1)
	b.Step = pb.IntField(m, 2)
	b.Steps = pb.IntField(m, 3)
	b.Multipliers = pb.PackedInts(raw, 4)
	b.Amounts = pb.PackedInts(raw, 5)
	b.CanClaim = pb.BoolField(m, 6)
	b.Base = pb.IntField(m, 7)
	b.Allowed = pb.PackedInts(raw, 8)
	b.Share = pb.IntField(m, 9)
	return b, nil
}

func decodeRecallInfo(raw []byte) (game.RecallInfo, error) {
	var r game.RecallInfo
	m, err := pb.FieldMap(raw)
	if err != nil {
		return r, err
	}
	r.Daily = pb.IntField(m, 1)
	r.DailyLimit = pb.IntField(m, 2)
	r.Total = pb.IntField(m, 3)
	r.Max = pb.IntField(m, 4)
	r.Pending = pb.IntField(m, 5)
	pending, err := decodeItemsAt(raw, 6)
	if err != nil {
		return r, err
	}
	r.PendingItems = pending
	r.Returner = pb.BoolField(m, 10)
	r.LoginAt = pb.IntField(m, 11)
	r.BuffExpire = pb.IntField(m, 12)
	r.BuffCD = pb.IntField(m, 13)
	r.ByGID = pb.IntField(m, 14)
	r.ReturnPending = pb.IntField(m, 15)
	ret, err := decodeItemsAt(raw, 16)
	if err != nil {
		return r, err
	}
	r.ReturnItems = ret
	return r, nil
}

func decodeInviteInfo(raw []byte) (game.InviteInfo, error) {
	var i game.InviteInfo
	m, err := pb.FieldMap(raw)
	if err != nil {
		return i, err
	}
	i.Limit = pb.IntField(m, 2)
	i.Invitee = pb.BoolField(m, 3)
	if len(m[4].Bytes) > 0 {
		inv, err := decodeInvitePerson(m[4].Bytes)
		if err != nil {
			return i, err
		}
		i.Inviter = &inv
	}
	i.InviteCount = pb.IntField(m, 8)
	i.GrowthCount = pb.IntField(m, 9)
	i.Expire = pb.IntField(m, 10)
	fields, err := pb.Walk(raw)
	if err != nil {
		return i, err
	}
	for _, f := range fields {
		switch f.Num {
		case 5:
			task, err := decodeInviteTask(f.Bytes)
			if err != nil {
				return i, err
			}
			i.Invite = append(i.Invite, task)
		case 6:
			task, err := decodeInviteTask(f.Bytes)
			if err != nil {
				return i, err
			}
			i.Growth = append(i.Growth, task)
		case 7:
			task, err := decodeInviteTask(f.Bytes)
			if err != nil {
				return i, err
			}
			i.Newcomer = append(i.Newcomer, task)
		}
	}
	return i, nil
}

func decodeInvitePerson(raw []byte) (game.Invitee, error) {
	var p game.Invitee
	m, err := pb.FieldMap(raw)
	if err != nil {
		return p, err
	}
	p.GID = pb.IntField(m, 1)
	p.Name = pb.StringField(m, 2)
	p.Avatar = pb.StringField(m, 3)
	p.Level = pb.IntField(m, 4)
	return p, nil
}

func decodeInviteTask(raw []byte) (game.InviteTask, error) {
	var t game.InviteTask
	m, err := pb.FieldMap(raw)
	if err != nil {
		return t, err
	}
	t.Stage = pb.IntField(m, 1)
	t.Desc = pb.StringField(m, 2)
	t.Target = pb.IntField(m, 3)
	t.Current = pb.IntField(m, 4)
	t.Completed = pb.BoolField(m, 5)
	t.Claimed = pb.BoolField(m, 6)
	t.Level = pb.IntField(m, 8)
	items, err := decodeItemsAt(raw, 9)
	if err != nil {
		return t, err
	}
	t.Rewards = items
	return t, nil
}

func decodeGiftState(raw []byte) (game.GiftState, error) {
	var g game.GiftState
	m, err := pb.FieldMap(raw)
	if err != nil {
		return g, err
	}
	g.Sent = pb.IntField(m, 1)
	g.SendLimit = pb.IntField(m, 2)
	g.ReceiveLimit = pb.IntField(m, 3)
	fields, err := pb.Walk(raw)
	if err != nil {
		return g, err
	}
	for _, f := range fields {
		if f.Num != 4 {
			continue
		}
		offer, err := decodeGiftOffer(f.Bytes)
		if err != nil {
			return g, err
		}
		g.Gifts = append(g.Gifts, offer)
	}
	return g, nil
}

func decodeGiftOffer(raw []byte) (game.GiftOffer, error) {
	var o game.GiftOffer
	m, err := pb.FieldMap(raw)
	if err != nil {
		return o, err
	}
	o.Type = pb.IntField(m, 3)
	o.Content = pb.IntField(m, 4)
	cost, err := decodeItemsAt(raw, 1)
	if err != nil {
		return o, err
	}
	o.Cost = cost
	recv, err := decodeItemsAt(raw, 2)
	if err != nil {
		return o, err
	}
	o.Receive = recv
	return o, nil
}

func decodeLotteryReply(body []byte) (game.LotteryOut, error) {
	var out game.LotteryOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 108 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		if code := pb.IntField(m, 5); code != 0 {
			msg := pb.StringField(m, 6)
			if msg == "" {
				msg = fmt.Sprintf("lottery error %d", code)
			}
			return out, fmt.Errorf("%s", msg)
		}
		items, err := decodeItemsAt(f.Bytes, 2)
		if err != nil {
			return out, err
		}
		out.Items = items
		costs, err := decodeItemsAt(f.Bytes, 3)
		if err != nil {
			return out, err
		}
		out.Costs = costs
		out.Partial = pb.BoolField(m, 4)
		hits, err := decodeLotteryHits(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Results = hits
	}
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	return out, nil
}

func decodeLotteryHits(raw []byte) ([]game.LotteryHit, error) {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var out []game.LotteryHit
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		items, err := decodeItemsAt(f.Bytes, 2)
		if err != nil {
			return nil, err
		}
		out = append(out, game.LotteryHit{
			GoodsID:   pb.IntField(m, 1),
			Items:     items,
			Quality:   pb.IntField(m, 3),
			Guarantee: pb.BoolField(m, 4),
		})
	}
	return out, nil
}

func decodeBrewStart(body []byte) (game.BrewStartOut, error) {
	var out game.BrewStartOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 113 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Value = pb.IntField(m, 1)
	}
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	return out, nil
}

func decodeBrewStep(body []byte) (game.BrewStepOut, error) {
	var out game.BrewStepOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 114 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Step = pb.IntField(m, 1)
		out.Multiplier = pb.IntField(m, 2)
		out.Amount = pb.IntField(m, 3)
		out.Finished = pb.BoolField(m, 4)
	}
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	return out, nil
}

func decodeSendCount(body []byte) (int64, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return 0, err
	}
	for _, f := range fields {
		if f.Num != 125 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return 0, err
		}
		return pb.IntField(m, 1), nil
	}
	return 0, nil
}

func decodeActivityShop(raw []byte) ([]game.ActivityGoods, error) {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var out []game.ActivityGoods
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		g, err := decodeActivityGoods(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

func decodeActivityGoods(raw []byte) (game.ActivityGoods, error) {
	var g game.ActivityGoods
	m, err := pb.FieldMap(raw)
	if err != nil {
		return g, err
	}
	g.ID = pb.IntField(m, 1)
	g.Limit = pb.IntField(m, 4)
	g.Bought = pb.IntField(m, 5)
	g.Name = pb.StringField(m, 7)
	g.Desc = pb.StringField(m, 8)
	g.Diamond = pb.IntField(m, 9)
	g.Order = pb.IntField(m, 6)
	g.Background = pb.IntField(m, 10)
	g.Restriction = pb.IntField(m, 11)
	g.Category = pb.StringField(m, 12)
	items, err := decodeItemsAt(raw, 2)
	if err != nil {
		return g, err
	}
	g.Items = items
	cost, err := decodeItemsAt(raw, 3)
	if err != nil {
		return g, err
	}
	g.Cost = cost
	return g, nil
}

func decodeRandShop(raw []byte) (game.RandShop, error) {
	var out game.RandShop
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			g, err := decodeRandGoods(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Goods = append(out.Goods, g)
		case 2:
			out.Next = int64(f.Varint)
		case 3:
			out.Cost = int64(f.Varint)
		case 4:
			out.CostID = int64(f.Varint)
		case 5:
			out.Diamond = int64(f.Varint)
		case 6:
			out.Limit = int64(f.Varint)
		case 7:
			out.Today = int64(f.Varint)
		case 8:
			out.FreeLimit = int64(f.Varint)
		case 9:
			out.FreeToday = int64(f.Varint)
		}
	}
	return out, nil
}

func decodeRandGoods(raw []byte) (game.RandGoods, error) {
	var g game.RandGoods
	m, err := pb.FieldMap(raw)
	if err != nil {
		return g, err
	}
	g.ID = pb.IntField(m, 1)
	g.Name = pb.StringField(m, 2)
	g.Desc = pb.StringField(m, 3)
	g.Limit = pb.IntField(m, 6)
	g.Bought = pb.IntField(m, 7)
	g.Available = pb.BoolField(m, 8)
	items, err := decodeItemsAt(raw, 4)
	if err != nil {
		return g, err
	}
	g.Items = items
	cost, err := decodeItemsAt(raw, 5)
	if err != nil {
		return g, err
	}
	g.Cost = cost
	return g, nil
}

func decodeTechNodes(raw []byte) ([]game.TechNode, error) {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var out []game.TechNode
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		tree, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		treeID := pb.IntField(tree, 1)
		nodes, err := pb.Walk(f.Bytes)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			if n.Num != 2 {
				continue
			}
			node, err := decodeTechNode(n.Bytes, treeID)
			if err != nil {
				return nil, err
			}
			out = append(out, node)
		}
	}
	return out, nil
}

func decodeTechNode(raw []byte, treeID int64) (game.TechNode, error) {
	var n game.TechNode
	m, err := pb.FieldMap(raw)
	if err != nil {
		return n, err
	}
	n.TreeID = treeID
	n.ID = pb.IntField(m, 1)
	n.Status = pb.IntField(m, 3)
	n.Progress = pb.IntField(m, 4)
	n.Target = pb.IntField(m, 5)
	cost, err := decodeItemsAt(raw, 6)
	if err != nil {
		return n, err
	}
	n.Cost = cost
	rewards, err := decodeItemsAt(raw, 7)
	if err != nil {
		return n, err
	}
	n.Rewards = rewards
	return n, nil
}

func decodeOperateItems(body []byte, msgField, itemField int32) ([]game.Item, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	for _, f := range fields {
		if int32(f.Num) != msgField {
			continue
		}
		return decodeItemsAt(f.Bytes, itemField)
	}
	return nil, nil
}

func decodeOperateActivity(body []byte) (*game.Activity, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	for _, f := range fields {
		if f.Num != 3 || len(f.Bytes) == 0 {
			continue
		}
		act, err := decodeActivityData(f.Bytes)
		if err != nil {
			return nil, err
		}
		return &act, nil
	}
	return nil, nil
}

func operateMsg(body []byte, num int32) []byte {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil
	}
	for _, f := range fields {
		if int32(f.Num) == num {
			return f.Bytes
		}
	}
	return nil
}

func decodeTechSubmit(body []byte) (game.ActivityOpOut, error) {
	out, err := decodeOperateOut(body, 140, 2, 0)
	if err != nil {
		return out, err
	}
	if raw := operateMsg(body, 140); len(raw) > 0 {
		out.Unlocked = pb.PackedInts(raw, 3)
	}
	return out, nil
}

func decodeDrawInfo(raw []byte) (game.DrawInfo, error) {
	var d game.DrawInfo
	m, err := pb.FieldMap(raw)
	if err != nil {
		return d, err
	}
	d.Today = pb.IntField(m, 1)
	d.Limit = pb.IntField(m, 2)
	d.Total = pb.IntField(m, 4)
	cost, err := decodeItemsAt(raw, 3)
	if err != nil {
		return d, err
	}
	d.Cost = cost
	return d, nil
}

func decodeCampCheer(raw []byte) (game.CampCheer, error) {
	var c game.CampCheer
	m, err := pb.FieldMap(raw)
	if err != nil {
		return c, err
	}
	c.CampID = pb.IntField(m, 1)
	c.Cheer = pb.IntField(m, 2)
	fields, err := pb.Walk(raw)
	if err != nil {
		return c, err
	}
	for _, f := range fields {
		if f.Num != 3 {
			continue
		}
		tier, err := decodeCheerTier(f.Bytes)
		if err != nil {
			return c, err
		}
		c.Tiers = append(c.Tiers, tier)
	}
	return c, nil
}

func decodeCheerTier(raw []byte) (game.CheerTier, error) {
	var t game.CheerTier
	m, err := pb.FieldMap(raw)
	if err != nil {
		return t, err
	}
	t.Index = pb.IntField(m, 1)
	t.Need = pb.IntField(m, 2)
	t.Claimed = pb.BoolField(m, 3)
	items, err := decodeItemsAt(raw, 4)
	if err != nil {
		return t, err
	}
	t.Rewards = items
	return t, nil
}

func decodeHuntBody(raw []byte) (*game.Hunt, error) {
	h := &game.Hunt{}
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	for _, f := range fields {
		if len(f.Bytes) == 0 {
			if f.Varint != 0 {
				h.UnreadPlunder = true
			}
			continue
		}
		if t := decodeHuntTreasure(f.Bytes); huntTreasureOK(t) {
			h.Treasures = append(h.Treasures, t)
			continue
		}
		if ts := collectHuntTreasures(f.Bytes); len(ts) > 0 {
			h.Treasures = append(h.Treasures, ts...)
			continue
		}
		pool, equipped := collectHuntCharms(f.Bytes)
		if len(pool) > 0 {
			h.DailyPool = pool
		}
		if len(equipped) > 0 {
			h.Equipped = equipped
		}
	}
	if len(h.Treasures) == 0 && len(h.DailyPool) == 0 && len(h.Equipped) == 0 && !h.UnreadPlunder {
		return nil, nil
	}
	return h, nil
}

func collectHuntTreasures(raw []byte) []game.HuntTreasure {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil
	}
	var out []game.HuntTreasure
	for _, f := range fields {
		if len(f.Bytes) == 0 {
			continue
		}
		t := decodeHuntTreasure(f.Bytes)
		if huntTreasureOK(t) {
			out = append(out, t)
			continue
		}
		out = append(out, collectHuntTreasures(f.Bytes)...)
	}
	return out
}

func decodeHuntTreasure(raw []byte) game.HuntTreasure {
	var t game.HuntTreasure
	m, err := pb.FieldMap(raw)
	if err != nil {
		return t
	}
	t.ID = pb.StringField(m, 1)
	if t.ID == "" {
		if n := pb.IntField(m, 1); n != 0 {
			t.ID = strconv.FormatInt(n, 10)
		}
	}
	t.ItemID = pb.IntField(m, 2)
	t.Count = pb.IntField(m, 3)
	t.Original = pb.IntField(m, 4)
	t.Created = pb.IntField(m, 5)
	t.Status = pb.IntField(m, 6)
	t.StartAt = pb.IntField(m, 7)
	t.EndAt = pb.IntField(m, 8)
	return t
}

func collectHuntCharms(raw []byte) (pool, equipped []game.HuntCharm) {
	byField := map[int32][]game.HuntCharm{}
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, nil
	}
	for _, f := range fields {
		if len(f.Bytes) == 0 {
			continue
		}
		if looksLikeHuntTreasure(f.Bytes) {
			continue
		}
		if c := decodeHuntCharm(f.Bytes); c.ID != 0 {
			byField[int32(f.Num)] = append(byField[int32(f.Num)], c)
			continue
		}
		innerPool, innerEq := collectHuntCharms(f.Bytes)
		if len(innerPool) > 0 {
			pool = innerPool
		}
		if len(innerEq) > 0 {
			equipped = innerEq
		}
	}
	var nums []int32
	for n := range byField {
		nums = append(nums, n)
	}
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[j] < nums[i] {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
	}
	if len(pool) == 0 && len(nums) >= 1 {
		pool = byField[nums[0]]
	}
	if len(equipped) == 0 && len(nums) >= 2 {
		equipped = byField[nums[1]]
	}
	return pool, equipped
}

func looksLikeHuntTreasure(raw []byte) bool {
	return huntTreasureOK(decodeHuntTreasure(raw))
}

func huntTreasureOK(t game.HuntTreasure) bool {
	return t.ItemID != 0 || t.Count != 0 || t.Status != 0 || t.Original != 0 || t.StartAt != 0
}

func decodeHuntCharm(raw []byte) game.HuntCharm {
	var c game.HuntCharm
	m, err := pb.FieldMap(raw)
	if err != nil {
		return c
	}
	c.ID = pb.IntField(m, 1)
	c.Total = pb.IntField(m, 2)
	c.Used = pb.IntField(m, 3)
	return c
}

func decodeHuntLogs(body []byte, cmd int64) (game.HuntLogOut, error) {
	var out game.HuntLogOut
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	for _, field := range []int32{int32(cmd), int32(cmd + 1)} {
		raw := operateMsg(body, field)
		if len(raw) == 0 {
			continue
		}
		logs, err := decodeHuntLogList(raw)
		if err != nil {
			return out, err
		}
		if len(logs) > 0 {
			out.Logs = logs
			break
		}
	}
	return out, nil
}

func decodeHuntLogList(raw []byte) ([]game.HuntLogEntry, error) {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var out []game.HuntLogEntry
	for _, f := range fields {
		if len(f.Bytes) == 0 {
			continue
		}
		entry, err := decodeHuntLogEntry(f.Bytes)
		if err != nil {
			return nil, err
		}
		if entry.Time != 0 || entry.Attacker != 0 || entry.Name != "" {
			out = append(out, entry)
			continue
		}
		inner, err := decodeHuntLogList(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, inner...)
	}
	return out, nil
}

func decodeHuntLogEntry(raw []byte) (game.HuntLogEntry, error) {
	var e game.HuntLogEntry
	fields, err := pb.Walk(raw)
	if err != nil {
		return e, err
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		return e, err
	}
	e.Time = pb.IntField(m, 1)
	e.Attacker = pb.IntField(m, 2)
	e.Name = pb.StringField(m, 3)
	e.Won = pb.BoolField(m, 4)
	e.Avatar = pb.StringField(m, 8)
	for _, f := range fields {
		switch f.Num {
		case 5:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return e, err
			}
			if it.ID != 0 || it.Count != 0 {
				e.Lost = append(e.Lost, it)
			}
		case 6:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return e, err
			}
			if it.ID != 0 || it.Count != 0 {
				e.Used = append(e.Used, it)
			}
		case 7:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return e, err
			}
			if it.ID != 0 || it.Count != 0 {
				e.Injected = append(e.Injected, it)
			}
		}
	}
	return e, nil
}

func decodeCharityFlower(raw []byte) (game.CharityFlower, error) {
	var c game.CharityFlower
	m, err := pb.FieldMap(raw)
	if err != nil {
		return c, err
	}
	c.LoveID = pb.IntField(m, 1)
	c.LoveCount = pb.IntField(m, 2)
	c.Personal = pb.IntField(m, 3)
	c.Global = pb.IntField(m, 4)
	c.MaxGlobal = pb.IntField(m, 5)
	c.Share = pb.IntField(m, 6)
	c.CanDonate = pb.BoolField(m, 16)
	c.Agreed = pb.BoolField(m, 19)
	share, err := decodeItemsAt(raw, 7)
	if err != nil {
		return c, err
	}
	c.ShareItems = share
	fields, err := pb.Walk(raw)
	if err != nil {
		return c, err
	}
	for _, f := range fields {
		switch f.Num {
		case 8:
			tier, err := decodeCharityTier(f.Bytes, true)
			if err != nil {
				return c, err
			}
			c.Tiers = append(c.Tiers, tier)
		case 10:
			tier, err := decodeCharityTier(f.Bytes, false)
			if err != nil {
				return c, err
			}
			c.GlobalTiers = append(c.GlobalTiers, tier)
		}
	}
	return c, nil
}

func decodeCharityTier(raw []byte, personal bool) (game.CharityTier, error) {
	var t game.CharityTier
	m, err := pb.FieldMap(raw)
	if err != nil {
		return t, err
	}
	t.Need = pb.IntField(m, 1)
	t.Reached = pb.BoolField(m, 3)
	if personal {
		t.Claimed = pb.BoolField(m, 4)
	}
	items, err := decodeItemsAt(raw, 2)
	if err != nil {
		return t, err
	}
	t.Rewards = items
	return t, nil
}

func decodeDropPreviews(raw []byte) ([]game.DropPreview, error) {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var out []game.DropPreview
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		d, err := decodeDropPreview(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func decodeDropPreview(raw []byte) (game.DropPreview, error) {
	var d game.DropPreview
	m, err := pb.FieldMap(raw)
	if err != nil {
		return d, err
	}
	d.ID = pb.IntField(m, 1)
	d.ItemID = pb.IntField(m, 2)
	d.Desc = pb.StringField(m, 4)
	d.Limit = pb.IntField(m, 5)
	d.Dropped = pb.IntField(m, 6)
	if raw3 := m[3].Bytes; len(raw3) > 0 {
		it, err := decodeItem(raw3)
		if err != nil {
			return d, err
		}
		d.Item = &it
	}
	return d, nil
}

func decodeDrawHistory(body []byte) (game.DrawHistoryOut, error) {
	var out game.DrawHistoryOut
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	raw := operateMsg(body, 106)
	if len(raw) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		rec, err := decodeDrawRecord(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Records = append(out.Records, rec)
	}
	return out, nil
}

func decodeDrawRecord(raw []byte) (game.DrawRecord, error) {
	var r game.DrawRecord
	m, err := pb.FieldMap(raw)
	if err != nil {
		return r, err
	}
	r.Time = pb.IntField(m, 1)
	r.Count = pb.IntField(m, 2)
	items, err := decodeItemsAt(raw, 3)
	if err != nil {
		return r, err
	}
	r.Rewards = items
	return r, nil
}

func decodeLotteryHistory(body []byte) (game.LotteryHistoryOut, error) {
	var out game.LotteryHistoryOut
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	raw := operateMsg(body, 109)
	if len(raw) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		rec, err := decodeLotteryRecord(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Records = append(out.Records, rec)
	}
	return out, nil
}

func decodeLotteryRecord(raw []byte) (game.LotteryRecord, error) {
	var r game.LotteryRecord
	m, err := pb.FieldMap(raw)
	if err != nil {
		return r, err
	}
	r.Time = pb.IntField(m, 1)
	r.CostType = pb.IntField(m, 3)
	r.CostCount = pb.IntField(m, 4)
	fields, err := pb.Walk(raw)
	if err != nil {
		return r, err
	}
	for _, f := range fields {
		if f.Num != 2 {
			continue
		}
		hm, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return r, err
		}
		items, err := decodeItemsAt(f.Bytes, 2)
		if err != nil {
			return r, err
		}
		r.Results = append(r.Results, game.LotteryHit{
			GoodsID:   pb.IntField(hm, 1),
			Items:     items,
			Quality:   pb.IntField(hm, 3),
			Guarantee: pb.BoolField(hm, 4),
		})
	}
	return r, nil
}

func decodeCheerSubmit(body []byte) (game.CheerSubmitOut, error) {
	var out game.CheerSubmitOut
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	raw := operateMsg(body, 111)
	if len(raw) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.Added = pb.IntField(m, 1)
	out.Cheer = pb.IntField(m, 2)
	out.Progress = pb.IntField(m, 3)
	return out, nil
}

func decodeRecallable(body []byte) (game.RecallListOut, error) {
	return decodeRecallList(body, 116, true)
}

func decodeRecalled(body []byte) (game.RecallListOut, error) {
	return decodeRecallList(body, 117, false)
}

func decodeRecallList(body []byte, msgField int32, recallable bool) (game.RecallListOut, error) {
	var out game.RecallListOut
	raw := operateMsg(body, msgField)
	if len(raw) == 0 {
		return out, nil
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		p, err := decodeRecallPerson(f.Bytes, recallable)
		if err != nil {
			return out, err
		}
		out.List = append(out.List, p)
	}
	return out, nil
}

func decodeRecallPerson(raw []byte, recallable bool) (game.RecallPerson, error) {
	var p game.RecallPerson
	m, err := pb.FieldMap(raw)
	if err != nil {
		return p, err
	}
	if recallable {
		p.Offline = pb.IntField(m, 2)
		p.LastLogin = pb.IntField(m, 3)
	} else {
		p.RecallAt = pb.IntField(m, 2)
	}
	player := m[1].Bytes
	if len(player) == 0 {
		return p, nil
	}
	pm, err := pb.FieldMap(player)
	if err != nil {
		return p, err
	}
	p.GID = pb.IntField(pm, 1)
	p.Name = pb.StringField(pm, 2)
	p.Avatar = pb.StringField(pm, 3)
	p.Level = pb.IntField(pm, 4)
	p.OpenID = pb.StringField(pm, 5)
	return p, nil
}

func decodeCharityDonate(body []byte) (game.CharityDonateOut, error) {
	var out game.CharityDonateOut
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	raw := operateMsg(body, 136)
	if len(raw) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.Consumed = pb.IntField(m, 1)
	out.Added = pb.IntField(m, 2)
	out.Personal = pb.IntField(m, 3)
	out.Global = pb.IntField(m, 4)
	return out, nil
}

func decodeCharityXhh(body []byte) (game.CharityXhhOut, error) {
	var out game.CharityXhhOut
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	raw := operateMsg(body, 138)
	if len(raw) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.Num = pb.IntField(m, 1)
	out.Code = pb.StringField(m, 2)
	out.Trans = pb.StringField(m, 4)
	out.BusinessID = pb.StringField(m, 5)
	items, err := decodeItemsAt(raw, 3)
	if err != nil {
		return out, err
	}
	out.Items = items
	return out, nil
}

func decodeOperateOut(body []byte, msgField, itemField, costField int32) (game.ActivityOpOut, error) {
	var out game.ActivityOpOut
	act, err := decodeOperateActivity(body)
	if err != nil {
		return out, err
	}
	out.Activity = act
	if itemField > 0 {
		items, err := decodeOperateItems(body, msgField, itemField)
		if err != nil {
			return out, err
		}
		out.Items = items
	}
	if costField > 0 {
		costs, err := decodeOperateItems(body, msgField, costField)
		if err != nil {
			return out, err
		}
		out.Costs = costs
	}
	return out, nil
}

func decodeMegaClaim(body []byte) (game.ActivityOpOut, error) {
	out, err := decodeOperateOut(body, 120, 2, 0)
	if err != nil {
		return out, err
	}
	if raw := operateMsg(body, 120); len(raw) > 0 {
		out.Days = pb.PackedInts(raw, 1)
	}
	return out, nil
}

func decodeProgressClaim(body []byte) (game.ActivityOpOut, error) {
	out, err := decodeOperateOut(body, 126, 2, 0)
	if err != nil {
		return out, err
	}
	if raw := operateMsg(body, 126); len(raw) > 0 {
		out.Steps = pb.PackedInts(raw, 1)
		m, err := pb.FieldMap(raw)
		if err != nil {
			return out, err
		}
		out.Completed = pb.BoolField(m, 3)
	}
	return out, nil
}

func decodeBrewClaim(body []byte) (game.ActivityOpOut, error) {
	out, err := decodeOperateOut(body, 115, 3, 0)
	if err != nil {
		return out, err
	}
	if raw := operateMsg(body, 115); len(raw) > 0 {
		m, err := pb.FieldMap(raw)
		if err != nil {
			return out, err
		}
		out.ClaimType = pb.IntField(m, 1)
		out.Gold = pb.IntField(m, 2)
	}
	return out, nil
}

func decodeMysteryShop(body []byte) (game.MysteryShop, error) {
	var out game.MysteryShop
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.Present = f.Varint != 0
		case 2:
			card, err := decodeMysteryGoods(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Goods = append(out.Goods, card)
		case 3:
			out.Start = int64(f.Varint)
		case 4:
			out.End = int64(f.Varint)
		}
	}
	return out, nil
}

func decodeMysteryGoods(raw []byte) (game.MysteryGoods, error) {
	var g game.MysteryGoods
	m, err := pb.FieldMap(raw)
	if err != nil {
		return g, err
	}
	g.ID = pb.IntField(m, 1)
	g.ItemID = pb.IntField(m, 2)
	g.Quality = pb.IntField(m, 3)
	g.Count = pb.IntField(m, 4)
	g.Currency = pb.IntField(m, 5)
	g.Price = pb.IntField(m, 6)
	g.Discount = pb.IntField(m, 7)
	g.Bought = pb.BoolField(m, 8)
	g.Original = pb.IntField(m, 9)
	return g, nil
}

func decodeMysteryBuy(body []byte) (game.MysteryBuyOut, error) {
	var out game.MysteryBuyOut
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return out, err
	}
	out.Items = items
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	if raw := m[2].Bytes; len(raw) > 0 {
		g, err := decodeMysteryGoods(raw)
		if err != nil {
			return out, err
		}
		out.Goods = &g
	}
	return out, nil
}

func decodeActivitySummary(raw []byte) (game.Activity, error) {
	var a game.Activity
	m, err := pb.FieldMap(raw)
	if err != nil {
		return a, err
	}
	a.ID = pb.IntField(m, 1)
	a.Name = pb.StringField(m, 2)
	a.Start = pb.IntField(m, 3)
	a.End = pb.IntField(m, 4)
	return a, nil
}

func decodeSeason(body []byte) (game.Season, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return game.Season{}, err
	}
	if len(top[1].Bytes) == 0 {
		return game.Season{}, nil
	}
	return decodeSeasonInfo(top[1].Bytes)
}

func decodeSeasonInfo(raw []byte) (game.Season, error) {
	var out game.Season
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.ID = int64(f.Varint)
		case 2:
			out.Name = string(f.Bytes)
		case 3:
			out.Phase = int64(f.Varint)
		case 4:
			out.Preheat = int64(f.Varint)
		case 5:
			out.Start = int64(f.Varint)
		case 6:
			out.End = int64(f.Varint)
		case 7:
			out.ServerTime = int64(f.Varint)
		case 8:
			a, err := decodeSeasonActivity(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Activities = append(out.Activities, a)
		case 9:
			next, err := decodeSeasonInfo(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Next = &next
		case 10:
			out.Pass, err = decodeBattlePass(f.Bytes)
			if err != nil {
				return out, err
			}
		}
	}
	return out, nil
}

func decodeSeasonActivity(raw []byte) (game.SeasonActivity, error) {
	var a game.SeasonActivity
	m, err := pb.FieldMap(raw)
	if err != nil {
		return a, err
	}
	a.ID = pb.IntField(m, 1)
	a.Type = pb.IntField(m, 2)
	a.Name = pb.StringField(m, 3)
	a.Start = pb.IntField(m, 4)
	a.End = pb.IntField(m, 5)
	return a, nil
}

func decodeBattlePass(raw []byte) (game.BattlePass, error) {
	var p game.BattlePass
	fields, err := pb.Walk(raw)
	if err != nil {
		return p, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			p.ID = int64(f.Varint)
		case 2:
			p.Level = int64(f.Varint)
		case 3:
			p.Exp = int64(f.Varint)
		case 4:
			p.LevelExp = int64(f.Varint)
		case 5:
			p.NextExp = int64(f.Varint)
		case 6:
			p.MaxLevel = int64(f.Varint)
		case 7:
			p.Premium = f.Varint != 0
		case 8:
			lv, err := decodePassLevel(f.Bytes)
			if err != nil {
				return p, err
			}
			p.Levels = append(p.Levels, lv)
		case 9:
			p.FreeClaimed = int64(f.Varint)
		case 10:
			p.PremiumClaimed = int64(f.Varint)
		case 15:
			p.Price = int64(f.Varint)
		case 16:
			p.Name = string(f.Bytes)
		case 17:
			p.Desc = string(f.Bytes)
		}
	}
	return p, nil
}

func decodePassLevel(raw []byte) (game.PassLevel, error) {
	var lv game.PassLevel
	m, err := pb.FieldMap(raw)
	if err != nil {
		return lv, err
	}
	lv.Level = pb.IntField(m, 1)
	lv.Key = pb.BoolField(m, 4)
	lv.Tag = pb.StringField(m, 5)
	lv.Free, err = decodeItemsAt(raw, 2)
	if err != nil {
		return lv, err
	}
	lv.Premium, err = decodeItemsAt(raw, 3)
	return lv, err
}

func decodePassClaim(body []byte) (game.PassClaimOut, error) {
	var out game.PassClaimOut
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return out, err
	}
	out.Items = items
	out.Levels = pb.PackedInts(body, 2)
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	if passRaw := m[3].Bytes; len(passRaw) > 0 {
		out.Pass, err = decodeBattlePass(passRaw)
		if err != nil {
			return out, err
		}
	}
	out.Overflow = pb.BoolField(m, 4)
	return out, nil
}

func decodeProducts(body []byte) ([]game.Product, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Product
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		p, err := decodeProduct(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func decodeProduct(raw []byte) (game.Product, error) {
	var p game.Product
	fields, err := pb.Walk(raw)
	if err != nil {
		return p, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			p.ID = int64(f.Varint)
		case 2:
			p.Name = string(f.Bytes)
		case 3:
			p.Num = int64(f.Varint)
		case 4:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return p, err
			}
			p.Rewards = append(p.Rewards, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		case 5:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return p, err
			}
			p.Price = game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)}
		case 6:
			p.Status = int64(f.Varint)
		case 7:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return p, err
			}
			p.Bought = pb.IntField(m, 2)
			p.Limit = pb.IntField(m, 3)
		case 8:
			p.Available = f.Varint != 0
		case 9:
			p.Discount = string(f.Bytes)
		case 10:
			p.Countdown = f.Varint != 0
		case 11:
			p.EndTime = int64(f.Varint)
		case 12:
			p.GoodsType = int64(f.Varint)
		case 13:
			p.Pic = string(f.Bytes)
		}
	}
	return p, nil
}

func decodeMonthCards(body []byte) ([]game.MonthCard, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.MonthCard
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		card, err := decodeMonthCard(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, card)
	}
	return out, nil
}

func decodeMonthCard(raw []byte) (game.MonthCard, error) {
	var c game.MonthCard
	fields, err := pb.Walk(raw)
	if err != nil {
		return c, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			c.ID = int64(f.Varint)
		case 2:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return c, err
			}
			c.Rewards = append(c.Rewards, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		case 3:
			c.Claimable = f.Varint != 0
		case 4:
			c.Days = int64(f.Varint)
		case 5:
			c.ClaimedAmount = int64(f.Varint)
		case 6:
			c.ExpireSeconds = int64(f.Varint)
		case 8:
			c.TotalCount = int64(f.Varint)
		case 9:
			c.TotalDays = int64(f.Varint)
		case 10:
			c.PayID = string(f.Bytes)
		case 11:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return c, err
			}
			c.PurchaseCost = &it
		case 12:
			c.Claimable2 = f.Varint != 0
		}
	}
	return c, nil
}

func decodeMonthCardClaim(body []byte) (game.MonthCardClaimOut, error) {
	var out game.MonthCardClaimOut
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return out, err
	}
	out.Items = items
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	if raw := m[2].Bytes; len(raw) > 0 {
		out.Card, err = decodeMonthCard(raw)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func decodeRedPackets(body []byte) ([]game.RedPacket, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.RedPacket
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		p, err := decodeRedPacket(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func decodeRedPacket(raw []byte) (game.RedPacket, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.RedPacket{}, err
	}
	p := game.RedPacket{
		ID:           pb.IntField(m, 1),
		Claimed:      pb.BoolField(m, 2),
		Status:       pb.IntField(m, 3),
		Desc:         pb.StringField(m, 4),
		Name:         pb.StringField(m, 5),
		Time:         pb.StringField(m, 6),
		Instructions: pb.StringField(m, 7),
	}
	p.CanClaim = !p.Claimed && p.Status != 0
	return p, nil
}

func decodeRedPacketClaim(body []byte) (game.RedPacketClaimOut, error) {
	var out game.RedPacketClaimOut
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Status = pb.IntField(m, 1)
	if raw := m[2].Bytes; len(raw) > 0 {
		it, err := decodeItem(raw)
		if err != nil {
			return out, err
		}
		if it.ID != 0 || it.Count != 0 {
			out.Items = []game.Item{it}
		}
	}
	return out, nil
}

func decodeVisitLogs(body []byte) ([]game.VisitLog, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.VisitLog
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		log, err := decodeVisitLog(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, log)
	}
	return out, nil
}

func decodeVisitLog(raw []byte) (game.VisitLog, error) {
	var log game.VisitLog
	m, err := pb.FieldMap(raw)
	if err != nil {
		return log, err
	}
	log.Time = pb.IntField(m, 1)
	log.Action = pb.IntField(m, 2)
	log.GID = pb.IntField(m, 3)
	log.Name = pb.StringField(m, 4)
	log.Avatar = pb.StringField(m, 5)
	log.CropID = pb.IntField(m, 6)
	log.CropCount = pb.IntField(m, 7)
	log.Times = pb.IntField(m, 8)
	log.FromType = pb.IntField(m, 9)
	log.Level = pb.IntField(m, 10)
	return log, nil
}

func decodeAlbum(body []byte) (game.Album, error) {
	var out game.Album
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			item, err := decodeAlbumItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Items = append(out.Items, item)
		case 2:
			out.Progress = int64(f.Varint)
		case 3:
			out.Level = int64(f.Varint)
		case 4:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Rewards = append(out.Rewards, it)
		case 6:
			out.Type = int64(f.Varint)
		case 7:
			out.Next = int64(f.Varint)
		case 8:
			out.Claimable = f.Varint != 0
		case 9:
			out.Claimed = f.Varint != 0
		case 10:
			b, err := decodeBuff(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Buffs = append(out.Buffs, b)
		case 11:
			b, err := decodeBuff(f.Bytes)
			if err != nil {
				return out, err
			}
			out.NextBuff = &b
		}
	}
	out.Rarities = pb.PackedInts(body, 5)
	return out, nil
}

func decodeBuff(raw []byte) (game.Buff, error) {
	var b game.Buff
	m, err := pb.FieldMap(raw)
	if err != nil {
		return b, err
	}
	b.ID = pb.IntField(m, 1)
	b.Value = pb.IntField(m, 2)
	b.Extra = pb.IntField(m, 3)
	return b, nil
}

func decodeAlbumClaim(body []byte) (game.AlbumClaimOut, error) {
	var out game.AlbumClaimOut
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return out, err
	}
	out.Items = items
	out.Levels = pb.PackedInts(body, 2)
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Level = pb.IntField(m, 3)
	out.Next, err = decodeItemsAt(body, 4)
	return out, err
}

func decodeAlbumItem(raw []byte) (game.AlbumItem, error) {
	var it game.AlbumItem
	fields, err := pb.Walk(raw)
	if err != nil {
		return it, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			it.FruitID = int64(f.Varint)
		case 2:
			it.Rarity = int64(f.Varint)
		case 3:
			it.Unlocked = f.Varint != 0
		case 4:
			it.Progress = int64(f.Varint)
		case 5:
			it.Layer = int64(f.Varint)
		case 6:
			rew, err := decodeItem(f.Bytes)
			if err != nil {
				return it, err
			}
			it.Rewards = append(it.Rewards, rew)
		case 7:
			it.New = f.Varint != 0
		}
	}
	return it, nil
}

func decodeDogYard(body []byte) (game.DogYard, error) {
	var out game.DogYard
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			d, err := decodeDog(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Dogs = append(out.Dogs, d)
		case 2:
			out.Deployed = int64(f.Varint)
		case 3:
			out.FoodLeft = int64(f.Varint)
		case 4:
			out.FoodMax = int64(f.Varint)
		case 5:
			food, err := decodeDogFood(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Foods = append(out.Foods, food)
		case 6:
			out.NewLog = f.Varint != 0
		case 7:
			out.Pending = int64(f.Varint)
		case 8:
			sk, err := decodeDogSkill(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Skills = append(out.Skills, sk)
		}
	}
	return out, nil
}

func decodeDog(raw []byte) (game.Dog, error) {
	var d game.Dog
	m, err := pb.FieldMap(raw)
	if err != nil {
		return d, err
	}
	d.ID = pb.IntField(m, 1)
	d.Name = pb.StringField(m, 2)
	d.Protect = pb.IntField(m, 3)
	d.LostMin = pb.IntField(m, 4)
	d.LostMax = pb.IntField(m, 5)
	d.Owned = pb.BoolField(m, 6)
	d.Activated = pb.BoolField(m, 7)
	d.Price = pb.IntField(m, 8)
	d.Expire = pb.IntField(m, 9)
	return d, nil
}

func decodeDogFood(raw []byte) (game.DogFood, error) {
	var f game.DogFood
	m, err := pb.FieldMap(raw)
	if err != nil {
		return f, err
	}
	f.ID = pb.IntField(m, 1)
	f.Seconds = pb.IntField(m, 2)
	f.Count = pb.IntField(m, 3)
	return f, nil
}

func decodeDogSkill(raw []byte) (game.DogSkill, error) {
	var s game.DogSkill
	m, err := pb.FieldMap(raw)
	if err != nil {
		return s, err
	}
	s.ID = pb.IntField(m, 1)
	s.Used = pb.IntField(m, 2)
	s.Max = pb.IntField(m, 3)
	s.DogID = pb.IntField(m, 4)
	return s, nil
}

func decodeProtectLogs(body []byte) (game.DogLogsOut, error) {
	var out game.DogLogsOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			log, err := decodeProtectLog(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Logs = append(out.Logs, log)
		case 2:
			out.Total = int64(f.Varint)
		}
	}
	return out, nil
}

func decodeProtectLog(raw []byte) (game.ProtectLog, error) {
	var l game.ProtectLog
	m, err := pb.FieldMap(raw)
	if err != nil {
		return l, err
	}
	l.Level = pb.IntField(m, 1)
	l.Name = pb.StringField(m, 2)
	l.Avatar = pb.StringField(m, 3)
	l.Online = pb.BoolField(m, 4)
	l.Time = pb.IntField(m, 5)
	l.Count = pb.IntField(m, 6)
	l.Gold = pb.IntField(m, 7)
	l.DogID = pb.IntField(m, 8)
	l.DogName = pb.StringField(m, 9)
	l.GID = pb.IntField(m, 10)
	l.LastOnline = pb.IntField(m, 11)
	l.Read = pb.BoolField(m, 12)
	l.Authorized = pb.IntField(m, 13)
	l.RecordType = pb.IntField(m, 15)
	l.SkillID = pb.IntField(m, 16)
	l.Skill = pb.StringField(m, 17)
	l.Trigger = pb.IntField(m, 18)
	return l, nil
}

func decodeDogGifts(body []byte) (game.DogGiftOut, error) {
	var out game.DogGiftOut
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return out, err
	}
	out.Items = items
	out.Compensated, err = decodeItemsAt(body, 2)
	if err != nil {
		return out, err
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Claimed = pb.IntField(m, 3)
	out.Pending = pb.IntField(m, 4)
	return out, nil
}

func decodeFoodLeft(body []byte) (int64, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return 0, err
	}
	return pb.IntField(m, 1), nil
}

func decodeActivateDog(body []byte, fallbackID int64) (game.Dog, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.Dog{}, err
	}
	if len(m[1].Bytes) == 0 {
		return game.Dog{ID: fallbackID, Activated: true}, nil
	}
	return decodeDog(m[1].Bytes)
}

func decodeWithdraw(body []byte) (game.DeployOut, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.DeployOut{}, err
	}
	return game.DeployOut{Withdrawn: pb.IntField(m, 1)}, nil
}

func decodeDeploy(body []byte) (game.DeployOut, error) {
	var out game.DeployOut
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Deployed = pb.IntField(m, 1)
	out.Previous = pb.IntField(m, 2)
	return out, nil
}

func decodeBulletins(body []byte) ([]game.Bulletin, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Bulletin
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.Bulletin{
			ID:      pb.IntField(m, 1),
			Title:   pb.StringField(m, 2),
			Read:    pb.BoolField(m, 3),
			Forced:  pb.BoolField(m, 4),
			PopType: pb.IntField(m, 5),
			Banner:  pb.StringField(m, 6),
		})
	}
	return out, nil
}

func decodeBulletinDetail(body []byte) (game.BulletinDetail, error) {
	var d game.BulletinDetail
	m, err := pb.FieldMap(body)
	if err != nil {
		return d, err
	}
	d.Title = pb.StringField(m, 1)
	d.Content = pb.StringField(m, 2)
	d.Start = pb.StringField(m, 3)
	d.End = pb.StringField(m, 4)
	d.Banner = pb.StringField(m, 5)
	d.DetailBanner = pb.StringField(m, 6)
	return d, nil
}

func decodeMutants(body []byte) ([]game.Mutant, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Mutant
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.Mutant{
			ID:         pb.IntField(m, 1),
			Start:      pb.IntField(m, 2),
			End:        pb.IntField(m, 3),
			ActivityID: pb.IntField(m, 4),
			Red:        pb.BoolField(m, 5),
		})
	}
	return out, nil
}

func decodeCareer(body []byte) (game.Career, error) {
	var c game.Career
	m, err := pb.FieldMap(body)
	if err != nil {
		return c, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return c, err
	}
	c.Items = items
	c.Harvested = pb.IntField(m, 2)
	c.Stolen = pb.IntField(m, 3)
	c.Name = pb.StringField(m, 4)
	c.Avatar = pb.StringField(m, 5)
	c.Remark = pb.StringField(m, 6)
	c.Signature = pb.StringField(m, 7)
	c.Gender = pb.IntField(m, 8)
	c.Level = pb.IntField(m, 9)
	c.Exp = pb.IntField(m, 10)
	c.GID = pb.IntField(m, 11)
	c.Authorized = pb.IntField(m, 13)
	return c, nil
}

func decodeRankBoard(body []byte) (game.RankBoard, error) {
	var out game.RankBoard
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
			out.Items = append(out.Items, game.RankItem{
				GID:    pb.IntField(m, 1),
				Name:   pb.StringField(m, 2),
				Value:  pb.IntField(m, 3),
				Rank:   pb.IntField(m, 4),
				Avatar: pb.StringField(m, 5),
				Level:  pb.IntField(m, 6),
			})
		case 2:
			out.Total = int64(f.Varint)
		}
	}
	return out, nil
}

func decodeAvatars(body []byte) ([]game.Avatar, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Avatar
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		a, err := decodeAvatar(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func decodeAvatar(raw []byte) (game.Avatar, error) {
	var a game.Avatar
	m, err := pb.FieldMap(raw)
	if err != nil {
		return a, err
	}
	a.ID = pb.IntField(m, 1)
	a.Type = pb.IntField(m, 2)
	a.Count = pb.IntField(m, 3)
	a.Priority = pb.IntField(m, 4)
	a.New = pb.BoolField(m, 5)
	a.Expire = pb.IntField(m, 6)
	return a, nil
}

func decodeEquippedAvatar(body []byte) (game.Avatar, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.Avatar{}, err
	}
	if len(m[1].Bytes) == 0 {
		return game.Avatar{}, nil
	}
	return decodeAvatar(m[1].Bytes)
}

func decodeSkins(body []byte) ([]game.Skin, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Skin
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.Skin{
			ID:       pb.IntField(m, 1),
			Slot:     pb.IntField(m, 2),
			Equipped: pb.IntField(m, 3) != 0,
			Expire:   pb.IntField(m, 4),
		})
	}
	return out, nil
}

func decodeDrops(body []byte) ([]game.Drop, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Drop
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		d, err := decodeDrop(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func decodeDrop(raw []byte) (game.Drop, error) {
	var d game.Drop
	m, err := pb.FieldMap(raw)
	if err != nil {
		return d, err
	}
	d.ID = pb.IntField(m, 1)
	d.Name = pb.StringField(m, 2)
	d.Status = pb.IntField(m, 3)
	d.Start = pb.IntField(m, 4)
	d.End = pb.IntField(m, 5)
	d.Dropped = pb.IntField(m, 7)
	d.Limit = pb.IntField(m, 8)
	rewards, err := pb.Walk(raw)
	if err != nil {
		return d, err
	}
	for _, f := range rewards {
		if f.Num != 6 {
			continue
		}
		rm, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return d, err
		}
		d.Rewards = append(d.Rewards, game.DropReward{
			ID:      pb.IntField(rm, 1),
			Count:   pb.IntField(rm, 2),
			Chance:  pb.IntField(rm, 3),
			Claimed: pb.BoolField(rm, 4),
		})
	}
	return d, nil
}

func decodeSolarTerms(body []byte) (game.SolarOut, error) {
	var out game.SolarOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			s, err := decodeSolarTerm(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Terms = append(out.Terms, s)
		case 2:
			out.ServerTime = int64(f.Varint)
		case 3:
			b, err := decodeSolarBasic(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Basic = &b
		case 4:
			b, err := decodeSolarBasic(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Basics = append(out.Basics, b)
		}
	}
	return out, nil
}

func decodeSolarBasic(raw []byte) (game.SolarBasic, error) {
	var b game.SolarBasic
	m, err := pb.FieldMap(raw)
	if err != nil {
		return b, err
	}
	b.ID = pb.IntField(m, 1)
	b.Season = pb.IntField(m, 2)
	b.Desc = pb.StringField(m, 3)
	b.Events = pb.PackedInts(raw, 4)
	return b, nil
}

func decodeSolarTerm(raw []byte) (game.SolarTerm, error) {
	var s game.SolarTerm
	m, err := pb.FieldMap(raw)
	if err != nil {
		return s, err
	}
	s.ID = pb.IntField(m, 1)
	s.Status = pb.IntField(m, 2)
	s.Start = pb.IntField(m, 3)
	s.End = pb.IntField(m, 4)
	s.Name = pb.StringField(m, 6)
	rewards, err := decodeItemsAt(raw, 5)
	if err != nil {
		return s, err
	}
	s.Rewards = rewards
	return s, nil
}

func decodeSolarClaim(body []byte) (game.SolarClaimOut, error) {
	var out game.SolarClaimOut
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return out, err
	}
	out.Items = items
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	if raw := m[2].Bytes; len(raw) > 0 {
		out.Event, err = decodeSolarTerm(raw)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func decodeAchieveScope(body []byte) (game.AchieveScope, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.AchieveScope{}, err
	}
	if len(m[1].Bytes) == 0 {
		return game.AchieveScope{}, nil
	}
	return decodeScopeView(m[1].Bytes)
}

func decodeAchieveGoalOut(body []byte) (game.AchieveGoalOut, error) {
	var out game.AchieveGoalOut
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Exp = pb.IntField(m, 1)
	out.Before = pb.IntField(m, 2)
	out.After = pb.IntField(m, 3)
	if raw := m[4].Bytes; len(raw) > 0 {
		out.Scope, err = decodeScopeView(raw)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func decodeScopeView(raw []byte) (game.AchieveScope, error) {
	var s game.AchieveScope
	top, err := pb.FieldMap(raw)
	if err != nil {
		return s, err
	}
	s.Kind = pb.IntField(top, 1)
	s.ID = pb.IntField(top, 2)
	s.Level = pb.IntField(top, 3)
	s.Exp = pb.IntField(top, 4)
	s.Next = pb.IntField(top, 5)
	s.Claimed = pb.IntField(top, 8)
	fields, err := pb.Walk(raw)
	if err != nil {
		return s, err
	}
	for _, f := range fields {
		switch f.Num {
		case 6:
			gm, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return s, err
			}
			need, err := decodeItemsAt(f.Bytes, 10)
			if err != nil {
				return s, err
			}
			s.Goals = append(s.Goals, game.AchieveGoal{
				ID:       pb.IntField(gm, 1),
				Cond:     pb.IntField(gm, 2),
				Progress: pb.IntField(gm, 3),
				Total:    pb.IntField(gm, 4),
				Claimed:  pb.BoolField(gm, 5),
				Unlock:   pb.IntField(gm, 6),
				Category: pb.IntField(gm, 7),
				Exp:      pb.IntField(gm, 8),
				Need:     need,
				Desc:     pb.StringField(gm, 11),
				Sort:     pb.IntField(gm, 9),
			})
		case 7:
			lm, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return s, err
			}
			rewards, err := decodeItemsAt(f.Bytes, 3)
			if err != nil {
				return s, err
			}
			s.Levels = append(s.Levels, game.AchieveLevel{
				Level:   pb.IntField(lm, 1),
				Need:    pb.IntField(lm, 2),
				Claimed: pb.BoolField(lm, 4),
				Rewards: rewards,
			})
		}
	}
	return s, nil
}

func decodeBlocked(raw []byte) (game.Blocked, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.Blocked{}, err
	}
	return game.Blocked{
		GID:    pb.IntField(m, 1),
		Name:   pb.StringField(m, 2),
		Avatar: pb.StringField(m, 3),
		Level:  pb.IntField(m, 4),
		Time:   pb.IntField(m, 5),
		OpenID: pb.StringField(m, 6),
	}, nil
}

func decodeBlockedList(body []byte) ([]game.Blocked, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Blocked
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		b, err := decodeBlocked(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func decodeBlockReply(body []byte) (game.Blocked, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.Blocked{}, err
	}
	if len(m[1].Bytes) == 0 {
		return game.Blocked{}, nil
	}
	return decodeBlocked(m[1].Bytes)
}

func decodeShareKey(body []byte) (string, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return "", err
	}
	return pb.StringField(m, 1), nil
}

func decodeUsersAt(body []byte, num protowire.Number) ([]game.User, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.User
	for _, f := range fields {
		if f.Num != num {
			continue
		}
		u, err := decodeBasic(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func decodeBrief(body []byte) (game.User, error) {
	m, err := pb.FieldMap(body)
	if err != nil {
		return game.User{}, err
	}
	if len(m[1].Bytes) == 0 {
		return game.User{}, nil
	}
	return decodeBasic(m[1].Bytes)
}

func decodeLockOut(body []byte) (game.LockOut, error) {
	var out game.LockOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.Newly = append(out.Newly, int64(f.Varint))
		case 2:
			out.Already = append(out.Already, int64(f.Varint))
		case 3:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Failures = append(out.Failures, game.LockFail{
				UID:    pb.IntField(m, 1),
				Code:   pb.IntField(m, 2),
				Reason: pb.StringField(m, 3),
			})
		}
	}
	return out, nil
}

func decodeShareInviteInfo(raw []byte) (game.ShareInviteInfo, error) {
	var info game.ShareInviteInfo
	m, err := pb.FieldMap(raw)
	if err != nil {
		return info, err
	}
	info.ID = pb.IntField(m, 1)
	info.ShareKey = pb.StringField(m, 2)
	info.ShareURL = pb.StringField(m, 3)
	info.InviteCount = pb.IntField(m, 4)
	info.RewardCount = pb.IntField(m, 5)
	info.CanClaim = pb.BoolField(m, 6)
	return info, nil
}

func decodeShareInvite(body []byte) (game.ShareInviteOut, error) {
	var out game.ShareInviteOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		info, err := decodeShareInviteInfo(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Infos = append(out.Infos, info)
	}
	return out, nil
}

func decodeShareAwardOut(body []byte) (game.ShareAwardOut, error) {
	var out game.ShareAwardOut
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			info, err := decodeShareInviteInfo(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Info = info
		case 2:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Awards = append(out.Awards, it)
		case 3:
			it, err := decodeItem(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Compensated = append(out.Compensated, it)
		case 4:
			out.Awarded = f.Varint != 0
		}
	}
	return out, nil
}

func decodeAlbumLevels(body []byte) (game.AlbumLevels, error) {
	var out game.AlbumLevels
	fields, err := pb.Walk(body)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			out.Level = int64(f.Varint)
		case 2:
			out.Progress = int64(f.Varint)
		case 3:
			lv, err := decodeAlbumLevel(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Levels = append(out.Levels, lv)
		case 4:
			lv, err := decodeAlbumLevel(f.Bytes)
			if err != nil {
				return out, err
			}
			out.Extra = append(out.Extra, lv)
		}
	}
	return out, nil
}

func decodeAlbumLevel(raw []byte) (game.AlbumLevel, error) {
	var lv game.AlbumLevel
	m, err := pb.FieldMap(raw)
	if err != nil {
		return lv, err
	}
	lv.Level = pb.IntField(m, 1)
	lv.Need = pb.IntField(m, 2)
	lv.CanClaim = pb.BoolField(m, 4)
	lv.Claimed = pb.BoolField(m, 5)
	rewards, err := decodeItemsAt(raw, 3)
	if err != nil {
		return lv, err
	}
	lv.Rewards = rewards
	return lv, nil
}

func decodeDogBuy(body []byte, fallbackID int64) (game.DogBuyOut, error) {
	var out game.DogBuyOut
	dog, err := decodeActivateDog(body, fallbackID)
	if err != nil {
		return out, err
	}
	cost, err := decodeItemsAt(body, 2)
	if err != nil {
		return out, err
	}
	out.Dog = dog
	out.Cost = cost
	return out, nil
}

func decodeInvitees(body []byte) ([]game.Invitee, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	for _, f := range fields {
		if f.Num == 123 {
			return decodeInviteeList(f.Bytes)
		}
	}
	return nil, nil
}

func decodeInviteeList(raw []byte) ([]game.Invitee, error) {
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var out []game.Invitee
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.Invitee{
			GID:        pb.IntField(m, 1),
			Level:      pb.IntField(m, 2),
			Name:       pb.StringField(m, 3),
			Avatar:     pb.StringField(m, 4),
			RegisterAt: pb.IntField(m, 5),
			InvitedAt:  pb.IntField(m, 6),
		})
	}
	return out, nil
}

func decodeVisitSummary(body []byte) (game.VisitSummary, error) {
	var out game.VisitSummary
	top, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	raw := top[1].Bytes
	if len(raw) == 0 {
		raw = body
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.StealCount = pb.IntField(m, 1)
	out.HelpCount = pb.IntField(m, 2)
	out.MischiefCount = pb.IntField(m, 3)
	out.StealItemNum = pb.IntField(m, 5)
	fields, err := pb.Walk(raw)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		if f.Num != 4 {
			continue
		}
		fm, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return out, err
		}
		out.Friends = append(out.Friends, game.VisitFriend{
			GID:           pb.IntField(fm, 1),
			Name:          pb.StringField(fm, 2),
			Avatar:        pb.StringField(fm, 3),
			StealCount:    pb.IntField(fm, 4),
			StealItemNum:  pb.IntField(fm, 5),
			HelpCount:     pb.IntField(fm, 6),
			MischiefCount: pb.IntField(fm, 7),
			Level:         pb.IntField(fm, 9),
		})
	}
	return out, nil
}

func decodeItemID(body []byte) (int64, error) {
	if len(body) == 0 {
		return 0, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return 0, err
	}
	return pb.IntField(m, 1), nil
}

func decodeRelation(body []byte) (game.Relation, error) {
	var out game.Relation
	if len(body) == 0 {
		return out, nil
	}
	top, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	raw := top[1].Bytes
	if len(raw) == 0 {
		raw = body
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.Friend = pb.BoolField(m, 1)
	out.Community = pb.BoolField(m, 2)
	out.Stranger = pb.BoolField(m, 3)
	return out, nil
}

func decodeOK(body []byte) (bool, error) {
	if len(body) == 0 {
		return true, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return false, err
	}
	if _, ok := m[1]; !ok {
		return true, nil
	}
	return pb.BoolField(m, 1) || pb.IntField(m, 1) != 0, nil
}

func decodeVisitPopup(body []byte) (game.VisitPopup, error) {
	var out game.VisitPopup
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Unread = pb.BoolField(m, 1)
	out.LastRead = pb.IntField(m, 2)
	out.NeedPopup = pb.BoolField(m, 3)
	return out, nil
}

func decodeMallProfiles(body []byte) ([]game.MallProfile, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.MallProfile
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.MallProfile{
			ID:   pb.IntField(m, 1),
			Type: pb.IntField(m, 2),
		})
	}
	return out, nil
}

func decodeRedDot(body []byte) (game.RedDot, error) {
	var out game.RedDot
	if len(body) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	out.Red = pb.BoolField(m, 1)
	return out, nil
}
