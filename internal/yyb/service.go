package yyb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/yyb/protocol"
	"github.com/alttab8520/qqfarm-sdk/internal/yyb/qr"
	"github.com/alttab8520/qqfarm-sdk/internal/yyb/store"
)

const FarmAppID = "wx5306c5978fdb76e4"

const (
	getCodeAttempts = 10
	getCodeRetry    = 2 * time.Second
	qrSessionTTL    = 10 * time.Minute
)

type expiredError struct{ openid string }

func (e expiredError) Error() string {
	return "应用宝登录态过期，需要重新扫码: " + e.openid
}

type Service struct {
	db   *store.DB
	pool *protocol.Pool
	qr   *qr.Client

	mu       sync.Mutex
	sessions map[string]*qr.Session
	images   map[string][]byte
}

func DirFromEnv() string {
	if v := os.Getenv("FARM_YYB_DIR"); v != "" {
		return v
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "yyb_data")
	}
	return "yyb_data"
}

func Open(dir string) (*Service, error) {
	if dir == "" {
		dir = DirFromEnv()
	}
	db, err := store.Open(filepath.Join(dir, "db", "yyb.db"))
	if err != nil {
		return nil, err
	}
	cfg := protocol.DefaultConfig()
	if p := os.Getenv("FARM_YYB_PROXY"); p != "" {
		cfg.TCPProxy = p
	}
	return &Service{
		db:       db,
		pool:     protocol.NewPool(cfg, db),
		qr:       qr.NewClient(8 * time.Second),
		sessions: map[string]*qr.Session{},
		images:   map[string][]byte{},
	}, nil
}

func (s *Service) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Service) Accounts(ctx context.Context) ([]game.YYBAccount, error) {
	rows, err := s.db.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]game.YYBAccount, 0, len(rows))
	for _, acc := range rows {
		out = append(out, publicAccount(acc))
	}
	return out, nil
}

func (s *Service) CreateQR(ctx context.Context) (game.YYBQROut, error) {
	s.pruneQR()
	img, err := s.qr.GetQRCodeImage(ctx)
	if err != nil {
		return game.YYBQROut{}, err
	}
	s.mu.Lock()
	s.sessions[img.Session.ID] = img.Session
	s.images[img.Session.ID] = img.ImageBytes
	s.mu.Unlock()
	return game.YYBQROut{
		SessionID: img.Session.ID,
		Status:    img.Session.Status,
		Image:     qr.DataURIJPEG(img.ImageBytes),
	}, nil
}

func (s *Service) QRImage(sessionID string) ([]byte, error) {
	s.mu.Lock()
	raw := s.images[sessionID]
	s.mu.Unlock()
	if len(raw) == 0 {
		return nil, fmt.Errorf("扫码会话不存在")
	}
	return raw, nil
}

func (s *Service) Poll(ctx context.Context, sessionID string) (game.YYBPollOut, error) {
	sess := s.getSession(sessionID)
	if sess == nil {
		return game.YYBPollOut{}, fmt.Errorf("扫码会话不存在")
	}
	result, err := s.qr.PollQRCode(ctx, sess)
	if err != nil {
		return game.YYBPollOut{Status: "pending", Msg: err.Error()}, nil
	}
	if result.Status == "expired" || result.Status == "cancelled" || result.Status == "unknown" {
		s.dropSession(sessionID)
	}
	out := game.YYBPollOut{Status: result.Status, Msg: result.Message}
	if result.ErrCode != nil {
		out.ErrCode = int64(*result.ErrCode)
	}
	return out, nil
}

func (s *Service) Confirm(ctx context.Context, sessionID string) (game.YYBAccount, error) {
	sess := s.getSession(sessionID)
	if sess == nil {
		return game.YYBAccount{}, fmt.Errorf("扫码会话不存在")
	}
	result, err := s.qr.GetLoginBuffer(ctx, sess)
	if err != nil {
		return game.YYBAccount{}, fmt.Errorf("还不能确认: %w", err)
	}
	var userInfo map[string]any
	if ui, err := s.qr.LoginBuffers().FetchUserInfo(ctx, result.Credentials); err == nil {
		userInfo = unwrapUserInfo(ui)
	}
	status := "alive"
	nick := pickNickname(userInfo, result.Credentials.Nickname)
	avatar := pickAvatarURL(userInfo)
	acc, err := s.db.UpsertAccount(ctx, result.Credentials.OpenID, result.LoginBuffer, stringPtr(nick), stringPtr(nick), stringPtr(avatar), userInfo, result.Credentials.ToMap(), &status)
	if err != nil {
		return game.YYBAccount{}, err
	}
	s.dropSession(sessionID)
	return publicAccount(acc), nil
}

