package game

import (
	"context"
	"errors"
)

var ErrNotLogin = errors.New("未登录")

type User struct {
	GID    int64  `json:"gid"`
	Name   string `json:"name"`
	Level  int64  `json:"level"`
	Exp    int64  `json:"exp"`
	Gold   int64  `json:"gold"`
	OpenID string `json:"open_id"`
	Avatar string `json:"avatar_url"`
}

type Land struct {
	ID          int64  `json:"id"`
	Unlocked    bool   `json:"unlocked"`
	Level       int64  `json:"level"`
	PlantID     int64  `json:"plant_id,omitempty"`
	PlantName   string `json:"plant_name,omitempty"`
	FruitID     int64  `json:"fruit_id,omitempty"`
	DryNum      int64  `json:"dry_num,omitempty"`
	HasWeed     bool   `json:"has_weed"`
	HasInsect   bool   `json:"has_insect"`
	Stealable   bool   `json:"stealable"`
	LeftFruit   int64  `json:"left_fruit,omitempty"`
	CouldUnlock bool   `json:"could_unlock,omitempty"`
}

type Item struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}

type Friend struct {
	GID       int64  `json:"gid"`
	Name      string `json:"name"`
	Level     int64  `json:"level"`
	OpenID    string `json:"open_id"`
	Avatar    string `json:"avatar_url"`
	DryNum    int64  `json:"dry_num,omitempty"`
	WeedNum   int64  `json:"weed_num,omitempty"`
	InsectNum int64  `json:"insect_num,omitempty"`
	StealNum  int64  `json:"steal_num,omitempty"`
}

type LoginIn struct {
	Code   string `json:"code"`
	OpenID string `json:"open_id"`
}

type HarvestIn struct {
	LandIDs []int64 `json:"land_ids"`
	HostGID int64   `json:"host_gid"`
	IsAll   bool    `json:"is_all"`
}

type PlantIn struct {
	SeedID  int64   `json:"seed_id"`
	LandIDs []int64 `json:"land_ids"`
}

type HelpIn struct {
	GID     int64   `json:"gid"`
	LandIDs []int64 `json:"land_ids"`
}

type LandOpIn struct {
	LandIDs []int64 `json:"land_ids"`
	HostGID int64   `json:"host_gid"`
}

type FertilizeIn struct {
	LandIDs      []int64 `json:"land_ids"`
	FertilizerID int64   `json:"fertilizer_id"`
}

type BagItem struct {
	ID        int64 `json:"id"`
	Count     int64 `json:"count"`
	UID       int64 `json:"uid,omitempty"`
	SellPrice int64 `json:"sell_price,omitempty"`
}

type SellIn struct {
	Items []BagItem `json:"items"`
}

type UseIn struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
	UID   int64 `json:"uid,omitempty"`
}

type Shop struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type int64  `json:"type"`
}

type ShopIn struct {
	ShopID int64 `json:"shop_id"`
}

type Goods struct {
	ID        int64 `json:"id"`
	Bought    int64 `json:"bought"`
	Price     int64 `json:"price"`
	Limit     int64 `json:"limit"`
	Unlocked  bool  `json:"unlocked"`
	ItemID    int64 `json:"item_id"`
	ItemCount int64 `json:"item_count,omitempty"`
}

type BuyIn struct {
	GoodsID int64 `json:"goods_id"`
	Num     int64 `json:"num"`
	Price   int64 `json:"price"`
}

type BuyOut struct {
	Items []Item `json:"items"`
	Cost  []Item `json:"cost"`
}

type EnterIn struct {
	GID int64 `json:"gid"`
}

type Visit struct {
	Host       User   `json:"host"`
	Lands      []Land `json:"lands"`
	DogID      int64  `json:"dog_id,omitempty"`
	DogFoodSec int64  `json:"dog_food_sec,omitempty"`
}

type Application struct {
	GID    int64  `json:"gid"`
	Time   int64  `json:"time,omitempty"`
	OpenID string `json:"open_id,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar_url,omitempty"`
}

type GIDsIn struct {
	GIDs []int64 `json:"gids"`
}

type RefreshIn struct {
	HostGID int64 `json:"host_gid"`
}

