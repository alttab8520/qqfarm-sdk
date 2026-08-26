package client

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/alttab8520/qqfarm-sdk/internal/ace"
	"github.com/alttab8520/qqfarm-sdk/internal/crypto"
	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/gate"
	"github.com/coder/websocket"
)

const wsBase = "wss://gate-obt.nqf.qq.com/prod/ws"

type Client struct {
	enc     crypto.Encryptor
	openErr error

	mu        sync.Mutex
	conn      *websocket.Conn
	pending   map[int64]chan gate.Message
	seq       int64
	serverSeq int64
	user      game.User
	bag       []game.BagItem
	loggedIn  bool
	closed    bool
	life      *ace.Life
}

func New() game.Session {
	enc, err := crypto.Open()
	c := NewWith(enc)
	c.openErr = err
	return c
}

func NewWith(enc crypto.Encryptor) *Client {
	if enc == nil {
		enc = crypto.Identity{}
	}
	return &Client{enc: enc, pending: map[int64]chan gate.Message{}}
}

func (c *Client) Login(ctx context.Context, in game.LoginIn) (game.User, error) {
	if c.openErr != nil {
		return game.User{}, c.openErr
	}
	if in.Code == "" {
		return game.User{}, fmt.Errorf("code 不能为空")
	}
	_ = c.closeConn()
	c.mu.Lock()
	c.pending = map[int64]chan gate.Message{}
	c.seq = 0
	c.serverSeq = 0
	c.closed = false
	c.loggedIn = false
	c.mu.Unlock()

	url := fmt.Sprintf("%s?platform=wx&os=Windows&ver=%s&code=%s&openID=%s",
		wsBase, gameVersion(), in.Code, in.OpenID)
	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPHeader: http.Header{"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"}},
	})
	if err != nil {
		return game.User{}, fmt.Errorf("连接网关失败: %w", err)
	}
	conn.SetReadLimit(1 << 20)
	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()
	go c.readLoop()

	body, err := c.call(ctx, "User", "Login", encodeLogin())
	if err != nil {
		_ = c.Close()
		return game.User{}, err
	}
	user, err := decodeUser(body)
	if err != nil {
		_ = c.Close()
		return game.User{}, err
	}
	bag, _ := decodeLoginBag(body)
	if openID := user.OpenID; openID != "" {
		if err := c.enc.BindUser(openID); err != nil {
			_ = c.Close()
			return game.User{}, fmt.Errorf("绑定用户失败: %w", err)
		}
	}
	c.mu.Lock()
	c.user = user
	c.bag = bag
	c.loggedIn = true
	c.life = ace.New(crypto.AsEngine(c.enc), c.uploadAnti)
	c.life.Start()
	c.mu.Unlock()
	return user, nil
}

func (c *Client) Info() (game.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.loggedIn {
		return game.User{}, game.ErrNotLogin
	}
	return c.user, nil
}

func (c *Client) Status() (game.Status, error) {
	user, err := c.Info()
	if err != nil {
		return game.Status{}, err
	}
	out := game.Status{LoggedIn: true, User: user}
	c.mu.Lock()
	life := c.life
	c.mu.Unlock()
	if life != nil {
		snap := life.Snapshot()
		out.ACE = game.ACE{Uploads: snap.Uploads, Reports: snap.Reports, Failures: snap.Failures, LastError: snap.LastError}
	}
	return out, nil
}

