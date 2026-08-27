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

const callTimeout = 15 * time.Second

const frogPrankBottle int64 = 5005

type Client struct {
	enc     crypto.Encryptor
	openErr error

	mu         sync.Mutex
	conn       *websocket.Conn
	pending    map[int64]chan gate.Message
	seq        int64
	serverSeq  int64
	user       game.User
	bag        []game.BagItem
	loggedIn   bool
	closed     bool
	life       *ace.Life
	hostGID    int64
	hbStop     chan struct{}
	hbCount    int
	serverMs   int64
	hbErr      string
	firstLogin bool
	weather    game.Weather
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
	first, weather := decodeLoginExtra(body)
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
	c.hostGID = user.GID
	c.firstLogin = first
	c.weather = weather
	c.life = ace.New(crypto.AsEngine(c.enc), c.uploadAnti)
	c.life.Start()
	c.mu.Unlock()
	c.startGateHeartbeat()
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
	hbErr := c.hbErr
	out.HostGID = c.hostGID
	out.Heartbeats = c.hbCount
	out.ServerMs = c.serverMs
	out.FirstLogin = c.firstLogin
	out.Weather = c.weather
	c.mu.Unlock()
	if life != nil {
		snap := life.Snapshot()
		out.ACE = game.ACE{Uploads: snap.Uploads, Reports: snap.Reports, Failures: snap.Failures, LastError: snap.LastError}
	}
	if out.ACE.LastError == "" {
		out.ACE.LastError = hbErr
	}
	return out, nil
}

func (c *Client) Heartbeat(ctx context.Context, in game.RefreshIn) (game.HeartbeatOut, error) {
	user, err := c.Info()
	if err != nil {
		return game.HeartbeatOut{}, err
	}
	host := in.HostGID
	if host == 0 {
		c.mu.Lock()
		host = c.hostGID
		c.mu.Unlock()
	}
	if host == 0 {
		host = user.GID
	}
	body, err := c.call(ctx, "User", "Heartbeat", encodeHeartbeat(host, gameVersion()))
	if err != nil {
		c.mu.Lock()
		c.hbErr = err.Error()
		c.mu.Unlock()
		return game.HeartbeatOut{}, err
	}
	out, err := decodeHeartbeat(body)
	if err != nil {
		return game.HeartbeatOut{}, err
	}
	out.HostGID = host
	c.mu.Lock()
	c.hbCount++
	c.serverMs = out.ServerMs
	c.hbErr = ""
	c.mu.Unlock()
	return out, nil
}

func (c *Client) Logout(ctx context.Context) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	_, _ = c.call(ctx, "User", "Logout", encodeLogout(""))
	return c.Close()
}

func (c *Client) Brief(ctx context.Context, in game.EnterIn) (game.User, error) {
	if _, err := c.Info(); err != nil {
		return game.User{}, err
	}
	if in.GID <= 0 {
		return game.User{}, fmt.Errorf("gid 不能为空")
	}
	body, err := c.call(ctx, "User", "GetBriefInfo", encodeLeave(in.GID))
	if err != nil {
		return game.User{}, err
	}
	return decodeBrief(body)
}

func (c *Client) BatchInfo(ctx context.Context, in game.GIDsIn) ([]game.User, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if len(in.GIDs) == 0 {
		return nil, fmt.Errorf("gids 不能为空")
	}
	body, err := c.call(ctx, "User", "BatchGetBasicInfo", encodeGIDs(in.GIDs))
	if err != nil {
		return nil, err
	}
	return decodeUsersAt(body, 1)
}

func (c *Client) ArkClick(ctx context.Context, in game.ArkIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 && in.OpenID == "" {
		return fmt.Errorf("gid 或 open_id 不能为空")
	}
	_, err := c.call(ctx, "User", "ReportArkClick", encodeArkClick(in))
	return err
}

func (c *Client) Refresh(ctx context.Context, in game.RefreshIn) (game.LandOpOut, error) {
	user, err := c.Info()
	if err != nil {
		return game.LandOpOut{}, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	c.setHost(host)
	body, err := c.call(ctx, "Plant", "AllLands", encodeAllLands(host))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeRefreshOut(body)
}

func (c *Client) RefreshLands(ctx context.Context, in game.RefreshLandsIn) ([]game.Land, error) {
	user, err := c.Info()
	if err != nil {
		return nil, err
	}
	if len(in.LandIDs) == 0 {
		return nil, nil
	}
	host := in.HostGID
	if host == user.GID {
		host = 0
	}
	if in.HostGID > 0 {
		c.setHost(in.HostGID)
	}
	body, err := c.call(ctx, "Plant", "RefreshLands", encodeRefreshLands(in.LandIDs, host, visitReason(user.GID, in.HostGID)))
	if err != nil {
		return nil, err
	}
	return decodeLands(body)
}

func (c *Client) CleanSocial(ctx context.Context, in game.CleanSocialIn) (game.LandOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LandOpOut{}, err
	}
	if len(in.LandIDs) == 0 {
		return game.LandOpOut{}, fmt.Errorf("land_ids 不能为空")
	}
	body, err := c.call(ctx, "Plant", "CleanSocialItems", encodeCleanSocial(in.LandIDs, in.ItemIDs))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeSocialOp(body)
}

func (c *Client) PutInsects(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	return c.putSocial(ctx, "PutInsects", in)
}

func (c *Client) PutWeeds(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	return c.putSocial(ctx, "PutWeeds", in)
}

func (c *Client) putSocial(ctx context.Context, method string, in game.LandOpIn) (game.LandOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LandOpOut{}, err
	}
	if in.HostGID <= 0 {
		return game.LandOpOut{}, fmt.Errorf("host_gid 不能为空")
	}
	if len(in.LandIDs) == 0 {
		return game.LandOpOut{}, fmt.Errorf("land_ids 不能为空")
	}
	body, err := c.call(ctx, "Plant", method, encodePutSocial(in.HostGID, in.LandIDs, 2))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeLandOp(body)
}

