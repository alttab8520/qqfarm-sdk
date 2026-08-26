package client

import (
	"fmt"
	"os"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/pb"
	"google.golang.org/protobuf/encoding/protowire"
)

const (
	fallbackGameVersion = "1.13.3.11_20260826"
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

func encodeHarvest(landIDs []int64, hostGID int64, isAll bool) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, hostGID)
	req.Bool(3, isAll)
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

func encodeEnter(hostGID int64) []byte {
	req := pb.NewEncoder()
	req.Int(1, hostGID)
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

func encodeWater(landIDs []int64, hostGID int64) []byte {
	return encodeLandOp(landIDs, hostGID)
}

func encodeLandOp(landIDs []int64, hostGID int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, hostGID)
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

func encodeFarming(landIDs []int64, hostGID int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, hostGID)
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
	req := pb.NewEncoder()
	req.Int(1, id)
	return req.Bytes()
}

func encodeBatchClaim(ids []int64) []byte {
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
	return u, nil
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

func decodeLand(raw []byte) (game.Land, error) {
	var land game.Land
	m, err := pb.FieldMap(raw)
	if err != nil {
		return land, err
	}
	land.ID = pb.IntField(m, 1)
	land.Unlocked = pb.BoolField(m, 2)
	land.Level = pb.IntField(m, 3)
	land.CouldUnlock = pb.BoolField(m, 5)
	if plantRaw := m[10].Bytes; len(plantRaw) > 0 {
		p, err := pb.FieldMap(plantRaw)
		if err != nil {
			return land, err
		}
		land.PlantID = pb.IntField(p, 1)
		land.PlantName = pb.StringField(p, 2)
		land.DryNum = pb.IntField(p, 6)
		land.FruitID = pb.IntField(p, 10)
		land.HasWeed = len(p[12].Bytes) > 0 || p[12].Varint > 0
		land.HasInsect = len(p[13].Bytes) > 0 || p[13].Varint > 0
		land.Stealable = pb.BoolField(p, 16)
		land.LeftFruit = pb.IntField(p, 18)
	}
	return land, nil
}

func decodeItems(body []byte) ([]game.Item, error) {
	return decodeItemsAt(body, 2)
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
	top, err := pb.FieldMap(body)
	if err != nil {
		return nil, err
	}
	return decodeItemBag(top[1].Bytes)
}

func decodeLoginBag(body []byte) ([]game.BagItem, error) {
	top, err := pb.FieldMap(body)
	if err != nil {
		return nil, err
	}
	return decodeItemBag(top[2].Bytes)
}

func decodeItemBag(raw []byte) ([]game.BagItem, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var items []game.BagItem
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		item, err := decodeBagItem(f.Bytes)
		if err != nil {
			return nil, err
		}
		if item.ID != 0 {
			items = append(items, item)
		}
	}
	return items, nil
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
	if show := m[100].Bytes; len(show) > 0 {
		sm, err := pb.FieldMap(show)
		if err != nil {
			return it, err
		}
		it.SellPrice = pb.IntField(sm, 1)
	}
	return it, nil
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
	return g, nil
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
		}
	}
	return out, nil
}

