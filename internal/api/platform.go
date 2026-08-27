package api

import (
	"net/http"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func registerPlatform(mux *http.ServeMux, hub *Hub) {
	mux.HandleFunc("POST /QQVip/DailyStatus", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.QQVipDailyStatus(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQVip/ClaimDaily", func(w http.ResponseWriter, r *http.Request) {
		var in game.QQVipClaimDailyIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.QQVipClaimDaily(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQVip/Refresh", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.QQVipRefresh(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQVip/ClaimRewards", func(w http.ResponseWriter, r *http.Request) {
		var in game.QQVipClaimRewardsIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.QQVipClaimRewards(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQVip/RewardsStatus", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.QQVipRewardsStatus(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQVip/MarkRedpoint", func(w http.ResponseWriter, r *http.Request) {
		if err := hub.QQVipMarkRedpoint(r.Context()); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Marquee/List", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.Marquee(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /SystemOpen/Unlocked", func(w http.ResponseWriter, r *http.Request) {
		var in game.SystemOpenIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.SystemUnlocked(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Mutant/OpenInfo", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.MutantOpenInfo(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Subscribe/QQ", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.QQSubscribe(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Subscribe/WX", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.WXSubscribe(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Subscribe/SetWX", func(w http.ResponseWriter, r *http.Request) {
		var in game.WXSubscribeIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.SetWXSubscribe(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Moderate/Text", func(w http.ResponseWriter, r *http.Request) {
		var in game.ModerateTextIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.ModerateText(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Moderate/BatchText", func(w http.ResponseWriter, r *http.Request) {
		var in game.ModerateTextBatchIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.BatchModerateText(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Moderate/Pic", func(w http.ResponseWriter, r *http.Request) {
		var in game.ModeratePicIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.ModeratePic(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Moderate/BatchPic", func(w http.ResponseWriter, r *http.Request) {
		var in game.ModeratePicBatchIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.BatchModeratePic(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
}
