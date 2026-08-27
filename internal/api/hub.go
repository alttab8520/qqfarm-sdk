package api

import (
	"context"
	"errors"
	"sync"

	"github.com/alttab8520/qqfarm-sdk/internal/client"
	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

type Hub struct {
	newS game.Factory
	mu   sync.Mutex
	sess game.Session
	yyb  yybAPI
	res  resourceAPI
}

func NewHub(newS game.Factory) *Hub {
	if newS == nil {
		newS = client.New
	}
	return &Hub{newS: newS}
}

func (h *Hub) Login(ctx context.Context, in game.LoginIn) (game.User, error) {
	s := h.newS()
	user, err := s.Login(ctx, in)
	if err != nil {
		_ = s.Close()
		return game.User{}, err
	}
	h.mu.Lock()
	if h.sess != nil {
		_ = h.sess.Close()
	}
	h.sess = s
	h.mu.Unlock()
	return user, nil
}

func (h *Hub) current() (game.Session, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.sess == nil {
		return nil, game.ErrNotLogin
	}
	return h.sess, nil
}

func (h *Hub) Info() (game.User, error) {
	s, err := h.current()
	if err != nil {
		return game.User{}, err
	}
	return s.Info()
}

func (h *Hub) Heartbeat(ctx context.Context, in game.RefreshIn) (game.HeartbeatOut, error) {
	s, err := h.current()
	if err != nil {
		return game.HeartbeatOut{}, err
	}
	return s.Heartbeat(ctx, in)
}

func (h *Hub) Refresh(ctx context.Context, in game.RefreshIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Refresh(ctx, in)
}

func (h *Hub) RefreshLands(ctx context.Context, in game.RefreshLandsIn) ([]game.Land, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.RefreshLands(ctx, in)
}

func (h *Hub) CleanSocial(ctx context.Context, in game.CleanSocialIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.CleanSocial(ctx, in)
}

func (h *Hub) Harvest(ctx context.Context, in game.HarvestIn) (game.HarvestOut, error) {
	s, err := h.current()
	if err != nil {
		return game.HarvestOut{}, err
	}
	return s.Harvest(ctx, in)
}

func (h *Hub) Plant(ctx context.Context, in game.PlantIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Plant(ctx, in)
}

func (h *Hub) Friends(ctx context.Context) (game.FriendsOut, error) {
	s, err := h.current()
	if err != nil {
		return game.FriendsOut{}, err
	}
	return s.Friends(ctx)
}

func (h *Hub) Applications(ctx context.Context) (game.ApplicationsOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ApplicationsOut{}, err
	}
	return s.Applications(ctx)
}

func (h *Hub) Accept(ctx context.Context, in game.GIDsIn) (game.AcceptOut, error) {
	s, err := h.current()
	if err != nil {
		return game.AcceptOut{}, err
	}
	return s.Accept(ctx, in)
}

func (h *Hub) Reject(ctx context.Context, in game.GIDsIn) (game.RejectOut, error) {
	s, err := h.current()
	if err != nil {
		return game.RejectOut{}, err
	}
	return s.Reject(ctx, in)
}

func (h *Hub) SetTags(ctx context.Context, in game.TagsIn) (game.Friend, error) {
	s, err := h.current()
	if err != nil {
		return game.Friend{}, err
	}
	return s.SetTags(ctx, in)
}

func (h *Hub) DeleteFriend(ctx context.Context, in game.EnterIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.DeleteFriend(ctx, in)
}

func (h *Hub) ShareCheck(ctx context.Context) (bool, error) {
	s, err := h.current()
	if err != nil {
		return false, err
	}
	return s.ShareCheck(ctx)
}

func (h *Hub) ShareClaim(ctx context.Context) (game.ShareClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ShareClaimOut{}, err
	}
	return s.ShareClaim(ctx)
}

func (h *Hub) Help(ctx context.Context, in game.HelpIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Help(ctx, in)
}

func (h *Hub) Status() (game.Status, error) {
	s, err := h.current()
	if err != nil {
		return game.Status{}, err
	}
	return s.Status()
}

func (h *Hub) Logout() error {
	h.mu.Lock()
	s := h.sess
	h.sess = nil
	h.mu.Unlock()
	if s == nil {
		return game.ErrNotLogin
	}
	return s.Logout(context.Background())
}