type LandIDIn struct {
	LandID int64 `json:"land_id"`
}

type IDIn struct {
	ID int64 `json:"id"`
}

type RemoveIn struct {
	LandIDs []int64 `json:"land_ids"`
}

type Weather struct {
	Type         int64  `json:"type"`
	Source       int64  `json:"source,omitempty"`
	Begin        int64  `json:"begin,omitempty"`
	End          int64  `json:"end,omitempty"`
	Active       bool   `json:"active"`
	Afterglow    int64  `json:"afterglow,omitempty"`
	AfterglowTyp int64  `json:"afterglow_type,omitempty"`
	CanCollect   bool   `json:"can_collect"`
	BlockReason  int64  `json:"block_reason,omitempty"`
	Name         string `json:"name,omitempty"`
}

type Task struct {
	ID       int64  `json:"id"`
	Progress int64  `json:"progress"`
	Total    int64  `json:"total,omitempty"`
	Claimed  bool   `json:"claimed"`
	Status   int64  `json:"status,omitempty"`
	Type     int64  `json:"type,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Rewards  []Item `json:"rewards,omitempty"`
}

type ActiveBox struct {
	ID      int64  `json:"id"`
	Need    int64  `json:"need"`
	Claimed bool   `json:"claimed"`
	Rewards []Item `json:"rewards,omitempty"`
}

type Active struct {
	Type     int64       `json:"type"`
	Progress int64       `json:"progress"`
	Boxes    []ActiveBox `json:"boxes,omitempty"`
}

type TaskBoard struct {
	Growth  []Task   `json:"growth"`
	Daily   []Task   `json:"daily"`
	Actives []Active `json:"actives"`
}

type ClaimIn struct {
	ID  int64   `json:"id"`
	IDs []int64 `json:"ids"`
}

type DailyIn struct {
	PointIDs []int64 `json:"point_ids"`
}

type EmailBoxIn struct {
	Box int64 `json:"box"`
}

type EmailIn struct {
	Box int64  `json:"box"`
	ID  string `json:"id"`
}

type Email struct {
	ID        string `json:"id"`
	Type      int64  `json:"type,omitempty"`
	Title     string `json:"title,omitempty"`
	Claimed   bool   `json:"claimed"`
	HasReward bool   `json:"has_reward"`
	Subtitle  string `json:"subtitle,omitempty"`
}

type EmailDetail struct {
	ID       string `json:"id"`
	Type     int64  `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	Rewards  []Item `json:"rewards,omitempty"`
	SendTime int64  `json:"send_time,omitempty"`
	Expire   int64  `json:"expire,omitempty"`
	Read     bool   `json:"read"`
	Claimed  bool   `json:"claimed"`
	Reason   string `json:"reason,omitempty"`
	Tips     string `json:"tips,omitempty"`
}

type Activity struct {
	ID             int64           `json:"id"`
	GroupID        int64           `json:"group_id,omitempty"`
	Type           int64           `json:"type,omitempty"`
	Name           string          `json:"name,omitempty"`
	Desc           string          `json:"desc,omitempty"`
	Start          int64           `json:"start,omitempty"`
	End            int64           `json:"end,omitempty"`
	ClientID       int64           `json:"client_id,omitempty"`
	Status         int64           `json:"status,omitempty"`
	RedDot         bool            `json:"red_dot,omitempty"`
	SigninClaimed  bool            `json:"signin_claimed,omitempty"`
	SigninRewardID int64           `json:"signin_reward_id,omitempty"`
	Shop           []ActivityGoods `json:"shop,omitempty"`
	Nodes          []TechNode      `json:"nodes,omitempty"`
	Lottery        *LotteryInfo    `json:"lottery,omitempty"`
	Brew           *BrewInfo       `json:"brew,omitempty"`
	Recall         *RecallInfo     `json:"recall,omitempty"`
	Invite         *InviteInfo     `json:"invite,omitempty"`
	Gift           *GiftState      `json:"gift,omitempty"`
}

type LotteryInfo struct {
	FreeLeft  int64 `json:"free_left"`
	FreeLimit int64 `json:"free_limit,omitempty"`
	PaidLeft  int64 `json:"paid_left"`
	PaidLimit int64 `json:"paid_limit,omitempty"`
	CostID    int64 `json:"cost_id,omitempty"`
	CostCount int64 `json:"cost_count,omitempty"`
	Diamond   int64 `json:"diamond,omitempty"`
}

