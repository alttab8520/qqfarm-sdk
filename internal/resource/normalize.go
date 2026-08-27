package resource

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

// keepKey wins over dropKey. Icon references say what a thing looks like, so
// they belong with the name even though they are asset paths. Note these are
// paths into the game's own bundles, not URLs: the CDN only serves the remote
// bundles, and item icons are not among them.
var keepKey = regexp.MustCompile(`(?i)icon|_pic$|^asset_name$|^moon_asset$`)

// dropKey matches the rest of the presentation layer: spine rigs, particle
// effects, animation clips, audio, layout geometry. None of it helps explain
// what an ID means, and it would several times the bundle.
var dropKey = regexp.MustCompile(`(?i)asset|_?res$|_res_|spine|effect|anim|aniname|offset|position|scale|audio|prefab|node|layer|sort|color|pic|img|texture|bundle|lod|priority|camera|skip`)

var nameKeys = []string{"name", "skin_name", "title"}
var descKeys = []string{"desc", "effectDesc", "description"}

// normalize turns raw config rows into lookup entries. A table survives only
// if something in it explains an ID: a name, a description, or a cross
// reference to another table's ID.
func normalize(tables []rawTable) []game.ResEntry {
	var out []game.ResEntry
	for _, t := range tables {
		rows := make([]game.ResEntry, 0, len(t.Rows))
		useful := false
		for _, row := range t.Rows {
			id, ok := asInt64(row["id"])
			if !ok {
				continue
			}
			e := game.ResEntry{ID: id, Table: t.Name}
			e.Name = firstString(row, nameKeys, true)
			e.Desc = firstString(row, descKeys, false)
			if v, ok := row["type"]; ok {
				e.Type = scalarString(v)
			}
			e.Extra = extraFields(row)
			if e.Name != "" || e.Desc != "" || hasCrossRef(e.Extra) {
				useful = true
			}
			rows = append(rows, e)
		}
		if useful {
			out = append(out, rows...)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Table != out[j].Table {
			return out[i].Table < out[j].Table
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func extraFields(row map[string]any) map[string]any {
	extra := map[string]any{}
	for k, v := range row {
		if k == "id" || k == "type" || contains(nameKeys, k) || contains(descKeys, k) {
			continue
		}
		if !keepKey.MatchString(k) && dropKey.MatchString(k) {
			continue
		}
		switch x := v.(type) {
		case bool:
			extra[k] = x
		case float64:
			extra[k] = x
		case string:
			if x != "" && len(x) <= 64 {
				extra[k] = x
			}
		case map[string]any:
			// {"id":..,"count":..} shows up for plant fruit and reward drops.
			if len(x) > 0 && len(x) <= 2 {
				for kk, vv := range x {
					if kk != "id" && kk != "count" {
						continue
					}
					if n, ok := asInt64(vv); ok {
						extra[k+"_"+kk] = n
					}
				}
			}
		}
	}
	if len(extra) == 0 {
		return nil
	}
	return extra
}

// iconFields are the config columns that point at a picture, best first.
var iconFields = []string{"icon_res", "icon_asset", "icon_path", "icon", "skin_icon", "entry_icon_path", "moon_asset", "land_icon_res"}

const (
	itemTypeSeed = "5"
	cropAssetDir = "model/v4/"
)

// resolveIcons turns the config's asset references into real CDN URLs.
//
// Two shapes show up. Most tables name a full asset path and only need the
// /spriteFrame suffix trimmed. Crops instead carry a bare asset_name like
// Crop_1, which addresses two different pictures: _6 is the 150x150 mature
// still, _Seed the 100x100 seed-bag closeup. Seeds want the bag, everything
// else wants the mature still.
func resolveIcons(entries []game.ResEntry, index map[string]string) {
	if len(index) == 0 {
		return
	}
	for i := range entries {
		e := &entries[i]
		for _, key := range iconFields {
			path, _ := e.Extra[key].(string)
			if path == "" {
				continue
			}
			if url := index[strings.TrimSuffix(path, "/spriteFrame")]; url != "" {
				e.IconURL = url
				break
			}
		}
		if e.IconURL != "" {
			continue
		}
		asset, _ := e.Extra["asset_name"].(string)
		if asset == "" {
			continue
		}
		order := []string{"_6", "_Seed"}
		if e.Type == itemTypeSeed {
			order = []string{"_Seed", "_6"}
		}
		for _, suffix := range order {
			if url := index[cropAssetDir+asset+suffix]; url != "" {
				e.IconURL = url
				break
			}
		}
	}
}

func hasCrossRef(extra map[string]any) bool {
	for k := range extra {
		if strings.HasSuffix(k, "_id") {
			return true
		}
	}
	return false
}

func firstString(row map[string]any, keys []string, alsoSuffix bool) string {
	for _, k := range keys {
		if s, ok := row[k].(string); ok && s != "" {
			return s
		}
	}
	if !alsoSuffix {
		return ""
	}
	// Tables that name their column after themselves, e.g. skin_name.
	suffixed := make([]string, 0, 4)
	for k, v := range row {
		if s, ok := v.(string); ok && s != "" && strings.HasSuffix(k, "_name") {
			suffixed = append(suffixed, k)
		}
	}
	if len(suffixed) == 0 {
		return ""
	}
	sort.Strings(suffixed)
	return row[suffixed[0]].(string)
}

func scalarString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case bool:
		return strconv.FormatBool(x)
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
	return ""
}

func asInt64(v any) (int64, bool) {
	switch x := v.(type) {
	case float64:
		return int64(x), true
	case int64:
		return x, true
	case int:
		return int64(x), true
	case string:
		n, err := strconv.ParseInt(x, 10, 64)
		return n, err == nil
	}
	return 0, false
}

func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}
