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
func (f *fakeSession) Heartbeat(context.Context, game.RefreshIn) (game.HeartbeatOut, error) {
	return game.HeartbeatOut{ServerMs: 1, HostGID: f.user.GID}, nil
}
func (f *fakeSession) Refresh(context.Context, game.RefreshIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) RefreshLands(context.Context, game.RefreshLandsIn) ([]game.Land, error) {
	return f.lands, nil
}
func (f *fakeSession) CleanSocial(context.Context, game.CleanSocialIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Harvest(context.Context, game.HarvestIn) (game.HarvestOut, error) {
	return game.HarvestOut{Items: f.items, Lands: f.lands}, nil
}
func (f *fakeSession) Plant(context.Context, game.PlantIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Friends(context.Context) (game.FriendsOut, error) {
	return game.FriendsOut{Friends: f.friends}, nil
}
func (f *fakeSession) Help(context.Context, game.HelpIn) error { return nil }
func (f *fakeSession) Applications(context.Context) (game.ApplicationsOut, error) {
	return game.ApplicationsOut{Applications: []game.Application{{GID: 3, Name: "申"}}}, nil
}
func (f *fakeSession) Accept(context.Context, game.GIDsIn) (game.AcceptOut, error) {
	return game.AcceptOut{Friends: f.friends, Success: 1}, nil
}
func (f *fakeSession) Reject(context.Context, game.GIDsIn) (game.RejectOut, error) {
	return game.RejectOut{Count: 1}, nil
}
func (f *fakeSession) SetTags(_ context.Context, in game.TagsIn) (game.Friend, error) {
	return game.Friend{GID: in.GID, New: in.New, Follow: in.Follow}, nil
}
func (f *fakeSession) DeleteFriend(context.Context, game.EnterIn) error {
	return nil
}
func (f *fakeSession) ShareCheck(context.Context) (bool, error) { return true, nil }
func (f *fakeSession) ShareClaim(context.Context) (game.ShareClaimOut, error) {
	return game.ShareClaimOut{Items: f.items}, nil
}
func (f *fakeSession) Water(context.Context, game.LandOpIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Weed(context.Context, game.LandOpIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Bug(context.Context, game.LandOpIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Fertilize(context.Context, game.FertilizeIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Enter(_ context.Context, in game.EnterIn) (game.Visit, error) {
	return game.Visit{Host: game.User{GID: in.GID}, Lands: f.lands}, nil
}
func (f *fakeSession) Leave(context.Context, game.EnterIn) error { return nil }
func (f *fakeSession) Bag(context.Context) (game.BagOut, error) {
	return game.BagOut{Items: []game.BagItem{{ID: 20002, Count: 3}}}, nil
}
func (f *fakeSession) Sell(context.Context, game.SellIn) (game.ItemOpOut, error) {
	return game.ItemOpOut{Items: []game.Item{{ID: 1001, Count: 10}}}, nil
}
func (f *fakeSession) Use(context.Context, game.UseIn) (game.ItemOpOut, error) {
	return game.ItemOpOut{Items: f.items}, nil
}
func (f *fakeSession) BatchUse(context.Context, game.BatchUseIn) (game.ItemOpOut, error) {
	return game.ItemOpOut{Items: f.items}, nil
}
func (f *fakeSession) CancelNew(_ context.Context, in game.IDIn) (int64, error) { return in.ID, nil }
func (f *fakeSession) Shops(context.Context) ([]game.Shop, error) {
	return []game.Shop{{ID: 2, Name: "种子"}}, nil
}
func (f *fakeSession) Goods(context.Context, game.ShopIn) ([]game.Goods, error) {
	return []game.Goods{{ID: 15, ItemID: 20002, Price: 50, Unlocked: true}}, nil
}
func (f *fakeSession) Buy(context.Context, game.BuyIn) (game.BuyOut, error) {
	return game.BuyOut{Items: []game.Item{{ID: 20002, Count: 1}}}, nil
}
func (f *fakeSession) Steal(context.Context, game.HarvestIn) (game.HarvestOut, error) {
	return game.HarvestOut{Items: f.items, Lands: f.lands}, nil
}
func (f *fakeSession) Remove(context.Context, game.RemoveIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Unlock(_ context.Context, in game.LandIDIn) (game.Land, error) {
	return game.Land{ID: in.LandID, Unlocked: true}, nil
}
func (f *fakeSession) Upgrade(_ context.Context, in game.LandIDIn) (game.Land, error) {
	return game.Land{ID: in.LandID, Level: 2}, nil
}
func (f *fakeSession) Farming(context.Context, game.LandOpIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) Weather(context.Context) (game.Weather, error) {
	return game.Weather{Type: 1, Active: true}, nil
}
func (f *fakeSession) TodayWeather(context.Context) ([]game.Weather, error) {
	return []game.Weather{{Type: 1, Name: "晴"}}, nil
}
func (f *fakeSession) CurrentWeather(context.Context) (game.Weather, error) {
	return game.Weather{Type: 1, Name: "晴"}, nil
}
func (f *fakeSession) Tasks(context.Context) (game.TaskBoard, error) {
	return game.TaskBoard{Daily: []game.Task{{ID: 1, Desc: "浇水"}}}, nil
}
func (f *fakeSession) ClaimTask(context.Context, game.ClaimIn) (game.TaskClaimOut, error) {
	return game.TaskClaimOut{Items: f.items}, nil
}
func (f *fakeSession) ReportTask(context.Context, game.TaskReportIn) (game.TaskBoard, error) {
	return game.TaskBoard{}, nil
}
func (f *fakeSession) ClaimDaily(context.Context, game.DailyIn) (game.TaskClaimOut, error) {
	return game.TaskClaimOut{Items: f.items}, nil
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
func (f *fakeSession) BatchReadEmail(context.Context, game.EmailIDsIn) error   { return nil }
func (f *fakeSession) BatchDeleteEmail(context.Context, game.EmailIDsIn) error { return nil }
func (f *fakeSession) ClaimAllEmail(context.Context, game.EmailBoxIn) (game.EmailClaimOut, error) {
	return game.EmailClaimOut{Items: f.items}, nil
}
func (f *fakeSession) Activities(context.Context) ([]game.Activity, error) {
	return []game.Activity{{ID: 8, Name: "签到"}}, nil
}
func (f *fakeSession) SetSplashed(context.Context, game.IDIn) (bool, error) { return true, nil }
func (f *fakeSession) Invitees(context.Context, game.IDIn) ([]game.Invitee, error) {
	return []game.Invitee{{GID: 2, Name: "客"}}, nil
}
func (f *fakeSession) ActivityGroup(_ context.Context, in game.GroupIn) (game.ActivityGroup, error) {
	return game.ActivityGroup{ID: in.Group(), Activities: []game.Activity{{ID: 8, GroupID: in.Group()}}}, nil
}
func (f *fakeSession) Season(context.Context) (game.Season, error) {
	return game.Season{ID: 1, Name: "春", Pass: game.BattlePass{Level: 3}}, nil
}
func (f *fakeSession) ClaimPass(context.Context) (game.PassClaimOut, error) {
	return game.PassClaimOut{Items: f.items}, nil
}
func (f *fakeSession) MallProfiles(context.Context) ([]game.MallProfile, error) {
	return []game.MallProfile{{ID: 1, Type: 1}}, nil
}
func (f *fakeSession) Mall(context.Context, game.MallIn) ([]game.Product, error) {
	return []game.Product{{ID: 1001, Name: "每日福利", Available: true}}, nil
}
func (f *fakeSession) MallDiamonds(context.Context) ([]game.Product, error) {
	return []game.Product{{ID: 2001, Name: "点券", Available: true}}, nil
}
func (f *fakeSession) MallBuy(context.Context, game.MallBuyIn) (game.BuyOut, error) {
	return game.BuyOut{Items: f.items}, nil
}
func (f *fakeSession) MonthCards(context.Context) ([]game.MonthCard, error) {
	return []game.MonthCard{{ID: 2001, Claimable: true}}, nil
}
func (f *fakeSession) ClaimMonthCard(context.Context, game.IDIn) (game.MonthCardClaimOut, error) {
	return game.MonthCardClaimOut{Items: f.items}, nil
}
func (f *fakeSession) RedPackets(context.Context) ([]game.RedPacket, error) {
	return []game.RedPacket{{ID: 1, CanClaim: true}}, nil
}
func (f *fakeSession) ClaimRedPacket(context.Context, game.IDIn) (game.RedPacketClaimOut, error) {
	return game.RedPacketClaimOut{Items: f.items}, nil
}
func (f *fakeSession) ClaimAllRedPackets(context.Context) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) VisitPopup(context.Context) (game.VisitPopup, error) {
	return game.VisitPopup{NeedPopup: true}, nil
}
func (f *fakeSession) VisitPage(context.Context, game.VisitPageIn) ([]game.VisitLog, error) {
	return nil, nil
}
func (f *fakeSession) DismissVisit(context.Context, game.EnterIn) error      { return nil }
func (f *fakeSession) DeleteVisit(context.Context, game.VisitDeleteIn) error { return nil }
func (f *fakeSession) VisitLogs(context.Context) ([]game.VisitLog, error) {
	return []game.VisitLog{{GID: 2, Action: 4, Name: "客"}}, nil
}
func (f *fakeSession) VisitSummary(context.Context) (game.VisitSummary, error) {
	return game.VisitSummary{StealCount: 1, HelpCount: 2}, nil
}
func (f *fakeSession) Album(context.Context, game.AlbumIn) (game.Album, error) {
	return game.Album{Items: []game.AlbumItem{{FruitID: 40061, Unlocked: true}}}, nil
}
func (f *fakeSession) MarkAlbum(context.Context, game.AlbumIn) error { return nil }
func (f *fakeSession) ClaimAlbum(context.Context, game.AlbumIn) (game.AlbumClaimOut, error) {
	return game.AlbumClaimOut{Items: f.items}, nil
}
func (f *fakeSession) Dog(context.Context) (game.DogYard, error) {
	return game.DogYard{Deployed: 1, FoodLeft: 60}, nil
}
func (f *fakeSession) Feed(context.Context, game.FeedIn) (int64, error) { return 120, nil }
func (f *fakeSession) ClaimDogGifts(context.Context) (game.DogGiftOut, error) {
	return game.DogGiftOut{Items: f.items}, nil
}
func (f *fakeSession) DogLogs(context.Context, game.PageIn) (game.DogLogsOut, error) {
	return game.DogLogsOut{Logs: []game.ProtectLog{{GID: 2, Name: "贼"}}}, nil
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
	return []game.Mutant{{ID: 1, Red: true}}, nil
}
func (f *fakeSession) Career(context.Context, game.EnterIn) (game.Career, error) {
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
func (f *fakeSession) MarkAvatar(context.Context, game.IDIn) error { return nil }
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
func (f *fakeSession) MarkSkin(context.Context, game.IDIn) error         { return nil }
func (f *fakeSession) Drops(context.Context) ([]game.Drop, error) {
	return []game.Drop{{ID: 1}}, nil
}
func (f *fakeSession) SolarRedDot(context.Context) (game.RedDot, error) {
	return game.RedDot{Red: true}, nil
}
func (f *fakeSession) SolarTerms(context.Context) (game.SolarOut, error) {
	return game.SolarOut{Terms: []game.SolarTerm{{ID: 1, Status: 2}}}, nil
}
func (f *fakeSession) ClaimSolar(context.Context, game.IDIn) (game.SolarClaimOut, error) {
	return game.SolarClaimOut{Items: f.items}, nil
}
func (f *fakeSession) ClaimAllSolar(context.Context) ([]game.Item, error) {
	return f.items, nil
}
func (f *fakeSession) AchieveView(context.Context, game.AchieveIn) (game.AchieveScope, error) {
	return game.AchieveScope{ID: 1}, nil
}
func (f *fakeSession) ClaimAchieveGoal(context.Context, game.AchieveGoalIn) (game.AchieveGoalOut, error) {
	return game.AchieveGoalOut{Exp: 10}, nil
}
func (f *fakeSession) ClaimAchieveLevel(context.Context, game.AchieveIn) (game.ItemOpOut, error) {
	return game.ItemOpOut{Items: f.items}, nil
}
func (f *fakeSession) Signin(context.Context, game.ActivityOpIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) ClaimProgress(context.Context, game.ActivityOpIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) ShopBuy(context.Context, game.ActivityShopIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) ShopBatchBuy(context.Context, game.ActivityBatchIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) RandBuy(context.Context, game.ActivityShopIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) RandRefresh(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Costs: f.items}, nil
}
func (f *fakeSession) ClaimMega(context.Context, game.ActivityOpIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) TechSubmit(context.Context, game.TechIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) Draw(context.Context, game.DrawIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) DrawHistory(context.Context, game.IDIn) (game.DrawHistoryOut, error) {
	return game.DrawHistoryOut{}, nil
}
func (f *fakeSession) MarkViewed(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) RandBatchBuy(context.Context, game.ActivityBatchIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) LotteryHistory(context.Context, game.IDIn) (game.LotteryHistoryOut, error) {
	return game.LotteryHistoryOut{}, nil
}
func (f *fakeSession) CheerJoin(context.Context, game.CheerJoinIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) CheerSubmit(context.Context, game.CheerSubmitIn) (game.CheerSubmitOut, error) {
	return game.CheerSubmitOut{Added: 1}, nil
}
func (f *fakeSession) CheerClaim(context.Context, game.CheerClaimIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) Recallable(context.Context, game.IDIn) (game.RecallListOut, error) {
	return game.RecallListOut{List: []game.RecallPerson{{GID: 9}}}, nil
}
func (f *fakeSession) Recalled(context.Context, game.IDIn) (game.RecallListOut, error) {
	return game.RecallListOut{}, nil
}
func (f *fakeSession) CharityShare(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) CharityDonate(context.Context, game.IDIn) (game.CharityDonateOut, error) {
	return game.CharityDonateOut{Consumed: 1}, nil
}
func (f *fakeSession) CharityClaim(context.Context, game.CharityClaimIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) CharityXhh(context.Context, game.IDIn) (game.CharityXhhOut, error) {
	return game.CharityXhhOut{}, nil
}
func (f *fakeSession) CharityAgree(context.Context, game.CharityAgreeIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntFinishCG(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntGuide(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntFeed(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntDraw(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntLog(context.Context, game.IDIn) (game.HuntLogOut, error) {
	return game.HuntLogOut{}, nil
}
func (f *fakeSession) HuntClaimStory(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntClaimSeed(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntRefreshCharm(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntEquip(context.Context, game.HuntEquipIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntBattle(context.Context, game.HuntBattleIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntPlunderedLog(context.Context, game.IDIn) (game.HuntLogOut, error) {
	return game.HuntLogOut{}, nil
}
func (f *fakeSession) HuntOpen(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntEscort(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntCompensate(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
}
func (f *fakeSession) HuntFriendInfo(context.Context, game.IDIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{}, nil
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
func (f *fakeSession) BrewClaim(context.Context, game.BrewClaimIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) ClaimRecall(context.Context, game.ActivityOpIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) ClaimReturn(context.Context, game.ActivityOpIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) ClaimInvite(context.Context, game.InviteIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) ClaimNewcomer(context.Context, game.ActivityOpIn) (game.ActivityOpOut, error) {
	return game.ActivityOpOut{Items: f.items}, nil
}
func (f *fakeSession) SendGift(context.Context, game.GiftIn) (int64, error) { return 1, nil }
func (f *fakeSession) Mystery(context.Context) (game.MysteryShop, error) {
	return game.MysteryShop{Present: true}, nil
}
func (f *fakeSession) MysteryBuy(context.Context, game.IDIn) (game.MysteryBuyOut, error) {
	return game.MysteryBuyOut{Items: f.items}, nil
}
func (f *fakeSession) MysteryLeave(context.Context) error { return nil }
func (f *fakeSession) Status() (game.Status, error) {
	if f.user.GID == 0 {
		return game.Status{}, game.ErrNotLogin
	}
	return game.Status{LoggedIn: true, User: f.user}, nil
}
func (f *fakeSession) Logout(context.Context) error { return f.Close() }
func (f *fakeSession) Brief(_ context.Context, in game.EnterIn) (game.User, error) {
	return game.User{GID: in.GID, Name: "友"}, nil
}
func (f *fakeSession) BatchInfo(context.Context, game.GIDsIn) ([]game.User, error) {
	return []game.User{{GID: 2, Name: "友"}}, nil
}
func (f *fakeSession) ArkClick(context.Context, game.ArkIn) error {
	return nil
}
func (f *fakeSession) PutInsects(context.Context, game.LandOpIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) PutWeeds(context.Context, game.LandOpIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) PutSocial(context.Context, game.PutSocialIn) (game.LandOpOut, error) {
	return game.LandOpOut{Lands: f.lands}, nil
}
func (f *fakeSession) CanOperate(context.Context, game.CanOperateIn) (game.CanOperateOut, error) {
	return game.CanOperateOut{OK: true, StealNum: 2}, nil
}
func (f *fakeSession) SyncFriends(context.Context, game.OpenIDsIn) ([]game.Friend, error) {
	return f.friends, nil
}
func (f *fakeSession) GameFriends(context.Context, game.GIDsIn) ([]game.Friend, error) {
	return f.friends, nil
}
func (f *fakeSession) BlockFriend(_ context.Context, in game.EnterIn) (game.Blocked, error) {
	return game.Blocked{GID: in.GID}, nil
}
func (f *fakeSession) UnblockFriend(context.Context, game.EnterIn) error { return nil }
func (f *fakeSession) BlockList(context.Context) ([]game.Blocked, error) {
	return []game.Blocked{{GID: 3}}, nil
}
func (f *fakeSession) ShareKey(context.Context) (string, error) { return "k", nil }
func (f *fakeSession) InviteInfo(context.Context, game.IDIn) (game.ShareInviteOut, error) {
	return game.ShareInviteOut{Infos: []game.ShareInviteInfo{{ID: 1, InviteCount: 1}}}, nil
}
func (f *fakeSession) InviteAward(context.Context, game.IDIn) (game.ShareAwardOut, error) {
	return game.ShareAwardOut{Awards: f.items, Awarded: true}, nil
}
func (f *fakeSession) PosterShown(context.Context, game.IDIn) (bool, error) { return true, nil }
func (f *fakeSession) ReportInvite(context.Context, game.InviteReportIn) (bool, error) {
	return true, nil
}
func (f *fakeSession) LockItems(context.Context, game.UIDsIn) (game.LockOut, error) {
	return game.LockOut{Newly: []int64{1}}, nil
}
func (f *fakeSession) UnlockItems(context.Context, game.UIDsIn) (game.LockOut, error) {
	return game.LockOut{Newly: []int64{1}}, nil
}
func (f *fakeSession) AutoBuy(context.Context, game.AutoBuyIn) (game.BuyOut, error) {
	return game.BuyOut{Items: f.items}, nil
}
func (f *fakeSession) MysteryAutoBuy(context.Context, game.MysteryAutoIn) (game.MysteryAutoOut, error) {
	return game.MysteryAutoOut{Items: f.items}, nil
}
func (f *fakeSession) BuyDog(_ context.Context, in game.DogBuyIn) (game.DogBuyOut, error) {
	return game.DogBuyOut{Dog: game.Dog{ID: in.ID, Activated: true}}, nil
}
func (f *fakeSession) AlbumLevels(context.Context, game.AlbumIn) (game.AlbumLevels, error) {
	return game.AlbumLevels{Level: 1}, nil
}
func (f *fakeSession) CleanFarmEvents(context.Context) (game.CleanFarmEventOut, error) {
	return game.CleanFarmEventOut{}, nil
}
func (f *fakeSession) BlockApplications(_ context.Context, in game.BlockAppsIn) (game.BlockAppsOut, error) {
	return game.BlockAppsOut{Block: in.Block}, nil
}
func (f *fakeSession) WXRecommend(context.Context, game.WXRecommendIn) (game.WXRecommendOut, error) {
	return game.WXRecommendOut{}, nil
}
func (f *fakeSession) WXRecommendPage(context.Context, game.WXRecommendPageIn) (game.WXRecommendOut, error) {
	return game.WXRecommendOut{}, nil
}
func (f *fakeSession) ApplyWXFriends(context.Context, game.GIDsIn) (game.WXApplyOut, error) {
	return game.WXApplyOut{Results: []game.WXApplyResult{{GID: 3, Success: true}}}, nil
}
func (f *fakeSession) EquipSkinSet(context.Context, game.SkinSetIn) error { return nil }
func (f *fakeSession) SetSkinSetEffect(context.Context, game.SkinSetEffectIn) error {
	return nil
}
func (f *fakeSession) SkinSets(context.Context) ([]game.SkinSetEffect, error) {
	return []game.SkinSetEffect{{SetID: 1, Active: true}}, nil
}
func (f *fakeSession) BuyPass(context.Context) (game.BuyPassOut, error) {
	return game.BuyPassOut{Success: true, Items: f.items}, nil
}
func (f *fakeSession) MarkSeasonOpening(context.Context) (bool, error) { return true, nil }
func (f *fakeSession) QQAuthGroups(context.Context, game.CookiesIn) (game.QQAuthGroupsOut, error) {
	return game.QQAuthGroupsOut{}, nil
}
func (f *fakeSession) QQRecommendGroups(context.Context, game.QQRecommendIn) (game.QQRecommendOut, error) {
	return game.QQRecommendOut{}, nil
}
func (f *fakeSession) QQBind(context.Context, game.QQBindIn) (game.QQBindOut, error) {
	return game.QQBindOut{}, nil
}
func (f *fakeSession) QQLeave(context.Context) (game.QQLeaveOut, error) {
	return game.QQLeaveOut{}, nil
}
func (f *fakeSession) QQCommunity(context.Context, game.PageIn) (game.QQCommunityOut, error) {
	return game.QQCommunityOut{}, nil
}
func (f *fakeSession) QQBindInfo(context.Context) (game.QQBindInfoOut, error) {
	return game.QQBindInfoOut{}, nil
}
func (f *fakeSession) QQClaimReward(context.Context) ([]game.Item, error)  { return f.items, nil }
func (f *fakeSession) QQRevokeAuth(context.Context, game.QQRevokeIn) error { return nil }
func (f *fakeSession) UseGiftToken(context.Context, game.UIDIn) (game.GiftTokenOut, error) {
	return game.GiftTokenOut{RedeemCode: "c"}, nil
}
func (f *fakeSession) GiftHistory(context.Context, game.GiftHistoryIn) (game.GiftHistoryOut, error) {
	return game.GiftHistoryOut{}, nil
}
func (f *fakeSession) TransferStatus(context.Context, game.TransferIn) (game.TransferOut, error) {
	return game.TransferOut{}, nil
}
func (f *fakeSession) CancelTransfer(context.Context, game.UIDIn) (game.TransferOut, error) {
	return game.TransferOut{}, nil
}
func (f *fakeSession) FollowGiftStatus(context.Context) (game.FollowGiftOut, error) {
	return game.FollowGiftOut{Followed: true}, nil
}
func (f *fakeSession) SetFollowGift(context.Context, game.FollowGiftIn) error { return nil }
func (f *fakeSession) ClaimFollowGift(context.Context) ([]game.Item, error)   { return f.items, nil }
func (f *fakeSession) RechargeBonus(context.Context) (game.RechargeBonusOut, error) {
	return game.RechargeBonusOut{Active: true}, nil
}
func (f *fakeSession) RechargeBonusData(context.Context) (game.RechargeDataOut, error) {
	return game.RechargeDataOut{}, nil
}
func (f *fakeSession) SetDisplay(context.Context, game.DisplayIn) (game.DisplayOut, error) {
	return game.DisplayOut{Name: "测"}, nil
}
func (f *fakeSession) GetSettings(context.Context, game.SettingsKeysIn) (game.UserSettings, error) {
	return game.UserSettings{}, nil
}
func (f *fakeSession) SetSettings(_ context.Context, in game.UserSettings) (game.UserSettings, error) {
	return in, nil
}
func (f *fakeSession) DeleteAccount(context.Context, game.DeleteAccountIn) (game.DeleteAccountOut, error) {
	return game.DeleteAccountOut{Success: true}, nil
}
func (f *fakeSession) DecryptOpenData(context.Context, game.DecryptIn) (string, error) {
	return "{}", nil
}
func (f *fakeSession) SetQQRecommendAuth(_ context.Context, in game.QQAuthIn) (game.QQAuthOut, error) {
	if in.Authorized {
		return game.QQAuthOut{Authorized: 1}, nil
	}
	return game.QQAuthOut{}, nil
}
func (f *fakeSession) ReportFlow(context.Context, game.ReportFlowIn) error { return nil }
func (f *fakeSession) BatchReportFlow(context.Context, game.BatchReportFlowIn) (game.BatchReportFlowOut, error) {
	return game.BatchReportFlowOut{Success: 1}, nil
}
func (f *fakeSession) ReportUser(context.Context, game.ReportUserIn) (game.ReportUserOut, error) {
	return game.ReportUserOut{Success: true}, nil
}
func (f *fakeSession) QQVipDailyStatus(context.Context) (game.QQVipDailyOut, error) {
	return game.QQVipDailyOut{IsQQVip: true, CanClaim: true}, nil
}
func (f *fakeSession) QQVipClaimDaily(context.Context, game.QQVipClaimDailyIn) (game.QQVipClaimDailyOut, error) {
	return game.QQVipClaimDailyOut{Rewards: f.items}, nil
}
func (f *fakeSession) QQVipRefresh(context.Context) (game.QQVipRefreshOut, error) {
	return game.QQVipRefreshOut{IsQQVip: true, VIPLevel: 1}, nil
}
func (f *fakeSession) QQVipClaimRewards(context.Context, game.QQVipClaimRewardsIn) (game.QQVipClaimRewardsOut, error) {
	return game.QQVipClaimRewardsOut{Rewards: f.items}, nil
}
func (f *fakeSession) QQVipRewardsStatus(context.Context) (game.QQVipRewardsStatusOut, error) {
	return game.QQVipRewardsStatusOut{IsQQVip: true, RewardsCanClaim: true}, nil
}
func (f *fakeSession) QQVipMarkRedpoint(context.Context) error { return nil }
func (f *fakeSession) Marquee(context.Context) (game.MarqueeOut, error) {
	return game.MarqueeOut{Msgs: []game.MarqueeMsg{{Content: "hi"}}}, nil
}
func (f *fakeSession) SystemUnlocked(context.Context, game.SystemOpenIn) (game.SystemOpenOut, error) {
	return game.SystemOpenOut{Unlocked: true}, nil
}
func (f *fakeSession) MutantOpenInfo(context.Context) (game.MutantOpenInfoOut, error) {
	return game.MutantOpenInfoOut{Tips: "ok", Rewards: f.items}, nil
}
func (f *fakeSession) QQSubscribe(context.Context) (game.QQSubscribeOut, error) {
	return game.QQSubscribeOut{Subscribed: true}, nil
}
func (f *fakeSession) WXSubscribe(context.Context) (game.WXSubscribeOut, error) {
	return game.WXSubscribeOut{Templates: []game.WXTemplateStatus{{TemplateID: "t1", Subscribed: true}}}, nil
}
func (f *fakeSession) SetWXSubscribe(_ context.Context, in game.WXSubscribeIn) (game.WXSubscribeOut, error) {
	return game.WXSubscribeOut{Templates: in.Templates}, nil
}
func (f *fakeSession) ModerateText(_ context.Context, in game.ModerateTextIn) (game.ModerateTextOut, error) {
	return game.ModerateTextOut{Text: in.Text}, nil
}
func (f *fakeSession) BatchModerateText(_ context.Context, in game.ModerateTextBatchIn) (game.ModerateTextBatchOut, error) {
	out := make([]game.ModerateTextOut, 0, len(in.Items))
	for _, it := range in.Items {
		out = append(out, game.ModerateTextOut{Text: it.Text})
	}
	return game.ModerateTextBatchOut{Items: out}, nil
}
func (f *fakeSession) ModeratePic(_ context.Context, in game.ModeratePicIn) (game.ModeratePicOut, error) {
	return game.ModeratePicOut{URL: in.URL}, nil
}
func (f *fakeSession) BatchModeratePic(_ context.Context, in game.ModeratePicBatchIn) (game.ModeratePicBatchOut, error) {
	out := make([]game.ModeratePicOut, 0, len(in.Items))
	for _, it := range in.Items {
		out = append(out, game.ModeratePicOut{URL: it.URL})
	}
	return game.ModeratePicBatchOut{Items: out}, nil
}
func (f *fakeSession) Close() error {
	f.user = game.User{}
	return nil
}

type fakeYYB struct{}

func (fakeYYB) Accounts(context.Context) ([]game.YYBAccount, error) {
	return []game.YYBAccount{{ID: 1, OpenID: "o1", Nickname: "测", Status: "alive"}}, nil
}
func (fakeYYB) CreateQR(context.Context) (game.YYBQROut, error) {
	return game.YYBQROut{SessionID: "s1", Status: "pending", Image: "data:image/jpeg;base64,AA"}, nil
}
func (fakeYYB) QRImage(string) ([]byte, error) { return []byte("jpeg"), nil }
func (fakeYYB) Poll(context.Context, string) (game.YYBPollOut, error) {
	return game.YYBPollOut{Status: "authorized"}, nil
}
func (fakeYYB) Confirm(context.Context, string) (game.YYBAccount, error) {
	return game.YYBAccount{ID: 1, OpenID: "o1", Status: "alive"}, nil
}
func (fakeYYB) Refresh(context.Context, string) (game.YYBRefreshOut, error) {
	return game.YYBRefreshOut{ID: 1, OpenID: "o1", Status: "alive"}, nil
}
func (fakeYYB) GetCode(context.Context, string, string) (game.YYBCodeOut, error) {
	return game.YYBCodeOut{Code: "abc", OpenID: "o1"}, nil
}
func (fakeYYB) Delete(context.Context, string) (game.YYBDeleteOut, error) {
	return game.YYBDeleteOut{Deleted: 1, OpenID: "o1"}, nil
}
func (fakeYYB) Profile(context.Context, string) (game.YYBAccount, error) {
	return game.YYBAccount{ID: 1, OpenID: "o1", Nickname: "测", Avatar: "http://a", Status: "alive"}, nil
}
func (fakeYYB) Phone(context.Context, string, string) (game.YYBRawOut, error) {
	return game.YYBRawOut{OpenID: "o1", Result: map[string]any{"ok": true}}, nil
}
func (fakeYYB) WXData(context.Context, game.YYBWXDataIn) (game.YYBRawOut, error) {
	return game.YYBRawOut{OpenID: "o1", Result: map[string]any{"ok": true}}, nil
}

func testClient(t *testing.T, sess *fakeSession) *http.ServeMux {
	t.Helper()
	return NewMux(NewHub(func() game.Session { return sess }).WithYYB(fakeYYB{}))
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
	out = post(t, mux, "/Friend/SetTags", game.TagsIn{GID: 3, Follow: true})
	if out.Code != 0 {
		t.Fatalf("tags %+v", out)
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
	out = post(t, mux, "/Farm/RefreshLands", game.RefreshLandsIn{LandIDs: []int64{1}})
	if out.Code != 0 {
		t.Fatalf("refreshlands %+v", out)
	}
	out = post(t, mux, "/Farm/CleanSocial", game.CleanSocialIn{LandIDs: []int64{1}})
	if out.Code != 0 {
		t.Fatalf("cleansocial %+v", out)
	}
	out = post(t, mux, "/Bag/BatchUse", game.BatchUseIn{Items: []game.UseIn{{ID: 1, Count: 1}}})
	if out.Code != 0 {
		t.Fatalf("batchuse %+v", out)
	}
	out = post(t, mux, "/User/Heartbeat", nil)
	if out.Code != 0 {
		t.Fatalf("heartbeat %+v", out)
	}
	out = post(t, mux, "/Activity/GetGroup", game.GroupIn{ID: 8})
	if out.Code != 0 {
		t.Fatalf("getgroup %+v", out)
	}
	out = post(t, mux, "/Farm/PutInsects", game.LandOpIn{HostGID: 9, LandIDs: []int64{1}})
	if out.Code != 0 {
		t.Fatalf("putinsects %+v", out)
	}
	out = post(t, mux, "/Shop/AutoBuy", game.AutoBuyIn{ItemID: 20002})
	if out.Code != 0 {
		t.Fatalf("autobuy %+v", out)
	}
	out = post(t, mux, "/Dog/Buy", game.DogBuyIn{ID: 1, Price: 10})
	if out.Code != 0 {
		t.Fatalf("buydog %+v", out)
	}
	out = post(t, mux, "/Activity/SetSplashed", game.IDIn{ID: 1})
	if out.Code != 0 {
		t.Fatalf("setsplashed %+v", out)
	}
	out = post(t, mux, "/Album/MarkViewed", game.AlbumIn{Type: 1})
	if out.Code != 0 {
		t.Fatalf("album viewed %+v", out)
	}
	out = post(t, mux, "/Avatar/MarkViewed", game.IDIn{ID: 1})
	if out.Code != 0 {
		t.Fatalf("avatar viewed %+v", out)
	}
	out = post(t, mux, "/Skin/MarkViewed", game.IDIn{ID: 1})
	if out.Code != 0 {
		t.Fatalf("skin viewed %+v", out)
	}
	out = post(t, mux, "/Visit/Popup", nil)
	if out.Code != 0 {
		t.Fatalf("visit popup %+v", out)
	}
	out = post(t, mux, "/Visit/Page", game.VisitPageIn{Page: 1})
	if out.Code != 0 {
		t.Fatalf("visit page %+v", out)
	}
	out = post(t, mux, "/Visit/Dismiss", game.EnterIn{GID: 9})
	if out.Code != 0 {
		t.Fatalf("visit dismiss %+v", out)
	}
	out = post(t, mux, "/Visit/Delete", game.VisitDeleteIn{IDs: []int64{1}})
	if out.Code != 0 {
		t.Fatalf("visit delete %+v", out)
	}
	out = post(t, mux, "/Email/BatchRead", game.EmailIDsIn{IDs: []string{"a"}})
	if out.Code != 0 {
		t.Fatalf("email batchread %+v", out)
	}
	out = post(t, mux, "/Email/BatchDelete", game.EmailIDsIn{IDs: []string{"a"}})
	if out.Code != 0 {
		t.Fatalf("email batchdelete %+v", out)
	}
	out = post(t, mux, "/Task/Report", game.TaskReportIn{ID: 1, Progress: 1})
	if out.Code != 0 {
		t.Fatalf("task report %+v", out)
	}
	out = post(t, mux, "/Mall/Profiles", nil)
	if out.Code != 0 {
		t.Fatalf("mall profiles %+v", out)
	}
	out = post(t, mux, "/Solar/RedDot", nil)
	if out.Code != 0 {
		t.Fatalf("solar reddot %+v", out)
	}
	out = post(t, mux, "/Weather/Current", nil)
	if out.Code != 0 {
		t.Fatalf("weather current %+v", out)
	}
	out = post(t, mux, "/Bag/CancelNew", game.IDIn{ID: 20002})
	if out.Code != 0 {
		t.Fatalf("cancelnew %+v", out)
	}
	out = post(t, mux, "/Activity/Invitees", game.IDIn{ID: 8})
	if out.Code != 0 {
		t.Fatalf("invitees %+v", out)
	}
	out = post(t, mux, "/Activity/Draw", game.DrawIn{ID: 8, Count: 1})
	if out.Code != 0 {
		t.Fatalf("draw %+v", out)
	}
	out = post(t, mux, "/Activity/CheerJoin", game.CheerJoinIn{ID: 8, CampID: 1})
	if out.Code != 0 {
		t.Fatalf("cheer %+v", out)
	}
	out = post(t, mux, "/Activity/CharityDonate", game.IDIn{ID: 8})
	if out.Code != 0 {
		t.Fatalf("charity %+v", out)
	}
	out = post(t, mux, "/Activity/RandBatchBuy", game.ActivityBatchIn{ID: 8, Items: []game.ShopBuyItem{{GoodsID: 2, Count: 1}}})
	if out.Code != 0 {
		t.Fatalf("rand batch %+v", out)
	}
	out = post(t, mux, "/Activity/HuntOpen", game.IDIn{ID: 2026090101})
	if out.Code != 0 {
		t.Fatalf("hunt open %+v", out)
	}
	out = post(t, mux, "/Activity/HuntEquip", game.HuntEquipIn{ID: 2026090101, CharmIDs: []int64{101}})
	if out.Code != 0 {
		t.Fatalf("hunt equip %+v", out)
	}
	out = post(t, mux, "/Activity/HuntBattle", game.HuntBattleIn{ID: 2026090101, DefenderGID: 9, TreasureID: "t1"})
	if out.Code != 0 {
		t.Fatalf("hunt battle %+v", out)
	}
	out = post(t, mux, "/Mall/Diamonds", nil)
	if out.Code != 0 {
		t.Fatalf("mall diamonds %+v", out)
	}
	out = post(t, mux, "/Share/PosterShown", game.IDIn{ID: 1})
	if out.Code != 0 {
		t.Fatalf("poster shown %+v", out)
	}
	out = post(t, mux, "/Visit/Summary", nil)
	if out.Code != 0 {
		t.Fatalf("visit summary %+v", out)
	}
	out = post(t, mux, "/Farm/PutSocial", game.PutSocialIn{HostGID: 9, LandIDs: []int64{1}, ItemID: 3})
	if out.Code != 0 {
		t.Fatalf("putsocial %+v", out)
	}
	out = post(t, mux, "/Share/ReportInvite", game.InviteReportIn{OpenID: "o", ShareKey: "k"})
	if out.Code != 0 {
		t.Fatalf("reportinvite %+v", out)
	}
	out = post(t, mux, "/Activity/ShopBatchBuy", game.ActivityBatchIn{ID: 1, Items: []game.ShopBuyItem{{GoodsID: 2, Count: 1}}})
	if out.Code != 0 {
		t.Fatalf("shopbatch %+v", out)
	}
	out = post(t, mux, "/Activity/RandBuy", game.ActivityShopIn{ID: 1, GoodsID: 2})
	if out.Code != 0 {
		t.Fatalf("randbuy %+v", out)
	}
	out = post(t, mux, "/Activity/RandRefresh", game.IDIn{ID: 1})
	if out.Code != 0 {
		t.Fatalf("randrefresh %+v", out)
	}
	out = post(t, mux, "/Career/Info", game.EnterIn{GID: 9})
	if out.Code != 0 {
		t.Fatalf("career gid %+v", out)
	}
	out = post(t, mux, "/Farm/CleanEvents", nil)
	if out.Code != 0 {
		t.Fatalf("cleanevents %+v", out)
	}
	out = post(t, mux, "/Friend/BlockApplications", game.BlockAppsIn{Block: true})
	if out.Code != 0 {
		t.Fatalf("blockapps %+v", out)
	}
	out = post(t, mux, "/Friend/ApplyWX", game.GIDsIn{GIDs: []int64{3}})
	if out.Code != 0 {
		t.Fatalf("applywx %+v", out)
	}
	out = post(t, mux, "/Skin/Sets", nil)
	if out.Code != 0 {
		t.Fatalf("skin sets %+v", out)
	}
	out = post(t, mux, "/Season/BuyPass", nil)
	if out.Code != 0 {
		t.Fatalf("buypass %+v", out)
	}
	out = post(t, mux, "/QQGroup/BindInfo", nil)
	if out.Code != 0 {
		t.Fatalf("qq bindinfo %+v", out)
	}
	out = post(t, mux, "/Gift/UseToken", game.UIDIn{UID: 1})
	if out.Code != 0 {
		t.Fatalf("gift token %+v", out)
	}
	out = post(t, mux, "/Follow/Status", nil)
	if out.Code != 0 {
		t.Fatalf("follow %+v", out)
	}
	out = post(t, mux, "/Recharge/Config", nil)
	if out.Code != 0 {
		t.Fatalf("recharge %+v", out)
	}
	out = post(t, mux, "/User/Settings", game.SettingsKeysIn{})
	if out.Code != 0 {
		t.Fatalf("settings %+v", out)
	}
	out = post(t, mux, "/User/Report", game.ReportUserIn{GID: 3})
	if out.Code != 0 {
		t.Fatalf("report %+v", out)
	}
	out = post(t, mux, "/QQVip/DailyStatus", nil)
	if out.Code != 0 {
		t.Fatalf("qqvip daily %+v", out)
	}
	out = post(t, mux, "/QQVip/ClaimDaily", game.QQVipClaimDailyIn{ConfigID: 1})
	if out.Code != 0 {
		t.Fatalf("qqvip claim daily %+v", out)
	}
	out = post(t, mux, "/QQVip/Refresh", nil)
	if out.Code != 0 {
		t.Fatalf("qqvip refresh %+v", out)
	}
	out = post(t, mux, "/QQVip/RewardsStatus", nil)
	if out.Code != 0 {
		t.Fatalf("qqvip rewards %+v", out)
	}
	out = post(t, mux, "/QQVip/ClaimRewards", game.QQVipClaimRewardsIn{ConfigIDs: []int64{1}})
	if out.Code != 0 {
		t.Fatalf("qqvip claim rewards %+v", out)
	}
	out = post(t, mux, "/QQVip/MarkRedpoint", nil)
	if out.Code != 0 {
		t.Fatalf("qqvip redpoint %+v", out)
	}
	out = post(t, mux, "/Marquee/List", nil)
	if out.Code != 0 {
		t.Fatalf("marquee %+v", out)
	}
	out = post(t, mux, "/SystemOpen/Unlocked", game.SystemOpenIn{SystemName: 1})
	if out.Code != 0 {
		t.Fatalf("systemopen %+v", out)
	}
	out = post(t, mux, "/Mutant/OpenInfo", nil)
	if out.Code != 0 {
		t.Fatalf("mutant open %+v", out)
	}
	out = post(t, mux, "/Subscribe/QQ", nil)
	if out.Code != 0 {
		t.Fatalf("subscribe qq %+v", out)
	}
	out = post(t, mux, "/Subscribe/WX", nil)
	if out.Code != 0 {
		t.Fatalf("subscribe wx %+v", out)
	}
	out = post(t, mux, "/Subscribe/SetWX", game.WXSubscribeIn{Templates: []game.WXTemplateStatus{{TemplateID: "t1", Subscribed: true}}})
	if out.Code != 0 {
		t.Fatalf("subscribe setwx %+v", out)
	}
	out = post(t, mux, "/Moderate/Text", game.ModerateTextIn{Text: "hi"})
	if out.Code != 0 {
		t.Fatalf("moderate text %+v", out)
	}
	out = post(t, mux, "/Moderate/BatchText", game.ModerateTextBatchIn{Items: []game.ModerateTextIn{{Text: "hi"}}})
	if out.Code != 0 {
		t.Fatalf("moderate batch text %+v", out)
	}
	out = post(t, mux, "/Moderate/Pic", game.ModeratePicIn{URL: "http://x"})
	if out.Code != 0 {
		t.Fatalf("moderate pic %+v", out)
	}
	out = post(t, mux, "/Moderate/BatchPic", game.ModeratePicBatchIn{Items: []game.ModeratePicIn{{URL: "http://x"}}})
	if out.Code != 0 {
		t.Fatalf("moderate batch pic %+v", out)
	}
	out = post(t, mux, "/YYB/Accounts", nil)
	if out.Code != 0 {
		t.Fatalf("yyb accounts %+v", out)
	}
	out = post(t, mux, "/YYB/QR", nil)
	if out.Code != 0 {
		t.Fatalf("yyb qr %+v", out)
	}
	out = post(t, mux, "/YYB/Poll", game.YYBSessionIn{SessionID: "s1"})
	if out.Code != 0 {
		t.Fatalf("yyb poll %+v", out)
	}
	out = post(t, mux, "/YYB/Confirm", game.YYBSessionIn{SessionID: "s1"})
	if out.Code != 0 {
		t.Fatalf("yyb confirm %+v", out)
	}
	out = post(t, mux, "/YYB/Code", game.YYBRefIn{Ref: "1"})
	if out.Code != 0 {
		t.Fatalf("yyb code %+v", out)
	}
	out = post(t, mux, "/YYB/Login", game.YYBRefIn{Ref: "1"})
	if out.Code != 0 {
		t.Fatalf("yyb login %+v", out)
	}
	out = post(t, mux, "/YYB/Delete", nil)
	if out.Code != 400 {
		t.Fatalf("yyb delete empty %+v", out)
	}
	out = post(t, mux, "/YYB/Delete", game.YYBRefIn{Ref: "1"})
	if out.Code != 0 {
		t.Fatalf("yyb delete %+v", out)
	}
	out = post(t, mux, "/YYB/Profile", game.YYBRefIn{Ref: "1"})
	if out.Code != 0 {
		t.Fatalf("yyb profile %+v", out)
	}
	out = post(t, mux, "/YYB/Phone", game.YYBRefIn{Ref: "1"})
	if out.Code != 0 {
		t.Fatalf("yyb phone %+v", out)
	}
	out = post(t, mux, "/YYB/WXData", game.YYBWXDataIn{Ref: "1"})
	if out.Code != 400 {
		t.Fatalf("yyb wxdata empty %+v", out)
	}
	out = post(t, mux, "/YYB/WXData", game.YYBWXDataIn{Ref: "1", Payload: map[string]any{"api_name": "x"}})
	if out.Code != 0 {
		t.Fatalf("yyb wxdata %+v", out)
	}
}

func TestDocs(t *testing.T) {
	mux := NewMux(NewHub(func() game.Session { return &fakeSession{} }))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))
	if rec.Code != 200 {
		t.Fatalf("openapi %d", rec.Code)
	}
	var spec map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatal(err)
	}
	comps, _ := spec["components"].(map[string]any)
	schemas, _ := comps["schemas"].(map[string]any)
	if schemas["User"] == nil || schemas["Land"] == nil || schemas["Item"] == nil {
		t.Fatalf("missing schemas %+v", schemas["User"])
	}
	paths, _ := spec["paths"].(map[string]any)
	harvest, _ := paths["/Farm/Harvest"].(map[string]any)
	if harvest == nil {
		t.Fatal("missing /Farm/Harvest")
	}
	tags, _ := spec["tags"].([]any)
	if len(tags) < 20 {
		t.Fatalf("tags %d", len(tags))
	}
	first, _ := tags[0].(map[string]any)
	if first["name"] != "System 系统" {
		t.Fatalf("first tag %+v", first)
	}
	if paths["/Farm/CleanEvents"] == nil || paths["/Friend/ApplyWX"] == nil || paths["/QQGroup/Bind"] == nil {
		t.Fatal("missing gap paths")
	}
	if paths["/QQVip/DailyStatus"] == nil || paths["/Marquee/List"] == nil || paths["/SystemOpen/Unlocked"] == nil || paths["/Subscribe/QQ"] == nil || paths["/Moderate/Text"] == nil {
		t.Fatal("missing platform paths")
	}
	if paths["/YYB/QR"] == nil || paths["/YYB/Code"] == nil || paths["/YYB/Login"] == nil ||
		paths["/YYB/Delete"] == nil || paths["/YYB/Profile"] == nil || paths["/YYB/Phone"] == nil || paths["/YYB/WXData"] == nil {
		t.Fatal("missing yyb paths")
	}
	raw := rec.Body.Bytes()
	ping := bytes.Index(raw, []byte(`"/System/Ping"`))
	huntPath := bytes.Index(raw, []byte(`"/Activity/HuntFinishCG"`))
	if ping < 0 || huntPath < 0 || ping > huntPath {
		t.Fatal("paths not classified")
	}
	hunt, _ := paths["/Activity/HuntFinishCG"].(map[string]any)
	post, _ := hunt["post"].(map[string]any)
	htags, _ := post["tags"].([]any)
	if len(htags) == 0 || htags[0] != "Activity 宠物寻宝" {
		t.Fatalf("hunt tags %+v", htags)
	}
}