type BrewInfo struct {
	Value    int64   `json:"value"`
	Step     int64   `json:"step"`
	Steps    int64   `json:"steps,omitempty"`
	CanClaim bool    `json:"can_claim"`
	Allowed  []int64 `json:"allowed,omitempty"`
}

type RecallInfo struct {
	Pending       int64 `json:"pending,omitempty"`
	Returner      bool  `json:"returner,omitempty"`
	ReturnPending int64 `json:"return_pending,omitempty"`
}

type InviteInfo struct {
	Invitee     bool  `json:"invitee,omitempty"`
	InviteCount int64 `json:"invite_count,omitempty"`
	Limit       int64 `json:"limit,omitempty"`
}

type GiftState struct {
	Sent      int64 `json:"sent"`
	SendLimit int64 `json:"send_limit,omitempty"`
}

type LotteryIn struct {
	ID      int64 `json:"id"`
	HostGID int64 `json:"host_gid"`
	Free    int64 `json:"free"`
	Paid    int64 `json:"paid"`
}

type LotteryOut struct {
	Items   []Item `json:"items"`
	Costs   []Item `json:"costs,omitempty"`
	Partial bool   `json:"partial,omitempty"`
}

type BrewItem struct {
	UID   int64 `json:"uid"`
	Count int64 `json:"count"`
}

type BrewStartIn struct {
	ID    int64      `json:"id"`
	Items []BrewItem `json:"items"`
}

type BrewStartOut struct {
	Value int64 `json:"value"`
}

type BrewStepOut struct {
	Step       int64 `json:"step"`
	Multiplier int64 `json:"multiplier,omitempty"`
	Amount     int64 `json:"amount,omitempty"`
	Finished   bool  `json:"finished"`
}

type BrewClaimIn struct {
	ID        int64 `json:"id"`
	ClaimType int64 `json:"claim_type"`
}

type InviteIn struct {
	ID         int64 `json:"id"`
	RewardType int64 `json:"reward_type"`
}

type GiftIn struct {
	ID        int64 `json:"id"`
	GID       int64 `json:"gid"`
	MsgTextID int64 `json:"msg_text_id"`
}

type ActivityGoods struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Limit   int64  `json:"limit,omitempty"`
	Bought  int64  `json:"bought,omitempty"`
	Items   []Item `json:"items,omitempty"`
	Cost    []Item `json:"cost,omitempty"`
	Diamond int64  `json:"diamond,omitempty"`
}

type TechNode struct {
	TreeID   int64  `json:"tree_id,omitempty"`
	ID       int64  `json:"id"`
	Status   int64  `json:"status"`
	Progress int64  `json:"progress,omitempty"`
	Target   int64  `json:"target,omitempty"`
	Cost     []Item `json:"cost,omitempty"`
	Rewards  []Item `json:"rewards,omitempty"`
}

type ActivityOpIn struct {
	ID       int64 `json:"id"`
	RewardID int64 `json:"reward_id"`
	Step     int64 `json:"step"`
}

type ActivityShopIn struct {
	ID      int64 `json:"id"`
	GoodsID int64 `json:"goods_id"`
	Count   int64 `json:"count"`
}

type TechIn struct {
	ID     int64 `json:"id"`
	NodeID int64 `json:"node_id"`
}

type MysteryGoods struct {
	ID       int64 `json:"id"`
	ItemID   int64 `json:"item_id"`
	Quality  int64 `json:"quality,omitempty"`
	Count    int64 `json:"count,omitempty"`
	Currency int64 `json:"currency,omitempty"`
	Price    int64 `json:"price,omitempty"`
	Discount int64 `json:"discount,omitempty"`
	Bought   bool  `json:"bought"`
	Original int64 `json:"original,omitempty"`
}

type MysteryShop struct {
	Present bool           `json:"present"`
	Start   int64          `json:"start,omitempty"`
	End     int64          `json:"end,omitempty"`
	Goods   []MysteryGoods `json:"goods"`
}

