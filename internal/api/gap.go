package api

import (
	"net/http"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func registerGap(mux *http.ServeMux, hub *Hub) {
	mux.HandleFunc("POST /Farm/CleanEvents", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.CleanFarmEvents(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Friend/BlockApplications", func(w http.ResponseWriter, r *http.Request) {
		var in game.BlockAppsIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.BlockApplications(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Friend/WXRecommend", func(w http.ResponseWriter, r *http.Request) {
		var in game.WXRecommendIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.WXRecommend(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Friend/WXRecommendPage", func(w http.ResponseWriter, r *http.Request) {
		var in game.WXRecommendPageIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.WXRecommendPage(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Friend/ApplyWX", func(w http.ResponseWriter, r *http.Request) {
		var in game.GIDsIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.ApplyWXFriends(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Skin/EquipSet", func(w http.ResponseWriter, r *http.Request) {
		var in game.SkinSetIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.EquipSkinSet(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Skin/SetEffect", func(w http.ResponseWriter, r *http.Request) {
		var in game.SkinSetEffectIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.SetSkinSetEffect(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Skin/Sets", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.SkinSets(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"sets": list}))
	})
	mux.HandleFunc("POST /Season/BuyPass", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.BuyPass(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Season/MarkOpening", func(w http.ResponseWriter, r *http.Request) {
		ok, err := hub.MarkSeasonOpening(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"success": ok}))
	})
	mux.HandleFunc("POST /QQGroup/AuthGroups", func(w http.ResponseWriter, r *http.Request) {
		var in game.CookiesIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.QQAuthGroups(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQGroup/Recommend", func(w http.ResponseWriter, r *http.Request) {
		var in game.QQRecommendIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.QQRecommendGroups(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQGroup/Bind", func(w http.ResponseWriter, r *http.Request) {
		var in game.QQBindIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.QQBind(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQGroup/Leave", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.QQLeave(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQGroup/Community", func(w http.ResponseWriter, r *http.Request) {
		var in game.PageIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.QQCommunity(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQGroup/BindInfo", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.QQBindInfo(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /QQGroup/ClaimReward", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.QQClaimReward(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /QQGroup/RevokeAuth", func(w http.ResponseWriter, r *http.Request) {
		var in game.QQRevokeIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.QQRevokeAuth(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Gift/UseToken", func(w http.ResponseWriter, r *http.Request) {
		var in game.UIDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.UseGiftToken(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Gift/History", func(w http.ResponseWriter, r *http.Request) {
		var in game.GiftHistoryIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.GiftHistory(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Gift/TransferStatus", func(w http.ResponseWriter, r *http.Request) {
		var in game.TransferIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.TransferStatus(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Gift/CancelTransfer", func(w http.ResponseWriter, r *http.Request) {
		var in game.UIDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.CancelTransfer(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Follow/Status", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.FollowGiftStatus(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Follow/Set", func(w http.ResponseWriter, r *http.Request) {
		var in game.FollowGiftIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.SetFollowGift(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Follow/Claim", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.ClaimFollowGift(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Recharge/Config", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.RechargeBonus(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Recharge/Data", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.RechargeBonusData(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /User/SetDisplay", func(w http.ResponseWriter, r *http.Request) {
		var in game.DisplayIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.SetDisplay(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /User/Settings", func(w http.ResponseWriter, r *http.Request) {
		var in game.SettingsKeysIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.GetSettings(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /User/SetSettings", func(w http.ResponseWriter, r *http.Request) {
		var in game.UserSettings
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.SetSettings(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /User/DeleteAccount", func(w http.ResponseWriter, r *http.Request) {
		var in game.DeleteAccountIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.DeleteAccount(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /User/DecryptOpenData", func(w http.ResponseWriter, r *http.Request) {
		var in game.DecryptIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		data, err := hub.DecryptOpenData(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"json_data": data}))
	})
	mux.HandleFunc("POST /User/QQRecommendAuth", func(w http.ResponseWriter, r *http.Request) {
		var in game.QQAuthIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.SetQQRecommendAuth(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /User/ReportFlow", func(w http.ResponseWriter, r *http.Request) {
		var in game.ReportFlowIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.ReportFlow(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /User/BatchReportFlow", func(w http.ResponseWriter, r *http.Request) {
		var in game.BatchReportFlowIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.BatchReportFlow(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /User/Report", func(w http.ResponseWriter, r *http.Request) {
		var in game.ReportUserIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.ReportUser(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
}
