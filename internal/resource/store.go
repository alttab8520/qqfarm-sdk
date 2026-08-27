package resource

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

// bundled is the offline fallback: the official config tables as of the game
// version this SDK was built against. Refresh replaces it at runtime.
//
//go:embed bundled.json.gz
var bundled []byte

const defaultListLimit = 200

type snapshot struct {
	Version string          `json:"version,omitempty"`
	Source  string          `json:"source"`
	Entries []game.ResEntry `json:"entries"`

	byID    map[int64][]int
	byTable map[string][]int
	tables  []game.ResTable
}

func (s *snapshot) index() {
	s.byID = make(map[int64][]int, len(s.Entries))
	s.byTable = map[string][]int{}
	for i, e := range s.Entries {
		s.byID[e.ID] = append(s.byID[e.ID], i)
		s.byTable[e.Table] = append(s.byTable[e.Table], i)
	}
	s.tables = make([]game.ResTable, 0, len(s.byTable))
	for name, idx := range s.byTable {
		s.tables = append(s.tables, game.ResTable{Name: name, Count: len(idx)})
	}
	sortTables(s.tables)
}

type Store struct {
	mu     sync.RWMutex
	snap   *snapshot
	client *http.Client
	cache  string
}

func Open() *Store {
	s := &Store{client: newHTTPClient(), cache: CachePath()}
	s.snap = loadInitial(s.cache)
	return s
}

// CachePath is where Refresh writes, and where Open looks before falling back
// to the embedded copy.
func CachePath() string {
	if p := os.Getenv("FARM_RES"); p != "" {
		return p
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "data", "resources.json")
	}
	return filepath.Join("data", "resources.json")
}

func loadInitial(cache string) *snapshot {
	if cache != "" {
		if raw, err := os.ReadFile(cache); err == nil {
			if snap, err := parseSnapshot(raw); err == nil && len(snap.Entries) > 0 {
				snap.Source = "cache"
				snap.index()
				return snap
			}
		}
	}
	if snap, err := parseBundled(); err == nil {
		return snap
	}
	empty := &snapshot{Source: "empty"}
	empty.index()
	return empty
}

func parseBundled() (*snapshot, error) {
	if len(bundled) == 0 {
		return nil, fmt.Errorf("内置资源表为空")
	}
	zr, err := gzip.NewReader(bytes.NewReader(bundled))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	raw, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	snap, err := parseSnapshot(raw)
	if err != nil {
		return nil, err
	}
	snap.Source = "bundled"
	snap.index()
	return snap, nil
}

func parseSnapshot(raw []byte) (*snapshot, error) {
	var snap snapshot
	if err := json.Unmarshal(raw, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

func (s *Store) current() *snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snap
}

func (s *Store) Lookup(in game.ResLookupIn) (game.ResListOut, error) {
	if len(in.IDs) == 0 {
		return game.ResListOut{}, fmt.Errorf("ids 不能为空")
	}
	snap := s.current()
	table := strings.ToLower(in.Table)
	out := game.ResListOut{Entries: []game.ResEntry{}, Source: snap.Source, Version: snap.Version}
	for _, id := range in.IDs {
		for _, i := range snap.byID[id] {
			e := snap.Entries[i]
			if table != "" && strings.ToLower(e.Table) != table {
				continue
			}
			out.Entries = append(out.Entries, e)
		}
	}
	out.Total = len(out.Entries)
	return out, nil
}

func (s *Store) List(in game.ResListIn) (game.ResListOut, error) {
	snap := s.current()
	table := strings.ToLower(in.Table)
	keyword := strings.ToLower(in.Keyword)

	pool := snap.byTable[in.Table]
	if pool == nil && table != "" {
		for name, idx := range snap.byTable {
			if strings.ToLower(name) == table {
				pool = idx
				break
			}
		}
		if pool == nil {
			return game.ResListOut{}, fmt.Errorf("没有这张配表: %s", in.Table)
		}
	}

	matched := make([]game.ResEntry, 0, 64)
	appendIfMatch := func(e game.ResEntry) {
		if in.Type != "" && e.Type != in.Type {
			return
		}
		if keyword != "" &&
			!strings.Contains(strings.ToLower(e.Name), keyword) &&
			!strings.Contains(strings.ToLower(e.Desc), keyword) {
			return
		}
		matched = append(matched, e)
	}
	if pool != nil {
		for _, i := range pool {
			appendIfMatch(snap.Entries[i])
		}
	} else {
		for _, e := range snap.Entries {
			appendIfMatch(e)
		}
	}

	out := game.ResListOut{Total: len(matched), Source: snap.Source, Version: snap.Version}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}
	limit := in.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}
	if offset >= len(matched) {
		out.Entries = []game.ResEntry{}
		return out, nil
	}
	end := offset + limit
	if end > len(matched) {
		end = len(matched)
	}
	out.Entries = matched[offset:end]
	return out, nil
}

func (s *Store) Tables() game.ResTablesOut {
	snap := s.current()
	return game.ResTablesOut{
		Tables:  snap.tables,
		Total:   len(snap.Entries),
		Source:  snap.Source,
		Version: snap.Version,
	}
}

func (s *Store) Refresh(ctx context.Context, in game.ResRefreshIn) (game.ResRefreshOut, error) {
	tables, err := fetchTables(ctx, s.client, in.BundleVers)
	if err != nil {
		return game.ResRefreshOut{}, err
	}
	entries := normalize(tables)
	if len(entries) == 0 {
		return game.ResRefreshOut{}, fmt.Errorf("CDN 配表解出来是空的")
	}
	resolveIcons(entries, fetchNativeIndex(ctx, s.client, in.BundleVers))
	snap := &snapshot{Source: "cdn", Entries: entries}
	snap.index()

	out := game.ResRefreshOut{
		Tables:  len(snap.tables),
		Entries: len(entries),
		Source:  snap.Source,
	}
	if err := s.writeCache(snap); err == nil {
		out.Cache = s.cache
	}
	s.mu.Lock()
	s.snap = snap
	s.mu.Unlock()
	return out, nil
}

// BuildBundle fetches the live tables and returns the gzipped snapshot that
// gets embedded as the offline fallback. Used by tools/genresources.
func BuildBundle(ctx context.Context, bundleVers map[string]string, version string) ([]byte, error) {
	client := newHTTPClient()
	tables, err := fetchTables(ctx, client, bundleVers)
	if err != nil {
		return nil, err
	}
	entries := normalize(tables)
	if len(entries) == 0 {
		return nil, fmt.Errorf("CDN 配表解出来是空的")
	}
	resolveIcons(entries, fetchNativeIndex(ctx, client, bundleVers))
	raw, err := json.Marshal(&snapshot{Source: "bundled", Version: version, Entries: entries})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := zw.Write(raw); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Store) writeCache(snap *snapshot) error {
	if s.cache == "" {
		return fmt.Errorf("没有缓存路径")
	}
	raw, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.cache), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.cache, raw, 0o644)
}

func sortTables(list []game.ResTable) {
	for i := 1; i < len(list); i++ {
		for j := i; j > 0 && list[j].Name < list[j-1].Name; j-- {
			list[j], list[j-1] = list[j-1], list[j]
		}
	}
}