func (h *Hub) Water(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Water(ctx, in)
}

func (h *Hub) Weed(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Weed(ctx, in)
}

func (h *Hub) Bug(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Bug(ctx, in)
}

func (h *Hub) Fertilize(ctx context.Context, in game.FertilizeIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Fertilize(ctx, in)
}

func (h *Hub) Enter(ctx context.Context, in game.EnterIn) (game.Visit, error) {
	s, err := h.current()
	if err != nil {
		return game.Visit{}, err
	}
	return s.Enter(ctx, in)
}

func (h *Hub) Leave(ctx context.Context, in game.EnterIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Leave(ctx, in)
}

func (h *Hub) Bag(ctx context.Context) (game.BagOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BagOut{}, err
	}
	return s.Bag(ctx)
}

func (h *Hub) Sell(ctx context.Context, in game.SellIn) (game.ItemOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ItemOpOut{}, err
	}
	return s.Sell(ctx, in)
}

func (h *Hub) Use(ctx context.Context, in game.UseIn) (game.ItemOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ItemOpOut{}, err
	}
	return s.Use(ctx, in)
}

func (h *Hub) BatchUse(ctx context.Context, in game.BatchUseIn) (game.ItemOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ItemOpOut{}, err
	}
	return s.BatchUse(ctx, in)
}

func (h *Hub) CancelNew(ctx context.Context, in game.IDIn) (int64, error) {
	s, err := h.current()
	if err != nil {
		return 0, err
	}
	return s.CancelNew(ctx, in)
}

func (h *Hub) Shops(ctx context.Context) ([]game.Shop, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Shops(ctx)
}

func (h *Hub) Goods(ctx context.Context, in game.ShopIn) ([]game.Goods, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Goods(ctx, in)
}

func (h *Hub) Buy(ctx context.Context, in game.BuyIn) (game.BuyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BuyOut{}, err
	}
	return s.Buy(ctx, in)
}

func (h *Hub) Steal(ctx context.Context, in game.HarvestIn) (game.HarvestOut, error) {
	s, err := h.current()
	if err != nil {
		return game.HarvestOut{}, err
	}
	return s.Steal(ctx, in)
}

func (h *Hub) Remove(ctx context.Context, in game.RemoveIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Remove(ctx, in)
}

func (h *Hub) Unlock(ctx context.Context, in game.LandIDIn) (game.Land, error) {
	s, err := h.current()
	if err != nil {
		return game.Land{}, err
	}
	return s.Unlock(ctx, in)
}

func (h *Hub) Upgrade(ctx context.Context, in game.LandIDIn) (game.Land, error) {
	s, err := h.current()
	if err != nil {
		return game.Land{}, err
	}
	return s.Upgrade(ctx, in)
}

func (h *Hub) Farming(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.Farming(ctx, in)
}

func (h *Hub) Weather(ctx context.Context) (game.Weather, error) {
	s, err := h.current()
	if err != nil {
		return game.Weather{}, err
	}
	return s.Weather(ctx)
}

func (h *Hub) CurrentWeather(ctx context.Context) (game.Weather, error) {
	s, err := h.current()
	if err != nil {
		return game.Weather{}, err
	}
	return s.CurrentWeather(ctx)
}

func (h *Hub) TodayWeather(ctx context.Context) ([]game.Weather, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.TodayWeather(ctx)
}

func (h *Hub) Tasks(ctx context.Context) (game.TaskBoard, error) {
	s, err := h.current()
	if err != nil {
		return game.TaskBoard{}, err
	}
	return s.Tasks(ctx)
}

func (h *Hub) ClaimTask(ctx context.Context, in game.ClaimIn) (game.TaskClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.TaskClaimOut{}, err
	}
	return s.ClaimTask(ctx, in)
}

func (h *Hub) ClaimDaily(ctx context.Context, in game.DailyIn) (game.TaskClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.TaskClaimOut{}, err
	}
	return s.ClaimDaily(ctx, in)
}

func (h *Hub) Emails(ctx context.Context, in game.EmailBoxIn) ([]game.Email, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Emails(ctx, in)
}

func (h *Hub) ReadEmail(ctx context.Context, in game.EmailIn) (game.EmailDetail, error) {
	s, err := h.current()
	if err != nil {
		return game.EmailDetail{}, err
	}
	return s.ReadEmail(ctx, in)
}

