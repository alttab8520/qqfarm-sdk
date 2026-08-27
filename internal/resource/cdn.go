package resource

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	cdnBase = "https://cdn-resource.nqf.qq.com/release/"
	xorKey  = "NQF_SHANGXIANDAMAI_#2026_SECURE"
)

// DefaultBundleVers are the asset hashes the packaged game shipped with, from
// its settings.*.json. They move with every game release, so Refresh takes an
// override for the case where the game updates before the SDK does.
//
// mainscene and delayRes carry the config tables; the rest carry artwork.
var DefaultBundleVers = map[string]string{
	"mainscene": "eab4a",
	"delayRes":  "30c5b",
	"extraRes":  "aca43",
	"plant":     "1d263",
	"petdog":    "f65f7",
	"weather":   "eb03a",
	"aiHead":    "3e212",
}

// configBundles hold the encrypted config tables. imageBundles hold PNGs
// addressed by native hash rather than import hash.
var configBundles = []string{"mainscene", "delayRes"}
var imageBundles = []string{"extraRes", "plant", "petdog", "weather", "aiHead", "mainscene", "delayRes"}

type rawTable struct {
	Name string
	Rows []map[string]any
}

// fetchTables pulls every config/* asset out of the remote bundles and
// decrypts it. A single asset that fails is skipped: the game ships assets the
// current client never reads, and one bad entry should not lose the rest.
func fetchTables(ctx context.Context, client *http.Client, bundleVers map[string]string) ([]rawTable, error) {
	if len(bundleVers) == 0 {
		bundleVers = DefaultBundleVers
	}
	var out []rawTable
	var firstErr error
	seen := map[string]bool{}
	for _, bundle := range configBundles {
		ver := bundleVers[bundle]
		if ver == "" {
			continue
		}
		tables, err := fetchBundle(ctx, client, bundle, ver)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		for _, t := range tables {
			if seen[t.Name] {
				continue
			}
			seen[t.Name] = true
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		if firstErr != nil {
			return nil, firstErr
		}
		return nil, fmt.Errorf("CDN 没有返回任何配表")
	}
	return out, nil
}

// fetchNativeIndex maps an asset path to its downloadable PNG on the CDN.
//
// Artwork is addressed differently from the config tables: the hash comes from
// versions.native instead of versions.import, the URL sits under native/ not
// import/, and it is a plain PNG with no XOR layer. Item icons live in
// extraRes, crop art in plant, dogs in petdog — none of them in the bundles
// that carry the configs.
func fetchNativeIndex(ctx context.Context, client *http.Client, bundleVers map[string]string) map[string]string {
	if len(bundleVers) == 0 {
		bundleVers = DefaultBundleVers
	}
	// A bundle can list a path but ship the bytes somewhere else: plant names
	// Crop_1_Seed and Crop_2_Seed without a native hash because the tutorial
	// needs those two before plant loads, so they ride along in extraRes. The
	// uuid is the asset's identity and survives that move, so resolve paths to
	// uuids first and hashes to uuids second, then join.
	pathUUID := map[string]string{}
	uuidURL := map[string]string{}

	for _, bundle := range imageBundles {
		ver := bundleVers[bundle]
		if ver == "" {
			continue
		}
		var cfg struct {
			UUIDs    []string         `json:"uuids"`
			Paths    map[string][]any `json:"paths"`
			Versions struct {
				Native []any `json:"native"`
			} `json:"versions"`
		}
		url := fmt.Sprintf("%sremote/%s/config.%s.json", cdnBase, bundle, ver)
		if err := getJSON(ctx, client, url, &cfg); err != nil {
			continue
		}

		for key, val := range cfg.Paths {
			if len(val) == 0 {
				continue
			}
			path, _ := val[0].(string)
			idx, err := strconv.Atoi(key)
			if path == "" || err != nil || idx < 0 || idx >= len(cfg.UUIDs) {
				continue
			}
			// "@" marks a sub-asset such as a texture or spriteFrame; only the
			// bare uuid addresses a standalone PNG.
			if compact := cfg.UUIDs[idx]; !strings.Contains(compact, "@") {
				if _, taken := pathUUID[path]; !taken {
					pathUUID[path] = compact
				}
			}
		}

		for i := 0; i+1 < len(cfg.Versions.Native); i += 2 {
			idx, ok := asInt(cfg.Versions.Native[i])
			if !ok || idx < 0 || idx >= len(cfg.UUIDs) {
				continue
			}
			hash, ok := cfg.Versions.Native[i+1].(string)
			if !ok || hash == "" {
				continue
			}
			compact := cfg.UUIDs[idx]
			if strings.Contains(compact, "@") {
				continue
			}
			if _, taken := uuidURL[compact]; taken {
				continue
			}
			uuid, err := decompressUUID(compact)
			if err != nil {
				continue
			}
			uuidURL[compact] = fmt.Sprintf("%sremote/%s/native/%s/%s.%s.png", cdnBase, bundle, uuid[:2], uuid, hash)
		}
	}

	index := make(map[string]string, len(pathUUID))
	for path, compact := range pathUUID {
		if url := uuidURL[compact]; url != "" {
			index[path] = url
		}
	}
	return index
}

func fetchBundle(ctx context.Context, client *http.Client, bundle, ver string) ([]rawTable, error) {
	var cfg struct {
		UUIDs    []string         `json:"uuids"`
		Paths    map[string][]any `json:"paths"`
		Versions struct {
			Import []any `json:"import"`
		} `json:"versions"`
	}
	url := fmt.Sprintf("%sremote/%s/config.%s.json", cdnBase, bundle, ver)
	if err := getJSON(ctx, client, url, &cfg); err != nil {
		return nil, fmt.Errorf("取 %s 清单失败: %w", bundle, err)
	}

	// versions.import is a flat [index, hash, index, hash, ...] array.
	importHash := map[int]string{}
	for i := 0; i+1 < len(cfg.Versions.Import); i += 2 {
		idx, ok := asInt(cfg.Versions.Import[i])
		if !ok {
			continue
		}
		if h, ok := cfg.Versions.Import[i+1].(string); ok {
			importHash[idx] = h
		}
	}

	var out []rawTable
	for key, val := range cfg.Paths {
		if len(val) == 0 {
			continue
		}
		path, _ := val[0].(string)
		if !strings.HasPrefix(path, "config/") {
			continue
		}
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(cfg.UUIDs) {
			continue
		}
		hash := importHash[idx]
		if hash == "" {
			continue
		}
		uuid, err := decompressUUID(cfg.UUIDs[idx])
		if err != nil {
			continue
		}
		assetURL := fmt.Sprintf("%sremote/%s/import/%s/%s.%s.json", cdnBase, bundle, uuid[:2], uuid, hash)
		rows, err := fetchAsset(ctx, client, assetURL)
		if err != nil || len(rows) == 0 {
			continue
		}
		out = append(out, rawTable{Name: strings.TrimPrefix(path, "config/"), Rows: rows})
	}
	return out, nil
}

// fetchAsset unwraps the Cocos asset envelope. The payload lives at [5][0][2]
// as base64 of the XOR-encrypted table JSON.
func fetchAsset(ctx context.Context, client *http.Client, url string) ([]map[string]any, error) {
	var envelope []any
	if err := getJSON(ctx, client, url, &envelope); err != nil {
		return nil, err
	}
	if len(envelope) < 6 {
		return nil, fmt.Errorf("资源结构不认识")
	}
	outer, ok := envelope[5].([]any)
	if !ok || len(outer) == 0 {
		return nil, fmt.Errorf("资源结构不认识")
	}
	inner, ok := outer[0].([]any)
	if !ok || len(inner) < 3 {
		return nil, fmt.Errorf("资源结构不认识")
	}
	payload, ok := inner[2].(string)
	if !ok || payload == "" {
		return nil, fmt.Errorf("资源没有载荷")
	}
	plain, err := xorDecrypt(payload)
	if err != nil {
		return nil, err
	}
	return decodeRows(plain)
}

func decodeRows(plain []byte) ([]map[string]any, error) {
	var list []map[string]any
	if err := json.Unmarshal(plain, &list); err == nil {
		return list, nil
	}
	var byKey map[string]map[string]any
	if err := json.Unmarshal(plain, &byKey); err != nil {
		return nil, fmt.Errorf("配表既不是数组也不是对象")
	}
	out := make([]map[string]any, 0, len(byKey))
	for _, row := range byKey {
		out = append(out, row)
	}
	return out, nil
}

func xorDecrypt(payload string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("载荷不是 base64: %w", err)
	}
	key := []byte(xorKey)
	for i := range raw {
		raw[i] ^= key[i%len(key)]
	}
	return raw, nil
}

// decompressUUID expands the Cocos short uuid: two literal hex chars followed
// by base64 of the remaining bytes.
func decompressUUID(compact string) (string, error) {
	if i := strings.IndexByte(compact, '@'); i >= 0 {
		compact = compact[:i]
	}
	if len(compact) < 3 {
		return "", fmt.Errorf("uuid 太短: %q", compact)
	}
	head, body := compact[:2], compact[2:]
	if pad := len(body) % 4; pad != 0 {
		body += strings.Repeat("=", 4-pad)
	}
	raw, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return "", fmt.Errorf("uuid 不是 base64: %w", err)
	}
	full := head + hex.EncodeToString(raw)
	if len(full) < 32 {
		return "", fmt.Errorf("uuid 长度不对: %q", full)
	}
	return full[:8] + "-" + full[8:12] + "-" + full[12:16] + "-" + full[16:20] + "-" + full[20:], nil
}

func getJSON(ctx context.Context, client *http.Client, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func asInt(v any) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case int:
		return x, true
	case json.Number:
		n, err := x.Int64()
		return int(n), err == nil
	case string:
		n, err := strconv.Atoi(x)
		return n, err == nil
	}
	return 0, false
}

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
