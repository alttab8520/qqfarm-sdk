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

var _ game.Session = (*fakeSession)(nil)

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
func (f *fakeSession) Refresh(context.Context, game.RefreshIn) ([]game.Land, error) {
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
func (f *fakeSession) Applications(context.Context) ([]game.Application, error) {
	return []game.Application{{GID: 3, Name: "申"}}, nil
}
func (f *fakeSession) Accept(context.Context, game.GIDsIn) ([]game.Friend, error) {
	return f.friends, nil
}
func (f *fakeSession) Reject(context.Context, game.GIDsIn) error { return nil }
func (f *fakeSession) DeleteFriend(context.Context, game.EnterIn) error {
	return nil
}
func (f *fakeSession) ShareCheck(context.Context) (bool, error) { return true, nil }
func (f *fakeSession) ShareClaim(context.Context) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Water(context.Context, game.LandOpIn) error {
	return nil
}
func (f *fakeSession) Weed(context.Context, game.LandOpIn) error { return nil }
func (f *fakeSession) Bug(context.Context, game.LandOpIn) error  { return nil }
func (f *fakeSession) Fertilize(context.Context, game.FertilizeIn) error {
	return nil
}
func (f *fakeSession) Enter(_ context.Context, in game.EnterIn) (game.Visit, error) {
	return game.Visit{Host: game.User{GID: in.GID}, Lands: f.lands}, nil
}
func (f *fakeSession) Leave(context.Context, game.EnterIn) error { return nil }
func (f *fakeSession) Bag(context.Context) ([]game.BagItem, error) {
	return []game.BagItem{{ID: 20002, Count: 3}}, nil
}
func (f *fakeSession) Sell(context.Context, game.SellIn) ([]game.Item, error) {
	return []game.Item{{ID: 1001, Count: 10}}, nil
}
func (f *fakeSession) Use(context.Context, game.UseIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Shops(context.Context) ([]game.Shop, error) {
	return []game.Shop{{ID: 2, Name: "种子"}}, nil
}
func (f *fakeSession) Goods(context.Context, game.ShopIn) ([]game.Goods, error) {
	return []game.Goods{{ID: 15, ItemID: 20002, Price: 50, Unlocked: true}}, nil
}
func (f *fakeSession) Buy(context.Context, game.BuyIn) (game.BuyOut, error) {
	return game.BuyOut{Items: []game.Item{{ID: 20002, Count: 1}}}, nil
}
func (f *fakeSession) Steal(context.Context, game.HarvestIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Remove(context.Context, game.RemoveIn) ([]game.Land, error) {
	return f.lands, nil
}
func (f *fakeSession) Unlock(_ context.Context, in game.LandIDIn) (game.Land, error) {
	return game.Land{ID: in.LandID, Unlocked: true}, nil
}
func (f *fakeSession) Upgrade(_ context.Context, in game.LandIDIn) (game.Land, error) {
	return game.Land{ID: in.LandID, Level: 2}, nil
}
func (f *fakeSession) Farming(context.Context, game.LandOpIn) ([]game.Land, error) {
	return f.lands, nil
}
func (f *fakeSession) Weather(context.Context) (game.Weather, error) {
	return game.Weather{Type: 1, Active: true}, nil
}
func (f *fakeSession) TodayWeather(context.Context) ([]game.Weather, error) {
	return []game.Weather{{Type: 1, Name: "晴"}}, nil
}
func (f *fakeSession) Tasks(context.Context) (game.TaskBoard, error) {
	return game.TaskBoard{Daily: []game.Task{{ID: 1, Desc: "浇水"}}}, nil
}
func (f *fakeSession) ClaimTask(context.Context, game.ClaimIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimDaily(context.Context, game.DailyIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Emails(context.Context, game.EmailBoxIn) ([]game.Email, error) {
	return []game.Email{{ID: "1", Title: "奖励", HasReward: true}}, nil
}
func (f *fakeSession) ReadEmail(_ context.Context, in game.EmailIn) (game.EmailDetail, error) {
	return game.EmailDetail{ID: in.ID, Title: "奖励"}, nil
}
func (f *fakeSession) ClaimEmail(context.Context, game.EmailIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimAllEmail(context.Context, game.EmailBoxIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Activities(context.Context) ([]game.Activity, error) {
	return []game.Activity{{ID: 8, Name: "签到"}}, nil
}
func (f *fakeSession) Season(context.Context) (game.Season, error) {
	return game.Season{ID: 1, Name: "春", Pass: game.BattlePass{Level: 3}}, nil
}
func (f *fakeSession) ClaimPass(context.Context) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Mall(context.Context, game.MallIn) ([]game.Product, error) {
	return []game.Product{{ID: 1001, Name: "每日福利", Available: true}}, nil
}
func (f *fakeSession) MallBuy(context.Context, game.MallBuyIn) (game.BuyOut, error) {
	return game.BuyOut{Items: f.items}, nil
}
func (f *fakeSession) MonthCards(context.Context) ([]game.MonthCard, error) {
	return []game.MonthCard{{ID: 2001, Claimable: true}}, nil
}
func (f *fakeSession) ClaimMonthCard(context.Context, game.IDIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) RedPackets(context.Context) ([]game.RedPacket, error) {
	return []game.RedPacket{{ID: 1, CanClaim: true}}, nil
}
func (f *fakeSession) ClaimRedPacket(context.Context, game.IDIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimAllRedPackets(context.Context) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) VisitLogs(context.Context) ([]game.VisitLog, error) {
	return []game.VisitLog{{GID: 2, Action: 4, Name: "客"}}, nil
}
func (f *fakeSession) Album(context.Context, game.AlbumIn) (game.Album, error) {
	return game.Album{Items: []game.AlbumItem{{FruitID: 40061, Unlocked: true}}}, nil
}
func (f *fakeSession) ClaimAlbum(context.Context, game.AlbumIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Dog(context.Context) (game.DogYard, error) {
	return game.DogYard{Deployed: 1, FoodLeft: 60}, nil
}
func (f *fakeSession) Feed(context.Context, game.FeedIn) (int64, error) { return 120, nil }
func (f *fakeSession) ClaimDogGifts(context.Context) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) DogLogs(context.Context, game.PageIn) ([]game.ProtectLog, error) {
	return []game.ProtectLog{{GID: 2, Name: "贼"}}, nil
}
func (f *fakeSession) DeployDog(context.Context, game.IDIn) (game.DeployOut, error) {
	return game.DeployOut{Deployed: 1}, nil
}
func (f *fakeSession) WithdrawDog(context.Context) (game.DeployOut, error) {
	return game.DeployOut{Previous: 1}, nil
}
func (f *fakeSession) ActivateDog(_ context.Context, in game.IDIn) (game.Dog, error) {
	return game.Dog{ID: in.ID, Activated: true}, nil
}
func (f *fakeSession) Bulletins(context.Context, game.PageIn) ([]game.Bulletin, error) {
	return []game.Bulletin{{ID: 1, Title: "告"}}, nil
}
func (f *fakeSession) ReadBulletin(context.Context, game.IDIn) (game.BulletinDetail, error) {
	return game.BulletinDetail{Title: "告", Content: "文"}, nil
}
func (f *fakeSession) Mutants(context.Context) ([]game.Mutant, error) {
	return []game.Mutant{{ID: 1, Unlocked: true}}, nil
}
func (f *fakeSession) Career(context.Context) (game.Career, error) {
	return game.Career{Name: "测", Harvested: 3}, nil
}
func (f *fakeSession) Ranks(context.Context, game.RankIn) (game.RankBoard, error) {
	return game.RankBoard{Items: []game.RankItem{{GID: 1, Rank: 1}}}, nil
}
func (f *fakeSession) Avatars(context.Context, game.TypeIn) ([]game.Avatar, error) {
	return []game.Avatar{{ID: 1}}, nil
}
func (f *fakeSession) EquippedAvatars(context.Context) ([]game.Avatar, error) {
	return []game.Avatar{{ID: 1}}, nil
}
func (f *fakeSession) EquipAvatar(_ context.Context, in game.AvatarEquipIn) (game.Avatar, error) {
	return game.Avatar{ID: in.ID}, nil
}
func (f *fakeSession) Skins(context.Context) ([]game.Skin, error) {
	return []game.Skin{{ID: 1}}, nil
}
func (f *fakeSession) EquippedSkins(context.Context) ([]game.Skin, error) {
	return []game.Skin{{ID: 1, Equipped: true}}, nil
}
func (f *fakeSession) EquipSkin(context.Context, game.SkinEquipIn) error { return nil }
func (f *fakeSession) Drops(context.Context) ([]game.Drop, error) {
	return []game.Drop{{ID: 1}}, nil
}
func (f *fakeSession) SolarTerms(context.Context) ([]game.SolarTerm, error) {
	return []game.SolarTerm{{ID: 1, Status: 2}}, nil
}
func (f *fakeSession) ClaimSolar(context.Context, game.IDIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimAllSolar(context.Context) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) AchieveView(context.Context, game.AchieveIn) (game.AchieveScope, error) {
	return game.AchieveScope{ID: 1}, nil
}
func (f *fakeSession) ClaimAchieveGoal(context.Context, game.AchieveGoalIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimAchieveLevel(context.Context, game.AchieveIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Signin(context.Context, game.ActivityOpIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimProgress(context.Context, game.ActivityOpIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ShopBuy(context.Context, game.ActivityShopIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimMega(context.Context, game.ActivityOpIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) TechSubmit(context.Context, game.TechIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) Lottery(context.Context, game.LotteryIn) (game.LotteryOut, error) {
	return game.LotteryOut{Items: f.items}, nil
}
func (f *fakeSession) BrewStart(context.Context, game.BrewStartIn) (game.BrewStartOut, error) {
	return game.BrewStartOut{Value: 10}, nil
}
func (f *fakeSession) BrewStep(context.Context, game.ActivityOpIn) (game.BrewStepOut, error) {
	return game.BrewStepOut{Step: 1}, nil
}
func (f *fakeSession) BrewClaim(context.Context, game.BrewClaimIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimRecall(context.Context, game.ActivityOpIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimReturn(context.Context, game.ActivityOpIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimInvite(context.Context, game.InviteIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) ClaimNewcomer(context.Context, game.ActivityOpIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) SendGift(context.Context, game.GiftIn) (int64, error) { return 1, nil }
func (f *fakeSession) Mystery(context.Context) (game.MysteryShop, error) {
	return game.MysteryShop{Present: true}, nil
}
func (f *fakeSession) MysteryBuy(context.Context, game.IDIn) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) MysteryLeave(context.Context) error { return nil }
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

func TestBagRequiresLogin(t *testing.T) {
	out := post(t, testClient(t, &fakeSession{}), "/Bag/Get", nil)
	if out.Code != 401 {
		t.Fatalf("%+v", out)
	}
}

func TestBagAfterLogin(t *testing.T) {
	mux := testClient(t, &fakeSession{})
	if out := post(t, mux, "/User/Login", game.LoginIn{Code: "abc"}); out.Code != 0 {
		t.Fatalf("login %+v", out)
	}
	out := post(t, mux, "/Bag/Get", nil)
	if out.Code != 0 {
		t.Fatalf("bag %+v", out)
	}
	out = post(t, mux, "/Shop/Goods", game.ShopIn{})
	if out.Code != 0 {
		t.Fatalf("goods %+v", out)
	}
	out = post(t, mux, "/Friend/Enter", game.EnterIn{GID: 9})
	if out.Code != 0 {
		t.Fatalf("enter %+v", out)
	}
	out = post(t, mux, "/Task/List", nil)
	if out.Code != 0 {
		t.Fatalf("task %+v", out)
	}
	out = post(t, mux, "/Weather/Status", nil)
	if out.Code != 0 {
		t.Fatalf("weather %+v", out)
	}
	out = post(t, mux, "/Farm/Unlock", game.LandIDIn{LandID: 4})
	if out.Code != 0 {
		t.Fatalf("unlock %+v", out)
	}
	out = post(t, mux, "/Email/List", nil)
	if out.Code != 0 {
		t.Fatalf("email %+v", out)
	}
	out = post(t, mux, "/Activity/List", nil)
	if out.Code != 0 {
		t.Fatalf("activity %+v", out)
	}
	out = post(t, mux, "/Season/Info", nil)
	if out.Code != 0 {
		t.Fatalf("season %+v", out)
	}
	out = post(t, mux, "/Mall/List", nil)
	if out.Code != 0 {
		t.Fatalf("mall %+v", out)
	}
	out = post(t, mux, "/Friend/Applications", nil)
	if out.Code != 0 {
		t.Fatalf("apps %+v", out)
	}
	out = post(t, mux, "/Share/Check", nil)
	if out.Code != 0 {
		t.Fatalf("share %+v", out)
	}
	out = post(t, mux, "/RedPacket/List", nil)
	if out.Code != 0 {
		t.Fatalf("red %+v", out)
	}
	out = post(t, mux, "/Visit/Logs", nil)
	if out.Code != 0 {
		t.Fatalf("logs %+v", out)
	}
	out = post(t, mux, "/Album/List", nil)
	if out.Code != 0 {
		t.Fatalf("album %+v", out)
	}
	out = post(t, mux, "/Mystery/Status", nil)
	if out.Code != 0 {
		t.Fatalf("mystery %+v", out)
	}
	out = post(t, mux, "/Activity/ShopBuy", game.ActivityShopIn{ID: 1, GoodsID: 2})
	if out.Code != 0 {
		t.Fatalf("shopbuy %+v", out)
	}
	out = post(t, mux, "/Activity/Lottery", game.LotteryIn{ID: 1, HostGID: 9})
	if out.Code != 0 {
		t.Fatalf("lottery %+v", out)
	}
	out = post(t, mux, "/Dog/Info", nil)
	if out.Code != 0 {
		t.Fatalf("dog %+v", out)
	}
	out = post(t, mux, "/Bulletin/List", nil)
	if out.Code != 0 {
		t.Fatalf("bulletin %+v", out)
	}
	out = post(t, mux, "/Career/Info", nil)
	if out.Code != 0 {
		t.Fatalf("career %+v", out)
	}
	out = post(t, mux, "/Rank/List", nil)
	if out.Code != 0 {
		t.Fatalf("rank %+v", out)
	}
	out = post(t, mux, "/Solar/List", nil)
	if out.Code != 0 {
		t.Fatalf("solar %+v", out)
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