func (h *Hub) ClaimEmail(ctx context.Context, in game.EmailIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimEmail(ctx, in)
}

func (h *Hub) ClaimAllEmail(ctx context.Context, in game.EmailBoxIn) (game.EmailClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.EmailClaimOut{}, err
	}
	return s.ClaimAllEmail(ctx, in)
}

func (h *Hub) Activities(ctx context.Context) ([]game.Activity, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Activities(ctx)
}

func (h *Hub) ActivityGroup(ctx context.Context, in game.GroupIn) (game.ActivityGroup, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityGroup{}, err
	}
	return s.ActivityGroup(ctx, in)
}

func (h *Hub) Signin(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.Signin(ctx, in)
}

func (h *Hub) ClaimProgress(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ClaimProgress(ctx, in)
}

func (h *Hub) ShopBuy(ctx context.Context, in game.ActivityShopIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ShopBuy(ctx, in)
}

func (h *Hub) ShopBatchBuy(ctx context.Context, in game.ActivityBatchIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ShopBatchBuy(ctx, in)
}

func (h *Hub) RandBuy(ctx context.Context, in game.ActivityShopIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.RandBuy(ctx, in)
}

func (h *Hub) RandRefresh(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.RandRefresh(ctx, in)
}

func (h *Hub) ClaimMega(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ClaimMega(ctx, in)
}

func (h *Hub) TechSubmit(ctx context.Context, in game.TechIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.TechSubmit(ctx, in)
}

func (h *Hub) Draw(ctx context.Context, in game.DrawIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.Draw(ctx, in)
}

func (h *Hub) DrawHistory(ctx context.Context, in game.IDIn) (game.DrawHistoryOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DrawHistoryOut{}, err
	}
	return s.DrawHistory(ctx, in)
}

func (h *Hub) MarkViewed(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.MarkViewed(ctx, in)
}

func (h *Hub) RandBatchBuy(ctx context.Context, in game.ActivityBatchIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.RandBatchBuy(ctx, in)
}

func (h *Hub) LotteryHistory(ctx context.Context, in game.IDIn) (game.LotteryHistoryOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LotteryHistoryOut{}, err
	}
	return s.LotteryHistory(ctx, in)
}

func (h *Hub) CheerJoin(ctx context.Context, in game.CheerJoinIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.CheerJoin(ctx, in)
}

func (h *Hub) CheerSubmit(ctx context.Context, in game.CheerSubmitIn) (game.CheerSubmitOut, error) {
	s, err := h.current()
	if err != nil {
		return game.CheerSubmitOut{}, err
	}
	return s.CheerSubmit(ctx, in)
}

func (h *Hub) CheerClaim(ctx context.Context, in game.CheerClaimIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.CheerClaim(ctx, in)
}

func (h *Hub) Recallable(ctx context.Context, in game.IDIn) (game.RecallListOut, error) {
	s, err := h.current()
	if err != nil {
		return game.RecallListOut{}, err
	}
	return s.Recallable(ctx, in)
}

func (h *Hub) Recalled(ctx context.Context, in game.IDIn) (game.RecallListOut, error) {
	s, err := h.current()
	if err != nil {
		return game.RecallListOut{}, err
	}
	return s.Recalled(ctx, in)
}

func (h *Hub) CharityShare(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.CharityShare(ctx, in)
}

func (h *Hub) CharityDonate(ctx context.Context, in game.IDIn) (game.CharityDonateOut, error) {
	s, err := h.current()
	if err != nil {
		return game.CharityDonateOut{}, err
	}
	return s.CharityDonate(ctx, in)
}

func (h *Hub) CharityClaim(ctx context.Context, in game.CharityClaimIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.CharityClaim(ctx, in)
}

func (h *Hub) CharityXhh(ctx context.Context, in game.IDIn) (game.CharityXhhOut, error) {
	s, err := h.current()
	if err != nil {
		return game.CharityXhhOut{}, err
	}
	return s.CharityXhh(ctx, in)
}

func (h *Hub) CharityAgree(ctx context.Context, in game.CharityAgreeIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.CharityAgree(ctx, in)
}

func (h *Hub) HuntFinishCG(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntFinishCG(ctx, in)
}

func (h *Hub) HuntGuide(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntGuide(ctx, in)
}

