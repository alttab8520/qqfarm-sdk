package resource

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	snap, err := parseBundled()
	if err != nil {
		t.Fatalf("内置资源表读不出来: %v", err)
	}
	if len(snap.Entries) == 0 {
		t.Fatal("内置资源表是空的")
	}
	return &Store{snap: snap, client: newHTTPClient()}
}

func TestBundledHasCoreTables(t *testing.T) {
	out := testStore(t).Tables()
	if out.Source != "bundled" {
		t.Fatalf("source = %q", out.Source)
	}
	counts := map[string]int{}
	for _, tbl := range out.Tables {
		counts[tbl.Name] = tbl.Count
	}
	for _, name := range []string{"ItemInfo", "Plant", "Goods"} {
		if counts[name] == 0 {
			t.Fatalf("内置表缺少 %s: %+v", name, counts)
		}
	}
}

func TestLookupKnownIDs(t *testing.T) {
	s := testStore(t)
	// README 里硬写过的几个 ID，名字必须对得上官方配表。
	want := map[int64]string{
		1001:  "金币",
		1002:  "点券",
		1005:  "金豆豆",
		5001:  "天气采集瓶",
		5005:  "青蛙使坏瓶",
		5006:  "乌云使坏瓶",
		20001: "草莓种子",
	}
	for id, name := range want {
		out, err := s.Lookup(game.ResLookupIn{IDs: []int64{id}, Table: "ItemInfo"})
		if err != nil {
			t.Fatalf("%d: %v", id, err)
		}
		if len(out.Entries) != 1 || out.Entries[0].Name != name {
			t.Fatalf("%d 查到 %+v，期望 %q", id, out.Entries, name)
		}
	}
}

func TestLookupRejectsEmptyIDs(t *testing.T) {
	if _, err := testStore(t).Lookup(game.ResLookupIn{}); err == nil {
		t.Fatal("空 ids 应当报错")
	}
}