func (c *Client) PutSocial(ctx context.Context, in game.PutSocialIn) (game.LandOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LandOpOut{}, err
	}
	if in.HostGID <= 0 {
		return game.LandOpOut{}, fmt.Errorf("host_gid 不能为空")
	}
	if in.ItemID <= 0 {
		return game.LandOpOut{}, fmt.Errorf("item_id 不能为空")
	}
	if in.ItemID != frogPrankBottle && len(in.LandIDs) == 0 {
		return game.LandOpOut{}, fmt.Errorf("地块道具必须带 land_ids；青蛙 5005 不带地")
	}
	use := game.UseIn{ID: in.ItemID, Count: 1, HostGID: in.HostGID, LandIDs: in.LandIDs}
	if in.ItemID == frogPrankBottle {
		use.LandIDs = nil
	}
	if uid := c.bagUID(in.ItemID); uid > 0 {
		use.UID = uid
	} else {
		_, _ = c.Bag(ctx)
		use.UID = c.bagUID(in.ItemID)
	}
	body, err := c.call(ctx, "Item", "Use", encodeUse(use))
	if err != nil {
		return game.LandOpOut{}, err
	}
	out, err := decodeUseLandOp(body)
	if err != nil {
		return game.LandOpOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) bagUID(itemID int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, it := range c.bag {
		if it.ID == itemID && it.Count > 0 && it.UID > 0 {
			return it.UID
		}
	}
	return 0
}

func (c *Client) CanOperate(ctx context.Context, in game.CanOperateIn) (game.CanOperateOut, error) {
	if _, err := c.Info(); err != nil {
		return game.CanOperateOut{}, err
	}
	if in.HostGID <= 0 || in.OperationID <= 0 {
		return game.CanOperateOut{}, fmt.Errorf("host_gid 和 operation_id 不能为空")
	}
	body, err := c.call(ctx, "Plant", "CheckCanOperate", encodeCheckCanOperate(in.HostGID, in.OperationID))
	if err != nil {
		return game.CanOperateOut{}, err
	}
	return decodeCanOperate(body)
}

func (c *Client) Harvest(ctx context.Context, in game.HarvestIn) (game.HarvestOut, error) {
	user, err := c.Info()
	if err != nil {
		return game.HarvestOut{}, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	if host != user.GID {
		if err := c.ensureCanOperate(ctx, host, 10004); err != nil {
			return game.HarvestOut{}, err
		}
	}
	body, err := c.call(ctx, "Plant", "Harvest", encodeHarvest(in.LandIDs, host, in.IsAll, visitReason(user.GID, host)))
	if err != nil {
		return game.HarvestOut{}, err
	}
	return decodeHarvest(body)
}

func (c *Client) Steal(ctx context.Context, in game.HarvestIn) (game.HarvestOut, error) {
	if in.HostGID <= 0 {
		return game.HarvestOut{}, fmt.Errorf("host_gid 不能为空")
	}
	return c.Harvest(ctx, in)
}

func (c *Client) Remove(ctx context.Context, in game.RemoveIn) (game.LandOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LandOpOut{}, err
	}
	if len(in.LandIDs) == 0 {
		return game.LandOpOut{}, fmt.Errorf("land_ids 不能为空")
	}
	body, err := c.call(ctx, "Plant", "RemovePlant", encodeRemove(in.LandIDs))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeLandOp(body)
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

func (c *Client) Farming(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	user, err := c.Info()
	if err != nil {
		return game.LandOpOut{}, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	body, err := c.call(ctx, "Plant", "Farming", encodeFarming(in.LandIDs, host, in.ItemIDs))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeFarmingOut(body)
}

func (c *Client) ensureCanOperate(ctx context.Context, hostGID, operationID int64) error {
	body, err := c.call(ctx, "Plant", "CheckCanOperate", encodeCheckCanOperate(hostGID, operationID))
	if err != nil {
		return err
	}
	out, err := decodeCanOperate(body)
	if err != nil {
		return err
	}
	if !out.OK {
		return fmt.Errorf("不能操作")
	}
	return nil
}

func (c *Client) Plant(ctx context.Context, in game.PlantIn) (game.LandOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LandOpOut{}, err
	}
	if in.SeedID <= 0 || len(in.LandIDs) == 0 {
		return game.LandOpOut{}, fmt.Errorf("seed_id 和 land_ids 不能为空")
	}
	body, err := c.call(ctx, "Plant", "Plant", encodePlant(in.SeedID, in.LandIDs))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeLandOp(body)
}

func (c *Client) Water(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	return c.landOp(ctx, "WaterLand", in)
}

func (c *Client) Weed(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	return c.landOp(ctx, "WeedOut", in)
}

func (c *Client) Bug(ctx context.Context, in game.LandOpIn) (game.LandOpOut, error) {
	return c.landOp(ctx, "Insecticide", in)
}

func (c *Client) Fertilize(ctx context.Context, in game.FertilizeIn) (game.LandOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LandOpOut{}, err
	}
	if in.FertilizerID <= 0 || len(in.LandIDs) == 0 {
		return game.LandOpOut{}, fmt.Errorf("fertilizer_id 和 land_ids 不能为空")
	}
	body, err := c.call(ctx, "Plant", "Fertilize", encodeFertilize(in.LandIDs, in.FertilizerID))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeFertilizeOut(body)
}

func (c *Client) landOp(ctx context.Context, method string, in game.LandOpIn) (game.LandOpOut, error) {
	user, err := c.Info()
	if err != nil {
		return game.LandOpOut{}, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	body, err := c.call(ctx, "Plant", method, encodeLandOp(in.LandIDs, host, visitReason(user.GID, host)))
	if err != nil {
		return game.LandOpOut{}, err
	}
	return decodeLandOp(body)
}

func (c *Client) Friends(ctx context.Context) (game.FriendsOut, error) {
	if _, err := c.Info(); err != nil {
		return game.FriendsOut{}, err
	}
	body, err := c.call(ctx, "Friend", "GetAll", nil)
	if err != nil {
		return game.FriendsOut{}, err
	}
	return decodeFriendsOut(body)
}

func (c *Client) Help(ctx context.Context, in game.HelpIn) error {
	user, err := c.Info()
	if err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	if _, err := c.call(ctx, "Visit", "Enter", encodeEnter(in.GID, visitReason(user.GID, in.GID))); err != nil {
		return err
	}
	c.setHost(in.GID)
	_, waterErr := c.call(ctx, "Plant", "WaterLand", encodeWater(in.LandIDs, in.GID, 2))
	leaveErr := c.Leave(ctx, game.EnterIn{GID: in.GID})
	if user, infoErr := c.Info(); infoErr == nil {
		c.setHost(user.GID)
	}
	if waterErr != nil {
		return waterErr
	}
	return leaveErr
}

func (c *Client) Enter(ctx context.Context, in game.EnterIn) (game.Visit, error) {
	user, err := c.Info()
	if err != nil {
		return game.Visit{}, err
	}
	if in.GID <= 0 {
		return game.Visit{}, fmt.Errorf("gid 不能为空")
	}
	body, err := c.call(ctx, "Visit", "Enter", encodeEnter(in.GID, visitReason(user.GID, in.GID)))
	if err != nil {
		return game.Visit{}, err
	}
	c.setHost(in.GID)
	return decodeVisit(body)
}

func (c *Client) Applications(ctx context.Context) (game.ApplicationsOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ApplicationsOut{}, err
	}
	body, err := c.call(ctx, "Friend", "GetApplications", nil)
	if err != nil {
		return game.ApplicationsOut{}, err
	}
	return decodeApplications(body)
}

