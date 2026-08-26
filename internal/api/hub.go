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

func (h *Hub) Refresh(ctx context.Context, in game.RefreshIn) ([]game.Land, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Refresh(ctx, in)
}

func (h *Hub) Harvest(ctx context.Context, in game.HarvestIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Harvest(ctx, in)
}

func (h *Hub) Plant(ctx context.Context, in game.PlantIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Plant(ctx, in)
}

func (h *Hub) Friends(ctx context.Context) ([]game.Friend, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Friends(ctx)
}

func (h *Hub) Applications(ctx context.Context) ([]game.Application, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Applications(ctx)
}

func (h *Hub) Accept(ctx context.Context, in game.GIDsIn) ([]game.Friend, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Accept(ctx, in)
}

func (h *Hub) Reject(ctx context.Context, in game.GIDsIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Reject(ctx, in)
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

func (h *Hub) ShareClaim(ctx context.Context) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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
	return s.Close()
}

func (h *Hub) Water(ctx context.Context, in game.LandOpIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Water(ctx, in)
}

func (h *Hub) Weed(ctx context.Context, in game.LandOpIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Weed(ctx, in)
}

func (h *Hub) Bug(ctx context.Context, in game.LandOpIn) error {
	s, err := h.current()
	if err != nil {
		return err
	}
	return s.Bug(ctx, in)
}

func (h *Hub) Fertilize(ctx context.Context, in game.FertilizeIn) error {
	s, err := h.current()
	if err != nil {
		return err
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

func (h *Hub) Bag(ctx context.Context) ([]game.BagItem, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Bag(ctx)
}

func (h *Hub) Sell(ctx context.Context, in game.SellIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Sell(ctx, in)
}

func (h *Hub) Use(ctx context.Context, in game.UseIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Use(ctx, in)
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

func (h *Hub) Steal(ctx context.Context, in game.HarvestIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Steal(ctx, in)
}

func (h *Hub) Remove(ctx context.Context, in game.RemoveIn) ([]game.Land, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) Farming(ctx context.Context, in game.LandOpIn) ([]game.Land, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) ClaimTask(ctx context.Context, in game.ClaimIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimTask(ctx, in)
}

func (h *Hub) ClaimDaily(ctx context.Context, in game.DailyIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) ClaimAllEmail(ctx context.Context, in game.EmailBoxIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) Signin(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.Signin(ctx, in)
}

func (h *Hub) ClaimProgress(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimProgress(ctx, in)
}

func (h *Hub) ShopBuy(ctx context.Context, in game.ActivityShopIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ShopBuy(ctx, in)
}

func (h *Hub) ClaimMega(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimMega(ctx, in)
}

func (h *Hub) TechSubmit(ctx context.Context, in game.TechIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.TechSubmit(ctx, in)
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

func (h *Hub) BrewClaim(ctx context.Context, in game.BrewClaimIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.BrewClaim(ctx, in)
}

func (h *Hub) ClaimRecall(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimRecall(ctx, in)
}

func (h *Hub) ClaimReturn(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimReturn(ctx, in)
}

func (h *Hub) ClaimInvite(ctx context.Context, in game.InviteIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimInvite(ctx, in)
}

func (h *Hub) ClaimNewcomer(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) MysteryBuy(ctx context.Context, in game.IDIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) ClaimPass(ctx context.Context) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimPass(ctx)
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

func (h *Hub) ClaimMonthCard(ctx context.Context, in game.IDIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) ClaimRedPacket(ctx context.Context, in game.IDIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
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

func (h *Hub) ClaimAlbum(ctx context.Context, in game.AlbumIn) ([]game.Item, error) {
	s, err := h.current()
	if err != nil {
		return nil, err
	}
	return s.ClaimAlbum(ctx, in)
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
