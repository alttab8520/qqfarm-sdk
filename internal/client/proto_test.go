package client

import (
	"testing"

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
	if string(dev[1].Bytes) != gameVersion {
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