func (h *Hub) HuntFeed(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntFeed(ctx, in)
}

func (h *Hub) HuntDraw(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntDraw(ctx, in)
}

func (h *Hub) HuntLog(ctx context.Context, in game.IDIn) (game.HuntLogOut, error) {
	s, err := h.current()
	if err != nil {
		return game.HuntLogOut{}, err
	}
	return s.HuntLog(ctx, in)
}

func (h *Hub) HuntClaimStory(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntClaimStory(ctx, in)
}

func (h *Hub) HuntClaimSeed(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntClaimSeed(ctx, in)
}

func (h *Hub) HuntRefreshCharm(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntRefreshCharm(ctx, in)
}

func (h *Hub) HuntEquip(ctx context.Context, in game.HuntEquipIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntEquip(ctx, in)
}

func (h *Hub) HuntBattle(ctx context.Context, in game.HuntBattleIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntBattle(ctx, in)
}

func (h *Hub) HuntPlunderedLog(ctx context.Context, in game.IDIn) (game.HuntLogOut, error) {
	s, err := h.current()
	if err != nil {
		return game.HuntLogOut{}, err
	}
	return s.HuntPlunderedLog(ctx, in)
}

func (h *Hub) HuntOpen(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntOpen(ctx, in)
}

func (h *Hub) HuntEscort(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntEscort(ctx, in)
}

func (h *Hub) HuntCompensate(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntCompensate(ctx, in)
}

func (h *Hub) HuntFriendInfo(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.HuntFriendInfo(ctx, in)
}

func (h *Hub) Lottery(ctx context.Context, in game.LotteryIn) (game.LotteryOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LotteryOut{}, err
	}
	return s.Lottery(ctx, in)
}

func (h *Hub) BrewStart(ctx context.Context, in game.BrewStartIn) (game.BrewStartOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BrewStartOut{}, err
	}
	return s.BrewStart(ctx, in)
}

func (h *Hub) BrewStep(ctx context.Context, in game.ActivityOpIn) (game.BrewStepOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BrewStepOut{}, err
	}
	return s.BrewStep(ctx, in)
}

func (h *Hub) BrewClaim(ctx context.Context, in game.BrewClaimIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.BrewClaim(ctx, in)
}

func (h *Hub) ClaimRecall(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ClaimRecall(ctx, in)
}

func (h *Hub) ClaimReturn(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ClaimReturn(ctx, in)
}

func (h *Hub) ClaimInvite(ctx context.Context, in game.InviteIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ClaimInvite(ctx, in)
}

func (h *Hub) ClaimNewcomer(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return s.ClaimNewcomer(ctx, in)
}

func (h *Hub) SendGift(ctx context.Context, in game.GiftIn) (int64, error) {
	s, err := h.current()
	if err != nil {
		return 0, err
	}
	return s.SendGift(ctx, in)
}

func (h *Hub) Mystery(ctx context.Context) (game.MysteryShop, error) {
	s, err := h.current()
	if err != nil {
		return game.MysteryShop{}, err
	}
	return s.Mystery(ctx)
}

func (h *Hub) MysteryBuy(ctx context.Context, in game.IDIn) (game.MysteryBuyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.MysteryBuyOut{}, err
	}
	return s.MysteryBuy(ctx, in)
}

func (h *Hub) MysteryLeave(ctx context.Context) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.MysteryLeave(ctx)
}

func (h *Hub) Season(ctx context.Context) (game.Season, error) {
	s, err := h.current()
	if err != nil {
		return game.Season{}, err
	}
	return s.Season(ctx)
}

func (h *Hub) ClaimPass(ctx context.Context) (game.PassClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.PassClaimOut{}, err
	}
	return s.ClaimPass(ctx)
}

func (h *Hub) MallDiamonds(ctx context.Context) ([]game.Product, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.MallDiamonds(ctx)
}

func (h *Hub) Mall(ctx context.Context, in game.MallIn) ([]game.Product, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Mall(ctx, in)
}

func (h *Hub) MallBuy(ctx context.Context, in game.MallBuyIn) (game.BuyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BuyOut{}, err
	}
	return s.MallBuy(ctx, in)
}

func (h *Hub) MonthCards(ctx context.Context) ([]game.MonthCard, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.MonthCards(ctx)
}