type BattlePass struct {
	ID             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Level          int64  `json:"level"`
	MaxLevel       int64  `json:"max_level,omitempty"`
	Exp            int64  `json:"exp,omitempty"`
	LevelExp       int64  `json:"level_exp,omitempty"`
	NextExp        int64  `json:"next_exp,omitempty"`
	Premium        bool   `json:"premium"`
	FreeClaimed    int64  `json:"free_claimed"`
	PremiumClaimed int64  `json:"premium_claimed,omitempty"`
}

type Season struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name,omitempty"`
	Phase      int64      `json:"phase,omitempty"`
	Start      int64      `json:"start,omitempty"`
	End        int64      `json:"end,omitempty"`
	ServerTime int64      `json:"server_time,omitempty"`
	Pass       BattlePass `json:"pass"`
}

type MallIn struct {
	Slot int64 `json:"slot"`
}

type Product struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Num       int64  `json:"num,omitempty"`
	Status    int64  `json:"status,omitempty"`
	Available bool   `json:"available"`
	GoodsType int64  `json:"goods_type,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
	Price     Item   `json:"price,omitempty"`
	Rewards   []Item `json:"rewards,omitempty"`
	Bought    int64  `json:"bought,omitempty"`
	Limit     int64  `json:"limit,omitempty"`
}

type MallBuyIn struct {
	ID  int64 `json:"id"`
	Num int64 `json:"num"`
}

type MonthCard struct {
	ID            int64  `json:"id"`
	Claimable     bool   `json:"claimable"`
	Days          int64  `json:"days,omitempty"`
	TotalDays     int64  `json:"total_days,omitempty"`
	ClaimedAmount int64  `json:"claimed_amount,omitempty"`
	Rewards       []Item `json:"rewards,omitempty"`
}

type RedPacket struct {
	ID       int64 `json:"id"`
	CanClaim bool  `json:"can_claim"`
}

type VisitLog struct {
	Time      int64  `json:"time"`
	Action    int64  `json:"action"`
	GID       int64  `json:"gid"`
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar_url,omitempty"`
	CropID    int64  `json:"crop_id,omitempty"`
	CropCount int64  `json:"crop_count,omitempty"`
	Times     int64  `json:"times,omitempty"`
	Level     int64  `json:"level,omitempty"`
}

type AlbumIn struct {
	Type   int64 `json:"type"`
	Rarity int64 `json:"rarity"`
}

type AlbumItem struct {
	FruitID  int64  `json:"fruit_id"`
	Rarity   int64  `json:"rarity,omitempty"`
	Unlocked bool   `json:"unlocked"`
	Progress int64  `json:"progress,omitempty"`
	New      bool   `json:"new,omitempty"`
	Rewards  []Item `json:"rewards,omitempty"`
}

type Album struct {
	Items     []AlbumItem `json:"items"`
	Progress  int64       `json:"progress,omitempty"`
	Claimable bool        `json:"claimable"`
}

type Status struct {
	LoggedIn bool `json:"logged_in"`
	User     User `json:"user,omitempty"`
	ACE      ACE  `json:"ace"`
}

type ACE struct {
	Uploads   int    `json:"uploads"`
	Reports   int    `json:"status_reports"`
	Failures  int    `json:"failures"`
	LastError string `json:"last_error,omitempty"`
}