func (c *Client) Refresh(ctx context.Context, in game.RefreshIn) ([]game.Land, error) {
	user, err := c.Info()
	if err != nil {
		return nil, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	body, err := c.call(ctx, "Plant", "AllLands", encodeAllLands(host))
	if err != nil {
		return nil, err
	}
	return decodeLands(body)
}

func (c *Client) Harvest(ctx context.Context, in game.HarvestIn) ([]game.Item, error) {
	user, err := c.Info()
	if err != nil {
		return nil, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	if host != user.GID {
		if err := c.ensureCanOperate(ctx, host, 10004); err != nil {
			return nil, err
		}
	}
	body, err := c.call(ctx, "Plant", "Harvest", encodeHarvest(in.LandIDs, host, in.IsAll))
	if err != nil {
		return nil, err
	}
	return decodeItems(body)
}

func (c *Client) Steal(ctx context.Context, in game.HarvestIn) ([]game.Item, error) {
	if in.HostGID <= 0 {
		return nil, fmt.Errorf("host_gid 不能为空")
	}
	return c.Harvest(ctx, in)
}

func (c *Client) Remove(ctx context.Context, in game.RemoveIn) ([]game.Land, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if len(in.LandIDs) == 0 {
		return nil, fmt.Errorf("land_ids 不能为空")
	}
	body, err := c.call(ctx, "Plant", "RemovePlant", encodeRemove(in.LandIDs))
	if err != nil {
		return nil, err
	}
	return decodeLands(body)
}

func (c *Client) Unlock(ctx context.Context, in game.LandIDIn) (game.Land, error) {
	if _, err := c.Info(); err != nil {
		return game.Land{}, err
	}
	if in.LandID <= 0 {
		return game.Land{}, fmt.Errorf("land_id 不能为空")
	}
	body, err := c.call(ctx, "Plant", "UnlockLand", encodeLandID(in.LandID))
	if err != nil {
		return game.Land{}, err
	}
	return decodeOneLand(body)
}

func (c *Client) Upgrade(ctx context.Context, in game.LandIDIn) (game.Land, error) {
	if _, err := c.Info(); err != nil {
		return game.Land{}, err
	}
	if in.LandID <= 0 {
		return game.Land{}, fmt.Errorf("land_id 不能为空")
	}
	body, err := c.call(ctx, "Plant", "UpgradeLand", encodeLandID(in.LandID))
	if err != nil {
		return game.Land{}, err
	}
	return decodeOneLand(body)
}

func (c *Client) Farming(ctx context.Context, in game.LandOpIn) ([]game.Land, error) {
	user, err := c.Info()
	if err != nil {
		return nil, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	body, err := c.call(ctx, "Plant", "Farming", encodeFarming(in.LandIDs, host))
	if err != nil {
		return nil, err
	}
	return decodeLands(body)
}

func (c *Client) ensureCanOperate(ctx context.Context, hostGID, operationID int64) error {
	body, err := c.call(ctx, "Plant", "CheckCanOperate", encodeCheckCanOperate(hostGID, operationID))
	if err != nil || len(body) == 0 {
		return nil
	}
	ok, err := decodeCanOperate(body)
	if err != nil || ok {
		return nil
	}
	return fmt.Errorf("不能操作")
}

func (c *Client) Plant(ctx context.Context, in game.PlantIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.SeedID <= 0 || len(in.LandIDs) == 0 {
		return fmt.Errorf("seed_id 和 land_ids 不能为空")
	}
	_, err := c.call(ctx, "Plant", "Plant", encodePlant(in.SeedID, in.LandIDs))
	return err
}

func (c *Client) Water(ctx context.Context, in game.LandOpIn) error {
	return c.landOp(ctx, "WaterLand", in)
}

func (c *Client) Weed(ctx context.Context, in game.LandOpIn) error {
	return c.landOp(ctx, "WeedOut", in)
}

func (c *Client) Bug(ctx context.Context, in game.LandOpIn) error {
	return c.landOp(ctx, "Insecticide", in)
}

func (c *Client) Fertilize(ctx context.Context, in game.FertilizeIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.FertilizerID <= 0 || len(in.LandIDs) == 0 {
		return fmt.Errorf("fertilizer_id 和 land_ids 不能为空")
	}
	_, err := c.call(ctx, "Plant", "Fertilize", encodeFertilize(in.LandIDs, in.FertilizerID))
	return err
}

func (c *Client) landOp(ctx context.Context, method string, in game.LandOpIn) error {
	user, err := c.Info()
	if err != nil {
		return err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	_, err = c.call(ctx, "Plant", method, encodeLandOp(in.LandIDs, host))
	return err
}

func (c *Client) Friends(ctx context.Context) ([]game.Friend, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Friend", "GetAll", nil)
	if err != nil {
		return nil, err
	}
	return decodeFriends(body)
}

func (c *Client) Help(ctx context.Context, in game.HelpIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	if _, err := c.call(ctx, "Visit", "Enter", encodeEnter(in.GID)); err != nil {
		return err
	}
	_, err := c.call(ctx, "Plant", "WaterLand", encodeWater(in.LandIDs, in.GID))
	return err
}

func (c *Client) Enter(ctx context.Context, in game.EnterIn) (game.Visit, error) {
	if _, err := c.Info(); err != nil {
		return game.Visit{}, err
	}
	if in.GID <= 0 {
		return game.Visit{}, fmt.Errorf("gid 不能为空")
	}
	body, err := c.call(ctx, "Visit", "Enter", encodeEnter(in.GID))
	if err != nil {
		return game.Visit{}, err
	}
	return decodeVisit(body)
}

func (c *Client) Applications(ctx context.Context) ([]game.Application, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Friend", "GetApplications", nil)
	if err != nil {
		return nil, err
	}
	return decodeApplications(body)
}

func (c *Client) Accept(ctx context.Context, in game.GIDsIn) ([]game.Friend, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if len(in.GIDs) == 0 {
		return nil, fmt.Errorf("gids 不能为空")
	}
	body, err := c.call(ctx, "Friend", "AcceptFriends", encodeGIDs(in.GIDs))
	if err != nil {
		return nil, err
	}
	return decodeFriends(body)
}

func (c *Client) Reject(ctx context.Context, in game.GIDsIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if len(in.GIDs) == 0 {
		return fmt.Errorf("gids 不能为空")
	}
	_, err := c.call(ctx, "Friend", "RejectFriends", encodeGIDs(in.GIDs))
	return err
}

func (c *Client) DeleteFriend(ctx context.Context, in game.EnterIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	_, err := c.call(ctx, "Friend", "DelFriend", encodeLeave(in.GID))
	return err
}

func (c *Client) ShareCheck(ctx context.Context) (bool, error) {
	if _, err := c.Info(); err != nil {
		return false, err
	}
	body, err := c.call(ctx, "Share", "CheckCanShare", nil)
	if err != nil {
		return false, err
	}
	return decodeCanShare(body)
}

func (c *Client) ShareClaim(ctx context.Context) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if _, err := c.call(ctx, "Share", "ReportShare", encodeShareFlag()); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Share", "ClaimShareReward", encodeShareFlag())
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 3)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) Leave(ctx context.Context, in game.EnterIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	_, err := c.call(ctx, "Visit", "Leave", encodeLeave(in.GID))
	return err
}

func (c *Client) Bag(ctx context.Context) ([]game.BagItem, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Item", "Bag", nil)
	if err != nil {
		return nil, err
	}
	items, err := decodeBagReply(body)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.bag = items
	c.mu.Unlock()
	return items, nil
}

func (c *Client) Sell(ctx context.Context, in game.SellIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if len(in.Items) == 0 {
		return nil, fmt.Errorf("items 不能为空")
	}
	for _, it := range in.Items {
		if it.ID <= 0 || it.Count <= 0 {
			return nil, fmt.Errorf("每个 item 需要 id 和 count")
		}
	}
	body, err := c.call(ctx, "Item", "Sell", encodeSell(in.Items))
	if err != nil {
		return nil, err
	}
	got, err := decodeItemsAt(body, 2)
	if err != nil {
		return nil, err
	}
	c.applyCurrency(got, nil)
	_, _ = c.Bag(ctx)
	return got, nil
}

func (c *Client) Use(ctx context.Context, in game.UseIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.Count <= 0 {
		in.Count = 1
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Item", "Use", encodeUse(in))
	if err != nil {
		return nil, err
	}
	got, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	c.applyCurrency(got, nil)
	_, _ = c.Bag(ctx)
	return got, nil
}

func (c *Client) Shops(ctx context.Context) ([]game.Shop, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Shop", "ShopProfiles", nil)
	if err != nil {
		return nil, err
	}
	return decodeShops(body)
}

func (c *Client) Goods(ctx context.Context, in game.ShopIn) ([]game.Goods, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ShopID <= 0 {
		in.ShopID = 2
	}
	body, err := c.call(ctx, "Shop", "ShopInfo", encodeShopInfo(in.ShopID))
	if err != nil {
		return nil, err
	}
	return decodeGoods(body)
}

func (c *Client) Buy(ctx context.Context, in game.BuyIn) (game.BuyOut, error) {
	var out game.BuyOut
	if _, err := c.Info(); err != nil {
		return out, err
	}
	if in.Num <= 0 {
		in.Num = 1
	}
	if in.GoodsID <= 0 {
		return out, fmt.Errorf("goods_id 不能为空")
	}
	body, err := c.call(ctx, "Shop", "BuyGoods", encodeBuy(in.GoodsID, in.Num, in.Price))
	if err != nil {
		return out, err
	}
	got, err := decodeItemsAt(body, 2)
	if err != nil {
		return out, err
	}
	cost, err := decodeItemsAt(body, 3)
	if err != nil {
		return out, err
	}
	out.Items = got
	out.Cost = cost
	c.applyCurrency(got, cost)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) Weather(ctx context.Context) (game.Weather, error) {
	if _, err := c.Info(); err != nil {
		return game.Weather{}, err
	}
	body, err := c.call(ctx, "Weather", "GetWeatherStatus", nil)
	if err != nil {
		return game.Weather{}, err
	}
	return decodeWeatherStatus(body)
}

func (c *Client) TodayWeather(ctx context.Context) ([]game.Weather, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Weather", "GetTodayWeather", nil)
	if err != nil {
		return nil, err
	}
	return decodeTodayWeather(body)
}

func (c *Client) Tasks(ctx context.Context) (game.TaskBoard, error) {
	if _, err := c.Info(); err != nil {
		return game.TaskBoard{}, err
	}
	body, err := c.call(ctx, "Task", "TaskInfo", nil)
	if err != nil {
		return game.TaskBoard{}, err
	}
	return decodeTaskBoard(body)
}

func (c *Client) ClaimTask(ctx context.Context, in game.ClaimIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	ids := in.IDs
	if in.ID > 0 {
		ids = append([]int64{in.ID}, ids...)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	var (
		body []byte
		err  error
	)
	if len(ids) == 1 {
		body, err = c.call(ctx, "Task", "ClaimTaskReward", encodeClaimTask(ids[0]))
	} else {
		body, err = c.call(ctx, "Task", "BatchClaimTaskReward", encodeBatchClaim(ids))
	}
	if err != nil {
		return nil, err
	}
	return decodeItemsAt(body, 1)
}

func (c *Client) ClaimDaily(ctx context.Context, in game.DailyIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if len(in.PointIDs) == 0 {
		return nil, fmt.Errorf("point_ids 不能为空")
	}
	body, err := c.call(ctx, "Task", "ClaimDailyReward", encodeBatchClaim(in.PointIDs))
	if err != nil {
		return nil, err
	}
	return decodeItemsAt(body, 1)
}

func emailBox(box int64) int64 {
	if box <= 0 {
		return 1
	}
	return box
}

func (c *Client) Emails(ctx context.Context, in game.EmailBoxIn) ([]game.Email, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Email", "GetEmailList", encodeBox(emailBox(in.Box)))
	if err != nil {
		return nil, err
	}
	return decodeEmails(body)
}

func (c *Client) ReadEmail(ctx context.Context, in game.EmailIn) (game.EmailDetail, error) {
	if _, err := c.Info(); err != nil {
		return game.EmailDetail{}, err
	}
	if in.ID == "" {
		return game.EmailDetail{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Email", "ReadEmail", encodeEmailRef(emailBox(in.Box), in.ID))
	if err != nil {
		return game.EmailDetail{}, err
	}
	return decodeEmailDetail(body)
}

func (c *Client) ClaimEmail(ctx context.Context, in game.EmailIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID == "" {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Email", "ClaimEmail", encodeEmailRef(emailBox(in.Box), in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimAllEmail(ctx context.Context, in game.EmailBoxIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	boxes := []int64{emailBox(in.Box)}
	if in.Box <= 0 {
		boxes = []int64{1, 2}
	}
	var all []game.Item
	for _, box := range boxes {
		body, err := c.call(ctx, "Email", "BatchClaimEmail", encodeBox(box))
		if err != nil {
			continue
		}
		items, err := decodeItemsAt(body, 1)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
	}
	_, _ = c.Bag(ctx)
	return all, nil
}

func (c *Client) Activities(ctx context.Context) ([]game.Activity, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Activity", "List", nil)
	if err != nil {
		return nil, err
	}
	return decodeActivities(body)
}

func (c *Client) Signin(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	if in.RewardID <= 0 {
		list, err := c.Activities(ctx)
		if err != nil {
			return nil, err
		}
		for _, a := range list {
			if a.ID == in.ID {
				in.RewardID = a.SigninRewardID
				break
			}
		}
	}
	if in.RewardID <= 0 {
		return nil, fmt.Errorf("reward_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeSignin(in.ID, in.RewardID))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 104, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimProgress(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeProgress(in.ID, in.Step))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 126, 2)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ShopBuy(ctx context.Context, in game.ActivityShopIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	if in.GoodsID <= 0 {
		return nil, fmt.Errorf("goods_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeShopBuy(in.ID, in.GoodsID, in.Count))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 101, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimMega(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeMegaClaim(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 120, 2)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) TechSubmit(ctx context.Context, in game.TechIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	if in.NodeID <= 0 {
		return nil, fmt.Errorf("node_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeTechSubmit(in.ID, in.NodeID))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 140, 2)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) Lottery(ctx context.Context, in game.LotteryIn) (game.LotteryOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LotteryOut{}, err
	}
	if in.ID <= 0 {
		return game.LotteryOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeLottery(in.ID, in.HostGID, in.Free, in.Paid))
	if err != nil {
		return game.LotteryOut{}, err
	}
	out, err := decodeLotteryReply(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) BrewStart(ctx context.Context, in game.BrewStartIn) (game.BrewStartOut, error) {
	if _, err := c.Info(); err != nil {
		return game.BrewStartOut{}, err
	}
	if in.ID <= 0 {
		return game.BrewStartOut{}, fmt.Errorf("id 不能为空")
	}
	if len(in.Items) == 0 {
		return game.BrewStartOut{}, fmt.Errorf("items 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeBrewStart(in.ID, in.Items))
	if err != nil {
		return game.BrewStartOut{}, err
	}
	return decodeBrewStart(body)
}

func (c *Client) BrewStep(ctx context.Context, in game.ActivityOpIn) (game.BrewStepOut, error) {
	if _, err := c.Info(); err != nil {
		return game.BrewStepOut{}, err
	}
	if in.ID <= 0 {
		return game.BrewStepOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeBrewStep(in.ID))
	if err != nil {
		return game.BrewStepOut{}, err
	}
	return decodeBrewStep(body)
}

func (c *Client) BrewClaim(ctx context.Context, in game.BrewClaimIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeBrewClaim(in.ID, in.ClaimType))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 115, 3)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimRecall(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeRecallClaim(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 118, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimReturn(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeReturnGift(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 119, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimInvite(ctx context.Context, in game.InviteIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	if in.RewardType <= 0 {
		return nil, fmt.Errorf("reward_type 不能为空，1 邀请 2 成长")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeInviteClaim(in.ID, in.RewardType))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 121, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimNewcomer(ctx context.Context, in game.ActivityOpIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeNewcomerClaim(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeOperateItems(body, 122, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) SendGift(ctx context.Context, in game.GiftIn) (int64, error) {
	if _, err := c.Info(); err != nil {
		return 0, err
	}
	if in.ID <= 0 {
		return 0, fmt.Errorf("id 不能为空")
	}
	if in.GID <= 0 {
		return 0, fmt.Errorf("gid 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeSendGift(in.ID, in.GID, in.MsgTextID))
	if err != nil {
		return 0, err
	}
	return decodeSendCount(body)
}

func (c *Client) Mystery(ctx context.Context) (game.MysteryShop, error) {
	if _, err := c.Info(); err != nil {
		return game.MysteryShop{}, err
	}
	body, err := c.call(ctx, "MysteryShop", "GetActiveNPC", nil)
	if err != nil {
		return game.MysteryShop{}, err
	}
	return decodeMysteryShop(body)
}

func (c *Client) MysteryBuy(ctx context.Context, in game.IDIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "MysteryShop", "Buy", encodeLandID(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) MysteryLeave(ctx context.Context) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	_, err := c.call(ctx, "MysteryShop", "Abandon", nil)
	return err
}

func (c *Client) Dog(ctx context.Context) (game.DogYard, error) {
	if _, err := c.Info(); err != nil {
		return game.DogYard{}, err
	}
	body, err := c.call(ctx, "Dog", "GetDogInfo", nil)
	if err != nil {
		return game.DogYard{}, err
	}
	return decodeDogYard(body)
}

func (c *Client) Feed(ctx context.Context, in game.FeedIn) (int64, error) {
	if _, err := c.Info(); err != nil {
		return 0, err
	}
	if in.FoodID <= 0 {
		return 0, fmt.Errorf("food_id 不能为空")
	}
	body, err := c.call(ctx, "Dog", "AddFood", encodeFeed(in.FoodID, in.Count))
	if err != nil {
		return 0, err
	}
	return decodeFoodLeft(body)
}

func (c *Client) ClaimDogGifts(ctx context.Context) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Dog", "ClaimSkillGifts", nil)
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) DogLogs(ctx context.Context, in game.PageIn) ([]game.ProtectLog, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Dog", "GetProtectLogs", encodePage(in.From, in.Count, 50))
	if err != nil {
		return nil, err
	}
	return decodeProtectLogs(body)
}

func (c *Client) DeployDog(ctx context.Context, in game.IDIn) (game.DeployOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DeployOut{}, err
	}
	if in.ID <= 0 {
		return game.DeployOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Dog", "DeployDog", encodeLandID(in.ID))
	if err != nil {
		return game.DeployOut{}, err
	}
	return decodeDeploy(body)
}

func (c *Client) WithdrawDog(ctx context.Context) (game.DeployOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DeployOut{}, err
	}
	body, err := c.call(ctx, "Dog", "WithdrawDog", nil)
	if err != nil {
		return game.DeployOut{}, err
	}
	return decodeWithdraw(body)
}

func (c *Client) ActivateDog(ctx context.Context, in game.IDIn) (game.Dog, error) {
	if _, err := c.Info(); err != nil {
		return game.Dog{}, err
	}
	if in.ID <= 0 {
		return game.Dog{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Dog", "ActivateDog", encodeLandID(in.ID))
	if err != nil {
		return game.Dog{}, err
	}
	return decodeActivateDog(body, in.ID)
}

func (c *Client) Bulletins(ctx context.Context, in game.PageIn) ([]game.Bulletin, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "BulletinBoard", "GetBulletinList", encodePage(in.From, in.Count, 20))
	if err != nil {
		return nil, err
	}
	return decodeBulletins(body)
}

func (c *Client) ReadBulletin(ctx context.Context, in game.IDIn) (game.BulletinDetail, error) {
	if _, err := c.Info(); err != nil {
		return game.BulletinDetail{}, err
	}
	if in.ID <= 0 {
		return game.BulletinDetail{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "BulletinBoard", "GetBulletinDetail", encodeLandID(in.ID))
	if err != nil {
		return game.BulletinDetail{}, err
	}
	return decodeBulletinDetail(body)
}

func (c *Client) Mutants(ctx context.Context) ([]game.Mutant, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Mutant", "ReadMutantBook", nil)
	if err != nil {
		return nil, err
	}
	return decodeMutants(body)
}

func (c *Client) Career(ctx context.Context) (game.Career, error) {
	if _, err := c.Info(); err != nil {
		return game.Career{}, err
	}
	body, err := c.call(ctx, "Career", "CareerInfoGet", nil)
	if err != nil {
		return game.Career{}, err
	}
	return decodeCareer(body)
}

func (c *Client) Ranks(ctx context.Context, in game.RankIn) (game.RankBoard, error) {
	if _, err := c.Info(); err != nil {
		return game.RankBoard{}, err
	}
	body, err := c.call(ctx, "Rank", "GetRankList", encodeRank(in.Type, in.Page))
	if err != nil {
		return game.RankBoard{}, err
	}
	return decodeRankBoard(body)
}

func (c *Client) Avatars(ctx context.Context, in game.TypeIn) ([]game.Avatar, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "AvatarFrame", "AvatarFramesOwned", encodeAvatarOwned(in.Type))
	if err != nil {
		return nil, err
	}
	return decodeAvatars(body)
}

func (c *Client) EquippedAvatars(ctx context.Context) ([]game.Avatar, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "AvatarFrame", "AvatarFramesEquiped", nil)
	if err != nil {
		return nil, err
	}
	return decodeAvatars(body)
}

func (c *Client) EquipAvatar(ctx context.Context, in game.AvatarEquipIn) (game.Avatar, error) {
	if _, err := c.Info(); err != nil {
		return game.Avatar{}, err
	}
	if !in.Off && in.ID <= 0 {
		return game.Avatar{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "AvatarFrame", "UpdateEquip", encodeAvatarEquip(in.ID, in.Off))
	if err != nil {
		return game.Avatar{}, err
	}
	return decodeEquippedAvatar(body)
}

func (c *Client) Skins(ctx context.Context) ([]game.Skin, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Skin", "SkinsOwned", nil)
	if err != nil {
		return nil, err
	}
	return decodeSkins(body)
}

func (c *Client) EquippedSkins(ctx context.Context) ([]game.Skin, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Skin", "SkinsEquipped", nil)
	if err != nil {
		return nil, err
	}
	return decodeSkins(body)
}

func (c *Client) EquipSkin(ctx context.Context, in game.SkinEquipIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	_, err := c.call(ctx, "Skin", "UpdateEquip", encodeSkinEquip(in.Current, in.ID))
	return err
}

func (c *Client) Drops(ctx context.Context) ([]game.Drop, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "RandomDrop", "GetActivityInfo", nil)
	if err != nil {
		return nil, err
	}
	return decodeDrops(body)
}

func (c *Client) SolarTerms(ctx context.Context) ([]game.SolarTerm, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "SolarTerms", "GetSolarTerms", nil)
	if err != nil {
		return nil, err
	}
	return decodeSolarTerms(body)
}

func (c *Client) ClaimSolar(ctx context.Context, in game.IDIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "SolarTerms", "ClaimSolarTerms", encodeLandID(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimAllSolar(ctx context.Context) ([]game.Item, error) {
	list, err := c.SolarTerms(ctx)
	if err != nil {
		return nil, err
	}
	var out []game.Item
	for _, s := range list {
		if s.Status != 2 || s.ID <= 0 {
			continue
		}
		items, err := c.ClaimSolar(ctx, game.IDIn{ID: s.ID})
		if err != nil {
			return out, err
		}
		out = append(out, items...)
	}
	return out, nil
}

func (c *Client) AchieveView(ctx context.Context, in game.AchieveIn) (game.AchieveScope, error) {
	if _, err := c.Info(); err != nil {
		return game.AchieveScope{}, err
	}
	body, err := c.call(ctx, "Achieve", "GetScopeView", encodeAchieve(in.Kind, in.ID))
	if err != nil {
		return game.AchieveScope{}, err
	}
	return decodeAchieveScope(body)
}

func (c *Client) ClaimAchieveGoal(ctx context.Context, in game.AchieveGoalIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.GoalID <= 0 {
		return nil, fmt.Errorf("goal_id 不能为空")
	}
	body, err := c.call(ctx, "Achieve", "ClaimGoalReward", encodeAchieveGoal(in.Kind, in.ID, in.GoalID))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimAchieveLevel(ctx context.Context, in game.AchieveIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Achieve", "ClaimScopeLevel", encodeAchieve(in.Kind, in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) Season(ctx context.Context) (game.Season, error) {
	if _, err := c.Info(); err != nil {
		return game.Season{}, err
	}
	body, err := c.call(ctx, "Season", "GetSeasonInfo", encodeSeasonInfo())
	if err != nil {
		return game.Season{}, err
	}
	return decodeSeason(body)
}

func (c *Client) ClaimPass(ctx context.Context) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Season", "ClaimBattlePassRewards", nil)
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) Mall(ctx context.Context, in game.MallIn) ([]game.Product, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.Slot <= 0 {
		in.Slot = 1
	}
	body, err := c.call(ctx, "Mall", "GetMallListBySlotType", encodeMallList(in.Slot))
	if err != nil {
		return nil, err
	}
	return decodeProducts(body)
}

func (c *Client) MallBuy(ctx context.Context, in game.MallBuyIn) (game.BuyOut, error) {
	var out game.BuyOut
	if _, err := c.Info(); err != nil {
		return out, err
	}
	if in.ID <= 0 {
		return out, fmt.Errorf("id 不能为空")
	}
	if in.Num <= 0 {
		in.Num = 1
	}
	body, err := c.call(ctx, "Mall", "Purchase", encodeMallBuy(in.ID, in.Num))
	if err != nil {
		return out, err
	}
	got, err := decodeItemsAt(body, 2)
	if err != nil {
		return out, err
	}
	cost, err := decodeItemsAt(body, 3)
	if err != nil {
		return out, err
	}
	out.Items = got
	out.Cost = cost
	c.applyCurrency(got, cost)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) MonthCards(ctx context.Context) ([]game.MonthCard, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Mall", "GetMonthCardInfos", nil)
	if err != nil {
		return nil, err
	}
	return decodeMonthCards(body)
}

func (c *Client) ClaimMonthCard(ctx context.Context, in game.IDIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Mall", "ClaimMonthCardReward", encodeLandID(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) RedPackets(ctx context.Context) ([]game.RedPacket, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "RedPacket", "GetTodayClaimStatus", nil)
	if err != nil {
		return nil, err
	}
	return decodeRedPackets(body)
}

func (c *Client) ClaimRedPacket(ctx context.Context, in game.IDIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "RedPacket", "ClaimRedPacket", encodeLandID(in.ID))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 2)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) ClaimAllRedPackets(ctx context.Context) ([]game.Item, error) {
	list, err := c.RedPackets(ctx)
	if err != nil {
		return nil, err
	}
	var all []game.Item
	for _, p := range list {
		if !p.CanClaim {
			continue
		}
		items, err := c.ClaimRedPacket(ctx, game.IDIn{ID: p.ID})
		if err != nil {
			continue
		}
		all = append(all, items...)
	}
	return all, nil
}

func (c *Client) VisitLogs(ctx context.Context) ([]game.VisitLog, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Interact", "InteractRecords", nil)
	if err != nil {
		return nil, err
	}
	return decodeVisitLogs(body)
}

func (c *Client) Album(ctx context.Context, in game.AlbumIn) (game.Album, error) {
	if _, err := c.Info(); err != nil {
		return game.Album{}, err
	}
	body, err := c.call(ctx, "Illustrated", "GetIllustratedListV2", encodeIllustratedList(in.Rarity, in.Type))
	if err != nil {
		return game.Album{}, err
	}
	return decodeAlbum(body)
}

func (c *Client) ClaimAlbum(ctx context.Context, in game.AlbumIn) ([]game.Item, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.Type <= 0 {
		in.Type = 1
	}
	body, err := c.call(ctx, "Illustrated", "ClaimAllRewardsV2", encodeLandID(in.Type))
	if err != nil {
		return nil, err
	}
	items, err := decodeItemsAt(body, 1)
	if err != nil {
		return nil, err
	}
	_, _ = c.Bag(ctx)
	return items, nil
}

func (c *Client) applyCurrency(got, cost []game.Item) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, it := range got {
		if it.ID == 1001 {
			c.user.Gold += it.Count
		}
	}
	for _, it := range cost {
		if it.ID == 1001 {
			c.user.Gold -= it.Count
		}
	}
}

func (c *Client) Close() error {
	err := c.closeConn()
	if c.enc != nil {
		if e := c.enc.Close(); e != nil && err == nil {
			err = e
		}
	}
	return err
}

func (c *Client) closeConn() error {
	c.mu.Lock()
	life := c.life
	c.life = nil
	c.closed = true
	conn := c.conn
	c.conn = nil
	c.loggedIn = false
	c.bag = nil
	for seq, ch := range c.pending {
		close(ch)
		delete(c.pending, seq)
	}
	c.mu.Unlock()
	if life != nil {
		life.Stop()
	}
	if conn != nil {
		return conn.Close(websocket.StatusNormalClosure, "")
	}
	return nil
}

func (c *Client) call(ctx context.Context, service, method string, body []byte) ([]byte, error) {
	if body == nil {
		body = []byte{}
	}
	sealed, token, err := c.enc.Seal(body)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	if c.conn == nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("未连接网关")
	}
	c.seq++
	seq := c.seq
	serverSeq := c.serverSeq
	ch := make(chan gate.Message, 1)
	c.pending[seq] = ch
	conn := c.conn
	raw := gate.EncodeRequest(service, method, sealed, token, seq, serverSeq)
	c.mu.Unlock()

	if err := conn.Write(ctx, websocket.MessageBinary, raw); err != nil {
		c.mu.Lock()
		delete(c.pending, seq)
		c.mu.Unlock()
		return nil, err
	}
	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return nil, fmt.Errorf("RPC %s.%s 超时", service, method)
	case msg, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("连接已关闭")
		}
		if msg.ErrorCode != 0 {
			return nil, fmt.Errorf("RPC %s: code=%d msg=%s", method, msg.ErrorCode, msg.ErrorMessage)
		}
		return msg.Body, nil
	}
}

func (c *Client) uploadAnti(data []byte) ([]byte, error) {
	body, err := c.call(context.Background(), "Ace", "AntiData", encodeAnti(data))
	if err != nil {
		return nil, err
	}
	return decodeAnti(body)
}

func (c *Client) readLoop() {
	ctx := context.Background()
	for {
		c.mu.Lock()
		conn := c.conn
		closed := c.closed
		c.mu.Unlock()
		if conn == nil || closed {
			return
		}
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		msg, err := gate.Decode(data)
		if err != nil {
			continue
		}
		c.mu.Lock()
		if msg.ServerSeq > c.serverSeq {
			c.serverSeq = msg.ServerSeq
		}
		if msg.MessageType == gate.TypeResponse {
			if ch, ok := c.pending[msg.ClientSeq]; ok {
				delete(c.pending, msg.ClientSeq)
				c.mu.Unlock()
				ch <- msg
				continue
			}
		}
		c.mu.Unlock()
	}
}
