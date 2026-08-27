package client

import (
	"testing"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/pb"
)

func TestEncodeRankAndAvatarEquip(t *testing.T) {
	raw := encodeRank(0, 0)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 1 || pb.IntField(m, 2) != 1 {
		t.Fatalf("%+v", m)
	}
	raw = encodeAvatarEquip(0, true)
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m[2]; !ok {
		t.Fatal("unequip flag omitted")
	}
	item := pb.NewEncoder()
	item.Int(1, 7)
	item.String(2, "甲")
	item.Int(4, 1)
	top := pb.NewEncoder()
	top.Message(1, item.Bytes())
	top.Int(2, 99)
	board, err := decodeRankBoard(top.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if board.Total != 99 || len(board.Items) != 1 || board.Items[0].GID != 7 {
		t.Fatalf("%+v", board)
	}
}

func TestDecodeDogYardAndBulletin(t *testing.T) {
	dog := pb.NewEncoder()
	dog.Int(1, 2)
	dog.String(2, "旺")
	dog.Bool(6, true)
	food := pb.NewEncoder()
	food.Int(1, 1)
	food.Int(2, 3600)
	food.Int(3, 4)
	yard := pb.NewEncoder()
	yard.Message(1, dog.Bytes())
	yard.Int(2, 2)
	yard.Int(3, 90)
	yard.Message(5, food.Bytes())
	got, err := decodeDogYard(yard.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if got.Deployed != 2 || got.FoodLeft != 90 || len(got.Dogs) != 1 || got.Dogs[0].Name != "旺" || len(got.Foods) != 1 {
		t.Fatalf("%+v", got)
	}

	raw := encodeFeed(1, 2)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 1 || pb.IntField(m, 2) != 2 {
		t.Fatalf("%+v", m)
	}

	brief := pb.NewEncoder()
	brief.Int(1, 9)
	brief.String(2, "更")
	list := pb.NewEncoder()
	list.Message(1, brief.Bytes())
	bs, err := decodeBulletins(list.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != 1 || bs[0].ID != 9 || bs[0].Title != "更" {
		t.Fatalf("%+v", bs)
	}
}

func TestEncodeLoginHasSceneAndDevice(t *testing.T) {
	raw := encodeLogin()
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(m[7].Bytes) != channelID {
		t.Fatalf("scene %q", m[7].Bytes)
	}
	dev, err := pb.FieldMap(m[5].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if string(dev[1].Bytes) != gameVersion() {
		t.Fatalf("ver %q", dev[1].Bytes)
	}
}

func TestDecodeUser(t *testing.T) {
	basic := pb.NewEncoder()
	basic.Int(1, 99)
	basic.String(2, "张三")
	basic.Int(3, 12)
	basic.Int(5, 100)
	basic.String(6, "oABC")
	top := pb.NewEncoder()
	top.Message(1, basic.Bytes())
	u, err := decodeUser(top.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if u.GID != 99 || u.Name != "张三" || u.OpenID != "oABC" {
		t.Fatalf("%+v", u)
	}
}

func TestEncodeAnti(t *testing.T) {
	raw := encodeAnti([]byte{7, 5, 1})
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(m[1].Bytes) != "\x07\x05\x01" {
		t.Fatalf("%x", m[1].Bytes)
	}
}

func TestDecodeBagAndVisit(t *testing.T) {
	sell := pb.NewEncoder()
	sell.Int(1, 1001)
	sell.Int(2, 12)
	show := pb.NewEncoder()
	show.Message(4, sell.Bytes())
	item := pb.NewEncoder()
	item.Int(1, 40061)
	item.Int(2, 8)
	item.Int(6, 77)
	item.Int(8, 2)
	item.Message(100, show.Bytes())
	bag := pb.NewEncoder()
	bag.Message(1, item.Bytes())
	reply := pb.NewEncoder()
	reply.Message(1, bag.Bytes())
	items, err := decodeBagReply(reply.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != 40061 || items[0].Count != 8 || items[0].UID != 77 || items[0].SellID != 1001 || items[0].SellPrice != 12 || len(items[0].Mutants) != 1 || items[0].Mutants[0] != 2 {
		t.Fatalf("%+v", items)
	}
	direct := pb.NewEncoder()
	direct.Int(1, 15)
	priced := pb.NewEncoder()
	priced.Int(1, 40003)
	priced.Int(2, 1)
	priced.Message(100, direct.Bytes())
	bag2 := pb.NewEncoder()
	bag2.Message(1, priced.Bytes())
	reply2 := pb.NewEncoder()
	reply2.Message(1, bag2.Bytes())
	directItems, err := decodeBagReply(reply2.Bytes())
	if err != nil || len(directItems) != 1 || directItems[0].SellPrice != 15 || directItems[0].SellID != 0 {
		t.Fatalf("f1 price %+v %v", directItems, err)
	}

	host := pb.NewEncoder()
	host.Int(1, 88)
	host.String(2, "友")
	land := pb.NewEncoder()
	land.Int(1, 3)
	land.Bool(2, true)
	dog := pb.NewEncoder()
	dog.Int(1, 90021)
	dog.Int(2, 60)
	enter := pb.NewEncoder()
	enter.Message(1, host.Bytes())
	enter.Message(2, land.Bytes())
	enter.Message(3, dog.Bytes())
	visit, err := decodeVisit(enter.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if visit.Host.GID != 88 || len(visit.Lands) != 1 || visit.Lands[0].ID != 3 || visit.DogID != 90021 {
		t.Fatalf("%+v", visit)
	}

	plant := pb.NewEncoder()
	plant.Int(6, 2)
	plant.Int(7, 1)
	fr := pb.NewEncoder()
	fr.Int(1, 5)
	fr.String(3, "友")
	fr.Message(9, plant.Bytes())
	list := pb.NewEncoder()
	list.Message(1, fr.Bytes())
	friends, err := decodeFriends(list.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(friends) != 1 || friends[0].GID != 5 || friends[0].StealNum != 2 || friends[0].DryNum != 1 {
		t.Fatalf("%+v", friends)
	}
}

func TestEncodeLotteryAndBrew(t *testing.T) {
	raw := encodeLottery(8, 99, 0, 0)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 2) != 9 {
		t.Fatalf("cmd %d", pb.IntField(m, 2))
	}
	inner, err := pb.FieldMap(m[107].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(inner, 3) != 99 {
		t.Fatalf("host %d", pb.IntField(inner, 3))
	}

	raw = encodeBrewStart(3, []game.BrewItem{{UID: 12, Count: 2}})
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 2) != 14 {
		t.Fatalf("cmd %d", pb.IntField(m, 2))
	}
	start, err := pb.FieldMap(m[112].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	sel, err := pb.FieldMap(start[1].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(sel, 1) != 12 || pb.IntField(sel, 2) != 2 {
		t.Fatalf("%+v", sel)
	}

	raw = encodeRecallClaim(7)
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m[117]; !ok || pb.IntField(m, 2) != 19 {
		t.Fatalf("%+v", m)
	}
}

func TestDecodeLotteryAndBrewBody(t *testing.T) {
	lot := pb.NewEncoder()
	lot.Int(1, 2)
	lot.Int(2, 5)
	head := pb.NewEncoder()
	head.Int(1, 6)
	data := pb.NewEncoder()
	data.Message(1, head.Bytes())
	data.Message(105, lot.Bytes())
	act, err := decodeActivityData(data.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if act.Lottery == nil || act.Lottery.FreeLeft != 2 || act.Lottery.FreeLimit != 5 {
		t.Fatalf("%+v", act.Lottery)
	}

	item := pb.NewEncoder()
	item.Int(1, 5001)
	item.Int(2, 1)
	rsp := pb.NewEncoder()
	rsp.Message(2, item.Bytes())
	top := pb.NewEncoder()
	top.Message(108, rsp.Bytes())
	out, err := decodeLotteryReply(top.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 1 || out.Items[0].ID != 5001 {
		t.Fatalf("%+v", out)
	}
}

func TestEncodeMegaEmptyBody(t *testing.T) {
	raw := encodeMegaClaim(12)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 12 || pb.IntField(m, 2) != 21 {
		t.Fatalf("head %+v", m)
	}
	if _, ok := m[119]; !ok {
		t.Fatal("empty mega body omitted")
	}
}

func TestDecodeActivityShopAndTech(t *testing.T) {
	item := pb.NewEncoder()
	item.Int(1, 80001)
	item.Int(2, 2)
	goods := pb.NewEncoder()
	goods.Int(1, 11)
	goods.Message(2, item.Bytes())
	goods.Int(4, 1)
	shop := pb.NewEncoder()
	shop.Message(1, goods.Bytes())
	head := pb.NewEncoder()
	head.Int(1, 5)
	head.String(4, "店")
	data := pb.NewEncoder()
	data.Message(1, head.Bytes())
	data.Message(102, shop.Bytes())
	act, err := decodeActivityData(data.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if act.ID != 5 || len(act.Shop) != 1 || act.Shop[0].ID != 11 || len(act.Shop[0].Items) != 1 {
		t.Fatalf("%+v", act)
	}

	node := pb.NewEncoder()
	node.Int(1, 3)
	node.Int(3, 2)
	tree := pb.NewEncoder()
	tree.Int(1, 1)
	tree.Message(2, node.Bytes())
	body := pb.NewEncoder()
	body.Message(1, tree.Bytes())
	nodes, err := decodeTechNodes(body.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 || nodes[0].ID != 3 || nodes[0].Status != 2 {
		t.Fatalf("%+v", nodes)
	}

	award := pb.NewEncoder()
	award.Int(1, 5001)
	award.Int(2, 1)
	rsp := pb.NewEncoder()
	rsp.Message(2, award.Bytes())
	rsp.PackedVarints(3, []int64{4, 5})
	top := pb.NewEncoder()
	top.Message(140, rsp.Bytes())
	out, err := decodeTechSubmit(top.Bytes())
	if err != nil || len(out.Items) != 1 || out.Items[0].ID != 5001 || len(out.Unlocked) != 2 || out.Unlocked[0] != 4 {
		t.Fatalf("tech submit %+v %v", out, err)
	}
}

func TestEncodeSigninAndProgress(t *testing.T) {
	raw := encodeProgress(9, 0)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 9 || pb.IntField(m, 2) != 25 {
		t.Fatalf("head %+v", m)
	}
	inner, err := pb.FieldMap(m[125].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if inner[1].Num != 1 {
		t.Fatalf("step omitted: %+v", inner)
	}

	raw = encodeSignin(3, 7)
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 2) != 4 {
		t.Fatalf("cmd %d", pb.IntField(m, 2))
	}
	claim, err := pb.FieldMap(m[103].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(claim, 1) != 7 {
		t.Fatalf("reward %d", pb.IntField(claim, 1))
	}
}

func TestDecodeMysteryAndSignin(t *testing.T) {
	card := pb.NewEncoder()
	card.Int(1, 11)
	card.Int(2, 20002)
	card.Int(6, 50)
	shop := pb.NewEncoder()
	shop.Bool(1, true)
	shop.Message(2, card.Bytes())
	s, err := decodeMysteryShop(shop.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !s.Present || len(s.Goods) != 1 || s.Goods[0].ID != 11 || s.Goods[0].Price != 50 {
		t.Fatalf("%+v", s)
	}

	reward := pb.NewEncoder()
	reward.Int(1, 7)
	signin := pb.NewEncoder()
	signin.Bool(1, true)
	signin.Message(2, reward.Bytes())
	data := pb.NewEncoder()
	head := pb.NewEncoder()
	head.Int(1, 88)
	head.String(4, "签到")
	data.Message(1, head.Bytes())
	data.Message(103, signin.Bytes())
	act, err := decodeActivityData(data.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if act.ID != 88 || !act.SigninClaimed || act.SigninRewardID != 7 {
		t.Fatalf("%+v", act)
	}
}

func TestDecodeRedPacketAndAlbum(t *testing.T) {
	info := pb.NewEncoder()
	info.Int(1, 3)
	info.Bool(3, true)
	st := pb.NewEncoder()
	st.Message(1, info.Bytes())
	ps, err := decodeRedPackets(st.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 || ps[0].ID != 3 || !ps[0].CanClaim {
		t.Fatalf("%+v", ps)
	}

	item := pb.NewEncoder()
	item.Int(1, 40061)
	item.Bool(3, true)
	item.Int(4, 2)
	album := pb.NewEncoder()
	album.Message(1, item.Bytes())
	album.Int(2, 8)
	album.Int(3, 4)
	album.Bool(8, true)
	a, err := decodeAlbum(album.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Items) != 1 || a.Items[0].FruitID != 40061 || a.Progress != 8 || a.Level != 4 || !a.Claimable {
		t.Fatalf("%+v", a)
	}
}

func TestEncodeShopAndSell(t *testing.T) {
	raw := encodeBuy(15, 2, 50)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 15 || pb.IntField(m, 2) != 2 || pb.IntField(m, 3) != 50 {
		t.Fatalf("%+v", m)
	}
	sell := encodeSell([]game.BagItem{{ID: 40061, Count: 3, UID: 9}})
	fields, err := pb.Walk(sell)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || fields[0].Num != 1 {
		t.Fatalf("%+v", fields)
	}
	item, err := pb.FieldMap(fields[0].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(item, 1) != 40061 || pb.IntField(item, 6) != 9 {
		t.Fatalf("%+v", item)
	}
}

func TestEncodeCheckCanOperateKeepsHostType(t *testing.T) {
	raw := encodeCheckCanOperate(9, 10004)
	m, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var hostType *uint64
	for _, f := range m {
		if f.Num == 3 {
			v := f.Varint
			hostType = &v
		}
	}
	if hostType == nil || *hostType != 0 {
		t.Fatalf("host_type %v", hostType)
	}
}

func TestDecodeWeatherAndTasks(t *testing.T) {
	st := pb.NewEncoder()
	st.Int(1, 3)
	st.Bool(5, true)
	st.Bool(8, true)
	wrap := pb.NewEncoder()
	wrap.Message(1, st.Bytes())
	w, err := decodeWeatherStatus(wrap.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if w.Type != 3 || !w.Active || !w.CanCollect {
		t.Fatalf("%+v", w)
	}

	task := pb.NewEncoder()
	task.Int(1, 11)
	task.Int(2, 2)
	task.Int(6, 5)
	task.String(9, "浇水")
	info := pb.NewEncoder()
	info.Message(2, task.Bytes())
	reply := pb.NewEncoder()
	reply.Message(1, info.Bytes())
	board, err := decodeTaskBoard(reply.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(board.Daily) != 1 || board.Daily[0].ID != 11 || board.Daily[0].Desc != "浇水" {
		t.Fatalf("%+v", board)
	}
}

func TestDecodeEmailAndSeason(t *testing.T) {
	item := pb.NewEncoder()
	item.String(1, "m1")
	item.String(3, "奖励")
	item.Bool(5, true)
	list := pb.NewEncoder()
	list.Message(1, item.Bytes())
	mails, err := decodeEmails(list.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(mails) != 1 || mails[0].ID != "m1" || !mails[0].HasReward {
		t.Fatalf("%+v", mails)
	}

	pass := pb.NewEncoder()
	pass.Int(1, 9)
	pass.Int(2, 4)
	pass.Int(6, 30)
	pass.String(16, "通行证")
	season := pb.NewEncoder()
	season.Int(1, 2)
	season.String(2, "春")
	season.Message(10, pass.Bytes())
	reply := pb.NewEncoder()
	reply.Message(1, season.Bytes())
	s, err := decodeSeason(reply.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != 2 || s.Name != "春" || s.Pass.Level != 4 || s.Pass.Name != "通行证" {
		t.Fatalf("%+v", s)
	}
}

func TestEncodeHarvest(t *testing.T) {
	raw := encodeHarvest([]int64{1, 2}, 8, true, 2)
	m, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var ids []int64
	for _, f := range m {
		if f.Num == 1 {
			ids = append(ids, int64(f.Varint))
		}
		if f.Num == 5 {
			t.Fatal("harvest has no reason")
		}
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("ids %v", ids)
	}
	fm, err := pb.FieldMap(raw)
	if err != nil || pb.IntField(fm, 2) != 8 {
		t.Fatalf("host %v %v", raw, err)
	}
}

func TestEncodeRefreshLandsAndCleanSocial(t *testing.T) {
	raw := encodeRefreshLands([]int64{3, 4}, 0, 0)
	fields, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var ids []int64
	for _, f := range fields {
		if f.Num == 1 {
			t.Fatal("own farm should omit host_gid")
		}
		if f.Num == 2 {
			ids = append(ids, int64(f.Varint))
		}
	}
	if len(ids) != 2 || ids[0] != 3 || ids[1] != 4 {
		t.Fatalf("ids %v", ids)
	}

	raw = encodeRefreshLands([]int64{1}, 9, 2)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 9 || pb.IntField(m, 4) != 2 {
		t.Fatalf("%+v", m)
	}

	raw = encodeCleanSocial([]int64{1, 2}, []int64{5005})
	fields, err = pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var lands, items []int64
	for _, f := range fields {
		if f.Num == 1 {
			lands = append(lands, int64(f.Varint))
		}
		if f.Num == 2 {
			items = append(items, int64(f.Varint))
		}
	}
	if len(lands) != 2 || items[0] != 5005 {
		t.Fatalf("lands %v items %v", lands, items)
	}
}

func TestEncodeBatchUseAndHeartbeat(t *testing.T) {
	raw := encodeBatchUse([]game.UseIn{{ID: 8, Count: 2, UID: 11}})
	fields, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || fields[0].Num != 1 {
		t.Fatalf("%+v", fields)
	}
	item, err := pb.FieldMap(fields[0].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(item, 1) != 8 || pb.IntField(item, 2) != 2 || pb.IntField(item, 6) != 11 {
		t.Fatalf("%+v", item)
	}

	raw = encodeHeartbeat(7, "1.13.3.11_20260826")
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 7 || pb.StringField(m, 2) != "1.13.3.11_20260826" {
		t.Fatalf("%+v", m)
	}

	reply := pb.NewEncoder()
	reply.Int(1, 1710000000000)
	hb, err := decodeHeartbeat(reply.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if hb.ServerMs != 1710000000000 {
		t.Fatalf("%+v", hb)
	}
}

func TestEncodePutSocialAndDogBuy(t *testing.T) {
	raw := encodePutSocial(9, []int64{1, 2}, 2)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 9 {
		t.Fatalf("%+v", m)
	}
	if _, ok := m[4]; ok {
		t.Fatal("put insects has no reason")
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var ids []int64
	for _, f := range fields {
		if f.Num == 2 {
			ids = append(ids, int64(f.Varint))
		}
	}
	if len(ids) != 2 {
		t.Fatalf("ids %v", ids)
	}
	raw = encodeDogBuy(3, 100)
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 3 || pb.IntField(m, 2) != 100 {
		t.Fatalf("%+v", m)
	}
}

func TestEncodeEmailIDsAndTaskReport(t *testing.T) {
	raw := encodeEmailIDs(2, []string{"a", "b"})
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 2 {
		t.Fatalf("box %+v", m)
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, f := range fields {
		if f.Num == 2 {
			ids = append(ids, string(f.Bytes))
		}
	}
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatalf("ids %v", ids)
	}
	raw = encodeTaskReport(9, 3)
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 9 || pb.IntField(m, 2) != 3 {
		t.Fatalf("%+v", m)
	}
	if encodeCareerGID() != nil {
		t.Fatal("career request is empty")
	}
	raw = encodeVisitPage(0)
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 1 {
		t.Fatalf("page %+v", m)
	}
	ok, err := decodeOK(nil)
	if err != nil || !ok {
		t.Fatalf("empty ok %v %v", ok, err)
	}
	flag := pb.NewEncoder()
	flag.Bool(1, true)
	dot, err := decodeRedDot(flag.Bytes())
	if err != nil || !dot.Red {
		t.Fatalf("red %+v %v", dot, err)
	}
	popup := pb.NewEncoder()
	popup.Bool(1, true)
	popup.Int(2, 10)
	popup.Bool(3, true)
	info, err := decodeVisitPopup(popup.Bytes())
	if err != nil || !info.Unread || info.LastRead != 10 || !info.NeedPopup {
		t.Fatalf("popup %+v", info)
	}
	prof := pb.NewEncoder()
	item := pb.NewEncoder()
	item.Int(1, 2)
	item.Int(2, 3)
	prof.Message(1, item.Bytes())
	list, err := decodeMallProfiles(prof.Bytes())
	if err != nil || len(list) != 1 || list[0].ID != 2 || list[0].Type != 3 {
		t.Fatalf("profiles %+v", list)
	}
}

func TestDecodeGetGroup(t *testing.T) {
	head := pb.NewEncoder()
	head.Int(1, 80)
	head.String(4, "气象")
	childHead := pb.NewEncoder()
	childHead.Int(1, 81)
	childHead.Int(2, 80)
	childHead.String(4, "商店")
	child := pb.NewEncoder()
	child.Message(1, childHead.Bytes())
	group := pb.NewEncoder()
	group.Message(1, head.Bytes())
	group.Message(2, child.Bytes())
	reply := pb.NewEncoder()
	reply.Message(1, group.Bytes())
	got, err := decodeGetGroup(reply.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 80 || got.Name != "气象" || len(got.Activities) != 1 || got.Activities[0].ID != 81 || got.Activities[0].GroupID != 80 {
		t.Fatalf("%+v", got)
	}
}

func TestDecodeLandEmailAndNewRPCs(t *testing.T) {
	phase := pb.NewEncoder()
	phase.Int(1, 6)
	phase.Int(2, 100)
	plant := pb.NewEncoder()
	plant.Int(1, 1020003)
	plant.String(2, "胡萝卜")
	plant.Message(4, phase.Bytes())
	plant.Int(10, 40003)
	plant.Int(11, 20)
	plant.Varint(12, 9)
	plant.Varint(13, 8)
	plant.Varint(14, 7)
	plant.Int(15, 3600)
	plant.Bool(16, true)
	plant.Int(18, 5)
	plant.Varint(20, 33)
	unlock := pb.NewEncoder()
	unlock.Int(1, 2)
	unlock.Int(2, 8)
	itemNeed := pb.NewEncoder()
	itemNeed.Int(1, 1001)
	itemNeed.Int(2, 500)
	unlock.Message(3, itemNeed.Bytes())
	plant.Int(28, 1001)
	plant.Int(29, 12)
	land := pb.NewEncoder()
	land.Int(1, 3)
	land.Bool(2, true)
	land.Int(3, 2)
	land.Int(4, 5)
	land.Message(6, unlock.Bytes())
	land.Message(10, plant.Bytes())
	land.Int(16, 4)
	got, err := decodeLand(land.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 3 || got.MaxLevel != 5 || got.FruitNum != 20 || !got.HasWeed || !got.HasInsect || got.GrowSec != 3600 || got.Phase != 6 || !got.Mutant || got.LandsLevel != 4 || got.SourcePrice != 12 || got.Unlock == nil || got.Unlock.LandID != 2 || len(got.Unlock.Items) != 1 {
		t.Fatalf("%+v", got)
	}

	item := pb.NewEncoder()
	item.String(1, "m1")
	item.String(3, "奖励")
	item.Bool(4, true)
	item.Bool(5, true)
	mail, err := decodeEmailItem(item.Bytes())
	if err != nil || !mail.Read || !mail.HasReward {
		t.Fatalf("%+v %v", mail, err)
	}

	info := pb.NewEncoder()
	info.Int(1, 2)
	info.String(4, "雨")
	cur := pb.NewEncoder()
	cur.Message(1, info.Bytes())
	w, err := decodeCurrentWeather(cur.Bytes())
	if err != nil || w.Type != 2 || w.Name != "雨" {
		t.Fatalf("%+v %v", w, err)
	}

	person := pb.NewEncoder()
	person.Int(1, 9)
	person.String(3, "客")
	rsp := pb.NewEncoder()
	rsp.Message(1, person.Bytes())
	op := pb.NewEncoder()
	op.Message(123, rsp.Bytes())
	invitees, err := decodeInvitees(op.Bytes())
	if err != nil || len(invitees) != 1 || invitees[0].GID != 9 {
		t.Fatalf("%+v %v", invitees, err)
	}

	stat := pb.NewEncoder()
	stat.Int(1, 9)
	stat.String(2, "友")
	stat.Int(4, 3)
	sum := pb.NewEncoder()
	sum.Int(1, 4)
	sum.Int(2, 5)
	sum.Int(3, 1)
	sum.Message(4, stat.Bytes())
	reply := pb.NewEncoder()
	reply.Message(1, sum.Bytes())
	summary, err := decodeVisitSummary(reply.Bytes())
	if err != nil || summary.StealCount != 4 || summary.HelpCount != 5 || len(summary.Friends) != 1 {
		t.Fatalf("%+v %v", summary, err)
	}

	raw := encodeInvitees(8)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 8 || pb.IntField(m, 2) != 24 {
		t.Fatalf("%+v", m)
	}
	if _, ok := m[122]; !ok {
		t.Fatal("missing invitees field")
	}

	if visitReason(7, 7) != 0 || visitReason(7, 9) != 2 || visitReason(7, 0) != 0 {
		t.Fatal("visitReason")
	}
}

func TestDecodeShareClaimPrefersRewards(t *testing.T) {
	reply := pb.NewEncoder()
	reply.Bool(1, true)
	reply.Bool(2, true)
	item := pb.NewEncoder()
	item.Int(1, 20002)
	item.Int(2, 3)
	reply.Message(3, item.Bytes())
	items, err := decodeShareClaim(reply.Bytes())
	if err != nil || len(items) != 1 || items[0].ID != 20002 || items[0].Count != 3 {
		t.Fatalf("%+v %v", items, err)
	}

	out, err := decodeShareClaimOut(reply.Bytes())
	if err != nil || !out.Success || !out.HasReward || len(out.Items) != 1 {
		t.Fatalf("out %+v %v", out, err)
	}
}

func TestDecodeActiveBoxAndTask(t *testing.T) {
	wait := pb.NewEncoder()
	wait.Int(1, 10)
	wait.Int(2, 20)
	wait.Int(3, 1)
	box, err := decodeActiveBox(wait.Bytes())
	if err != nil || box.ID != 10 || !box.CanClaim || box.Claimed {
		t.Fatalf("wait %+v %v", box, err)
	}

	done := pb.NewEncoder()
	done.Int(1, 11)
	done.Int(3, 2)
	box, err = decodeActiveBox(done.Bytes())
	if err != nil || !box.Claimed || box.CanClaim {
		t.Fatalf("done %+v %v", box, err)
	}

	task := pb.NewEncoder()
	task.Int(1, 8)
	task.Int(2, 1)
	task.Bool(4, true)
	got, err := decodeTask(task.Bytes())
	if err != nil || got.ID != 8 || !got.Unlocked || got.Claimed {
		t.Fatalf("task %+v %v", got, err)
	}
}

func TestDecodeCanOperateAndHarvest(t *testing.T) {
	reply := pb.NewEncoder()
	reply.Bool(1, true)
	reply.Int(2, 4)
	out, err := decodeCanOperate(reply.Bytes())
	if err != nil || !out.OK || out.StealNum != 4 {
		t.Fatalf("%+v %v", out, err)
	}

	item := pb.NewEncoder()
	item.Int(1, 20002)
	item.Int(2, 5)
	lost := pb.NewEncoder()
	lost.Int(1, 20002)
	lost.Int(2, 1)
	land := pb.NewEncoder()
	land.Int(1, 3)
	land.Bool(2, true)
	harvest := pb.NewEncoder()
	harvest.Message(1, land.Bytes())
	harvest.Message(2, item.Bytes())
	harvest.Message(3, lost.Bytes())
	got, err := decodeHarvest(harvest.Bytes())
	if err != nil || len(got.Items) != 1 || got.Items[0].Count != 5 || len(got.Lost) != 1 || len(got.Lands) != 1 {
		t.Fatalf("%+v %v", got, err)
	}

	lim := pb.NewEncoder()
	lim.Int(1, 10001)
	lim.Int(2, 3)
	lim.Int(3, 20)
	rew := pb.NewEncoder()
	rew.Int(1, 1001)
	rew.Int(2, 2)
	drop := pb.NewEncoder()
	drop.Int(1, 3)
	drop.Message(2, rew.Bytes())
	harvest.Message(4, lim.Bytes())
	harvest.Message(7, drop.Bytes())
	got, err = decodeHarvest(harvest.Bytes())
	if err != nil || len(got.Limits) != 1 || got.Limits[0].Used != 3 || len(got.Drops) != 1 || got.Drops[0].Rewards[0].Count != 2 {
		t.Fatalf("limits/drops %+v %v", got, err)
	}
}

func TestDecodeRandShopAndEncodeShopOps(t *testing.T) {
	goods := pb.NewEncoder()
	goods.Int(1, 22)
	goods.String(2, "种")
	it := pb.NewEncoder()
	it.Int(1, 20002)
	it.Int(2, 1)
	goods.Message(4, it.Bytes())
	goods.Bool(8, true)
	shop := pb.NewEncoder()
	shop.Message(1, goods.Bytes())
	shop.Int(2, 99)
	shop.Int(6, 5)
	got, err := decodeRandShop(shop.Bytes())
	if err != nil || len(got.Goods) != 1 || got.Goods[0].ID != 22 || !got.Goods[0].Available || got.Next != 99 || got.Limit != 5 {
		t.Fatalf("%+v %v", got, err)
	}

	raw := encodeRandBuy(8, 22, 0)
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 8 || pb.IntField(m, 2) != 2 {
		t.Fatalf("rand buy %+v", m)
	}
	inner, err := pb.FieldMap(m[102].Bytes)
	if err != nil || pb.IntField(inner, 1) != 22 || pb.IntField(inner, 2) != 1 {
		t.Fatalf("rand buy inner %+v %v", inner, err)
	}

	raw = encodeRandRefresh(8)
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 8 || pb.IntField(m, 2) != 3 {
		t.Fatalf("rand refresh %+v", m)
	}
	if _, ok := m[102]; ok {
		t.Fatal("refresh should have no payload")
	}

	raw = encodeShopBatch(8, []game.ShopBuyItem{{GoodsID: 3, Count: 2}})
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 1) != 8 || pb.IntField(m, 2) != 34 {
		t.Fatalf("batch %+v", m)
	}
	batch, err := pb.FieldMap(m[133].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	sel, err := pb.FieldMap(batch[1].Bytes)
	if err != nil || pb.IntField(sel, 1) != 3 || pb.IntField(sel, 2) != 2 {
		t.Fatalf("batch item %+v %v", sel, err)
	}

	raw = encodeUse(game.UseIn{ID: 5005, Count: 1, UID: 11, HostGID: 9})
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	item, err := pb.FieldMap(m[1].Bytes)
	if err != nil || pb.IntField(item, 1) != 5005 || pb.IntField(item, 2) != 1 || pb.IntField(item, 6) != 11 {
		t.Fatalf("use item %+v %v", item, err)
	}
	target, err := pb.FieldMap(m[2].Bytes)
	if err != nil || pb.IntField(target, 1) != 9 {
		t.Fatalf("use target %+v %v", target, err)
	}
	if ids := pb.PackedInts(m[2].Bytes, 2); len(ids) != 0 {
		t.Fatalf("frog has no land_ids %v", ids)
	}

	raw = encodeUse(game.UseIn{ID: 5006, Count: 1, HostGID: 9, LandIDs: []int64{1, 2}})
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	target, err = pb.FieldMap(m[2].Bytes)
	if err != nil {
		t.Fatal(err)
	}
	ids := pb.PackedInts(m[2].Bytes, 2)
	if pb.IntField(target, 1) != 9 || len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("cloud target %+v ids %v", target, ids)
	}

	raw = encodeUse(game.UseIn{ID: 8, Count: 2, UID: 11})
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m[2]; ok {
		t.Fatal("warehouse use has no target")
	}

	raw = encodeReportInvite("oid", "key")
	m, err = pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.StringField(m, 1) != "oid" || pb.StringField(m, 2) != "key" {
		t.Fatalf("invite %+v", m)
	}
}

func TestDecodeActivityBodiesAndOperateOut(t *testing.T) {
	rew := pb.NewEncoder()
	rew.Int(1, 3)
	rew.Bool(4, true)
	item := pb.NewEncoder()
	item.Int(1, 20002)
	item.Int(2, 1)
	rew.Message(5, item.Bytes())
	mega := pb.NewEncoder()
	mega.Int(1, 8)
	mega.Int(2, 14)
	mega.Message(4, rew.Bytes())
	gotMega, err := decodeMegaEvent(mega.Bytes())
	if err != nil || gotMega.Day != 8 || gotMega.Total != 14 || len(gotMega.Rewards) != 1 || !gotMega.Rewards[0].Claimable {
		t.Fatalf("mega %+v %v", gotMega, err)
	}

	step := pb.NewEncoder()
	step.Int(1, 0)
	step.Int(4, 1)
	step.Message(3, item.Bytes())
	prog := pb.NewEncoder()
	prog.Message(2, step.Bytes())
	prog.Int(3, 1)
	gotProg, err := decodeProgressReward(prog.Bytes())
	if err != nil || gotProg.Unlocked != 1 || len(gotProg.Steps) != 1 || gotProg.Steps[0].Status != 1 {
		t.Fatalf("progress %+v %v", gotProg, err)
	}

	signin := pb.NewEncoder()
	signin.Bool(1, true)
	day := pb.NewEncoder()
	day.Int(1, 9)
	day.String(2, "第1天")
	day.Message(3, item.Bytes())
	signin.Message(2, day.Bytes())
	body := pb.NewEncoder()
	head := pb.NewEncoder()
	head.Int(1, 88)
	head.String(4, "签到")
	body.Message(1, head.Bytes())
	body.Message(103, signin.Bytes())
	act, err := decodeActivityData(body.Bytes())
	if err != nil || act.ID != 88 || !act.SigninClaimed || act.SigninRewardID != 9 || len(act.Signin) != 1 {
		t.Fatalf("signin %+v %v", act, err)
	}

	award := pb.NewEncoder()
	award.Int(1, 20002)
	award.Int(2, 2)
	cost := pb.NewEncoder()
	cost.Int(1, 7)
	cost.Int(2, 1)
	shop := pb.NewEncoder()
	shop.Message(1, award.Bytes())
	shop.Message(2, cost.Bytes())
	reply := pb.NewEncoder()
	reply.Message(101, shop.Bytes())
	reply.Message(3, body.Bytes())
	out, err := decodeOperateOut(reply.Bytes(), 101, 1, 2)
	if err != nil || len(out.Items) != 1 || out.Items[0].Count != 2 || len(out.Costs) != 1 || out.Activity == nil || out.Activity.ID != 88 {
		t.Fatalf("operate %+v %v", out, err)
	}

	days := pb.NewEncoder()
	days.PackedVarints(1, []int64{1, 2})
	days.Message(2, award.Bytes())
	megaReply := pb.NewEncoder()
	megaReply.Message(120, days.Bytes())
	megaOut, err := decodeMegaClaim(megaReply.Bytes())
	if err != nil || len(megaOut.Days) != 2 || len(megaOut.Items) != 1 {
		t.Fatalf("mega claim %+v %v", megaOut, err)
	}
}

func TestDecodeLandOpItemOpFriendsAlbum(t *testing.T) {
	land := pb.NewEncoder()
	land.Int(1, 2)
	land.Bool(2, true)
	lim := pb.NewEncoder()
	lim.Int(1, 10001)
	lim.Int(2, 1)
	lim.Int(3, 20)
	rew := pb.NewEncoder()
	rew.Int(1, 1001)
	rew.Int(2, 3)
	drop := pb.NewEncoder()
	drop.Int(1, 2)
	drop.Message(2, rew.Bytes())
	evt := pb.NewEncoder()
	evt.Int(1, 8001)
	evt.Int(2, 9)
	reply := pb.NewEncoder()
	reply.Message(1, land.Bytes())
	reply.Message(2, lim.Bytes())
	reply.Message(3, drop.Bytes())
	op, err := decodeLandOp(reply.Bytes())
	if err != nil || len(op.Lands) != 1 || op.Lands[0].ID != 2 || len(op.Limits) != 1 || len(op.Drops) != 1 {
		t.Fatalf("landop %+v %v", op, err)
	}

	useLand := pb.NewEncoder()
	useLand.Message(4, land.Bytes())
	useLand.Message(5, drop.Bytes())
	farmRew := pb.NewEncoder()
	farmRew.Int(1, 5005)
	useLand.Message(6, farmRew.Bytes())
	uo, err := decodeUseLandOp(useLand.Bytes())
	if err != nil || len(uo.Lands) != 1 || uo.Lands[0].ID != 2 || len(uo.Drops) != 1 || len(uo.Rewards) != 1 || uo.Rewards[0].SourceItemID != 5005 {
		t.Fatalf("use landop %+v %v", uo, err)
	}

	fert := pb.NewEncoder()
	fert.Int(1, 30001)
	fert.Int(2, 1)
	fr := pb.NewEncoder()
	fr.Message(1, land.Bytes())
	fr.Message(2, lim.Bytes())
	fr.Message(3, fert.Bytes())
	fr.Message(4, drop.Bytes())
	fo, err := decodeFertilizeOut(fr.Bytes())
	if err != nil || len(fo.Costs) != 1 || fo.Costs[0].ID != 30001 || len(fo.Drops) != 1 {
		t.Fatalf("fertilize %+v %v", fo, err)
	}

	all := pb.NewEncoder()
	all.Message(1, land.Bytes())
	all.Message(2, lim.Bytes())
	all.Message(3, evt.Bytes())
	ro, err := decodeRefreshOut(all.Bytes())
	if err != nil || len(ro.Events) != 1 || ro.Events[0].ItemID != 8001 || len(ro.Drops) != 0 {
		t.Fatalf("refresh %+v %v", ro, err)
	}

	used := pb.NewEncoder()
	used.Int(1, 20002)
	used.Int(2, 1)
	got := pb.NewEncoder()
	got.Int(1, 1001)
	got.Int(2, 50)
	comp := pb.NewEncoder()
	comp.Int(1, 1001)
	comp.Int(2, 10)
	use := pb.NewEncoder()
	use.Message(1, used.Bytes())
	use.Message(2, got.Bytes())
	use.Message(3, comp.Bytes())
	io, err := decodeItemOp(use.Bytes(), 1, 2, 3, true)
	if err != nil || len(io.Used) != 1 || io.Used[0].ID != 20002 || len(io.Items) != 1 || io.Items[0].Count != 50 || len(io.Compensated) != 1 {
		t.Fatalf("use %+v %v", io, err)
	}

	frd := pb.NewEncoder()
	frd.Int(1, 5)
	frd.String(3, "友")
	blk := pb.NewEncoder()
	blk.Int(1, 8)
	blk.String(2, "黑")
	list := pb.NewEncoder()
	list.Message(1, frd.Bytes())
	list.Int(3, 2)
	list.Message(4, blk.Bytes())
	fo2, err := decodeFriendsOut(list.Bytes())
	if err != nil || len(fo2.Friends) != 1 || fo2.ApplicationCount != 2 || len(fo2.Blocked) != 1 {
		t.Fatalf("friends %+v %v", fo2, err)
	}

	aitem := pb.NewEncoder()
	aitem.Int(1, 40061)
	aitem.Int(5, 2)
	rew2 := pb.NewEncoder()
	rew2.Int(1, 1001)
	rew2.Int(2, 1)
	aitem.Message(6, rew2.Bytes())
	al := pb.NewEncoder()
	al.Message(1, aitem.Bytes())
	al.Int(2, 8)
	al.Int(3, 4)
	al.Int(6, 1)
	al.Int(7, 12)
	al.Bool(8, true)
	a, err := decodeAlbum(al.Bytes())
	if err != nil || a.Level != 4 || a.Type != 1 || a.Next != 12 || !a.Claimable || a.Items[0].Layer != 2 || len(a.Items[0].Rewards) != 1 {
		t.Fatalf("album %+v %v", a, err)
	}
}

func TestEncodeClaimDailyAndDecodeMallMutantVisit(t *testing.T) {
	raw := encodeClaimDaily(2, []int64{10, 20})
	if got := pb.PackedInts(raw, 1); len(got) != 2 || got[0] != 10 || got[1] != 20 {
		t.Fatalf("ids %v", got)
	}
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m[2]; ok {
		t.Fatal("claim daily has no type field")
	}

	rew := pb.NewEncoder()
	rew.Int(1, 1002)
	rew.Int(2, 120)
	mall := pb.NewEncoder()
	mall.Message(2, rew.Bytes())
	lim := pb.NewEncoder()
	lim.Int(2, 1)
	lim.Int(3, 1)
	mall.Message(4, lim.Bytes())
	bo, err := decodeMallBuy(mall.Bytes())
	if err != nil || len(bo.Items) != 1 || bo.Items[0].Count != 120 || bo.Bought != 1 || bo.Limit != 1 {
		t.Fatalf("mall %+v %v", bo, err)
	}

	ent := pb.NewEncoder()
	ent.Int(1, 7)
	ent.Int(2, 100)
	ent.Int(3, 200)
	ent.Int(4, 88)
	ent.Bool(5, true)
	book := pb.NewEncoder()
	book.Message(1, ent.Bytes())
	ms, err := decodeMutants(book.Bytes())
	if err != nil || len(ms) != 1 || ms[0].ID != 7 || ms[0].ActivityID != 88 || !ms[0].Red {
		t.Fatalf("mutant %+v %v", ms, err)
	}

	rel := pb.NewEncoder()
	rel.Bool(2, true)
	w := pb.NewEncoder()
	w.Int(1, 2)
	w.Bool(5, true)
	w.Bool(8, true)
	ol := pb.NewEncoder()
	ol.Int(1, 10004)
	ol.Int(2, 1)
	ol.Int(3, 5)
	enter := pb.NewEncoder()
	enter.Bool(6, true)
	enter.Message(7, rel.Bytes())
	enter.Int(8, 123)
	enter.Message(11, ol.Bytes())
	enter.Message(13, w.Bytes())
	visit, err := decodeVisit(enter.Bytes())
	if err != nil || !visit.AtHome || !visit.Friend || visit.ServerMs != 123 || len(visit.Limits) != 1 || !visit.Weather.CanCollect {
		t.Fatalf("visit %+v %v", visit, err)
	}

	item := pb.NewEncoder()
	item.Int(1, 20002)
	item.Int(2, 1)
	info := pb.NewEncoder()
	info.Message(1, item.Bytes())
	claim := pb.NewEncoder()
	claim.Message(1, item.Bytes())
	claim.Message(2, item.Bytes())
	tc, err := decodeTaskClaim(claim.Bytes())
	if err != nil || len(tc.Items) != 1 || len(tc.Compensated) != 1 {
		t.Fatalf("task claim %+v %v", tc, err)
	}

	sc := pb.NewEncoder()
	sc.Bool(1, true)
	sc.Bool(2, true)
	sc.Message(3, item.Bytes())
	so, err := decodeShareClaimOut(sc.Bytes())
	if err != nil || !so.Success || !so.HasReward || len(so.Items) != 1 || so.Items[0].ID != 20002 {
		t.Fatalf("share %+v %v", so, err)
	}

	bagItem := pb.NewEncoder()
	bagItem.Int(1, 40061)
	bagItem.Int(2, 2)
	bagItem.Int(3, 99)
	bagItem.Bool(7, true)
	bagItem.Bool(9, true)
	bag := pb.NewEncoder()
	bag.Message(1, bagItem.Bytes())
	bag.Int(2, 200)
	bag.Int(3, 12)
	reply := pb.NewEncoder()
	reply.Message(1, bag.Bytes())
	bo2, err := decodeBagOut(reply.Bytes())
	if err != nil || len(bo2.Items) != 1 || bo2.Items[0].ID != 40061 || bo2.Items[0].Count != 2 {
		t.Fatalf("bag %+v %v", bo2, err)
	}

	cond := pb.NewEncoder()
	cond.Int(1, 2)
	cond.Int(2, 10)
	goods := pb.NewEncoder()
	goods.Int(1, 15)
	goods.Int(6, 20002)
	goods.Message(8, cond.Bytes())
	shop := pb.NewEncoder()
	shop.Message(1, goods.Bytes())
	shop.Message(2, item.Bytes())
	sb, err := decodeShopBuy(shop.Bytes())
	if err != nil || sb.Goods == nil || sb.Goods.ID != 15 || len(sb.Goods.Conds) != 1 || sb.Goods.Conds[0].Param != 10 || len(sb.Items) != 1 {
		t.Fatalf("shop buy %+v %v", sb, err)
	}

	log := pb.NewEncoder()
	log.Int(1, 1)
	log.Int(2, 5)
	log.Int(3, 9)
	log.Int(9, 2)
	vl, err := decodeVisitLog(log.Bytes())
	if err != nil || vl.Action != 5 || vl.FromType != 2 {
		t.Fatalf("visit log %+v %v", vl, err)
	}

	ec := pb.NewEncoder()
	ec.Message(1, item.Bytes())
	ec.String(2, "m1")
	eo, err := decodeEmailClaim(ec.Bytes())
	if err != nil || len(eo.Items) != 1 || len(eo.Unclaimed) != 1 || eo.Unclaimed[0] != "m1" {
		t.Fatalf("email %+v %v", eo, err)
	}
}

func TestDecodeAchieveGoalAndDogAndPass(t *testing.T) {
	goal := pb.NewEncoder()
	goal.Int(1, 50)
	goal.Int(2, 3)
	goal.Int(3, 4)
	scope := pb.NewEncoder()
	scope.Int(1, 1)
	scope.Int(2, 9)
	scope.Int(3, 4)
	goal.Message(4, scope.Bytes())
	out, err := decodeAchieveGoalOut(goal.Bytes())
	if err != nil || out.Exp != 50 || out.Before != 3 || out.After != 4 || out.Scope.ID != 9 {
		t.Fatalf("goal %+v %v", out, err)
	}

	item := pb.NewEncoder()
	item.Int(1, 1001)
	item.Int(2, 2)
	gifts := pb.NewEncoder()
	gifts.Message(1, item.Bytes())
	gifts.Int(3, 2)
	gifts.Int(4, 1)
	dg, err := decodeDogGifts(gifts.Bytes())
	if err != nil || len(dg.Items) != 1 || dg.Claimed != 2 || dg.Pending != 1 {
		t.Fatalf("gifts %+v %v", dg, err)
	}

	log := pb.NewEncoder()
	log.Int(1, 12)
	log.String(2, "贼")
	log.Bool(4, true)
	log.Int(10, 8)
	log.Int(16, 3)
	logs := pb.NewEncoder()
	logs.Message(1, log.Bytes())
	logs.Int(2, 7)
	lo, err := decodeProtectLogs(logs.Bytes())
	if err != nil || lo.Total != 7 || len(lo.Logs) != 1 || !lo.Logs[0].Online || lo.Logs[0].SkillID != 3 {
		t.Fatalf("logs %+v %v", lo, err)
	}

	lv := pb.NewEncoder()
	lv.Int(1, 5)
	lv.Message(2, item.Bytes())
	lv.Bool(4, true)
	pass := pb.NewEncoder()
	pass.Int(2, 6)
	pass.Message(8, lv.Bytes())
	pass.Int(15, 180)
	pass.String(17, "战令")
	claim := pb.NewEncoder()
	claim.Message(1, item.Bytes())
	claim.Int(2, 5)
	claim.Message(3, pass.Bytes())
	claim.Bool(4, true)
	pc, err := decodePassClaim(claim.Bytes())
	if err != nil || !pc.Overflow || pc.Pass.Price != 180 || pc.Pass.Desc != "战令" || len(pc.Pass.Levels) != 1 || pc.Pass.Levels[0].Level != 5 {
		t.Fatalf("pass %+v %v", pc, err)
	}

	app := pb.NewEncoder()
	app.Int(1, 3)
	app.String(4, "甲")
	apps := pb.NewEncoder()
	apps.Message(1, app.Bytes())
	apps.Bool(2, true)
	ao, err := decodeApplications(apps.Bytes())
	if err != nil || !ao.Blocked || len(ao.Applications) != 1 || ao.Applications[0].GID != 3 || ao.Applications[0].Name != "甲" {
		t.Fatalf("apps %+v %v", ao, err)
	}

	info := pb.NewEncoder()
	info.Int(1, 3)
	info.Bool(2, true)
	info.Int(3, 1)
	rp, err := decodeRedPacket(info.Bytes())
	if err != nil || !rp.Claimed || rp.Status != 1 || rp.CanClaim {
		t.Fatalf("claimed packet %+v %v", rp, err)
	}

	if encodeSetTags() != nil {
		t.Fatal("set tags request is empty")
	}

	warn := pb.NewEncoder()
	warn.Int(1, 3)
	warn.String(2, "满了")
	hv := pb.NewEncoder()
	hv.Message(5, warn.Bytes())
	ho, err := decodeHarvest(hv.Bytes())
	if err != nil || len(ho.Warnings) != 1 || ho.Warnings[0].LandID != 3 {
		t.Fatalf("harvest %+v %v", ho, err)
	}

	basic := pb.NewEncoder()
	basic.Int(1, 1)
	basic.Int(2, 9)
	basic.String(3, "夏")
	solar := pb.NewEncoder()
	solar.Int(2, 100)
	solar.Message(3, basic.Bytes())
	so, err := decodeSolarTerms(solar.Bytes())
	if err != nil || so.ServerTime != 100 || so.Basic == nil || so.Basic.Season != 9 {
		t.Fatalf("solar %+v %v", so, err)
	}
}

func TestEncodeFarmingAndShareAndMonthCard(t *testing.T) {
	raw := encodeFarming([]int64{1, 2}, 8, []int64{5005})
	m, err := pb.FieldMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	if pb.IntField(m, 2) != 8 {
		t.Fatalf("host %+v", m)
	}
	fields, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var lands, items []int64
	for _, f := range fields {
		switch f.Num {
		case 1:
			lands = append(lands, int64(f.Varint))
		case 3, 4:
			t.Fatalf("farming should not write host_type/reason %d", f.Num)
		case 5:
			items = append(items, int64(f.Varint))
		}
	}
	if len(lands) != 2 || lands[0] != 1 || items[0] != 5005 {
		t.Fatalf("farming %v %v", lands, items)
	}

	raw = encodeFarming(nil, 8, []int64{5005})
	fields, err = pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var clear []int64
	for _, f := range fields {
		switch f.Num {
		case 1:
			t.Fatal("clear frog has no land_ids")
		case 3, 4:
			t.Fatalf("farming should not write host_type/reason %d", f.Num)
		case 5:
			clear = append(clear, int64(f.Varint))
		}
	}
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 8 || len(clear) != 1 || clear[0] != 5005 {
		t.Fatalf("clear frog farming %+v %v", m, err)
	}

	raw = encodeEnter(9, 2)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 9 || pb.IntField(m, 2) != 2 {
		t.Fatalf("enter %+v %v", m, err)
	}
	raw = encodeEnter(9, 0)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 9 {
		t.Fatalf("enter self %+v %v", m, err)
	}
	if _, ok := m[2]; ok {
		t.Fatal("own farm enter has no reason")
	}

	raw = encodeGIDs([]int64{7, 9})
	fields, err = pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var gids []int64
	for _, f := range fields {
		if f.Num == 1 {
			gids = append(gids, int64(f.Varint))
		}
	}
	if len(gids) != 2 || gids[0] != 7 || gids[1] != 9 {
		t.Fatalf("gids %v", gids)
	}

	info := pb.NewEncoder()
	info.Int(1, 11)
	info.String(2, "k")
	info.String(3, "https://s")
	info.Int(4, 2)
	info.Int(5, 1)
	info.Bool(6, true)
	list := pb.NewEncoder()
	list.Message(1, info.Bytes())
	inv, err := decodeShareInvite(list.Bytes())
	if err != nil || len(inv.Infos) != 1 || inv.Infos[0].ID != 11 || inv.Infos[0].ShareKey != "k" || !inv.Infos[0].CanClaim {
		t.Fatalf("invite %+v %v", inv, err)
	}

	item := pb.NewEncoder()
	item.Int(1, 1002)
	item.Int(2, 10)
	award := pb.NewEncoder()
	award.Message(1, info.Bytes())
	award.Message(2, item.Bytes())
	award.Bool(4, true)
	ao, err := decodeShareAwardOut(award.Bytes())
	if err != nil || !ao.Awarded || ao.Info.ID != 11 || len(ao.Awards) != 1 {
		t.Fatalf("award %+v %v", ao, err)
	}

	cost := pb.NewEncoder()
	cost.Int(1, 1004)
	cost.Int(2, 180)
	card := pb.NewEncoder()
	card.Int(1, 2001)
	card.Int(3, 1)
	card.Int(4, 29)
	card.Int(6, 3600)
	card.Int(8, 180)
	card.Int(9, 30)
	card.String(10, "2001")
	card.Message(11, cost.Bytes())
	card.Int(12, 1)
	mc, err := decodeMonthCard(card.Bytes())
	if err != nil || mc.ID != 2001 || !mc.Claimable || mc.Days != 29 || mc.ExpireSeconds != 3600 || mc.TotalCount != 180 || mc.PurchaseCost == nil || mc.PurchaseCost.ID != 1004 || !mc.Claimable2 {
		t.Fatalf("card %+v %v", mc, err)
	}

	ref := encodeEmailRef(2, "m9")
	m, err = pb.FieldMap(ref)
	if err != nil || pb.IntField(m, 1) != 2 || pb.StringField(m, 2) != "m9" {
		t.Fatalf("email ref %+v %v", m, err)
	}
}

func TestDecodeV59DogSolarAchieve(t *testing.T) {
	dog := pb.NewEncoder()
	dog.Int(1, 3)
	dog.String(2, "旺财")
	dog.Int(8, 180)
	d, err := decodeDog(dog.Bytes())
	if err != nil || d.ID != 3 || d.Name != "旺财" || d.Price != 180 {
		t.Fatalf("dog %+v %v", d, err)
	}

	food := pb.NewEncoder()
	food.Int(1, 101)
	food.Int(2, 3600)
	food.Int(3, 5)
	f, err := decodeDogFood(food.Bytes())
	if err != nil || f.ID != 101 || f.Seconds != 3600 || f.Count != 5 {
		t.Fatalf("food %+v %v", f, err)
	}

	wd := pb.NewEncoder()
	wd.Int(1, 7)
	wo, err := decodeWithdraw(wd.Bytes())
	if err != nil || wo.Withdrawn != 7 || wo.Previous != 0 {
		t.Fatalf("withdraw %+v %v", wo, err)
	}

	term := pb.NewEncoder()
	term.Int(1, 12)
	term.Int(2, 2)
	term.String(6, "立秋")
	s, err := decodeSolarTerm(term.Bytes())
	if err != nil || s.ID != 12 || s.Status != 2 || s.Name != "立秋" {
		t.Fatalf("solar %+v %v", s, err)
	}

	goal := pb.NewEncoder()
	goal.Int(1, 9)
	goal.Int(9, 4)
	goal.String(11, "种")
	scope := pb.NewEncoder()
	scope.Int(1, 1)
	scope.Int(2, 2)
	scope.Message(6, goal.Bytes())
	sv, err := decodeScopeView(scope.Bytes())
	if err != nil || sv.ID != 2 || len(sv.Goals) != 1 || sv.Goals[0].Sort != 4 || sv.Goals[0].Desc != "种" {
		t.Fatalf("scope %+v %v", sv, err)
	}
}

func TestEncodeV59ActivityOperate(t *testing.T) {
	raw := encodeDraw(8, 2)
	m, err := pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 8 || pb.IntField(m, 2) != 5 {
		t.Fatalf("draw head %+v %v", m, err)
	}
	inner, err := pb.FieldMap(m[104].Bytes)
	if err != nil || pb.IntField(inner, 1) != 2 {
		t.Fatalf("draw inner %+v %v", inner, err)
	}

	raw = encodeRandBatch(8, []game.ShopBuyItem{{GoodsID: 3, Count: 2}})
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 8 {
		t.Fatalf("rand batch cmd %+v %v", m, err)
	}

	raw = encodeCheerJoin(8, 2)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 11 {
		t.Fatalf("cheer join %+v %v", m, err)
	}

	raw = encodeRecallable(8)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 17 {
		t.Fatalf("recallable %+v %v", m, err)
	}

	raw = encodeCharityDonate(8)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 36 {
		t.Fatalf("donate %+v %v", m, err)
	}

	raw = encodeMarkViewed(8)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 7 || m[100].Num != 0 {
		t.Fatalf("mark %+v %v", m, err)
	}

	cost := pb.NewEncoder()
	cost.Int(1, 1005)
	cost.Int(2, 10)
	draw := pb.NewEncoder()
	draw.Int(1, 1)
	draw.Int(2, 5)
	draw.Message(3, cost.Bytes())
	data := pb.NewEncoder()
	data.Message(104, draw.Bytes())
	act, err := decodeActivityData(data.Bytes())
	if err != nil || act.Draw == nil || act.Draw.Today != 1 || act.Draw.Limit != 5 || len(act.Draw.Cost) != 1 {
		t.Fatalf("draw body %+v %v", act.Draw, err)
	}

	tier := pb.NewEncoder()
	tier.Int(1, 2)
	tier.Int(2, 100)
	cheer := pb.NewEncoder()
	cheer.Int(1, 1)
	cheer.Int(2, 30)
	cheer.Message(3, tier.Bytes())
	data = pb.NewEncoder()
	data.Message(107, cheer.Bytes())
	act, err = decodeActivityData(data.Bytes())
	if err != nil || act.Cheer == nil || act.Cheer.CampID != 1 || act.Cheer.Cheer != 30 || len(act.Cheer.Tiers) != 1 {
		t.Fatalf("cheer body %+v %v", act.Cheer, err)
	}
}

func TestEncodeHuntOperate(t *testing.T) {
	raw := encodeHunt(2026090101, cmdHuntOpen)
	m, err := pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 2026090101 || pb.IntField(m, 2) != 913 {
		t.Fatalf("open %+v %v", m, err)
	}

	raw = encodeHuntEquip(2026090101, []int64{101, 102})
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 909 {
		t.Fatalf("equip cmd %+v %v", m, err)
	}
	inner, err := pb.FieldMap(m[909].Bytes)
	if err != nil || pb.IntField(inner, 1) != 102 {
		t.Fatalf("equip inner %+v %v", inner, err)
	}

	raw = encodeHuntBattle(2026090101, 9, "t1")
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 2) != 910 {
		t.Fatalf("battle cmd %+v %v", m, err)
	}
	inner, err = pb.FieldMap(m[910].Bytes)
	if err != nil || pb.IntField(inner, 1) != 9 || pb.StringField(inner, 2) != "t1" {
		t.Fatalf("battle inner %+v %v", inner, err)
	}

	tr := pb.NewEncoder()
	tr.String(1, "tb1")
	tr.Int(2, 5001)
	tr.Int(3, 80)
	tr.Int(6, 2)
	pool := pb.NewEncoder()
	pool.Message(1, tr.Bytes())
	data := pb.NewEncoder()
	data.Message(115, pool.Bytes())
	act, err := decodeActivityData(data.Bytes())
	if err != nil || act.Hunt == nil || len(act.Hunt.Treasures) != 1 || act.Hunt.Treasures[0].ID != "tb1" || act.Hunt.Treasures[0].Status != 2 {
		t.Fatalf("hunt body %+v %v", act.Hunt, err)
	}
}

func TestGapEncodeDecode(t *testing.T) {
	ev := pb.NewEncoder()
	ev.Int(1, 5005)
	ev.Int(2, 9)
	rw := pb.NewEncoder()
	it := pb.NewEncoder()
	it.Int(1, 1001)
	it.Int(2, 2)
	rw.Int(1, 5005)
	rw.Message(2, it.Bytes())
	body := pb.NewEncoder()
	body.Message(1, ev.Bytes())
	body.Message(2, rw.Bytes())
	out, err := decodeCleanFarmEvents(body.Bytes())
	if err != nil || len(out.Events) != 1 || out.Events[0].ItemID != 5005 || len(out.Rewards) != 1 {
		t.Fatalf("clean events %+v %v", out, err)
	}

	raw := encodeBlockApps(true)
	m, err := pb.FieldMap(raw)
	if err != nil || !pb.BoolField(m, 1) {
		t.Fatalf("block %+v %v", m, err)
	}
	got, err := decodeBlockApps(raw)
	if err != nil || !got.Block {
		t.Fatalf("block decode %+v %v", got, err)
	}

	raw = encodeGIDs([]int64{3, 4})
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 4 {
		t.Fatalf("apply wx gids %+v %v", m, err)
	}
	res := pb.NewEncoder()
	res.Int(1, 3)
	res.Bool(2, true)
	apply := pb.NewEncoder()
	apply.Message(1, res.Bytes())
	wx, err := decodeWXApply(apply.Bytes())
	if err != nil || len(wx.Results) != 1 || wx.Results[0].GID != 3 || !wx.Results[0].Success {
		t.Fatalf("apply wx %+v %v", wx, err)
	}

	raw = encodeUID(88)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 88 {
		t.Fatalf("uid %+v %v", m, err)
	}
	tok := pb.NewEncoder()
	tok.String(3, "CODE")
	tok.String(7, "bill1")
	gt, err := decodeGiftToken(tok.Bytes())
	if err != nil || gt.RedeemCode != "CODE" || gt.OutBillNo != "bill1" {
		t.Fatalf("gift %+v %v", gt, err)
	}

	set := game.UserSettings{DisableNudge: true, AllowArkVisit: true}
	raw = encodeUserSettings(set)
	flat, err := decodeUserSettings(raw)
	if err != nil || !flat.DisableNudge || !flat.AllowArkVisit {
		t.Fatalf("settings flat %+v %v", flat, err)
	}
	wrap := pb.NewEncoder()
	wrap.Message(1, raw)
	nested, err := decodeUserSettings(wrap.Bytes())
	if err != nil || !nested.DisableNudge || !nested.AllowArkVisit {
		t.Fatalf("settings nested %+v %v", nested, err)
	}

	raw = encodeCommunityID("c1")
	m, err = pb.FieldMap(raw)
	if err != nil || pb.StringField(m, 1) != "c1" {
		t.Fatalf("bind %+v %v", m, err)
	}
}

func TestPlatformEncodeDecode(t *testing.T) {
	raw := encodeQQVipClaimDaily(9)
	m, err := pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 9 {
		t.Fatalf("claim daily %+v %v", m, err)
	}

	it := pb.NewEncoder()
	it.Int(1, 1001)
	it.Int(2, 2)
	daily := pb.NewEncoder()
	daily.Bool(1, true)
	daily.Bool(2, true)
	daily.Message(4, it.Bytes())
	got, err := decodeQQVipDaily(daily.Bytes())
	if err != nil || !got.IsQQVip || !got.CanClaim || len(got.Rewards) != 1 || got.Rewards[0].ID != 1001 {
		t.Fatalf("daily %+v %v", got, err)
	}

	raw = encodeQQVipClaimRewards([]int64{3, 4})
	ids := pb.PackedInts(raw, 1)
	if len(ids) != 2 || ids[0] != 3 || ids[1] != 4 {
		t.Fatalf("claim rewards encode %+v", ids)
	}
	rew := pb.NewEncoder()
	rew.PackedVarints(1, []int64{11, 12})
	rew.PackedVarints(2, []int64{21})
	rew.Message(3, it.Bytes())
	cr, err := decodeQQVipClaimRewards(rew.Bytes())
	if err != nil || len(cr.SkinIDs) != 2 || cr.SkinIDs[0] != 11 || len(cr.FrameIDs) != 1 || len(cr.Rewards) != 1 {
		t.Fatalf("claim rewards %+v %v", cr, err)
	}

	cfg := pb.NewEncoder()
	cfg.Int(1, 2)
	cfg.Message(2, it.Bytes())
	cfg.Bool(3, true)
	cfg.Int(5, 88)
	st := pb.NewEncoder()
	st.Bool(1, true)
	st.Bool(4, true)
	st.Message(5, cfg.Bytes())
	st.Bool(6, true)
	rs, err := decodeQQVipRewardsStatus(st.Bytes())
	if err != nil || !rs.IsQQVip || !rs.RewardsCanClaim || !rs.HasRedpoint || len(rs.Configs) != 1 || rs.Configs[0].ID != 88 {
		t.Fatalf("rewards status %+v %v", rs, err)
	}

	msg := pb.NewEncoder()
	msg.Int(1, 7)
	msg.Int(2, 1)
	msg.String(5, "hello")
	msg.Int(6, 9)
	mq := pb.NewEncoder()
	mq.Message(1, msg.Bytes())
	mar, err := decodeMarquee(mq.Bytes())
	if err != nil || len(mar.Msgs) != 1 || mar.Msgs[0].UUID != 7 || mar.Msgs[0].Content != "hello" {
		t.Fatalf("marquee %+v %v", mar, err)
	}

	raw = encodeSystemName(3)
	m, err = pb.FieldMap(raw)
	if err != nil || pb.IntField(m, 1) != 3 {
		t.Fatalf("system name %+v %v", m, err)
	}
	un := pb.NewEncoder()
	un.Bool(1, true)
	so, err := decodeSystemOpen(un.Bytes())
	if err != nil || !so.Unlocked {
		t.Fatalf("system open %+v %v", so, err)
	}

	open := pb.NewEncoder()
	open.String(1, "tip")
	open.Message(2, it.Bytes())
	oi, err := decodeMutantOpenInfo(open.Bytes())
	if err != nil || oi.Tips != "tip" || len(oi.Rewards) != 1 {
		t.Fatalf("open info %+v %v", oi, err)
	}

	entry := pb.NewEncoder()
	entry.Int(1, 5)
	entry.Bool(2, true)
	sub := pb.NewEncoder()
	sub.Message(1, entry.Bytes())
	qs, err := decodeQQSubscribe(sub.Bytes())
	if err != nil || !qs.Subscribed || len(qs.Items) != 1 || qs.Items[0].ID != 5 {
		t.Fatalf("qq subscribe %+v %v", qs, err)
	}
	scalar := pb.NewEncoder()
	scalar.Int(1, 2)
	qs, err = decodeQQSubscribe(scalar.Bytes())
	if err != nil || qs.Status != 2 || !qs.Subscribed {
		t.Fatalf("qq subscribe scalar %+v %v", qs, err)
	}

	raw = encodeWXSubscribe([]game.WXTemplateStatus{{TemplateID: "t1", Subscribed: true}})
	wx, err := decodeWXSubscribe(raw)
	if err != nil || len(wx.Templates) != 1 || wx.Templates[0].TemplateID != "t1" || !wx.Templates[0].Subscribed {
		t.Fatalf("wx subscribe %+v %v", wx, err)
	}

	raw = encodeModerateText(game.ModerateTextIn{Text: "hi", Reason: "name"})
	m, err = pb.FieldMap(raw)
	if err != nil || pb.StringField(m, 1) != "hi" || pb.StringField(m, 2) != "name" {
		t.Fatalf("mod text %+v %v", m, err)
	}
	text := pb.NewEncoder()
	text.String(2, "h*")
	text.Bool(3, true)
	text.String(4, "dirty")
	rep := pb.NewEncoder()
	rep.Message(1, text.Bytes())
	rep.String(2, "name")
	mt, err := decodeModerateText(rep.Bytes())
	if err != nil || mt.Text != "h*" || !mt.Dirty || mt.Reason != "dirty" {
		t.Fatalf("mod text decode %+v %v", mt, err)
	}

	raw = encodeBatchModerateText([]game.ModerateTextIn{{Text: "a"}, {Text: "b"}})
	batch := pb.NewEncoder()
	batch.Message(1, text.Bytes())
	bt, err := decodeBatchModerateText(batch.Bytes())
	if err != nil || len(bt.Items) != 1 || !bt.Items[0].Dirty {
		t.Fatalf("batch text %+v %v", bt, err)
	}

	raw = encodeModeratePic(game.ModeratePicIn{URL: "http://x", Reason: "avatar"})
	m, err = pb.FieldMap(raw)
	if err != nil || pb.StringField(m, 1) != "http://x" {
		t.Fatalf("mod pic %+v %v", m, err)
	}
	pic := pb.NewEncoder()
	pic.String(1, "http://y")
	pic.Bool(2, true)
	pic.Int(3, 8)
	pr := pb.NewEncoder()
	pr.Message(1, pic.Bytes())
	mp, err := decodeModeratePic(pr.Bytes())
	if err != nil || mp.URL != "http://y" || !mp.Dirty || mp.DirtyType != 8 {
		t.Fatalf("mod pic decode %+v %v", mp, err)
	}
}
