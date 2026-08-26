package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var stubs = []string{
	"/User/Login",
	"/User/GetInfo",
	"/Farm/Refresh",
	"/Farm/Harvest",
	"/Farm/Plant",
	"/Friend/GetList",
	"/Friend/Help",
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) Reply {
	t.Helper()
	var body Reply
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v body=%s", err, rec.Body.String())
	}
	return body
}

func TestPing(t *testing.T) {
	rec := httptest.NewRecorder()
	NewMux().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/System/Ping", nil))
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	body := decode(t, rec)
	if body.Code != 0 {
		t.Fatalf("code %d", body.Code)
	}
	data, _ := body.Data.(map[string]any)
	if data["pong"] != true {
		t.Fatalf("data %#v", body.Data)
	}
}

func TestStubsNotReady(t *testing.T) {
	mux := NewMux()
	for _, path := range stubs {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, nil))
		body := decode(t, rec)
		if rec.Code != 200 || body.Code != 501 || body.Msg != notReadyMsg {
			t.Fatalf("%s status=%d body=%+v", path, rec.Code, body)
		}
	}
}

func TestDocs(t *testing.T) {
	rec := httptest.NewRecorder()
	NewMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
}