func (s *Service) Refresh(ctx context.Context, ref string) (game.YYBRefreshOut, error) {
	acc, err := s.resolve(ctx, ref)
	if err != nil {
		return game.YYBRefreshOut{}, err
	}
	status := s.refreshLiveness(ctx, acc)
	return game.YYBRefreshOut{ID: acc.ID, OpenID: acc.OpenID, Status: status}, nil
}

func (s *Service) GetCode(ctx context.Context, ref, appID string) (game.YYBCodeOut, error) {
	if appID == "" {
		appID = FarmAppID
	}
	var last error
	for i := 0; i < getCodeAttempts; i++ {
		acc, err := s.resolve(ctx, ref)
		if err != nil {
			return game.YYBCodeOut{}, err
		}
		out, err := s.getCodeOnce(ctx, acc, appID)
		if err == nil {
			return out, nil
		}
		last = err
		var expired expiredError
		if errors.As(err, &expired) {
			if st := s.refreshLiveness(ctx, acc); st != "alive" {
				return game.YYBCodeOut{}, expired
			}
		}
		if i+1 < getCodeAttempts {
			select {
			case <-ctx.Done():
				return game.YYBCodeOut{}, ctx.Err()
			case <-time.After(getCodeRetry):
			}
		}
	}
	return game.YYBCodeOut{}, fmt.Errorf("取 code 失败: %w", last)
}

func (s *Service) getCodeOnce(ctx context.Context, acc *store.WechatAccount, appID string) (game.YYBCodeOut, error) {
	res, live, err := s.withLive(ctx, acc, func(acc *store.WechatAccount) (map[string]any, error) {
		return s.pool.GetCode(ctx, acc.LoginBuffer, appID, acc.ID, os.Getenv("FARM_YYB_PROXY"))
	})
	if err != nil {
		return game.YYBCodeOut{}, err
	}
	return codeOut(live.OpenID, res)
}

func (s *Service) Delete(ctx context.Context, ref string) (game.YYBDeleteOut, error) {
	if ref == "" {
		return game.YYBDeleteOut{}, fmt.Errorf("ref 不能为空")
	}
	acc, err := s.resolve(ctx, ref)
	if err != nil {
		return game.YYBDeleteOut{}, err
	}
	if err := s.db.DeleteAccount(ctx, acc.ID); err != nil {
		return game.YYBDeleteOut{}, err
	}
	return game.YYBDeleteOut{Deleted: acc.ID, OpenID: acc.OpenID}, nil
}

func (s *Service) Profile(ctx context.Context, ref string) (game.YYBAccount, error) {
	acc, err := s.resolve(ctx, ref)
	if err != nil {
		return game.YYBAccount{}, err
	}
	if acc.Credentials == nil {
		return game.YYBAccount{}, fmt.Errorf("账号没有登录凭据，重新扫码")
	}
	creds := protocol.CredentialsFromMap(acc.Credentials)
	raw, err := s.qr.LoginBuffers().FetchUserInfo(ctx, creds)
	if err != nil {
		return game.YYBAccount{}, err
	}
	userInfo := unwrapUserInfo(raw)
	nick := pickNickname(userInfo, deref(acc.Nickname))
	avatar := pickAvatarURL(userInfo)
	if avatar == "" {
		avatar = deref(acc.Avatar)
	}
	if err := s.db.SetAccountProfile(ctx, acc.ID, stringPtr(nick), stringPtr(avatar), userInfo); err != nil {
		return game.YYBAccount{}, err
	}
	fresh, err := s.db.GetAccount(ctx, acc.ID)
	if err != nil {
		return game.YYBAccount{}, err
	}
	return publicAccount(fresh), nil
}

func (s *Service) Phone(ctx context.Context, ref, appID string) (game.YYBRawOut, error) {
	if appID == "" {
		appID = FarmAppID
	}
	acc, err := s.resolve(ctx, ref)
	if err != nil {
		return game.YYBRawOut{}, err
	}
	res, live, err := s.withLive(ctx, acc, func(acc *store.WechatAccount) (map[string]any, error) {
		return s.pool.GetPhoneNumber(ctx, acc.LoginBuffer, appID, acc.ID, os.Getenv("FARM_YYB_PROXY"))
	})
	if err != nil {
		return game.YYBRawOut{}, err
	}
	return game.YYBRawOut{OpenID: live.OpenID, Result: res}, nil
}