func (h *Hub) ClaimMonthCard(ctx context.Context, in game.IDIn) (game.MonthCardClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.MonthCardClaimOut{}, err
	}
	return s.ClaimMonthCard(ctx, in)
}

func (h *Hub) RedPackets(ctx context.Context) ([]game.RedPacket, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.RedPackets(ctx)
}

func (h *Hub) ClaimRedPacket(ctx context.Context, in game.IDIn) (game.RedPacketClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.RedPacketClaimOut{}, err
	}
	return s.ClaimRedPacket(ctx, in)
}

func (h *Hub) ClaimAllRedPackets(ctx context.Context) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimAllRedPackets(ctx)
}

func (h *Hub) VisitSummary(ctx context.Context) (game.VisitSummary, error) {
	s, err := h.current()
	if err != nil {
		return game.VisitSummary{}, err
	}
	return s.VisitSummary(ctx)
}

func (h *Hub) VisitLogs(ctx context.Context) ([]game.VisitLog, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.VisitLogs(ctx)
}

func (h *Hub) Album(ctx context.Context, in game.AlbumIn) (game.Album, error) {
	s, err := h.current()
	if err != nil {
		return game.Album{}, err
	}
	return s.Album(ctx, in)
}

func (h *Hub) ClaimAlbum(ctx context.Context, in game.AlbumIn) (game.AlbumClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.AlbumClaimOut{}, err
	}
	return s.ClaimAlbum(ctx, in)
}

func (h *Hub) Dog(ctx context.Context) (game.DogYard, error) {
	s, err := h.current()
	if err != nil {
		return game.DogYard{}, err
	}
	return s.Dog(ctx)
}

func (h *Hub) Feed(ctx context.Context, in game.FeedIn) (int64, error) {
	s, err := h.current()
	if err != nil {
		return 0, err
	}
	return s.Feed(ctx, in)
}

func (h *Hub) ClaimDogGifts(ctx context.Context) (game.DogGiftOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DogGiftOut{}, err
	}
	return s.ClaimDogGifts(ctx)
}

func (h *Hub) DogLogs(ctx context.Context, in game.PageIn) (game.DogLogsOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DogLogsOut{}, err
	}
	return s.DogLogs(ctx, in)
}

func (h *Hub) DeployDog(ctx context.Context, in game.IDIn) (game.DeployOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DeployOut{}, err
	}
	return s.DeployDog(ctx, in)
}

func (h *Hub) WithdrawDog(ctx context.Context) (game.DeployOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DeployOut{}, err
	}
	return s.WithdrawDog(ctx)
}

func (h *Hub) ActivateDog(ctx context.Context, in game.IDIn) (game.Dog, error) {
	s, err := h.current()
	if err != nil {
		return game.Dog{}, err
	}
	return s.ActivateDog(ctx, in)
}

func (h *Hub) Bulletins(ctx context.Context, in game.PageIn) ([]game.Bulletin, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Bulletins(ctx, in)
}

func (h *Hub) ReadBulletin(ctx context.Context, in game.IDIn) (game.BulletinDetail, error) {
	s, err := h.current()
	if err != nil {
		return game.BulletinDetail{}, err
	}
	return s.ReadBulletin(ctx, in)
}

func (h *Hub) Mutants(ctx context.Context) ([]game.Mutant, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Mutants(ctx)
}

func (h *Hub) Career(ctx context.Context, in game.EnterIn) (game.Career, error) {
	s, err := h.current()
	if err != nil {
		return game.Career{}, err
	}
	return s.Career(ctx, in)
}

func (h *Hub) Ranks(ctx context.Context, in game.RankIn) (game.RankBoard, error) {
	s, err := h.current()
	if err != nil {
		return game.RankBoard{}, err
	}
	return s.Ranks(ctx, in)
}

func (h *Hub) Avatars(ctx context.Context, in game.TypeIn) ([]game.Avatar, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Avatars(ctx, in)
}

func (h *Hub) EquippedAvatars(ctx context.Context) ([]game.Avatar, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.EquippedAvatars(ctx)
}

func (h *Hub) EquipAvatar(ctx context.Context, in game.AvatarEquipIn) (game.Avatar, error) {
	s, err := h.current()
	if err != nil {
		return game.Avatar{}, err
	}
	return s.EquipAvatar(ctx, in)
}

func (h *Hub) Skins(ctx context.Context) ([]game.Skin, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Skins(ctx)
}

