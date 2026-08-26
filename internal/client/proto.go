package client

import (
	"fmt"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/pb"
)

const (
	gameVersion = "1.13.3.11_20260826"
	channelID   = "1256"
)

func encodeLogin() []byte {
	dev := pb.NewEncoder()
	dev.String(1, gameVersion)
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

func encodeWater(landIDs []int64, hostGID int64) []byte {
	req := pb.NewEncoder()
	req.RepeatedVarint(1, landIDs)
	req.Int(2, hostGID)
	return req.Bytes()
}

func decodeUser(body []byte) (game.User, error) {
	var u game.User
	top, err := pb.FieldMap(body)
	if err != nil {
		return u, err
	}
	if len(top[1].Bytes) == 0 {
		return u, fmt.Errorf("登录回包没有用户资料")
	}
	basic, err := pb.FieldMap(top[1].Bytes)
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
	fields, err := pb.Walk(body)
	if err != nil {
		return nil, err
	}
	var items []game.Item
	for _, f := range fields {
		if f.Num != 2 {
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
		m, err := pb.FieldMap(f.Bytes)
		if err != nil {
			return nil, err
		}
		friends = append(friends, game.Friend{
			GID:    pb.IntField(m, 1),
			OpenID: pb.StringField(m, 2),
			Name:   pb.StringField(m, 3),
			Avatar: pb.StringField(m, 4),
			Level:  pb.IntField(m, 6),
		})
	}
	return friends, nil
}
