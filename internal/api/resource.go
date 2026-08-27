package api

import (
	"net/http"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func registerResource(mux *http.ServeMux, hub *Hub) {
	mux.HandleFunc("POST /Resource/Lookup", func(w http.ResponseWriter, r *http.Request) {
		var in game.ResLookupIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if len(in.IDs) == 0 {
			writeReply(w, Fail(400, "ids 不能为空"))
			return
		}
		out, err := hub.ResLookup(in)
		if err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Resource/Items", func(w http.ResponseWriter, r *http.Request) {
		var in game.ResListIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.ResList(in)
		if err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Resource/Tables", func(w http.ResponseWriter, r *http.Request) {
		writeReply(w, OK(hub.ResTables()))
	})
	mux.HandleFunc("POST /Resource/Refresh", func(w http.ResponseWriter, r *http.Request) {
		var in game.ResRefreshIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.ResRefresh(r.Context(), in)
		if err != nil {
			writeReply(w, Fail(502, err.Error()))
			return
		}
		writeReply(w, OK(out))
	})
}