func decodeFriends(body []byte) ([]game.Friend, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var friends []game.Friend
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		friend, err := decodeFriend(f.Bytes)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func decodeFriend(raw []byte) (game.Friend, error) {
	m, err := pb.FieldMap(raw)
	if err != nil {
		return game.Friend{}, err
	}
	f := game.Friend{
		GID:    pb.IntField(m, 1),
		OpenID: pb.StringField(m, 2),
		Name:   pb.StringField(m, 3),
		Avatar: pb.StringField(m, 4),
		Level:  pb.IntField(m, 6),
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
	return f, nil
}

func decodeApplications(body []byte) ([]game.Application, error) {
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var out []game.Application
	for _, f := range fields {
		if f.Num != 1 {
			continue
		}
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.Application{
			GID:    pb.IntField(m, 1),
			Time:   pb.IntField(m, 2),
			OpenID: pb.StringField(m, 3),
			Name:   pb.StringField(m, 4),
			Avatar: pb.StringField(m, 5),
		})
	}
	return out, nil
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

func decodeCanOperate(body []byte) (bool, error) {
	if len(body) == 0 {
		return true, nil
	}
	m, err := pb.FieldMap(body)
	if err != nil {
		return true, err
	}
	return pb.BoolField(m, 1), nil
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
	var out game.TaskBoard
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
			t.Status = int64(f.Varint)
		case 5:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return t, err
			}
			t.Rewards = append(t.Rewards, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		case 6:
			t.Total = int64(f.Varint)
		case 9:
			t.Desc = string(f.Bytes)
		case 10:
			t.Type = int64(f.Varint)
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
			b.Claimed = f.Varint != 0
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

func decodeEmailItem(raw []byte) (game.Email, error) {
	var e game.Email
	m, err := pb.FieldMap(raw)
	if err != nil {
		return e, err
	}
	e.ID = pb.StringField(m, 1)
	e.Type = pb.IntField(m, 2)
	e.Title = pb.StringField(m, 3)
	e.Claimed = pb.BoolField(m, 4)
	e.HasReward = pb.BoolField(m, 5)
	e.Subtitle = pb.StringField(m, 7)
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
	fields, err := pb.Walk(raw)
	if err != nil {
		return nil, err
	}
	var groupID int64
	var out []game.Activity
	for _, f := range fields {
		switch f.Num {
		case 1:
			head, err := decodeActivityHead(f.Bytes)
			if err != nil {
				return nil, err
			}
			groupID = head.ID
		case 2:
			act, err := decodeActivityData(f.Bytes)
			if err != nil {
				return nil, err
			}
			if act.ID == 0 {
				continue
			}
			if act.GroupID == 0 {
				act.GroupID = groupID
			}
			out = append(out, act)
		}
	}
	return out, nil
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
	a.Start = pb.IntField(m, 6)
	a.End = pb.IntField(m, 7)
	a.ClientID = pb.IntField(m, 8)
	a.Status = pb.IntField(m, 20)
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
			head.Shop = a.Shop
			head.Nodes = a.Nodes
			head.Lottery = a.Lottery
			head.Brew = a.Brew
			head.Recall = a.Recall
			head.Invite = a.Invite
			head.Gift = a.Gift
			a = head
		case 102:
			goods, err := decodeActivityShop(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Shop = goods
		case 105:
			lot, err := decodeLotteryInfo(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Lottery = &lot
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
		case 113:
			gift, err := decodeGiftState(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Gift = &gift
		case 118:
			nodes, err := decodeTechNodes(f.Bytes)
			if err != nil {
				return a, err
			}
			a.Nodes = nodes
		case 103:
			signin, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return a, err
			}
			a.SigninClaimed = pb.BoolField(signin, 1)
			rewards, err := pb.Walk(f.Bytes)
			if err != nil {
				return a, err
			}
			for _, r := range rewards {
				if r.Num != 2 {
					continue
				}
				rm, err := pb.FieldMap(r.Bytes)
				if err != nil {
					return a, err
				}
				if id := pb.IntField(rm, 1); id > 0 {
					a.SigninRewardID = id
					break
				}
			}
		}
	}
	return a, nil
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
	b.CanClaim = pb.BoolField(m, 6)
	b.Allowed = pb.PackedInts(raw, 8)
	return b, nil
}

func decodeRecallInfo(raw []byte) (game.RecallInfo, error) {
	var r game.RecallInfo
	m, err := pb.FieldMap(raw)
	if err != nil {
		return r, err
	}
	r.Pending = pb.IntField(m, 5)
	r.Returner = pb.BoolField(m, 10)
	r.ReturnPending = pb.IntField(m, 15)
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
	i.InviteCount = pb.IntField(m, 8)
	return i, nil
}

func decodeGiftState(raw []byte) (game.GiftState, error) {
	var g game.GiftState
	m, err := pb.FieldMap(raw)
	if err != nil {
		return g, err
	}
	g.Sent = pb.IntField(m, 1)
	g.SendLimit = pb.IntField(m, 2)
	return g, nil
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
		return out, nil
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
		return out, nil
	}
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
		return out, nil
	}
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
	var out game.Season
	top, err := pb.FieldMap(body)
	if err != nil {
		return out, err
	}
	raw := top[1].Bytes
	if len(raw) == 0 {
		return out, nil
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		return out, err
	}
	out.ID = pb.IntField(m, 1)
	out.Name = pb.StringField(m, 2)
	out.Phase = pb.IntField(m, 3)
	out.Start = pb.IntField(m, 5)
	out.End = pb.IntField(m, 6)
	out.ServerTime = pb.IntField(m, 7)
	if passRaw := m[10].Bytes; len(passRaw) > 0 {
		out.Pass, err = decodeBattlePass(passRaw)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func decodeBattlePass(raw []byte) (game.BattlePass, error) {
	var p game.BattlePass
	m, err := pb.FieldMap(raw)
	if err != nil {
		return p, err
	}
	p.ID = pb.IntField(m, 1)
	p.Level = pb.IntField(m, 2)
	p.Exp = pb.IntField(m, 3)
	p.LevelExp = pb.IntField(m, 4)
	p.NextExp = pb.IntField(m, 5)
	p.MaxLevel = pb.IntField(m, 6)
	p.Premium = pb.BoolField(m, 7)
	p.FreeClaimed = pb.IntField(m, 9)
	p.PremiumClaimed = pb.IntField(m, 10)
	p.Name = pb.StringField(m, 16)
	return p, nil
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
		case 11:
			p.EndTime = int64(f.Varint)
		case 12:
			p.GoodsType = int64(f.Varint)
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
		case 9:
			c.TotalDays = int64(f.Varint)
		}
	}
	return c, nil
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
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, game.RedPacket{
			ID:       pb.IntField(m, 1),
			CanClaim: pb.BoolField(m, 3),
		})
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
		case 6, 8:
			out.Claimable = out.Claimable || f.Varint != 0
		}
	}
	return out, nil
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
		case 5, 6:
			m, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return it, err
			}
			it.Rewards = append(it.Rewards, game.Item{ID: pb.IntField(m, 1), Count: pb.IntField(m, 2)})
		case 7:
			it.New = f.Varint != 0
		}
	}
	return it, nil
}