func (h *Hub) EquippedSkins(ctx context.Context) ([]game.Skin, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.EquippedSkins(ctx)
}

func (h *Hub) EquipSkin(ctx context.Context, in game.SkinEquipIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.EquipSkin(ctx, in)
}

func (h *Hub) Drops(ctx context.Context) ([]game.Drop, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Drops(ctx)
}

func (h *Hub) SolarTerms(ctx context.Context) (game.SolarOut, error) {
	s, err := h.current()
	if err != nil {
		return game.SolarOut{}, err
	}
	return s.SolarTerms(ctx)
}

func (h *Hub) ClaimSolar(ctx context.Context, in game.IDIn) (game.SolarClaimOut, error) {
	s, err := h.current()
	if err != nil {
		return game.SolarClaimOut{}, err
	}
	return s.ClaimSolar(ctx, in)
}

func (h *Hub) ClaimAllSolar(ctx context.Context) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimAllSolar(ctx)
}

func (h *Hub) AchieveView(ctx context.Context, in game.AchieveIn) (game.AchieveScope, error) {
	s, err := h.current()
	if err != nil {
		return game.AchieveScope{}, err
	}
	return s.AchieveView(ctx, in)
}

func (h *Hub) ClaimAchieveGoal(ctx context.Context, in game.AchieveGoalIn) (game.AchieveGoalOut, error) {
	s, err := h.current()
	if err != nil {
		return game.AchieveGoalOut{}, err
	}
	return s.ClaimAchieveGoal(ctx, in)
}

func (h *Hub) ClaimAchieveLevel(ctx context.Context, in game.AchieveIn) (game.ItemOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ItemOpOut{}, err
	}
	return s.ClaimAchieveLevel(ctx, in)
}

func (h *Hub) Brief(ctx context.Context, in game.EnterIn) (game.User, error) {
	s, err := h.current()
	if err != nil {
		return game.User{}, err
	}
	return s.Brief(ctx, in)
}

func (h *Hub) BatchInfo(ctx context.Context, in game.GIDsIn) ([]game.User, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.BatchInfo(ctx, in)
}

func (h *Hub) ArkClick(ctx context.Context, in game.ArkIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.ArkClick(ctx, in)
}

func (h *Hub) PutInsects(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.PutInsects(ctx, in)
}

func (h *Hub) PutWeeds(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.PutWeeds(ctx, in)
}

func (h *Hub) PutSocial(ctx context.Context, in game.PutSocialIn) (game.LandOpOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LandOpOut{}, err
	}
	return s.PutSocial(ctx, in)
}

func (h *Hub) CanOperate(ctx context.Context, in game.CanOperateIn) (game.CanOperateOut, error) {
	s, err := h.current()
	if err != nil {
		return game.CanOperateOut{}, err
	}
	return s.CanOperate(ctx, in)
}

func (h *Hub) SyncFriends(ctx context.Context, in game.OpenIDsIn) ([]game.Friend, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.SyncFriends(ctx, in)
}

func (h *Hub) GameFriends(ctx context.Context, in game.GIDsIn) ([]game.Friend, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.GameFriends(ctx, in)
}

func (h *Hub) BlockFriend(ctx context.Context, in game.EnterIn) (game.Blocked, error) {
	s, err := h.current()
	if err != nil {
		return game.Blocked{}, err
	}
	return s.BlockFriend(ctx, in)
}

func (h *Hub) UnblockFriend(ctx context.Context, in game.EnterIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.UnblockFriend(ctx, in)
}

func (h *Hub) BlockList(ctx context.Context) ([]game.Blocked, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.BlockList(ctx)
}

func (h *Hub) ShareKey(ctx context.Context) (string, error) {
	s, err := h.current()
	if err != nil {
		return "", err
	}
	return s.ShareKey(ctx)
}

func (h *Hub) InviteInfo(ctx context.Context, in game.IDIn) (game.ShareInviteOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ShareInviteOut{}, err
	}
	return s.InviteInfo(ctx, in)
}

func (h *Hub) InviteAward(ctx context.Context, in game.IDIn) (game.ShareAwardOut, error) {
	s, err := h.current()
	if err != nil {
		return game.ShareAwardOut{}, err
	}
	return s.InviteAward(ctx, in)
}

