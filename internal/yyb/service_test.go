package yyb

import "testing"

func TestPickNicknameAndAvatar(t *testing.T) {
	info := map[string]any{"nick_name": "甲", "head_img_url": "http://a"}
	if got := pickNickname(info, "乙"); got != "甲" {
		t.Fatalf("nick %q", got)
	}
	if got := pickAvatarURL(info); got != "http://a" {
		t.Fatalf("avatar %q", got)
	}
	if got := pickNickname(nil, "乙"); got != "乙" {
		t.Fatalf("fallback %q", got)
	}
}

func TestUnwrapUserInfo(t *testing.T) {
	raw := map[string]any{"ret": 0, "data": map[string]any{"nickname": "甲", "avatar": "http://a"}}
	inner := unwrapUserInfo(raw)
	if pickNickname(inner, "") != "甲" || pickAvatarURL(inner) != "http://a" {
		t.Fatalf("%+v", inner)
	}
	if unwrapUserInfo(nil) != nil {
		t.Fatal("nil")
	}
}