// Session is the live game connection used by HTTP handlers.
type Session interface {
	Login(ctx context.Context, in LoginIn) (User, error)
	Info() (User, error)
	Status() (Status, error)
	Refresh(ctx context.Context, in RefreshIn) ([]Land, error)
	Harvest(ctx context.Context, in HarvestIn) ([]Item, error)
	Steal(ctx context.Context, in HarvestIn) ([]Item, error)
	Plant(ctx context.Context, in PlantIn) error
	Remove(ctx context.Context, in RemoveIn) ([]Land, error)
	Unlock(ctx context.Context, in LandIDIn) (Land, error)
	Upgrade(ctx context.Context, in LandIDIn) (Land, error)
	Farming(ctx context.Context, in LandOpIn) ([]Land, error)
	Water(ctx context.Context, in LandOpIn) error
	Weed(ctx context.Context, in LandOpIn) error
	Bug(ctx context.Context, in LandOpIn) error
	Fertilize(ctx context.Context, in FertilizeIn) error
	Friends(ctx context.Context) ([]Friend, error)
	Applications(ctx context.Context) ([]Application, error)
	Accept(ctx context.Context, in GIDsIn) ([]Friend, error)
	Reject(ctx context.Context, in GIDsIn) error
	DeleteFriend(ctx context.Context, in EnterIn) error
	Help(ctx context.Context, in HelpIn) error
	Enter(ctx context.Context, in EnterIn) (Visit, error)
	Leave(ctx context.Context, in EnterIn) error
	ShareCheck(ctx context.Context) (bool, error)
	ShareClaim(ctx context.Context) ([]Item, error)
	Bag(ctx context.Context) ([]BagItem, error)
	Sell(ctx context.Context, in SellIn) ([]Item, error)
	Use(ctx context.Context, in UseIn) ([]Item, error)
	Shops(ctx context.Context) ([]Shop, error)
	Goods(ctx context.Context, in ShopIn) ([]Goods, error)
	Buy(ctx context.Context, in BuyIn) (BuyOut, error)
	Weather(ctx context.Context) (Weather, error)
	TodayWeather(ctx context.Context) ([]Weather, error)
	Tasks(ctx context.Context) (TaskBoard, error)
	ClaimTask(ctx context.Context, in ClaimIn) ([]Item, error)
	ClaimDaily(ctx context.Context, in DailyIn) ([]Item, error)
	Emails(ctx context.Context, in EmailBoxIn) ([]Email, error)
	ReadEmail(ctx context.Context, in EmailIn) (EmailDetail, error)
	ClaimEmail(ctx context.Context, in EmailIn) ([]Item, error)
	ClaimAllEmail(ctx context.Context, in EmailBoxIn) ([]Item, error)
	Activities(ctx context.Context) ([]Activity, error)
	Signin(ctx context.Context, in ActivityOpIn) ([]Item, error)
	ClaimProgress(ctx context.Context, in ActivityOpIn) ([]Item, error)
	ShopBuy(ctx context.Context, in ActivityShopIn) ([]Item, error)
	ClaimMega(ctx context.Context, in ActivityOpIn) ([]Item, error)
	TechSubmit(ctx context.Context, in TechIn) ([]Item, error)
	Lottery(ctx context.Context, in LotteryIn) (LotteryOut, error)
	BrewStart(ctx context.Context, in BrewStartIn) (BrewStartOut, error)
	BrewStep(ctx context.Context, in ActivityOpIn) (BrewStepOut, error)
	BrewClaim(ctx context.Context, in BrewClaimIn) ([]Item, error)
	ClaimRecall(ctx context.Context, in ActivityOpIn) ([]Item, error)
	ClaimReturn(ctx context.Context, in ActivityOpIn) ([]Item, error)
	ClaimInvite(ctx context.Context, in InviteIn) ([]Item, error)
	ClaimNewcomer(ctx context.Context, in ActivityOpIn) ([]Item, error)
	SendGift(ctx context.Context, in GiftIn) (int64, error)
	Mystery(ctx context.Context) (MysteryShop, error)
	MysteryBuy(ctx context.Context, in IDIn) ([]Item, error)
	MysteryLeave(ctx context.Context) error
	Season(ctx context.Context) (Season, error)
	ClaimPass(ctx context.Context) ([]Item, error)
	Mall(ctx context.Context, in MallIn) ([]Product, error)
	MallBuy(ctx context.Context, in MallBuyIn) (BuyOut, error)
	MonthCards(ctx context.Context) ([]MonthCard, error)
	ClaimMonthCard(ctx context.Context, in IDIn) ([]Item, error)
	RedPackets(ctx context.Context) ([]RedPacket, error)
	ClaimRedPacket(ctx context.Context, in IDIn) ([]Item, error)
	ClaimAllRedPackets(ctx context.Context) ([]Item, error)
	VisitLogs(ctx context.Context) ([]VisitLog, error)
	Album(ctx context.Context, in AlbumIn) (Album, error)
	ClaimAlbum(ctx context.Context, in AlbumIn) ([]Item, error)
	Close() error
}

type Factory func() Session
