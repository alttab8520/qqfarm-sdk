package api

import (
	"encoding/json"
	"net/http"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func NewMux(hub *Hub) *http.ServeMux {
	if hub == nil {
		hub = NewHub(nil)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /docs", serveDocs)
	mux.HandleFunc("GET /openapi.yaml", serveOpenAPI)
	mux.HandleFunc("POST /System/Ping", writeOK(map[string]any{"pong": true}))
	mux.HandleFunc("POST /User/Login", func(w http.ResponseWriter, r *http.Request) {
		var in game.LoginIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if in.Code == "" {
			writeReply(w, Fail(400, "code 不能为空"))
			return
		}
		user, err := hub.Login(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(user))
	})
	mux.HandleFunc("POST /System/Status", func(w http.ResponseWriter, r *http.Request) {
		st, err := hub.Status()
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(st))
	})
	mux.HandleFunc("POST /User/Logout", func(w http.ResponseWriter, r *http.Request) {
		if err := hub.Logout(); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /User/GetInfo", func(w http.ResponseWriter, r *http.Request) {
		user, err := hub.Info()
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(user))
	})
	mux.HandleFunc("POST /Farm/Refresh", func(w http.ResponseWriter, r *http.Request) {
		lands, err := hub.Refresh(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"lands": lands}))
	})
	mux.HandleFunc("POST /Farm/Harvest", func(w http.ResponseWriter, r *http.Request) {
		var in game.HarvestIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.Harvest(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Farm/Plant", func(w http.ResponseWriter, r *http.Request) {
		var in game.PlantIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Plant(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Farm/Water", func(w http.ResponseWriter, r *http.Request) {
		var in game.LandOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Water(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Farm/Weed", func(w http.ResponseWriter, r *http.Request) {
		var in game.LandOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Weed(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Farm/Bug", func(w http.ResponseWriter, r *http.Request) {
		var in game.LandOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Bug(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Farm/Fertilize", func(w http.ResponseWriter, r *http.Request) {
		var in game.FertilizeIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Fertilize(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Friend/GetList", func(w http.ResponseWriter, r *http.Request) {
		friends, err := hub.Friends(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"friends": friends}))
	})
	mux.HandleFunc("POST /Friend/Help", func(w http.ResponseWriter, r *http.Request) {
		var in game.HelpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Help(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	return mux
}

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil || r.ContentLength == 0 {
		return nil
	}
	dec := json.NewDecoder(r.Body)
	return dec.Decode(dst)
}

func writeOK(data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeReply(w, OK(data))
	}
}

func writeReply(w http.ResponseWriter, body Reply) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(body)
}