func (s *Service) WXData(ctx context.Context, in game.YYBWXDataIn) (game.YYBRawOut, error) {
	if in.Payload == nil {
		return game.YYBRawOut{}, fmt.Errorf("payload 不能为空")
	}
	appID := in.AppID
	if appID == "" {
		appID = FarmAppID
	}
	acc, err := s.resolve(ctx, in.Ref)
	if err != nil {
		return game.YYBRawOut{}, err
	}
	res, live, err := s.withLive(ctx, acc, func(acc *store.WechatAccount) (map[string]any, error) {
		return s.pool.OperateWXData(ctx, acc.LoginBuffer, appID, in.Payload, acc.ID, os.Getenv("FARM_YYB_PROXY"))
	})
	if err != nil {
		return game.YYBRawOut{}, err
	}
	return game.YYBRawOut{OpenID: live.OpenID, Result: res}, nil
}

func (s *Service) withLive(ctx context.Context, acc *store.WechatAccount, fn func(*store.WechatAccount) (map[string]any, error)) (map[string]any, *store.WechatAccount, error) {
	proxy := os.Getenv("FARM_YYB_PROXY")
	if _, err := s.db.GetSession(ctx, acc.ID, proxy); err == nil {
		res, err := fn(acc)
		if err == nil {
			return res, acc, nil
		}
		_ = s.db.InvalidateSession(ctx, acc.ID, proxy)
	}
	if st := s.refreshLiveness(ctx, acc); st != "alive" {
		return nil, acc, expiredError{openid: acc.OpenID}
	}
	fresh, err := s.db.GetAccount(ctx, acc.ID)
	if err == nil && fresh != nil {
		acc = fresh
	}
	res, err := fn(acc)
	if err != nil {
		return nil, acc, err
	}
	return res, acc, nil
}

func (s *Service) refreshLiveness(ctx context.Context, acc *store.WechatAccount) string {
	if acc.Credentials == nil {
		_ = s.db.SetAccountStatus(ctx, acc.ID, "unknown")
		return "unknown"
	}
	creds := protocol.CredentialsFromMap(acc.Credentials)
	result, err := s.qr.RefreshLoginBuffer(ctx, creds)
	if err != nil {
		_ = s.db.SetAccountStatus(ctx, acc.ID, "expired")
		return "expired"
	}
	_ = s.db.SetAccountCredential(ctx, acc.ID, result.LoginBuffer, result.Credentials.ToMap())
	_ = s.db.SetAccountStatus(ctx, acc.ID, "alive")
	return "alive"
}

func (s *Service) resolve(ctx context.Context, ref string) (*store.WechatAccount, error) {
	if ref == "" {
		rows, err := s.db.ListAccounts(ctx)
		if err != nil {
			return nil, err
		}
		if len(rows) == 0 {
			return nil, fmt.Errorf("没有应用宝账号，先扫码")
		}
		return rows[0], nil
	}
	acc, err := s.db.ResolveAccount(ctx, ref)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("账号不存在: %s", ref)
	}
	return acc, err
}

func (s *Service) getSession(id string) *qr.Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sessions[id]
}

func (s *Service) dropSession(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	delete(s.images, id)
	s.mu.Unlock()
}

func (s *Service) pruneQR() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, sess := range s.sessions {
		if sess.Age() > qrSessionTTL {
			delete(s.sessions, id)
			delete(s.images, id)
		}
	}
}

func publicAccount(acc *store.WechatAccount) game.YYBAccount {
	out := game.YYBAccount{ID: acc.ID, OpenID: acc.OpenID, UpdatedAt: acc.UpdatedAt}
	if acc.UIN != nil {
		out.UIN = *acc.UIN
	}
	if acc.Nickname != nil {
		out.Nickname = *acc.Nickname
	}
	if acc.Avatar != nil {
		out.Avatar = *acc.Avatar
	}
	if acc.Status != nil {
		out.Status = *acc.Status
	}
	return out
}

func unwrapUserInfo(raw map[string]any) map[string]any {
	if raw == nil {
		return nil
	}
	for _, k := range []string{"data", "result", "user_info", "userInfo"} {
		inner, ok := raw[k].(map[string]any)
		if !ok {
			continue
		}
		if pickNickname(inner, "") != "" || pickAvatarURL(inner) != "" {
			return inner
		}
	}
	return raw
}

func pickNickname(userInfo map[string]any, fallback string) string {
	for _, k := range []string{"nick_name", "nickname", "nickName", "name"} {
		if s := stringFromAny(userInfo[k]); s != "" {
			return s
		}
	}
	return fallback
}

func pickAvatarURL(userInfo map[string]any) string {
	for _, k := range []string{"head_img_url", "head_url", "headimgurl", "avatar", "avatar_url"} {
		if s := stringFromAny(userInfo[k]); s != "" {
			return s
		}
	}
	return ""
}

func stringFromAny(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func codeOut(openid string, res map[string]any) (game.YYBCodeOut, error) {
	code, _ := res["code"].(string)
	if code == "" {
		return game.YYBCodeOut{}, fmt.Errorf("getCode 回了空 code")
	}
	return game.YYBCodeOut{Code: code, OpenID: openid}, nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
