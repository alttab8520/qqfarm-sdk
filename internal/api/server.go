package api

import (
	"encoding/json"
	"net/http"
)

const notReadyMsg = "未接入游戏会话"

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /docs", serveDocs)
	mux.HandleFunc("GET /openapi.yaml", serveOpenAPI)
	mux.HandleFunc("POST /System/Ping", writeJSON(OK(map[string]any{"pong": true})))
	for _, path := range []string{
		"/User/Login",
		"/User/GetInfo",
		"/Farm/Refresh",
		"/Farm/Harvest",
		"/Farm/Plant",
		"/Friend/GetList",
		"/Friend/Help",
	} {
		mux.HandleFunc("POST "+path, writeJSON(Fail(501, notReadyMsg)))
	}
	return mux
}

func writeJSON(body Reply) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(body)
	}
}