func (c *Client) Accept(ctx context.Context, in game.GIDsIn) (game.AcceptOut, error) {
	if _, err := c.Info(); err != nil {
		return game.AcceptOut{}, err
	}
	if len(in.GIDs) == 0 {
		return game.AcceptOut{}, fmt.Errorf("gids 不能为空")
	}
	body, err := c.call(ctx, "Friend", "AcceptFriends", encodeGIDs(in.GIDs))
	if err != nil {
		return game.AcceptOut{}, err
	}
	return decodeAccept(body)
}

func (c *Client) Reject(ctx context.Context, in game.GIDsIn) (game.RejectOut, error) {
	if _, err := c.Info(); err != nil {
		return game.RejectOut{}, err
	}
	if len(in.GIDs) == 0 {
		return game.RejectOut{}, fmt.Errorf("gids 不能为空")
	}
	body, err := c.call(ctx, "Friend", "RejectFriends", encodeGIDs(in.GIDs))
	if err != nil {
		return game.RejectOut{}, err
	}
	return decodeReject(body)
}

func (c *Client) SetTags(ctx context.Context, in game.TagsIn) (game.Friend, error) {
	if _, err := c.Info(); err != nil {
		return game.Friend{}, err
	}
	if in.GID <= 0 {
		return game.Friend{}, fmt.Errorf("gid 不能为空")
	}
	body, err := c.call(ctx, "Friend", "SetTags", encodeSetTags())
	if err != nil {
		return game.Friend{}, err
	}
	return decodeSetTags(body)
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

func (c *Client) SyncFriends(ctx context.Context, in game.OpenIDsIn) ([]game.Friend, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if len(in.OpenIDs) == 0 {
		return nil, fmt.Errorf("open_ids 不能为空")
	}
	body, err := c.call(ctx, "Friend", "SyncAll", encodeOpenIDs(2, in.OpenIDs))
	if err != nil {
		return nil, err
	}
	return decodeFriends(body)
}

func (c *Client) GameFriends(ctx context.Context, in game.GIDsIn) ([]game.Friend, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if len(in.GIDs) == 0 {
		return nil, fmt.Errorf("gids 不能为空")
	}
	body, err := c.call(ctx, "Friend", "GetGameFriends", encodeGIDsPacked(in.GIDs))
	if err != nil {
		return nil, err
	}
	return decodeFriends(body)
}

func (c *Client) BlockFriend(ctx context.Context, in game.EnterIn) (game.Blocked, error) {
	if _, err := c.Info(); err != nil {
		return game.Blocked{}, err
	}
	if in.GID <= 0 {
		return game.Blocked{}, fmt.Errorf("gid 不能为空")
	}
	body, err := c.call(ctx, "Friend", "BlockFriend", encodeLeave(in.GID))
	if err != nil {
		return game.Blocked{}, err
	}
	out, err := decodeBlockReply(body)
	if err != nil {
		return game.Blocked{}, err
	}
	if out.GID == 0 {
		out.GID = in.GID
	}
	return out, nil
}

func (c *Client) UnblockFriend(ctx context.Context, in game.EnterIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	_, err := c.call(ctx, "Friend", "UnblockFriend", encodeLeave(in.GID))
	return err
}

func (c *Client) BlockList(ctx context.Context) ([]game.Blocked, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Friend", "GetBlockList", nil)
	if err != nil {
		return nil, err
	}
	return decodeBlockedList(body)
}

func (c *Client) ShareKey(ctx context.Context) (string, error) {
	if _, err := c.Info(); err != nil {
		return "", err
	}
	body, err := c.call(ctx, "Friend", "GetShareKey", nil)
	if err != nil {
		return "", err
	}
	return decodeShareKey(body)
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

func (c *Client) InviteInfo(ctx context.Context, in game.IDIn) (game.ShareInviteOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ShareInviteOut{}, err
	}
	var req []byte
	if in.ID > 0 {
		req = encodeLandID(in.ID)
	}
	body, err := c.call(ctx, "Share", "GetInviteInfo", req)
	if err != nil {
		return game.ShareInviteOut{}, err
	}
	return decodeShareInvite(body)
}

func (c *Client) InviteAward(ctx context.Context, in game.IDIn) (game.ShareAwardOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ShareAwardOut{}, err
	}
	if in.ID <= 0 {
		return game.ShareAwardOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Share", "GetInviteAward", encodeLandID(in.ID))
	if err != nil {
		return game.ShareAwardOut{}, err
	}
	return decodeShareAwardOut(body)
}

func (c *Client) PosterShown(ctx context.Context, in game.IDIn) (bool, error) {
	if _, err := c.Info(); err != nil {
		return false, err
	}
	if in.ID <= 0 {
		return false, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Share", "SetPosterShown", encodeLandID(in.ID))
	if err != nil {
		return false, err
	}
	return decodeOK(body)
}

func (c *Client) ShareClaim(ctx context.Context) (game.ShareClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ShareClaimOut{}, err
	}
	if _, err := c.call(ctx, "Share", "ReportShare", encodeShareFlag()); err != nil {
		return game.ShareClaimOut{}, err
	}
	body, err := c.call(ctx, "Share", "ClaimShareReward", encodeShareFlag())
	if err != nil {
		return game.ShareClaimOut{}, err
	}
	out, err := decodeShareClaimOut(body)
	if err != nil {
		return game.ShareClaimOut{}, err
	}
	c.applyCurrency(out.Items, nil)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ReportInvite(ctx context.Context, in game.InviteReportIn) (bool, error) {
	if _, err := c.Info(); err != nil {
		return false, err
	}
	if in.OpenID == "" || in.ShareKey == "" {
		return false, fmt.Errorf("open_id 和 share_key 不能为空")
	}
	body, err := c.call(ctx, "Share", "ReportInvitation", encodeReportInvite(in.OpenID, in.ShareKey))
	if err != nil {
		return false, err
	}
	return decodeOK(body)
}

func (c *Client) Leave(ctx context.Context, in game.EnterIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	_, err := c.call(ctx, "Visit", "Leave", encodeLeave(in.GID))
	if err == nil {
		if user, infoErr := c.Info(); infoErr == nil {
			c.setHost(user.GID)
		}
	}
	return err
}

func (c *Client) Bag(ctx context.Context) (game.BagOut, error) {
	if _, err := c.Info(); err != nil {
		return game.BagOut{}, err
	}
	body, err := c.call(ctx, "Item", "Bag", nil)
	if err != nil {
		return game.BagOut{}, err
	}
	out, err := decodeBagOut(body)
	if err != nil {
		return game.BagOut{}, err
	}
	c.mu.Lock()
	c.bag = out.Items
	c.mu.Unlock()
	return out, nil
}

func (c *Client) Sell(ctx context.Context, in game.SellIn) (game.ItemOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ItemOpOut{}, err
	}
	if len(in.Items) == 0 {
		return game.ItemOpOut{}, fmt.Errorf("items 不能为空")
	}
	for _, it := range in.Items {
		if it.ID <= 0 || it.Count <= 0 {
			return game.ItemOpOut{}, fmt.Errorf("每个 item 需要 id 和 count")
		}
	}
	body, err := c.call(ctx, "Item", "Sell", encodeSell(in.Items))
	if err != nil {
		return game.ItemOpOut{}, err
	}
	out, err := decodeItemOp(body, 1, 2, 0, false)
	if err != nil {
		return game.ItemOpOut{}, err
	}
	c.applyCurrency(out.Items, nil)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) Use(ctx context.Context, in game.UseIn) (game.ItemOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ItemOpOut{}, err
	}
	if in.Count <= 0 {
		in.Count = 1
	}
	if in.ID <= 0 {
		return game.ItemOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Item", "Use", encodeUse(in))
	if err != nil {
		return game.ItemOpOut{}, err
	}
	out, err := decodeItemOp(body, 1, 2, 3, true)
	if err != nil {
		return game.ItemOpOut{}, err
	}
	c.applyCurrency(out.Items, nil)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) BatchUse(ctx context.Context, in game.BatchUseIn) (game.ItemOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ItemOpOut{}, err
	}
	if len(in.Items) == 0 {
		return game.ItemOpOut{}, fmt.Errorf("items 不能为空")
	}
	for i := range in.Items {
		if in.Items[i].Count <= 0 {
			in.Items[i].Count = 1
		}
		if in.Items[i].ID <= 0 {
			return game.ItemOpOut{}, fmt.Errorf("id 不能为空")
		}
	}
	body, err := c.call(ctx, "Item", "BatchUse", encodeBatchUse(in.Items))
	if err != nil {
		return game.ItemOpOut{}, err
	}
	out, err := decodeItemOp(body, 1, 2, 3, false)
	if err != nil {
		return game.ItemOpOut{}, err
	}
	c.applyCurrency(out.Items, nil)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) CancelNew(ctx context.Context, in game.IDIn) (int64, error) {
	if _, err := c.Info(); err != nil {
		return 0, err
	}
	if in.ID <= 0 {
		return 0, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Item", "CannelNew", encodeLandID(in.ID))
	if err != nil {
		return 0, err
	}
	return decodeItemID(body)
}

func (c *Client) LockItems(ctx context.Context, in game.UIDsIn) (game.LockOut, error) {
	return c.lockItems(ctx, "LockItems", in)
}

func (c *Client) UnlockItems(ctx context.Context, in game.UIDsIn) (game.LockOut, error) {
	return c.lockItems(ctx, "UnlockItems", in)
}

func (c *Client) lockItems(ctx context.Context, method string, in game.UIDsIn) (game.LockOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LockOut{}, err
	}
	if len(in.UIDs) == 0 {
		return game.LockOut{}, fmt.Errorf("uids 不能为空")
	}
	body, err := c.call(ctx, "Item", method, encodeUIDs(in.UIDs))
	if err != nil {
		return game.LockOut{}, err
	}
	return decodeLockOut(body)
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
	out, err = decodeShopBuy(body)
	if err != nil {
		return out, err
	}
	c.applyCurrency(out.Items, out.Cost)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) AutoBuy(ctx context.Context, in game.AutoBuyIn) (game.BuyOut, error) {
	if in.ItemID <= 0 {
		return game.BuyOut{}, fmt.Errorf("item_id 不能为空")
	}
	if in.Num <= 0 {
		in.Num = 1
	}
	goods, err := c.Goods(ctx, game.ShopIn{ShopID: in.ShopID})
	if err != nil {
		return game.BuyOut{}, err
	}
	for _, g := range goods {
		if g.ItemID != in.ItemID {
			continue
		}
		if !g.Unlocked {
			return game.BuyOut{}, fmt.Errorf("商品未解锁")
		}
		if g.Limit > 0 {
			remain := g.Limit - g.Bought
			if remain <= 0 {
				return game.BuyOut{}, fmt.Errorf("已达限购 %d", g.Limit)
			}
			if in.Num > remain {
				in.Num = remain
			}
		}
		return c.Buy(ctx, game.BuyIn{GoodsID: g.ID, Num: in.Num, Price: g.Price})
	}
	return game.BuyOut{}, fmt.Errorf("货架上没有 item_id=%d", in.ItemID)
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

func (c *Client) CurrentWeather(ctx context.Context) (game.Weather, error) {
	if _, err := c.Info(); err != nil {
		return game.Weather{}, err
	}
	body, err := c.call(ctx, "Weather", "GetCurrentWeather", nil)
	if err != nil {
		return game.Weather{}, err
	}
	return decodeCurrentWeather(body)
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

func (c *Client) ClaimTask(ctx context.Context, in game.ClaimIn) (game.TaskClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.TaskClaimOut{}, err
	}
	ids := in.IDs
	if in.ID > 0 {
		ids = append([]int64{in.ID}, ids...)
	}
	if len(ids) == 0 {
		return game.TaskClaimOut{}, fmt.Errorf("id 不能为空")
	}
	var (
		body []byte
		err  error
	)
	if len(ids) == 1 {
		body, err = c.call(ctx, "Task", "ClaimTaskReward", encodeClaimTaskShared(ids[0], in.Shared))
	} else {
		body, err = c.call(ctx, "Task", "BatchClaimTaskReward", encodeBatchClaim(ids))
	}
	if err != nil {
		return game.TaskClaimOut{}, err
	}
	out, err := decodeTaskClaim(body)
	if err != nil {
		return game.TaskClaimOut{}, err
	}
	c.applyCurrency(out.Items, nil)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimDaily(ctx context.Context, in game.DailyIn) (game.TaskClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.TaskClaimOut{}, err
	}
	if len(in.PointIDs) == 0 {
		return game.TaskClaimOut{}, fmt.Errorf("point_ids 不能为空")
	}
	body, err := c.call(ctx, "Task", "ClaimDailyReward", encodeClaimDaily(in.Type, in.PointIDs))
	if err != nil {
		return game.TaskClaimOut{}, err
	}
	out, err := decodeTaskClaim(body)
	if err != nil {
		return game.TaskClaimOut{}, err
	}
	c.applyCurrency(out.Items, nil)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ReportTask(ctx context.Context, in game.TaskReportIn) (game.TaskBoard, error) {
	if _, err := c.Info(); err != nil {
		return game.TaskBoard{}, err
	}
	if in.ID <= 0 {
		return game.TaskBoard{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Task", "ClientReportProgress", encodeTaskReport(in.ID, in.Progress))
	if err != nil {
		return game.TaskBoard{}, err
	}
	return decodeTaskBoard(body)
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

func (c *Client) ClaimAllEmail(ctx context.Context, in game.EmailBoxIn) (game.EmailClaimOut, error) {
	var out game.EmailClaimOut
	if _, err := c.Info(); err != nil {
		return out, err
	}
	boxes := []int64{emailBox(in.Box)}
	if in.Box <= 0 {
		boxes = []int64{1, 2}
	}
	for _, box := range boxes {
		body, err := c.call(ctx, "Email", "BatchClaimEmail", encodeEmailRef(box, in.ID))
		if err != nil {
			continue
		}
		part, err := decodeEmailClaim(body)
		if err != nil {
			return out, err
		}
		out.Items = append(out.Items, part.Items...)
		out.Unclaimed = append(out.Unclaimed, part.Unclaimed...)
	}
	c.applyCurrency(out.Items, nil)
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) BatchReadEmail(ctx context.Context, in game.EmailIDsIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if len(in.IDs) == 0 {
		return fmt.Errorf("ids 不能为空")
	}
	_, err := c.call(ctx, "Email", "BatchReadEmail", encodeEmailIDs(emailBox(in.Box), in.IDs))
	return err
}

func (c *Client) BatchDeleteEmail(ctx context.Context, in game.EmailIDsIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if len(in.IDs) == 0 {
		return fmt.Errorf("ids 不能为空")
	}
	_, err := c.call(ctx, "Email", "BatchDeleteEmail", encodeEmailIDs(emailBox(in.Box), in.IDs))
	return err
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

func (c *Client) ActivityGroup(ctx context.Context, in game.GroupIn) (game.ActivityGroup, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityGroup{}, err
	}
	id := in.Group()
	if id <= 0 {
		return game.ActivityGroup{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "GetGroup", encodeGroupID(id))
	if err != nil {
		return game.ActivityGroup{}, err
	}
	return decodeGetGroup(body)
}

func (c *Client) SetSplashed(ctx context.Context, in game.IDIn) (bool, error) {
	if _, err := c.Info(); err != nil {
		return false, err
	}
	if in.ID <= 0 {
		return false, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "SetSplashed", encodeLandID(in.ID))
	if err != nil {
		return false, err
	}
	return decodeOK(body)
}

func (c *Client) Invitees(ctx context.Context, in game.IDIn) ([]game.Invitee, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeInvitees(in.ID))
	if err != nil {
		return nil, err
	}
	return decodeInvitees(body)
}

func (c *Client) Signin(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if in.RewardID <= 0 {
		list, err := c.Activities(ctx)
		if err != nil {
			return game.ActivityOpOut{}, err
		}
		for _, a := range list {
			if a.ID == in.ID {
				in.RewardID = a.SigninRewardID
				break
			}
		}
	}
	if in.RewardID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("reward_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeSignin(in.ID, in.RewardID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 104, 1, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimProgress(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeProgress(in.ID, in.Step))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeProgressClaim(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ShopBuy(ctx context.Context, in game.ActivityShopIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if in.GoodsID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("goods_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeShopBuy(in.ID, in.GoodsID, in.Count))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 101, 1, 2)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ShopBatchBuy(ctx context.Context, in game.ActivityBatchIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if len(in.Items) == 0 {
		return game.ActivityOpOut{}, fmt.Errorf("items 不能为空")
	}
	for _, it := range in.Items {
		if it.GoodsID <= 0 {
			return game.ActivityOpOut{}, fmt.Errorf("goods_id 不能为空")
		}
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeShopBatch(in.ID, in.Items))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 134, 1, 2)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) RandBuy(ctx context.Context, in game.ActivityShopIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if in.GoodsID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("goods_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeRandBuy(in.ID, in.GoodsID, in.Count))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 102, 1, 2)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) RandRefresh(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeRandRefresh(in.ID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 103, 0, 1)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimMega(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeMegaClaim(in.ID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeMegaClaim(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) TechSubmit(ctx context.Context, in game.TechIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if in.NodeID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("node_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeTechSubmit(in.ID, in.NodeID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeTechSubmit(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) Draw(ctx context.Context, in game.DrawIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeDraw(in.ID, in.Count))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 105, 1, 2)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) DrawHistory(ctx context.Context, in game.IDIn) (game.DrawHistoryOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DrawHistoryOut{}, err
	}
	if in.ID <= 0 {
		return game.DrawHistoryOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeDrawHistory(in.ID))
	if err != nil {
		return game.DrawHistoryOut{}, err
	}
	return decodeDrawHistory(body)
}

func (c *Client) MarkViewed(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeMarkViewed(in.ID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return decodeOperateOut(body, 0, 0, 0)
}

func (c *Client) RandBatchBuy(ctx context.Context, in game.ActivityBatchIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if len(in.Items) == 0 {
		return game.ActivityOpOut{}, fmt.Errorf("items 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeRandBatch(in.ID, in.Items))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 107, 1, 2)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) LotteryHistory(ctx context.Context, in game.IDIn) (game.LotteryHistoryOut, error) {
	if _, err := c.Info(); err != nil {
		return game.LotteryHistoryOut{}, err
	}
	if in.ID <= 0 {
		return game.LotteryHistoryOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeLotteryHistory(in.ID))
	if err != nil {
		return game.LotteryHistoryOut{}, err
	}
	return decodeLotteryHistory(body)
}

func (c *Client) CheerJoin(ctx context.Context, in game.CheerJoinIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 || in.CampID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 和 camp_id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCheerJoin(in.ID, in.CampID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return decodeOperateOut(body, 0, 0, 0)
}

func (c *Client) CheerSubmit(ctx context.Context, in game.CheerSubmitIn) (game.CheerSubmitOut, error) {
	if _, err := c.Info(); err != nil {
		return game.CheerSubmitOut{}, err
	}
	if in.ID <= 0 {
		return game.CheerSubmitOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCheerSubmit(in.ID, in.Count))
	if err != nil {
		return game.CheerSubmitOut{}, err
	}
	out, err := decodeCheerSubmit(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) CheerClaim(ctx context.Context, in game.CheerClaimIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCheerClaim(in.ID, in.Tier))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 112, 1, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) Recallable(ctx context.Context, in game.IDIn) (game.RecallListOut, error) {
	if _, err := c.Info(); err != nil {
		return game.RecallListOut{}, err
	}
	if in.ID <= 0 {
		return game.RecallListOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeRecallable(in.ID))
	if err != nil {
		return game.RecallListOut{}, err
	}
	return decodeRecallable(body)
}

func (c *Client) Recalled(ctx context.Context, in game.IDIn) (game.RecallListOut, error) {
	if _, err := c.Info(); err != nil {
		return game.RecallListOut{}, err
	}
	if in.ID <= 0 {
		return game.RecallListOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeRecalled(in.ID))
	if err != nil {
		return game.RecallListOut{}, err
	}
	return decodeRecalled(body)
}

func (c *Client) CharityShare(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCharityShare(in.ID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 135, 1, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) CharityDonate(ctx context.Context, in game.IDIn) (game.CharityDonateOut, error) {
	if _, err := c.Info(); err != nil {
		return game.CharityDonateOut{}, err
	}
	if in.ID <= 0 {
		return game.CharityDonateOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCharityDonate(in.ID))
	if err != nil {
		return game.CharityDonateOut{}, err
	}
	out, err := decodeCharityDonate(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) CharityClaim(ctx context.Context, in game.CharityClaimIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCharityClaim(in.ID, in.Score))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 137, 2, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) CharityXhh(ctx context.Context, in game.IDIn) (game.CharityXhhOut, error) {
	if _, err := c.Info(); err != nil {
		return game.CharityXhhOut{}, err
	}
	if in.ID <= 0 {
		return game.CharityXhhOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCharityXhh(in.ID))
	if err != nil {
		return game.CharityXhhOut{}, err
	}
	out, err := decodeCharityXhh(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) CharityAgree(ctx context.Context, in game.CharityAgreeIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeCharityAgree(in.ID, in.Agreed))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return decodeOperateOut(body, 139, 0, 0)
}

func (c *Client) huntOperate(ctx context.Context, id, cmd int64) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if id <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeHunt(id, cmd))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return decodeOperateOut(body, 0, 0, 0)
}

func (c *Client) HuntFinishCG(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntFinishCG)
}

func (c *Client) HuntGuide(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntGuide)
}

func (c *Client) HuntFeed(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntFeed)
}

func (c *Client) HuntDraw(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntDraw)
}

func (c *Client) huntLogs(ctx context.Context, id, cmd int64) (game.HuntLogOut, error) {
	if _, err := c.Info(); err != nil {
		return game.HuntLogOut{}, err
	}
	if id <= 0 {
		return game.HuntLogOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeHunt(id, cmd))
	if err != nil {
		return game.HuntLogOut{}, err
	}
	return decodeHuntLogs(body, cmd)
}

func (c *Client) HuntLog(ctx context.Context, in game.IDIn) (game.HuntLogOut, error) {
	return c.huntLogs(ctx, in.ID, cmdHuntLog)
}

func (c *Client) HuntClaimStory(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntClaimStory)
}

func (c *Client) HuntClaimSeed(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntClaimSeed)
}

func (c *Client) HuntRefreshCharm(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntRefreshCharm)
}

