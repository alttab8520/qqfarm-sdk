package api

import (
	"net/http"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func registerYYB(mux *http.ServeMux, hub *Hub) {
	mux.HandleFunc("POST /YYB/Accounts", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.YYBAccounts(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"accounts": list}))
	})
	mux.HandleFunc("POST /YYB/QR", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.YYBCreateQR(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("GET /YYB/Image", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("session_id")
		if id == "" {
			writeReply(w, Fail(400, "session_id 不能为空"))
			return
		}
		raw, err := hub.YYBImage(id)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write(raw)
	})
	mux.HandleFunc("POST /YYB/Poll", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBSessionIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if in.SessionID == "" {
			writeReply(w, Fail(400, "session_id 不能为空"))
			return
		}
		out, err := hub.YYBPoll(r.Context(), in.SessionID)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /YYB/Confirm", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBSessionIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if in.SessionID == "" {
			writeReply(w, Fail(400, "session_id 不能为空"))
			return
		}
		out, err := hub.YYBConfirm(r.Context(), in.SessionID)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /YYB/Refresh", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBRefIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.YYBRefresh(r.Context(), in.Ref)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /YYB/Code", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBRefIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.YYBCode(r.Context(), in.Ref, in.AppID)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /YYB/Login", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBRefIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		user, err := hub.YYBLogin(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(user))
	})
	mux.HandleFunc("POST /YYB/Delete", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBRefIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if in.Ref == "" {
			writeReply(w, Fail(400, "ref 不能为空"))
			return
		}
		out, err := hub.YYBDelete(r.Context(), in.Ref)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /YYB/Profile", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBRefIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.YYBProfile(r.Context(), in.Ref)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /YYB/Phone", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBRefIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.YYBPhone(r.Context(), in.Ref, in.AppID)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /YYB/WXData", func(w http.ResponseWriter, r *http.Request) {
		var in game.YYBWXDataIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if in.Payload == nil {
			writeReply(w, Fail(400, "payload 不能为空"))
			return
		}
		out, err := hub.YYBWXData(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
}
