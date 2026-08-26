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
		var in game.RefreshIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		lands, err := hub.Refresh(r.Context(), in)
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
	mux.HandleFunc("POST /Farm/Steal", func(w http.ResponseWriter, r *http.Request) {
		var in game.HarvestIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.Steal(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Farm/Remove", func(w http.ResponseWriter, r *http.Request) {
		var in game.RemoveIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		lands, err := hub.Remove(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"lands": lands}))
	})
	mux.HandleFunc("POST /Farm/Unlock", func(w http.ResponseWriter, r *http.Request) {
		var in game.LandIDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		land, err := hub.Unlock(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(land))
	})
	mux.HandleFunc("POST /Farm/Upgrade", func(w http.ResponseWriter, r *http.Request) {
		var in game.LandIDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		land, err := hub.Upgrade(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(land))
	})
	mux.HandleFunc("POST /Farm/Farming", func(w http.ResponseWriter, r *http.Request) {
		var in game.LandOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		lands, err := hub.Farming(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"lands": lands}))
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
	mux.HandleFunc("POST /Friend/Enter", func(w http.ResponseWriter, r *http.Request) {
		var in game.EnterIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		visit, err := hub.Enter(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(visit))
	})
	mux.HandleFunc("POST /Friend/Applications", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.Applications(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"applications": list}))
	})
	mux.HandleFunc("POST /Friend/Accept", func(w http.ResponseWriter, r *http.Request) {
		var in game.GIDsIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		friends, err := hub.Accept(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"friends": friends}))
	})
	mux.HandleFunc("POST /Friend/Reject", func(w http.ResponseWriter, r *http.Request) {
		var in game.GIDsIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Reject(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Friend/Delete", func(w http.ResponseWriter, r *http.Request) {
		var in game.EnterIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.DeleteFriend(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Share/Check", func(w http.ResponseWriter, r *http.Request) {
		ok, err := hub.ShareCheck(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"can_share": ok}))
	})
	mux.HandleFunc("POST /Share/Claim", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.ShareClaim(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Friend/Leave", func(w http.ResponseWriter, r *http.Request) {
		var in game.EnterIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.Leave(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Bag/Get", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.Bag(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Bag/Sell", func(w http.ResponseWriter, r *http.Request) {
		var in game.SellIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.Sell(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Bag/Use", func(w http.ResponseWriter, r *http.Request) {
		var in game.UseIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.Use(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Shop/List", func(w http.ResponseWriter, r *http.Request) {
		shops, err := hub.Shops(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"shops": shops}))
	})
	mux.HandleFunc("POST /Shop/Goods", func(w http.ResponseWriter, r *http.Request) {
		var in game.ShopIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		goods, err := hub.Goods(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"goods": goods}))
	})
	mux.HandleFunc("POST /Shop/Buy", func(w http.ResponseWriter, r *http.Request) {
		var in game.BuyIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.Buy(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Weather/Status", func(w http.ResponseWriter, r *http.Request) {
		weather, err := hub.Weather(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(weather))
	})
	mux.HandleFunc("POST /Weather/Today", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.TodayWeather(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"weathers": list}))
	})
	mux.HandleFunc("POST /Task/List", func(w http.ResponseWriter, r *http.Request) {
		board, err := hub.Tasks(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(board))
	})
	mux.HandleFunc("POST /Task/Claim", func(w http.ResponseWriter, r *http.Request) {
		var in game.ClaimIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimTask(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Task/ClaimDaily", func(w http.ResponseWriter, r *http.Request) {
		var in game.DailyIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimDaily(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Email/List", func(w http.ResponseWriter, r *http.Request) {
		var in game.EmailBoxIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		list, err := hub.Emails(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"emails": list}))
	})
	mux.HandleFunc("POST /Email/Read", func(w http.ResponseWriter, r *http.Request) {
		var in game.EmailIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		mail, err := hub.ReadEmail(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(mail))
	})
	mux.HandleFunc("POST /Email/Claim", func(w http.ResponseWriter, r *http.Request) {
		var in game.EmailIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimEmail(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Email/ClaimAll", func(w http.ResponseWriter, r *http.Request) {
		var in game.EmailBoxIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimAllEmail(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/Signin", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.Signin(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/ClaimProgress", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimProgress(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/ShopBuy", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityShopIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ShopBuy(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/ClaimMega", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimMega(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/TechSubmit", func(w http.ResponseWriter, r *http.Request) {
		var in game.TechIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.TechSubmit(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/Lottery", func(w http.ResponseWriter, r *http.Request) {
		var in game.LotteryIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.Lottery(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Activity/BrewStart", func(w http.ResponseWriter, r *http.Request) {
		var in game.BrewStartIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.BrewStart(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Activity/BrewStep", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.BrewStep(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Activity/BrewClaim", func(w http.ResponseWriter, r *http.Request) {
		var in game.BrewClaimIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.BrewClaim(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/ClaimRecall", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimRecall(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/ClaimReturn", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimReturn(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/ClaimInvite", func(w http.ResponseWriter, r *http.Request) {
		var in game.InviteIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimInvite(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/ClaimNewcomer", func(w http.ResponseWriter, r *http.Request) {
		var in game.ActivityOpIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimNewcomer(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Activity/SendGift", func(w http.ResponseWriter, r *http.Request) {
		var in game.GiftIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		n, err := hub.SendGift(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"total_send_count": n}))
	})
	mux.HandleFunc("POST /Mystery/Status", func(w http.ResponseWriter, r *http.Request) {
		shop, err := hub.Mystery(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(shop))
	})
	mux.HandleFunc("POST /Mystery/Buy", func(w http.ResponseWriter, r *http.Request) {
		var in game.IDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.MysteryBuy(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Mystery/Leave", func(w http.ResponseWriter, r *http.Request) {
		if err := hub.MysteryLeave(r.Context()); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Activity/List", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.Activities(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"activities": list}))
	})
	mux.HandleFunc("POST /Season/Info", func(w http.ResponseWriter, r *http.Request) {
		season, err := hub.Season(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(season))
	})
	mux.HandleFunc("POST /Season/ClaimPass", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.ClaimPass(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Mall/List", func(w http.ResponseWriter, r *http.Request) {
		var in game.MallIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		list, err := hub.Mall(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"products": list}))
	})
	mux.HandleFunc("POST /Mall/Buy", func(w http.ResponseWriter, r *http.Request) {
		var in game.MallBuyIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.MallBuy(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Mall/MonthCard", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.MonthCards(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"cards": list}))
	})
	mux.HandleFunc("POST /Mall/ClaimMonthCard", func(w http.ResponseWriter, r *http.Request) {
		var in game.IDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimMonthCard(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /RedPacket/List", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.RedPackets(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"packets": list}))
	})
	mux.HandleFunc("POST /RedPacket/Claim", func(w http.ResponseWriter, r *http.Request) {
		var in game.IDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimRedPacket(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /RedPacket/ClaimAll", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.ClaimAllRedPackets(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Visit/Logs", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.VisitLogs(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"logs": list}))
	})
	mux.HandleFunc("POST /Album/List", func(w http.ResponseWriter, r *http.Request) {
		var in game.AlbumIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		album, err := hub.Album(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(album))
	})
	mux.HandleFunc("POST /Album/Claim", func(w http.ResponseWriter, r *http.Request) {
		var in game.AlbumIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimAlbum(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Dog/Info", func(w http.ResponseWriter, r *http.Request) {
		yard, err := hub.Dog(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(yard))
	})
	mux.HandleFunc("POST /Dog/Feed", func(w http.ResponseWriter, r *http.Request) {
		var in game.FeedIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		left, err := hub.Feed(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"food_left": left}))
	})
	mux.HandleFunc("POST /Dog/ClaimGifts", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.ClaimDogGifts(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Dog/Logs", func(w http.ResponseWriter, r *http.Request) {
		var in game.PageIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		logs, err := hub.DogLogs(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"logs": logs}))
	})
	mux.HandleFunc("POST /Dog/Deploy", func(w http.ResponseWriter, r *http.Request) {
		var in game.IDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		out, err := hub.DeployDog(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Dog/Withdraw", func(w http.ResponseWriter, r *http.Request) {
		out, err := hub.WithdrawDog(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(out))
	})
	mux.HandleFunc("POST /Dog/Activate", func(w http.ResponseWriter, r *http.Request) {
		var in game.IDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		dog, err := hub.ActivateDog(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(dog))
	})
	mux.HandleFunc("POST /Bulletin/List", func(w http.ResponseWriter, r *http.Request) {
		var in game.PageIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		list, err := hub.Bulletins(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"bulletins": list}))
	})
	mux.HandleFunc("POST /Bulletin/Read", func(w http.ResponseWriter, r *http.Request) {
		var in game.IDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		detail, err := hub.ReadBulletin(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(detail))
	})
	mux.HandleFunc("POST /Mutant/List", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.Mutants(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"mutants": list}))
	})
	mux.HandleFunc("POST /Career/Info", func(w http.ResponseWriter, r *http.Request) {
		info, err := hub.Career(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(info))
	})
	mux.HandleFunc("POST /Rank/List", func(w http.ResponseWriter, r *http.Request) {
		var in game.RankIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		board, err := hub.Ranks(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(board))
	})
	mux.HandleFunc("POST /Avatar/Owned", func(w http.ResponseWriter, r *http.Request) {
		var in game.TypeIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		list, err := hub.Avatars(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"avatars": list}))
	})
	mux.HandleFunc("POST /Avatar/Equipped", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.EquippedAvatars(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"avatars": list}))
	})
	mux.HandleFunc("POST /Avatar/Equip", func(w http.ResponseWriter, r *http.Request) {
		var in game.AvatarEquipIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		av, err := hub.EquipAvatar(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(av))
	})
	mux.HandleFunc("POST /Skin/Owned", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.Skins(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"skins": list}))
	})
	mux.HandleFunc("POST /Skin/Equipped", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.EquippedSkins(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"skins": list}))
	})
	mux.HandleFunc("POST /Skin/Equip", func(w http.ResponseWriter, r *http.Request) {
		var in game.SkinEquipIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		if err := hub.EquipSkin(r.Context(), in); err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(nil))
	})
	mux.HandleFunc("POST /Drop/List", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.Drops(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"drops": list}))
	})
	mux.HandleFunc("POST /Solar/List", func(w http.ResponseWriter, r *http.Request) {
		list, err := hub.SolarTerms(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"terms": list}))
	})
	mux.HandleFunc("POST /Solar/Claim", func(w http.ResponseWriter, r *http.Request) {
		var in game.IDIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimSolar(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Solar/ClaimAll", func(w http.ResponseWriter, r *http.Request) {
		items, err := hub.ClaimAllSolar(r.Context())
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Achieve/View", func(w http.ResponseWriter, r *http.Request) {
		var in game.AchieveIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		view, err := hub.AchieveView(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(view))
	})
	mux.HandleFunc("POST /Achieve/ClaimGoal", func(w http.ResponseWriter, r *http.Request) {
		var in game.AchieveGoalIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimAchieveGoal(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
	})
	mux.HandleFunc("POST /Achieve/ClaimLevel", func(w http.ResponseWriter, r *http.Request) {
		var in game.AchieveIn
		if err := decodeJSON(r, &in); err != nil {
			writeReply(w, Fail(400, err.Error()))
			return
		}
		items, err := hub.ClaimAchieveLevel(r.Context(), in)
		if err != nil {
			code, msg := failCode(err)
			writeReply(w, Fail(code, msg))
			return
		}
		writeReply(w, OK(map[string]any{"items": items}))
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