func (c *Client) HuntEquip(ctx context.Context, in game.HuntEquipIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if len(in.CharmIDs) == 0 {
		return game.ActivityOpOut{}, fmt.Errorf("charm_ids 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeHuntEquip(in.ID, in.CharmIDs))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return decodeOperateOut(body, 0, 0, 0)
}

func (c *Client) HuntBattle(ctx context.Context, in game.HuntBattleIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if in.DefenderGID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("defender_gid 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeHuntBattle(in.ID, in.DefenderGID, in.TreasureID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	return decodeOperateOut(body, 0, 0, 0)
}

func (c *Client) HuntPlunderedLog(ctx context.Context, in game.IDIn) (game.HuntLogOut, error) {
	return c.huntLogs(ctx, in.ID, cmdHuntPlunderedLog)
}

func (c *Client) HuntOpen(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntOpen)
}

func (c *Client) HuntEscort(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntEscort)
}

func (c *Client) HuntCompensate(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntCompensate)
}

func (c *Client) HuntFriendInfo(ctx context.Context, in game.IDIn) (game.ActivityOpOut, error) {
	return c.huntOperate(ctx, in.ID, cmdHuntFriendInfo)
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

func (c *Client) BrewClaim(ctx context.Context, in game.BrewClaimIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeBrewClaim(in.ID, in.ClaimType))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeBrewClaim(body)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimRecall(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeRecallClaim(in.ID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 118, 1, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimReturn(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeReturnGift(in.ID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 119, 1, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimInvite(ctx context.Context, in game.InviteIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	if in.RewardType <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("reward_type 不能为空，1 邀请 2 成长")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeInviteClaim(in.ID, in.RewardType))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 121, 1, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimNewcomer(ctx context.Context, in game.ActivityOpIn) (game.ActivityOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ActivityOpOut{}, err
	}
	if in.ID <= 0 {
		return game.ActivityOpOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Activity", "Operate", encodeNewcomerClaim(in.ID))
	if err != nil {
		return game.ActivityOpOut{}, err
	}
	out, err := decodeOperateOut(body, 122, 1, 0)
	if err != nil {
		return out, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
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

func (c *Client) MysteryBuy(ctx context.Context, in game.IDIn) (game.MysteryBuyOut, error) {
	if _, err := c.Info(); err != nil {
		return game.MysteryBuyOut{}, err
	}
	if in.ID <= 0 {
		return game.MysteryBuyOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "MysteryShop", "Buy", encodeLandID(in.ID))
	if err != nil {
		return game.MysteryBuyOut{}, err
	}
	out, err := decodeMysteryBuy(body)
	if err != nil {
		return game.MysteryBuyOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) MysteryAutoBuy(ctx context.Context, in game.MysteryAutoIn) (game.MysteryAutoOut, error) {
	var out game.MysteryAutoOut
	shop, err := c.Mystery(ctx)
	if err != nil {
		return out, err
	}
	if !shop.Present {
		return out, nil
	}
	for _, g := range shop.Goods {
		if g.Bought || g.ID <= 0 {
			continue
		}
		if in.Currency > 0 && g.Currency != in.Currency {
			continue
		}
		got, err := c.MysteryBuy(ctx, game.IDIn{ID: g.ID})
		if err != nil {
			out.Failed = append(out.Failed, game.BuyFail{ID: g.ID, Reason: err.Error()})
			continue
		}
		out.Items = append(out.Items, got.Items...)
	}
	return out, nil
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

func (c *Client) ClaimDogGifts(ctx context.Context) (game.DogGiftOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DogGiftOut{}, err
	}
	body, err := c.call(ctx, "Dog", "ClaimSkillGifts", nil)
	if err != nil {
		return game.DogGiftOut{}, err
	}
	out, err := decodeDogGifts(body)
	if err != nil {
		return game.DogGiftOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) DogLogs(ctx context.Context, in game.PageIn) (game.DogLogsOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DogLogsOut{}, err
	}
	body, err := c.call(ctx, "Dog", "GetProtectLogs", encodePage(in.From, in.Count, 50))
	if err != nil {
		return game.DogLogsOut{}, err
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

func (c *Client) BuyDog(ctx context.Context, in game.DogBuyIn) (game.DogBuyOut, error) {
	if _, err := c.Info(); err != nil {
		return game.DogBuyOut{}, err
	}
	if in.ID <= 0 {
		return game.DogBuyOut{}, fmt.Errorf("id 不能为空")
	}
	if in.Price <= 0 {
		yard, err := c.Dog(ctx)
		if err != nil {
			return game.DogBuyOut{}, err
		}
		for _, d := range yard.Dogs {
			if d.ID == in.ID {
				in.Price = d.Price
				break
			}
		}
	}
	if in.Price <= 0 {
		return game.DogBuyOut{}, fmt.Errorf("price 不能为空")
	}
	body, err := c.call(ctx, "Dog", "BuyAndActivateDog", encodeDogBuy(in.ID, in.Price))
	if err != nil {
		return game.DogBuyOut{}, err
	}
	out, err := decodeDogBuy(body, in.ID)
	if err != nil {
		return game.DogBuyOut{}, err
	}
	c.applyCurrency(nil, out.Cost)
	return out, nil
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

func (c *Client) Career(ctx context.Context, in game.EnterIn) (game.Career, error) {
	if _, err := c.Info(); err != nil {
		return game.Career{}, err
	}
	body, err := c.call(ctx, "Career", "CareerInfoGet", encodeCareerGID())
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

func (c *Client) MarkAvatar(ctx context.Context, in game.IDIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.ID <= 0 {
		return fmt.Errorf("id 不能为空")
	}
	_, err := c.call(ctx, "AvatarFrame", "MarkAsViewed", encodeLandID(in.ID))
	return err
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

func (c *Client) MarkSkin(ctx context.Context, in game.IDIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.ID <= 0 {
		return fmt.Errorf("id 不能为空")
	}
	_, err := c.call(ctx, "Skin", "MarkAsViewed", encodeLandID(in.ID))
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

func (c *Client) SolarTerms(ctx context.Context) (game.SolarOut, error) {
	if _, err := c.Info(); err != nil {
		return game.SolarOut{}, err
	}
	body, err := c.call(ctx, "SolarTerms", "GetSolarTerms", nil)
	if err != nil {
		return game.SolarOut{}, err
	}
	return decodeSolarTerms(body)
}

func (c *Client) SolarRedDot(ctx context.Context) (game.RedDot, error) {
	if _, err := c.Info(); err != nil {
		return game.RedDot{}, err
	}
	body, err := c.call(ctx, "SolarTerms", "GetSolarTermsRedDot", nil)
	if err != nil {
		return game.RedDot{}, err
	}
	return decodeRedDot(body)
}

func (c *Client) ClaimSolar(ctx context.Context, in game.IDIn) (game.SolarClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.SolarClaimOut{}, err
	}
	if in.ID <= 0 {
		return game.SolarClaimOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "SolarTerms", "ClaimSolarTerms", encodeLandID(in.ID))
	if err != nil {
		return game.SolarClaimOut{}, err
	}
	out, err := decodeSolarClaim(body)
	if err != nil {
		return game.SolarClaimOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) ClaimAllSolar(ctx context.Context) ([]game.Item, error) {
	list, err := c.SolarTerms(ctx)
	if err != nil {
		return nil, err
	}
	var out []game.Item
	for _, s := range list.Terms {
		if s.Status != 2 || s.ID <= 0 {
			continue
		}
		got, err := c.ClaimSolar(ctx, game.IDIn{ID: s.ID})
		if err != nil {
			return out, err
		}
		out = append(out, got.Items...)
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

func (c *Client) ClaimAchieveGoal(ctx context.Context, in game.AchieveGoalIn) (game.AchieveGoalOut, error) {
	if _, err := c.Info(); err != nil {
		return game.AchieveGoalOut{}, err
	}
	if in.GoalID <= 0 {
		return game.AchieveGoalOut{}, fmt.Errorf("goal_id 不能为空")
	}
	body, err := c.call(ctx, "Achieve", "ClaimGoalReward", encodeAchieveGoal(in.Kind, in.ID, in.GoalID))
	if err != nil {
		return game.AchieveGoalOut{}, err
	}
	return decodeAchieveGoalOut(body)
}

func (c *Client) ClaimAchieveLevel(ctx context.Context, in game.AchieveIn) (game.ItemOpOut, error) {
	if _, err := c.Info(); err != nil {
		return game.ItemOpOut{}, err
	}
	body, err := c.call(ctx, "Achieve", "ClaimScopeLevel", encodeAchieve(in.Kind, in.ID))
	if err != nil {
		return game.ItemOpOut{}, err
	}
	out, err := decodeItemOp(body, 0, 1, 2, false)
	if err != nil {
		return game.ItemOpOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
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

func (c *Client) ClaimPass(ctx context.Context) (game.PassClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.PassClaimOut{}, err
	}
	body, err := c.call(ctx, "Season", "ClaimBattlePassRewards", nil)
	if err != nil {
		return game.PassClaimOut{}, err
	}
	out, err := decodePassClaim(body)
	if err != nil {
		return game.PassClaimOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
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

func (c *Client) MallDiamonds(ctx context.Context) ([]game.Product, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Mall", "GetDiamondItems", nil)
	if err != nil {
		return nil, err
	}
	return decodeProducts(body)
}

func (c *Client) MallProfiles(ctx context.Context) ([]game.MallProfile, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Mall", "GetMallProfiles", nil)
	if err != nil {
		return nil, err
	}
	return decodeMallProfiles(body)
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
	out, err = decodeMallBuy(body)
	if err != nil {
		return out, err
	}
	c.applyCurrency(out.Items, nil)
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

func (c *Client) ClaimMonthCard(ctx context.Context, in game.IDIn) (game.MonthCardClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.MonthCardClaimOut{}, err
	}
	if in.ID <= 0 {
		return game.MonthCardClaimOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "Mall", "ClaimMonthCardReward", encodeLandID(in.ID))
	if err != nil {
		return game.MonthCardClaimOut{}, err
	}
	out, err := decodeMonthCardClaim(body)
	if err != nil {
		return game.MonthCardClaimOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
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

func (c *Client) ClaimRedPacket(ctx context.Context, in game.IDIn) (game.RedPacketClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.RedPacketClaimOut{}, err
	}
	if in.ID <= 0 {
		return game.RedPacketClaimOut{}, fmt.Errorf("id 不能为空")
	}
	body, err := c.call(ctx, "RedPacket", "ClaimRedPacket", encodeLandID(in.ID))
	if err != nil {
		return game.RedPacketClaimOut{}, err
	}
	out, err := decodeRedPacketClaim(body)
	if err != nil {
		return game.RedPacketClaimOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
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
		got, err := c.ClaimRedPacket(ctx, game.IDIn{ID: p.ID})
		if err != nil {
			continue
		}
		all = append(all, got.Items...)
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

func (c *Client) VisitSummary(ctx context.Context) (game.VisitSummary, error) {
	if _, err := c.Info(); err != nil {
		return game.VisitSummary{}, err
	}
	body, err := c.call(ctx, "Interact", "GetInteractSummary", nil)
	if err != nil {
		return game.VisitSummary{}, err
	}
	return decodeVisitSummary(body)
}

func (c *Client) VisitPopup(ctx context.Context) (game.VisitPopup, error) {
	if _, err := c.Info(); err != nil {
		return game.VisitPopup{}, err
	}
	body, err := c.call(ctx, "Interact", "GetInteractInfo", nil)
	if err != nil {
		return game.VisitPopup{}, err
	}
	return decodeVisitPopup(body)
}

func (c *Client) VisitPage(ctx context.Context, in game.VisitPageIn) ([]game.VisitLog, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Interact", "GetInteractions", encodeVisitPage(in.Page))
	if err != nil {
		return nil, err
	}
	return decodeVisitLogs(body)
}

func (c *Client) DismissVisit(ctx context.Context, in game.EnterIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	_, err := c.call(ctx, "Interact", "DismissInteractPopup", encodeLeave(in.GID))
	return err
}

func (c *Client) DeleteVisit(ctx context.Context, in game.VisitDeleteIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if len(in.IDs) == 0 {
		return fmt.Errorf("ids 不能为空")
	}
	_, err := c.call(ctx, "Interact", "DeleteInteractions", encodeGIDs(in.IDs))
	return err
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

func (c *Client) AlbumLevels(ctx context.Context, in game.AlbumIn) (game.AlbumLevels, error) {
	if _, err := c.Info(); err != nil {
		return game.AlbumLevels{}, err
	}
	body, err := c.call(ctx, "Illustrated", "GetIllustratedLevelListV2", encodeIllustratedLevels(in.Type))
	if err != nil {
		return game.AlbumLevels{}, err
	}
	return decodeAlbumLevels(body)
}

func (c *Client) ClaimAlbum(ctx context.Context, in game.AlbumIn) (game.AlbumClaimOut, error) {
	if _, err := c.Info(); err != nil {
		return game.AlbumClaimOut{}, err
	}
	if in.Type <= 0 {
		in.Type = 1
	}
	body, err := c.call(ctx, "Illustrated", "ClaimAllRewardsV2", encodeLandID(in.Type))
	if err != nil {
		return game.AlbumClaimOut{}, err
	}
	out, err := decodeAlbumClaim(body)
	if err != nil {
		return game.AlbumClaimOut{}, err
	}
	_, _ = c.Bag(ctx)
	return out, nil
}

func (c *Client) MarkAlbum(ctx context.Context, in game.AlbumIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.Type <= 0 {
		in.Type = 1
	}
	_, err := c.call(ctx, "Illustrated", "ClearNewUnlockedFruitsV2", encodeIllustratedLevels(in.Type))
	return err
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

func (c *Client) setHost(gid int64) {
	c.mu.Lock()
	c.hostGID = gid
	c.mu.Unlock()
}

func (c *Client) startGateHeartbeat() {
	c.mu.Lock()
	if c.hbStop != nil {
		close(c.hbStop)
	}
	stop := make(chan struct{})
	c.hbStop = stop
	c.mu.Unlock()
	go c.heartbeatLoop(stop)
}

func (c *Client) heartbeatLoop(stop <-chan struct{}) {
	ticker := time.NewTicker(ace.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, _ = c.Heartbeat(ctx, game.RefreshIn{})
			cancel()
		}
	}
}

func (c *Client) closeConn() error {
	c.mu.Lock()
	life := c.life
	c.life = nil
	if c.hbStop != nil {
		close(c.hbStop)
		c.hbStop = nil
	}
	c.closed = true
	conn := c.conn
	c.conn = nil
	c.loggedIn = false
	c.bag = nil
	c.hostGID = 0
	c.hbCount = 0
	c.serverMs = 0
	c.hbErr = ""
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
	if _, ok := gate.ServiceName(service); !ok {
		return nil, fmt.Errorf("服务 %q 没有登记官方名字，拒绝发包", service)
	}
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
	timer := time.NewTimer(callTimeout)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		c.dropPending(seq)
		return nil, ctx.Err()
	case <-timer.C:
		c.dropPending(seq)
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

// dropPending forgets a request the caller gave up on. readLoop deletes the
// entry itself when a reply lands, so a timed-out call must clean up or the
// map grows for the life of the session.
func (c *Client) dropPending(seq int64) {
	c.mu.Lock()
	delete(c.pending, seq)
	c.mu.Unlock()
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
