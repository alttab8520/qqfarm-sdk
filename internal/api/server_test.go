package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

type fakeSession struct {
	user    game.User
	lands   []game.Land
	items   []game.Item
	friends []game.Friend
}

func (f *fakeSession) Login(_ context.Context, in game.LoginIn) (game.User, error) {
	f.user = game.User{GID: 7, Name: "测", OpenID: in.OpenID}
	return f.user, nil
}
func (f *fakeSession) Info() (game.User, error) {
	if f.user.GID == 0 {
		return game.User{}, game.ErrNotLogin
	}
	return f.user, nil
}
func (f *fakeSession) Refresh(context.Context) ([]game.Land, error) {
	return f.lands, nil
}
func (f *fakeSession) Harvest(context.Context, game.HarvestIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Plant(context.Context, game.PlantIn) error { return nil }
func (f *fakeSession) Friends(context.Context) ([]game.Friend, error) {
	return f.friends, nil
}
func (f *fakeSession) Help(context.Context, game.HelpIn) error { return nil }
func (f *fakeSession) Water(context.Context, game.LandOpIn) error {
	return nil
}
func (f *fakeSession) Weed(context.Context, game.LandOpIn) error { return nil }
func (f *fakeSession) Bug(context.Context, game.LandOpIn) error  { return nil }
func (f *fakeSession) Fertilize(context.Context, game.FertilizeIn) error {
	return nil
}
func (f *fakeSession) Status() (game.Status, error) {
	if f.user.GID == 0 {
		return game.Status{}, game.ErrNotLogin
	}
	return game.Status{LoggedIn: true, User: f.user}, nil
}
func (f *fakeSession) Close() error {
	f.user = game.User{}
	return nil
}

func testClient(t *testing.T, sess *fakeSession) *http.ServeMux {
	t.Helper()
	return NewMux(NewHub(func() game.Session { return sess }))
}

func post(t *testing.T, mux *http.ServeMux, path string, body any) Reply {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, &buf))
	var out Reply
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("json %v %s", err, rec.Body.String())
	}
	return out
}

func TestPing(t *testing.T) {
	out := post(t, testClient(t, &fakeSession{}), "/System/Ping", nil)
	if out.Code != 0 {
		t.Fatalf("%+v", out)
	}
}

func TestLoginRequiresCode(t *testing.T) {
	out := post(t, testClient(t, &fakeSession{}), "/User/Login", map[string]string{})
	if out.Code != 400 {
		t.Fatalf("%+v", out)
	}
}

func TestGetInfoRequiresLogin(t *testing.T) {
	out := post(t, testClient(t, &fakeSession{}), "/User/GetInfo", nil)
	if out.Code != 401 {
		t.Fatalf("%+v", out)
	}
}

func TestLoginThenInfo(t *testing.T) {
	mux := testClient(t, &fakeSession{})
	out := post(t, mux, "/User/Login", game.LoginIn{Code: "abc", OpenID: "o1"})
	if out.Code != 0 {
		t.Fatalf("login %+v", out)
	}
	out = post(t, mux, "/User/GetInfo", nil)
	if out.Code != 0 {
		t.Fatalf("info %+v", out)
	}
}

func TestStatusRequiresLogin(t *testing.T) {
	out := post(t, testClient(t, &fakeSession{}), "/System/Status", nil)
	if out.Code != 401 {
		t.Fatalf("%+v", out)
	}
}

func TestLogout(t *testing.T) {
	mux := testClient(t, &fakeSession{})
	if out := post(t, mux, "/User/Login", game.LoginIn{Code: "abc", OpenID: "o1"}); out.Code != 0 {
		t.Fatalf("login %+v", out)
	}
	if out := post(t, mux, "/User/Logout", nil); out.Code != 0 {
		t.Fatalf("logout %+v", out)
	}
	if out := post(t, mux, "/User/GetInfo", nil); out.Code != 401 {
		t.Fatalf("after logout %+v", out)
	}
}

func TestDocs(t *testing.T) {
	rec := httptest.NewRecorder()
	NewMux(NewHub(func() game.Session { return &fakeSession{} })).ServeHTTP(
		rec, httptest.NewRequest(http.MethodGet, "/docs", nil),
	)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
}
