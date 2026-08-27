package gate

import (
	"strings"
	"testing"
)

func TestRequestRoundTrip(t *testing.T) {
	raw := EncodeRequest("User", "Login", []byte("hi"), "tok", 7, 3)
	msg, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ServiceName != Services["User"] || msg.MethodName != "Login" {
		t.Fatalf("meta %+v", msg)
	}
	if string(msg.Body) != "hi" || msg.Token != "tok" || msg.ClientSeq != 7 || msg.ServerSeq != 3 {
		t.Fatalf("msg %+v", msg)
	}
	if msg.MessageType != TypeRequest {
		t.Fatalf("type %d", msg.MessageType)
	}
}

func TestServiceNamesAreOfficial(t *testing.T) {
	for key, name := range Services {
		if !strings.HasPrefix(name, "gamepb.") || !strings.Contains(strings.TrimPrefix(name, "gamepb."), ".") {
			t.Fatalf("%s 的服务名不是 gamepb.<pkg>.<Service> 形态: %q", key, name)
		}
	}
	// 这些短名在 client 里发过 RPC。漏一条网关就收到短名，静默打不进服务。
	for _, key := range []string{
		"User", "Plant", "Item", "Friend", "Visit", "Share", "Shop", "Task",
		"Email", "Activity", "Season", "Mall", "Weather", "RedPacket",
		"Interact", "Illustrated", "MysteryShop", "Dog", "BulletinBoard",
		"Mutant", "Career", "Rank", "AvatarFrame", "Skin", "RandomDrop",
		"SolarTerms", "Achieve", "Ace", "QQVip", "Marquee", "SystemOpen",
		"SubscribeQQ", "SubscribeWX", "Uicproxy",
		"Gift", "Misc", "QQGroup", "RechargeBonus",
	} {
		if _, ok := ServiceName(key); !ok {
			t.Fatalf("服务短名 %q 没有登记官方名字", key)
		}
	}
}

func TestServiceNameRejectsUnknownShortKey(t *testing.T) {
	if _, ok := ServiceName("Nope"); ok {
		t.Fatal("未知短名应当判为未登记")
	}
	full := "gamepb.nudgepb.NudgeService"
	if name, ok := ServiceName(full); !ok || name != full {
		t.Fatalf("完整服务名应当原样放行: %q %v", name, ok)
	}
}
