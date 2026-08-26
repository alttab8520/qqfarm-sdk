package client

import (
	"testing"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/pb"
)

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
	show := pb.NewEncoder()
	show.Int(1, 12)
	item := pb.NewEncoder()
	item.Int(1, 40061)
	item.Int(2, 8)
	item.Int(6, 77)
	item.Message(100, show.Bytes())
	bag := pb.NewEncoder()
	bag.Message(1, item.Bytes())
	reply := pb.NewEncoder()
	reply.Message(1, bag.Bytes())
	items, err := decodeBagReply(reply.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != 40061 || items[0].Count != 8 || items[0].UID != 77 || items[0].SellPrice != 12 {
		t.Fatalf("%+v", items)
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
	album.Bool(6, true)
	a, err := decodeAlbum(album.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Items) != 1 || a.Items[0].FruitID != 40061 || a.Progress != 8 || !a.Claimable {
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
	raw := encodeHarvest([]int64{1, 2}, 8, true)
	m, err := pb.Walk(raw)
	if err != nil {
		t.Fatal(err)
	}
	var ids []int64
	for _, f := range m {
		if f.Num == 1 {
			ids = append(ids, int64(f.Varint))
		}
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("ids %v", ids)
	}
}