func (h *Hub) PosterShown(ctx context.Context, in game.IDIn) (bool, error) {
	s, err := h.current()
	if err != nil {
		return false, err
	}
	return s.PosterShown(ctx, in)
}

func (h *Hub) ReportInvite(ctx context.Context, in game.InviteReportIn) (bool, error) {
	s, err := h.current()
	if err != nil {
		return false, err
	}
	return s.ReportInvite(ctx, in)
}

func (h *Hub) LockItems(ctx context.Context, in game.UIDsIn) (game.LockOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LockOut{}, err
	}
	return s.LockItems(ctx, in)
}

func (h *Hub) UnlockItems(ctx context.Context, in game.UIDsIn) (game.LockOut, error) {
	s, err := h.current()
	if err != nil {
		return game.LockOut{}, err
	}
	return s.UnlockItems(ctx, in)
}

func (h *Hub) AutoBuy(ctx context.Context, in game.AutoBuyIn) (game.BuyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.BuyOut{}, err
	}
	return s.AutoBuy(ctx, in)
}

func (h *Hub) MysteryAutoBuy(ctx context.Context, in game.MysteryAutoIn) (game.MysteryAutoOut, error) {
	s, err := h.current()
	if err != nil {
		return game.MysteryAutoOut{}, err
	}
	return s.MysteryAutoBuy(ctx, in)
}

func (h *Hub) BuyDog(ctx context.Context, in game.DogBuyIn) (game.DogBuyOut, error) {
	s, err := h.current()
	if err != nil {
		return game.DogBuyOut{}, err
	}
	return s.BuyDog(ctx, in)
}

func (h *Hub) AlbumLevels(ctx context.Context, in game.AlbumIn) (game.AlbumLevels, error) {
	s, err := h.current()
	if err != nil {
		return game.AlbumLevels{}, err
	}
	return s.AlbumLevels(ctx, in)
}

func (h *Hub) SetSplashed(ctx context.Context, in game.IDIn) (bool, error) {
	s, err := h.current()
	if err != nil {
		return false, err
	}
	return s.SetSplashed(ctx, in)
}

func (h *Hub) Invitees(ctx context.Context, in game.IDIn) ([]game.Invitee, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Invitees(ctx, in)
}

func (h *Hub) MarkAlbum(ctx context.Context, in game.AlbumIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.MarkAlbum(ctx, in)
}

func (h *Hub) MarkAvatar(ctx context.Context, in game.IDIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.MarkAvatar(ctx, in)
}

func (h *Hub) MarkSkin(ctx context.Context, in game.IDIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.MarkSkin(ctx, in)
}

func (h *Hub) VisitPopup(ctx context.Context) (game.VisitPopup, error) {
	s, err := h.current()
	if err != nil {
		return game.VisitPopup{}, err
	}
	return s.VisitPopup(ctx)
}

func (h *Hub) VisitPage(ctx context.Context, in game.VisitPageIn) ([]game.VisitLog, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.VisitPage(ctx, in)
}

func (h *Hub) DismissVisit(ctx context.Context, in game.EnterIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.DismissVisit(ctx, in)
}

func (h *Hub) DeleteVisit(ctx context.Context, in game.VisitDeleteIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.DeleteVisit(ctx, in)
}

func (h *Hub) BatchReadEmail(ctx context.Context, in game.EmailIDsIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.BatchReadEmail(ctx, in)
}

func (h *Hub) BatchDeleteEmail(ctx context.Context, in game.EmailIDsIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.BatchDeleteEmail(ctx, in)
}

func (h *Hub) ReportTask(ctx context.Context, in game.TaskReportIn) (game.TaskBoard, error) {
	s, err := h.current()
	if err != nil {
		return game.TaskBoard{}, err
	}
	return s.ReportTask(ctx, in)
}

func (h *Hub) MallProfiles(ctx context.Context) ([]game.MallProfile, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.MallProfiles(ctx)
}

func (h *Hub) SolarRedDot(ctx context.Context) (game.RedDot, error) {
	s, err := h.current()
	if err != nil {
		return game.RedDot{}, err
	}
	return s.SolarRedDot(ctx)
}

func failCode(err error) (int, string) {
	if err == nil {
		return 0, "ok"
	}
	if errors.Is(err, game.ErrNotLogin) {
		return 401, err.Error()
	}
	return 502, err.Error()
}