func TestListSeedsByType(t *testing.T) {
	s := testStore(t)
	out, err := s.List(game.ResListIn{Table: "ItemInfo", Type: "5", Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if out.Total < 100 {
		t.Fatalf("种子只有 %d 个，太少了", out.Total)
	}
	if len(out.Entries) != 5 {
		t.Fatalf("limit 没生效: %d", len(out.Entries))
	}
	for _, e := range out.Entries {
		if e.Type != "5" {
			t.Fatalf("type 过滤漏了: %+v", e)
		}
	}
}

func TestListKeywordAndOffset(t *testing.T) {
	s := testStore(t)
	all, err := s.List(game.ResListIn{Table: "ItemInfo", Keyword: "草莓"})
	if err != nil {
		t.Fatal(err)
	}
	if all.Total == 0 {
		t.Fatal("按“草莓”一条都没搜到")
	}
	page, err := s.List(game.ResListIn{Table: "ItemInfo", Keyword: "草莓", Offset: 1, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != all.Total {
		t.Fatalf("total 应当不受分页影响: %d vs %d", page.Total, all.Total)
	}
	if all.Total > 1 && (len(page.Entries) != 1 || page.Entries[0].ID == all.Entries[0].ID) {
		t.Fatalf("offset 没生效: %+v", page.Entries)
	}
}

func TestListUnknownTable(t *testing.T) {
	if _, err := testStore(t).List(game.ResListIn{Table: "Nope"}); err == nil {
		t.Fatal("不存在的表应当报错")
	}
}

func TestPlantCarriesSeedAndFruit(t *testing.T) {
	out, err := testStore(t).List(game.ResListIn{Table: "Plant", Keyword: "草莓", Limit: 1})
	if err != nil || len(out.Entries) == 0 {
		t.Fatalf("没查到草莓: %v %+v", err, out)
	}
	extra := out.Entries[0].Extra
	for _, k := range []string{"seed_id", "fruit_id", "land_level_need"} {
		if _, ok := extra[k]; !ok {
			t.Fatalf("Plant 缺少 %s: %+v", k, extra)
		}
	}
}

func TestDecompressUUID(t *testing.T) {
	full, err := decompressUUID("28dsRacyBAxYIeb+ecKpwB@6c48a")
	if err != nil {
		t.Fatal(err)
	}
	if len(full) != 36 || full[8] != '-' || full[13] != '-' {
		t.Fatalf("uuid 形态不对: %q", full)
	}
	if _, err := decompressUUID("ab"); err == nil {
		t.Fatal("过短的 uuid 应当报错")
	}
}

func TestXORDecryptRoundTrip(t *testing.T) {
	plain := []byte(`[{"id":1,"name":"甲"}]`)
	key := []byte(xorKey)
	enc := make([]byte, len(plain))
	for i := range plain {
		enc[i] = plain[i] ^ key[i%len(key)]
	}
	got, err := xorDecrypt(base64.StdEncoding.EncodeToString(enc))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(plain) {
		t.Fatalf("解出来是 %q", got)
	}
}

func TestNormalizeDropsPresentationOnlyTables(t *testing.T) {
	entries := normalize([]rawTable{
		{Name: "Sprites", Rows: []map[string]any{{"id": 1.0, "icon_res": "a", "spine_asset": "b"}}},
		{Name: "Named", Rows: []map[string]any{{"id": 2.0, "name": "甲", "seed_id": 7.0}}},
	})
	if len(entries) != 1 || entries[0].Table != "Named" {
		t.Fatalf("只应留下有名字的表: %+v", entries)
	}
	if entries[0].Extra["seed_id"] != 7.0 {
		t.Fatalf("交叉引用字段应当保留: %+v", entries[0].Extra)
	}
}

func TestNormalizeKeepsIconDropsRig(t *testing.T) {
	entries := normalize([]rawTable{{Name: "T", Rows: []map[string]any{{
		"id":         1.0,
		"name":       "甲",
		"icon_res":   "gui/texture/icon/golds/spriteFrame",
		"asset_name": "Crop_2101",
		"icon_asset": "gui/texture/skinDetail/a/spriteFrame",
		// 纯渲染，继续丢
		"spine_asset":  "spine/v2/x",
		"effect_res":   "effect/prefab/y",
		"unlocked_res": "model/v3/z",
	}}}})
	extra := entries[0].Extra
	for _, k := range []string{"icon_res", "asset_name", "icon_asset"} {
		if _, ok := extra[k]; !ok {
			t.Fatalf("图标引用 %s 应当保留: %+v", k, extra)
		}
	}
	for _, k := range []string{"spine_asset", "effect_res", "unlocked_res"} {
		if _, ok := extra[k]; ok {
			t.Fatalf("渲染字段 %s 不该留下: %+v", k, extra)
		}
	}
}

func TestResolveIconsPrefersSeedBagForSeeds(t *testing.T) {
	index := map[string]string{
		"model/v4/Crop_1_6":         "https://cdn/mature.png",
		"model/v4/Crop_1_Seed":      "https://cdn/bag.png",
		"gui/texture/icon/golds":    "https://cdn/gold.png",
		"gui/texture/skinDetail/im": "https://cdn/skin.png",
	}
	entries := []game.ResEntry{
		{ID: 20001, Type: "5", Extra: map[string]any{"asset_name": "Crop_1"}},
		{ID: 40001, Type: "6", Extra: map[string]any{"asset_name": "Crop_1"}},
		// 完整路径带 /spriteFrame 后缀，要先去掉再查
		{ID: 1001, Extra: map[string]any{"icon_res": "gui/texture/icon/golds/spriteFrame"}},
		{ID: 9, Extra: map[string]any{"asset_name": "Crop_没有"}},
	}
	resolveIcons(entries, index)
	if entries[0].IconURL != "https://cdn/bag.png" {
		t.Fatalf("种子应当拿种子袋图: %q", entries[0].IconURL)
	}
	if entries[1].IconURL != "https://cdn/mature.png" {
		t.Fatalf("果实应当拿成熟图: %q", entries[1].IconURL)
	}
	if entries[2].IconURL != "https://cdn/gold.png" {
		t.Fatalf("/spriteFrame 后缀没去掉: %q", entries[2].IconURL)
	}
	if entries[3].IconURL != "" {
		t.Fatalf("查不到就该留空: %q", entries[3].IconURL)
	}
}

func TestBundledIconURLsPointAtCDN(t *testing.T) {
	out, err := testStore(t).List(game.ResListIn{Table: "ItemInfo", Limit: 2000})
	if err != nil {
		t.Fatal(err)
	}
	withIcon := 0
	for _, e := range out.Entries {
		if e.IconURL == "" {
			continue
		}
		withIcon++
		if !strings.HasPrefix(e.IconURL, "https://cdn-resource.nqf.qq.com/release/remote/") ||
			!strings.HasSuffix(e.IconURL, ".png") {
			t.Fatalf("%d %s 的图不像官方 CDN 地址: %q", e.ID, e.Name, e.IconURL)
		}
	}
	if withIcon < 600 {
		t.Fatalf("带图的物品只有 %d 个，太少了", withIcon)
	}
}

// 种子袋图和成熟图是两张不同的图。官方只做了种子袋图的那几株除外，
// 它们游戏里本来就共用一张，不是解析错。
func TestSeedAndFruitIconsDiffer(t *testing.T) {
	s := testStore(t)
	items, err := s.List(game.ResListIn{Table: "ItemInfo", Limit: 2000})
	if err != nil {
		t.Fatal(err)
	}
	byID := map[int64]game.ResEntry{}
	for _, e := range items.Entries {
		byID[e.ID] = e
	}
	plants, err := s.List(game.ResListIn{Table: "Plant", Limit: 2000})
	if err != nil {
		t.Fatal(err)
	}
	sharedArt := map[int64]bool{21032: true, 21542: true, 29003: true}
	pairs, collisions := 0, 0
	for _, p := range plants.Entries {
		seedID, _ := asInt64(p.Extra["seed_id"])
		fruitID, _ := asInt64(p.Extra["fruit_id"])
		seed, fruit := byID[seedID], byID[fruitID]
		if seed.IconURL == "" || fruit.IconURL == "" {
			continue
		}
		pairs++
		if seed.IconURL != fruit.IconURL {
			continue
		}
		collisions++
		if !sharedArt[seedID] {
			t.Fatalf("%d %s 和 %d %s 指到了同一张图: %s",
				seedID, seed.Name, fruitID, fruit.Name, seed.IconURL)
		}
	}
	if pairs < 150 {
		t.Fatalf("只对上 %d 组种子/果实，太少了", pairs)
	}
	if collisions != len(sharedArt) {
		t.Fatalf("共用图的作物有 %d 株，期望 %d 株", collisions, len(sharedArt))
	}
}

func TestBundledItemsCarryIcons(t *testing.T) {
	out, err := testStore(t).List(game.ResListIn{Table: "ItemInfo", Limit: 1000})
	if err != nil {
		t.Fatal(err)
	}
	icons := 0
	for _, e := range out.Entries {
		if e.Extra["icon_res"] != nil || e.Extra["asset_name"] != nil {
			icons++
		}
	}
	if icons < 400 {
		t.Fatalf("带图标引用的物品只有 %d 个，太少了", icons)
	}
}

func TestNormalizeKeepsCrossRefOnlyTable(t *testing.T) {
	entries := normalize([]rawTable{
		{Name: "Goods", Rows: []map[string]any{{"id": 1.0, "item_id": 20001.0, "price_id": 1001.0}}},
	})
	if len(entries) != 1 {
		t.Fatalf("只有交叉引用的表也要留下: %+v", entries)
	}
}
