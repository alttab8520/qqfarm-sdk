package game

import (
	"context"
	"errors"
)

var ErrNotLogin = errors.New("未登录")

type User struct {
	GID        int64  `json:"gid"`
	Name       string `json:"name"`
	Level      int64  `json:"level"`
	Exp        int64  `json:"exp"`
	Gold       int64  `json:"gold"`
	OpenID     string `json:"open_id"`
	Avatar     string `json:"avatar_url"`
	Remark     string `json:"remark,omitempty"`
	Signature  string `json:"signature,omitempty"`
	Gender     int64  `json:"gender,omitempty"`
	Authorized int64  `json:"authorized,omitempty"`
	LastOnline int64  `json:"last_online,omitempty"`
}

type Land struct {
	ID             int64      `json:"id"`
	Unlocked       bool       `json:"unlocked"`
	Level          int64      `json:"level"`
	MaxLevel       int64      `json:"max_level,omitempty"`
	CouldUnlock    bool       `json:"could_unlock,omitempty"`
	CouldUpgrade   bool       `json:"could_upgrade,omitempty"`
	Shared         bool       `json:"shared,omitempty"`
	CanShare       bool       `json:"can_share,omitempty"`
	MasterLandID   int64      `json:"master_land_id,omitempty"`
	SlaveLandIDs   []int64    `json:"slave_land_ids,omitempty"`
	LandSize       int64      `json:"land_size,omitempty"`
	PlantID        int64      `json:"plant_id,omitempty"`
	PlantName      string     `json:"plant_name,omitempty"`
	FruitID        int64      `json:"fruit_id,omitempty"`
	FruitNum       int64      `json:"fruit_num,omitempty"`
	DryNum         int64      `json:"dry_num,omitempty"`
	HasWeed        bool       `json:"has_weed"`
	HasInsect      bool       `json:"has_insect"`
	Stealable      bool       `json:"stealable"`
	LeftFruit      int64      `json:"left_fruit,omitempty"`
	GrowSec        int64      `json:"grow_sec,omitempty"`
	Season         int64      `json:"season,omitempty"`
	Phase          int64      `json:"phase,omitempty"`
	StoleNum       int64      `json:"stole_num,omitempty"`
	Stealers       []int64    `json:"stealers,omitempty"`
	Mutant         bool       `json:"mutant,omitempty"`
	MutantIDs      []int64    `json:"mutant_ids,omitempty"`
	LandsLevel     int64      `json:"lands_level,omitempty"`
	FertLeft       int64      `json:"fert_left,omitempty"`
	Intimacy       int64      `json:"intimacy,omitempty"`
	PlantExp       int64      `json:"plant_exp,omitempty"`
	LastFruit      int64      `json:"last_fruit,omitempty"`
	LastFruitNum   int64      `json:"last_fruit_num,omitempty"`
	SourceFruit    int64      `json:"source_fruit,omitempty"`
	SourceFruitNum int64      `json:"source_fruit_num,omitempty"`
	SourcePriceID  int64      `json:"source_price_id,omitempty"`
	SourcePrice    int64      `json:"source_price,omitempty"`
	MutantPriceID  int64      `json:"mutant_price_id,omitempty"`
	MutantPrice    int64      `json:"mutant_price,omitempty"`
	Unlock         *LandNeed  `json:"unlock,omitempty"`
	Upgrade        *LandNeed  `json:"upgrade,omitempty"`
	Buff           *LandBuff  `json:"buff,omitempty"`
	SocialItems    []CropItem `json:"social_items,omitempty"`
}

type LandNeed struct {
	LandID int64  `json:"land_id,omitempty"`
	Lands  int64  `json:"lands,omitempty"`
	Level  int64  `json:"level,omitempty"`
	Gold   int64  `json:"gold,omitempty"`
	Items  []Item `json:"items,omitempty"`
}

type LandBuff struct {
	Yield  int64 `json:"yield,omitempty"`
	Time   int64 `json:"time,omitempty"`
	Exp    int64 `json:"exp,omitempty"`
	Mutant int64 `json:"mutant,omitempty"`
	Pass   int64 `json:"pass,omitempty"`
}

type CropItem struct {
	ItemID  int64 `json:"item_id"`
	Count   int64 `json:"count,omitempty"`
	Type    int64 `json:"type,omitempty"`
	Owner   int64 `json:"owner_id,omitempty"`
	PutTime int64 `json:"put_time,omitempty"`
	End     int64 `json:"end,omitempty"`
	LandID  int64 `json:"land_id,omitempty"`
}

type Item struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}

type Friend struct {
	GID          int64    `json:"gid"`
	Name         string   `json:"name"`
	Level        int64    `json:"level"`
	OpenID       string   `json:"open_id"`
	Avatar       string   `json:"avatar_url"`
	Remark       string   `json:"remark,omitempty"`
	Gold         int64    `json:"gold,omitempty"`
	Authorized   int64    `json:"authorized,omitempty"`
	LastLogin    int64    `json:"last_login,omitempty"`
	DryNum       int64    `json:"dry_num,omitempty"`
	WeedNum      int64    `json:"weed_num,omitempty"`
	InsectNum    int64    `json:"insect_num,omitempty"`
	StealNum     int64    `json:"steal_num,omitempty"`
	New          bool     `json:"new,omitempty"`
	Follow       bool     `json:"follow,omitempty"`
	Exp          int64    `json:"exp,omitempty"`
	Banned       bool     `json:"banned,omitempty"`
	HostType     int64    `json:"host_type,omitempty"`
	Returner     bool     `json:"returner,omitempty"`
	AlbumNormal  int64    `json:"album_normal,omitempty"`
	AlbumPremium int64    `json:"album_premium,omitempty"`
	PassSeason   int64    `json:"pass_season,omitempty"`
	PassLevel    int64    `json:"pass_level,omitempty"`
	Weather      *Weather `json:"weather,omitempty"`
}

type LoginIn struct {
	Code   string `json:"code"`
	OpenID string `json:"open_id"`
}

type YYBAccount struct {
	ID        int64  `json:"id"`
	OpenID    string `json:"openid"`
	UIN       int64  `json:"uin,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Status    string `json:"status,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}

type YYBQROut struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	Image     string `json:"image,omitempty"`
}

type YYBSessionIn struct {
	SessionID string `json:"session_id"`
}

type YYBPollOut struct {
	Status  string `json:"status"`
	ErrCode int64  `json:"errcode,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

type YYBRefIn struct {
	Ref   string `json:"ref,omitempty"`
	AppID string `json:"app_id,omitempty"`
}

type YYBCodeOut struct {
	Code   string `json:"code"`
	OpenID string `json:"openid"`
}

type YYBRefreshOut struct {
	ID     int64  `json:"id"`
	OpenID string `json:"openid"`
	Status string `json:"status"`
}

type YYBDeleteOut struct {
	Deleted int64  `json:"deleted"`
	OpenID  string `json:"openid"`
}

type YYBWXDataIn struct {
	Ref     string         `json:"ref,omitempty"`
	AppID   string         `json:"app_id,omitempty"`
	Payload map[string]any `json:"payload"`
}

type YYBRawOut struct {
	OpenID string         `json:"openid"`
	Result map[string]any `json:"result"`
}

type HarvestIn struct {
	LandIDs []int64 `json:"land_ids"`
	HostGID int64   `json:"host_gid"`
	IsAll   bool    `json:"is_all"`
}

type HarvestOut struct {
	Items    []Item        `json:"items"`
	Lost     []Item        `json:"lost,omitempty"`
	Extra    []Item        `json:"extra,omitempty"`
	Lands    []Land        `json:"lands,omitempty"`
	Limits   []OpLimit     `json:"limits,omitempty"`
	Drops    []LandDrop    `json:"drops,omitempty"`
	Warnings []HarvestWarn `json:"warnings,omitempty"`
	Buffs    []Buff        `json:"buffs,omitempty"`
}

type HarvestWarn struct {
	LandID int64  `json:"land_id"`
	Text   string `json:"text,omitempty"`
}

type OpLimit struct {
	ID         int64 `json:"id"`
	Used       int64 `json:"used"`
	Limit      int64 `json:"limit,omitempty"`
	ShareID    int64 `json:"share_id,omitempty"`
	ExpUsed    int64 `json:"exp_used,omitempty"`
	ExpLimit   int64 `json:"exp_limit,omitempty"`
	ExpShareID int64 `json:"exp_share_id,omitempty"`
}

type LandDrop struct {
	LandID  int64   `json:"land_id"`
	Rewards []Item  `json:"rewards,omitempty"`
	Costs   []Item  `json:"costs,omitempty"`
	Skills  []int64 `json:"skills,omitempty"`
}

type FarmSocial struct {
	ItemID  int64 `json:"item_id"`
	OwnerID int64 `json:"owner_id,omitempty"`
	PutTime int64 `json:"put_time,omitempty"`
}

type LandSocial struct {
	LandID int64  `json:"land_id"`
	Items  []Item `json:"items,omitempty"`
}

type FarmReward struct {
	SourceItemID int64  `json:"source_item_id,omitempty"`
	Items        []Item `json:"items,omitempty"`
}

type LandOpOut struct {
	Lands   []Land       `json:"lands"`
	Limits  []OpLimit    `json:"limits,omitempty"`
	Drops   []LandDrop   `json:"drops,omitempty"`
	Costs   []Item       `json:"costs,omitempty"`
	Social  []LandSocial `json:"social,omitempty"`
	Events  []FarmSocial `json:"events,omitempty"`
	Rewards []FarmReward `json:"rewards,omitempty"`
}

type ItemOpOut struct {
	Used        []Item `json:"used,omitempty"`
	Items       []Item `json:"items,omitempty"`
	Compensated []Item `json:"compensated,omitempty"`
}

type FriendsOut struct {
	Friends          []Friend  `json:"friends"`
	ApplicationCount int64     `json:"application_count,omitempty"`
	Blocked          []Blocked `json:"blocked,omitempty"`
	BlockedBy        []Blocked `json:"blocked_by,omitempty"`
}

type PutSocialIn struct {
	LandIDs []int64 `json:"land_ids"`
	HostGID int64   `json:"host_gid"`
	ItemID  int64   `json:"item_id"`
}

type InviteReportIn struct {
	OpenID   string `json:"open_id"`
	ShareKey string `json:"share_key"`
}

type CanOperateOut struct {
	OK       bool  `json:"ok"`
	StealNum int64 `json:"steal_num,omitempty"`
}

type ActivityBatchIn struct {
	ID    int64         `json:"id"`
	Items []ShopBuyItem `json:"items"`
}

type ShopBuyItem struct {
	GoodsID int64 `json:"goods_id"`
	Count   int64 `json:"count"`
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
	ItemIDs []int64 `json:"item_ids,omitempty"`
}

type FertilizeIn struct {
	LandIDs      []int64 `json:"land_ids"`
	FertilizerID int64   `json:"fertilizer_id"`
}

type BagItem struct {
	ID        int64   `json:"id"`
	Count     int64   `json:"count"`
	UID       int64   `json:"uid,omitempty"`
	SellID    int64   `json:"sell_id,omitempty"`
	SellPrice int64   `json:"sell_price,omitempty"`
	Expire    int64   `json:"expire,omitempty"`
	New       bool    `json:"new,omitempty"`
	Locked    bool    `json:"locked,omitempty"`
	Mutants   []int64 `json:"mutants,omitempty"`
}

type BagOut struct {
	Items []BagItem `json:"items"`
	Max   int64     `json:"max,omitempty"`
	Used  int64     `json:"used,omitempty"`
}

type SellIn struct {
	Items []BagItem `json:"items"`
}

type UseIn struct {
	ID      int64   `json:"id"`
	Count   int64   `json:"count"`
	UID     int64   `json:"uid,omitempty"`
	HostGID int64   `json:"host_gid,omitempty"`
	LandIDs []int64 `json:"land_ids,omitempty"`
}

type Shop struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type int64  `json:"type"`
}

type ShopIn struct {
	ShopID int64 `json:"shop_id"`
}

type ShopCond struct {
	Type  int64 `json:"type"`
	Param int64 `json:"param,omitempty"`
}

type Goods struct {
	ID        int64      `json:"id"`
	Bought    int64      `json:"bought"`
	Price     int64      `json:"price"`
	Limit     int64      `json:"limit"`
	Unlocked  bool       `json:"unlocked"`
	ItemID    int64      `json:"item_id"`
	ItemCount int64      `json:"item_count,omitempty"`
	Conds     []ShopCond `json:"conds,omitempty"`
	Countdown bool       `json:"countdown,omitempty"`
	EndTime   int64      `json:"end_time,omitempty"`
}

type BuyIn struct {
	GoodsID int64 `json:"goods_id"`
	Num     int64 `json:"num"`
	Price   int64 `json:"price"`
}

type BuyOut struct {
	Items    []Item `json:"items"`
	Cost     []Item `json:"cost,omitempty"`
	Goods    *Goods `json:"goods,omitempty"`
	Success  bool   `json:"success,omitempty"`
	FirstBuy bool   `json:"first_buy,omitempty"`
	Bought   int64  `json:"bought,omitempty"`
	Limit    int64  `json:"limit,omitempty"`
}

type EnterIn struct {
	GID int64 `json:"gid"`
}

type Visit struct {
	Host       User      `json:"host"`
	Lands      []Land    `json:"lands"`
	DogID      int64     `json:"dog_id,omitempty"`
	DogFoodSec int64     `json:"dog_food_sec,omitempty"`
	ApplyToken string    `json:"apply_token,omitempty"`
	AtHome     bool      `json:"at_home,omitempty"`
	Tourist    bool      `json:"tourist,omitempty"`
	ServerMs   int64     `json:"server_ms,omitempty"`
	Friend     bool      `json:"friend,omitempty"`
	Community  bool      `json:"community,omitempty"`
	Limits     []OpLimit `json:"limits,omitempty"`
	Weather    Weather   `json:"weather,omitempty"`
}

type OpenIDsIn struct {
	OpenIDs []string `json:"open_ids"`
}

type ArkIn struct {
	GID     int64  `json:"gid"`
	OpenID  string `json:"open_id"`
	Scene   string `json:"scene"`
	ShareID int64  `json:"share_id"`
	Key     string `json:"key"`
}

type Blocked struct {
	GID    int64  `json:"gid"`
	Name   string `json:"name,omitempty"`
	OpenID string `json:"open_id,omitempty"`
	Avatar string `json:"avatar_url,omitempty"`
	Level  int64  `json:"level,omitempty"`
	Time   int64  `json:"time,omitempty"`
}

type UIDsIn struct {
	UIDs []int64 `json:"uids"`
}

type LockFail struct {
	UID    int64  `json:"uid"`
	Code   int64  `json:"code,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type LockOut struct {
	Newly    []int64    `json:"newly,omitempty"`
	Already  []int64    `json:"already,omitempty"`
	Failures []LockFail `json:"failures,omitempty"`
}

type CanOperateIn struct {
	HostGID     int64 `json:"host_gid"`
	OperationID int64 `json:"operation_id"`
}

type AutoBuyIn struct {
	ItemID int64 `json:"item_id"`
	Num    int64 `json:"num"`
	ShopID int64 `json:"shop_id"`
}

type MysteryAutoIn struct {
	Currency int64 `json:"currency"`
}

type BuyFail struct {
	ID     int64  `json:"id"`
	Reason string `json:"reason"`
}

type MysteryAutoOut struct {
	Items  []Item    `json:"items"`
	Failed []BuyFail `json:"failed,omitempty"`
}

type DogBuyIn struct {
	ID    int64 `json:"id"`
	Price int64 `json:"price"`
}

type DogBuyOut struct {
	Dog  Dog    `json:"dog"`
	Cost []Item `json:"cost,omitempty"`
}

type ShareInviteInfo struct {
	ID          int64  `json:"id"`
	ShareKey    string `json:"share_key,omitempty"`
	ShareURL    string `json:"share_url,omitempty"`
	InviteCount int64  `json:"invite_count"`
	RewardCount int64  `json:"reward_count"`
	CanClaim    bool   `json:"can_claim"`
}

type ShareInviteOut struct {
	Infos []ShareInviteInfo `json:"infos"`
}

type ShareAwardOut struct {
	Info        ShareInviteInfo `json:"info,omitempty"`
	Awards      []Item          `json:"awards,omitempty"`
	Compensated []Item          `json:"compensated,omitempty"`
	Awarded     bool            `json:"awarded"`
}

type ShareClaimOut struct {
	Success   bool   `json:"success"`
	HasReward bool   `json:"has_reward"`
	Items     []Item `json:"items"`
}

type AlbumLevel struct {
	Level    int64  `json:"level"`
	Need     int64  `json:"need,omitempty"`
	Rewards  []Item `json:"rewards,omitempty"`
	CanClaim bool   `json:"can_claim"`
	Claimed  bool   `json:"claimed"`
}

type AlbumLevels struct {
	Level    int64        `json:"level"`
	Progress int64        `json:"progress,omitempty"`
	Levels   []AlbumLevel `json:"levels"`
	Extra    []AlbumLevel `json:"extra,omitempty"`
}

type Application struct {
	GID    int64  `json:"gid"`
	Time   int64  `json:"time,omitempty"`
	OpenID string `json:"open_id,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar_url,omitempty"`
}

type ApplicationsOut struct {
	Applications []Application `json:"applications"`
	Blocked      bool          `json:"blocked,omitempty"`
}

type AcceptOut struct {
	Friends []Friend `json:"friends"`
	Success int64    `json:"success,omitempty"`
	Failed  int64    `json:"failed,omitempty"`
	Full    int64    `json:"full,omitempty"`
}

type RejectOut struct {
	Count int64 `json:"count"`
}

type TagsIn struct {
	GID    int64 `json:"gid"`
	New    bool  `json:"new"`
	Follow bool  `json:"follow"`
}

type GIDsIn struct {
	GIDs []int64 `json:"gids"`
}

type RefreshIn struct {
	HostGID int64 `json:"host_gid"`
}

type RefreshLandsIn struct {
	LandIDs []int64 `json:"land_ids"`
	HostGID int64   `json:"host_gid"`
}

type CleanSocialIn struct {
	LandIDs []int64 `json:"land_ids"`
	ItemIDs []int64 `json:"item_ids"`
}

type BatchUseIn struct {
	Items []UseIn `json:"items"`
}

type GroupIn struct {
	ID      int64 `json:"id"`
	GroupID int64 `json:"group_id"`
}

func (in GroupIn) Group() int64 {
	if in.GroupID > 0 {
		return in.GroupID
	}
	return in.ID
}

type ActivityGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
	Type int64  `json:"type,omitempty"`
	// Start 组开始时间，Unix 秒。
	Start int64 `json:"start,omitempty"`
	// End 组结束时间，Unix 秒。本地配置没有结束日，以这里为准。
	End        int64      `json:"end,omitempty"`
	Status     int64      `json:"status,omitempty"`
	Activities []Activity `json:"activities"`
}

type HeartbeatOut struct {
	ServerMs int64 `json:"server_ms,omitempty"`
	HostGID  int64 `json:"host_gid,omitempty"`
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
	Unlocked bool   `json:"unlocked"`
	Type     int64  `json:"type,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Rewards  []Item `json:"rewards,omitempty"`
	Share    int64  `json:"share,omitempty"`
	Group    int64  `json:"group,omitempty"`
	Cond     int64  `json:"cond,omitempty"`
	Extra    []Item `json:"extra,omitempty"`
}

type TaskClaimOut struct {
	Items       []Item    `json:"items"`
	Compensated []Item    `json:"compensated,omitempty"`
	Board       TaskBoard `json:"board,omitempty"`
}

type ActiveBox struct {
	ID       int64  `json:"id"`
	Need     int64  `json:"need"`
	Status   int64  `json:"status"`
	CanClaim bool   `json:"can_claim"`
	Claimed  bool   `json:"claimed"`
	Rewards  []Item `json:"rewards,omitempty"`
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
	ID     int64   `json:"id"`
	IDs    []int64 `json:"ids"`
	Shared bool    `json:"shared,omitempty"`
}

type DailyIn struct {
	Type     int64   `json:"type"`
	PointIDs []int64 `json:"point_ids"`
}

type EmailClaimOut struct {
	Items     []Item   `json:"items"`
	Unclaimed []string `json:"unclaimed,omitempty"`
}

type EmailBoxIn struct {
	Box int64  `json:"box"`
	ID  string `json:"id,omitempty"`
}

type EmailIn struct {
	Box int64  `json:"box"`
	ID  string `json:"id"`
}

type EmailIDsIn struct {
	Box int64    `json:"box"`
	IDs []string `json:"ids"`
}

type TaskReportIn struct {
	ID       int64 `json:"id"`
	Progress int64 `json:"progress"`
}

type Email struct {
	ID        string `json:"id"`
	Type      int64  `json:"type,omitempty"`
	Title     string `json:"title,omitempty"`
	Read      bool   `json:"read"`
	HasReward bool   `json:"has_reward"`
	Subtitle  string `json:"subtitle,omitempty"`
	SendTime  int64  `json:"send_time,omitempty"`
	Expire    int64  `json:"expire,omitempty"`
	Tips      string `json:"tips,omitempty"`
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
	ID      int64  `json:"id"`
	GroupID int64  `json:"group_id,omitempty"`
	Type    int64  `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Desc    string `json:"desc,omitempty"`
	// Start 开始时间，Unix 秒。活动号不是开日，以这里为准。
	Start int64 `json:"start,omitempty"`
	// End 结束时间，Unix 秒。本地配置没有结束日，以这里为准。
	End int64 `json:"end,omitempty"`
	// ClientID 客户端展示 ID。采集 103、建设 104、故友 61、青酿日签 71、青酿 72、筑桥 111、香囊 112、寻宝 18。
	ClientID int64 `json:"client_id,omitempty"`
	// Status 1 进行中，2 已结束。值为 0 时 JSON 会省略。
	Status         int64           `json:"status,omitempty"`
	SplashOrder    int64           `json:"splash_order,omitempty"`
	Splashed       bool            `json:"splashed,omitempty"`
	RedDot         bool            `json:"red_dot,omitempty"`
	SigninClaimed  bool            `json:"signin_claimed,omitempty"`
	SigninRewardID int64           `json:"signin_reward_id,omitempty"`
	Signin         []SigninReward  `json:"signin,omitempty"`
	Shop           []ActivityGoods `json:"shop,omitempty"`
	Nodes          []TechNode      `json:"nodes,omitempty"`
	Lottery        *LotteryInfo    `json:"lottery,omitempty"`
	Brew           *BrewInfo       `json:"brew,omitempty"`
	Recall         *RecallInfo     `json:"recall,omitempty"`
	Invite         *InviteInfo     `json:"invite,omitempty"`
	Gift           *GiftState      `json:"gift,omitempty"`
	RandShop       *RandShop       `json:"rand_shop,omitempty"`
	Mega           *MegaEvent      `json:"mega,omitempty"`
	Progress       *ProgressReward `json:"progress,omitempty"`
	Draw           *DrawInfo       `json:"draw,omitempty"`
	Cheer          *CampCheer      `json:"cheer,omitempty"`
	Charity        *CharityFlower  `json:"charity,omitempty"`
	Drops          []DropPreview   `json:"drops,omitempty"`
	Hunt           *Hunt           `json:"hunt,omitempty"`
}

type SigninReward struct {
	ID    int64  `json:"id"`
	Desc  string `json:"desc,omitempty"`
	Items []Item `json:"items,omitempty"`
}

type MegaEvent struct {
	Day      int64         `json:"day,omitempty"`
	Total    int64         `json:"total,omitempty"`
	Lookback int64         `json:"lookback,omitempty"`
	Rewards  []MegaReward  `json:"rewards,omitempty"`
	Events   []MegaSegment `json:"events,omitempty"`
}

type MegaReward struct {
	Day       int64  `json:"day,omitempty"`
	Unlocked  bool   `json:"unlocked"`
	Claimed   bool   `json:"claimed"`
	Claimable bool   `json:"claimable"`
	Rewards   []Item `json:"rewards,omitempty"`
}

type MegaSegment struct {
	Level    int64  `json:"level,omitempty"`
	Unlocked bool   `json:"unlocked"`
	Name     string `json:"name,omitempty"`
	Art      string `json:"art,omitempty"`
	Desc     string `json:"desc,omitempty"`
}

type ProgressReward struct {
	Preview   []Item         `json:"preview,omitempty"`
	Steps     []ProgressStep `json:"steps,omitempty"`
	Unlocked  int64          `json:"unlocked,omitempty"`
	Completed bool           `json:"completed,omitempty"`
}

type ProgressStep struct {
	Step    int64  `json:"step"`
	Cost    []Item `json:"cost,omitempty"`
	Rewards []Item `json:"rewards,omitempty"`
	Status  int64  `json:"status,omitempty"`
}

type ActivityOpOut struct {
	Items     []Item    `json:"items"`
	Costs     []Item    `json:"costs,omitempty"`
	Days      []int64   `json:"days,omitempty"`
	Steps     []int64   `json:"steps,omitempty"`
	Completed bool      `json:"completed,omitempty"`
	Gold      int64     `json:"gold,omitempty"`
	ClaimType int64     `json:"claim_type,omitempty"`
	Unlocked  []int64   `json:"unlocked,omitempty"`
	Activity  *Activity `json:"activity,omitempty"`
}

type RandShop struct {
	Goods     []RandGoods `json:"goods"`
	Next      int64       `json:"next_refresh,omitempty"`
	Cost      int64       `json:"refresh_cost,omitempty"`
	CostID    int64       `json:"refresh_cost_id,omitempty"`
	Diamond   int64       `json:"diamond_cost,omitempty"`
	Limit     int64       `json:"daily_limit,omitempty"`
	Today     int64       `json:"today_count,omitempty"`
	FreeLimit int64       `json:"free_limit,omitempty"`
	FreeToday int64       `json:"free_today,omitempty"`
}

type RandGoods struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Desc      string `json:"desc,omitempty"`
	Items     []Item `json:"items,omitempty"`
	Cost      []Item `json:"cost,omitempty"`
	Limit     int64  `json:"limit,omitempty"`
	Bought    int64  `json:"bought,omitempty"`
	Available bool   `json:"available"`
}

type LotteryInfo struct {
	FreeLeft  int64   `json:"free_left"`
	FreeLimit int64   `json:"free_limit,omitempty"`
	PaidLeft  int64   `json:"paid_left"`
	PaidLimit int64   `json:"paid_limit,omitempty"`
	CostID    int64   `json:"cost_id,omitempty"`
	CostCount int64   `json:"cost_count,omitempty"`
	Diamond   int64   `json:"diamond,omitempty"`
	Preview   []int64 `json:"preview,omitempty"`
	BizType   int64   `json:"biz_type,omitempty"`
}

type BrewInfo struct {
	Value       int64   `json:"value"`
	Step        int64   `json:"step"`
	Steps       int64   `json:"steps,omitempty"`
	Multipliers []int64 `json:"multipliers,omitempty"`
	Amounts     []int64 `json:"amounts,omitempty"`
	CanClaim    bool    `json:"can_claim"`
	Base        int64   `json:"base,omitempty"`
	Allowed     []int64 `json:"allowed,omitempty"`
	Share       int64   `json:"share,omitempty"`
}

type RecallInfo struct {
	Daily         int64  `json:"daily,omitempty"`
	DailyLimit    int64  `json:"daily_limit,omitempty"`
	Total         int64  `json:"total,omitempty"`
	Max           int64  `json:"max,omitempty"`
	Pending       int64  `json:"pending,omitempty"`
	PendingItems  []Item `json:"pending_items,omitempty"`
	Returner      bool   `json:"returner,omitempty"`
	LoginAt       int64  `json:"login_at,omitempty"`
	BuffExpire    int64  `json:"buff_expire,omitempty"`
	BuffCD        int64  `json:"buff_cd,omitempty"`
	ByGID         int64  `json:"by_gid,omitempty"`
	ReturnPending int64  `json:"return_pending,omitempty"`
	ReturnItems   []Item `json:"return_items,omitempty"`
}

type InviteTask struct {
	Stage     int64  `json:"stage,omitempty"`
	Desc      string `json:"desc,omitempty"`
	Target    int64  `json:"target,omitempty"`
	Current   int64  `json:"current,omitempty"`
	Completed bool   `json:"completed,omitempty"`
	Claimed   bool   `json:"claimed,omitempty"`
	Level     int64  `json:"growth_level,omitempty"`
	Rewards   []Item `json:"rewards,omitempty"`
}

type InviteInfo struct {
	Invitee     bool         `json:"invitee,omitempty"`
	InviteCount int64        `json:"invite_count,omitempty"`
	GrowthCount int64        `json:"growth_count,omitempty"`
	Limit       int64        `json:"limit,omitempty"`
	Expire      int64        `json:"expire,omitempty"`
	Inviter     *Invitee     `json:"inviter,omitempty"`
	Invite      []InviteTask `json:"invite,omitempty"`
	Growth      []InviteTask `json:"growth,omitempty"`
	Newcomer    []InviteTask `json:"newcomer,omitempty"`
}

type GiftOffer struct {
	Cost    []Item `json:"cost,omitempty"`
	Receive []Item `json:"receive,omitempty"`
	Type    int64  `json:"type,omitempty"`
	Content int64  `json:"content,omitempty"`
}

type GiftState struct {
	Sent         int64       `json:"sent"`
	SendLimit    int64       `json:"send_limit,omitempty"`
	ReceiveLimit int64       `json:"receive_limit,omitempty"`
	Gifts        []GiftOffer `json:"gifts,omitempty"`
}

type LotteryIn struct {
	ID      int64 `json:"id"`
	HostGID int64 `json:"host_gid"`
	Free    int64 `json:"free"`
	Paid    int64 `json:"paid"`
}

type LotteryHit struct {
	GoodsID   int64  `json:"goods_id,omitempty"`
	Items     []Item `json:"items,omitempty"`
	Quality   int64  `json:"quality,omitempty"`
	Guarantee bool   `json:"guarantee,omitempty"`
}

type LotteryOut struct {
	Items    []Item       `json:"items"`
	Costs    []Item       `json:"costs,omitempty"`
	Partial  bool         `json:"partial,omitempty"`
	Results  []LotteryHit `json:"results,omitempty"`
	Activity *Activity    `json:"activity,omitempty"`
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
	Value    int64     `json:"value"`
	Activity *Activity `json:"activity,omitempty"`
}

type BrewStepOut struct {
	Step       int64     `json:"step"`
	Multiplier int64     `json:"multiplier,omitempty"`
	Amount     int64     `json:"amount,omitempty"`
	Finished   bool      `json:"finished"`
	Activity   *Activity `json:"activity,omitempty"`
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
	ID          int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	Desc        string `json:"desc,omitempty"`
	Limit       int64  `json:"limit,omitempty"`
	Bought      int64  `json:"bought,omitempty"`
	Order       int64  `json:"order,omitempty"`
	Items       []Item `json:"items,omitempty"`
	Cost        []Item `json:"cost,omitempty"`
	Diamond     int64  `json:"diamond,omitempty"`
	Background  int64  `json:"background,omitempty"`
	Restriction int64  `json:"restriction,omitempty"`
	Category    string `json:"category,omitempty"`
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

type DrawInfo struct {
	Today int64  `json:"today,omitempty"`
	Limit int64  `json:"daily_limit,omitempty"`
	Cost  []Item `json:"cost,omitempty"`
	Total int64  `json:"total,omitempty"`
}

type DrawIn struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}

type DrawRecord struct {
	Time    int64  `json:"time,omitempty"`
	Count   int64  `json:"count,omitempty"`
	Rewards []Item `json:"rewards,omitempty"`
}

type DrawHistoryOut struct {
	Records  []DrawRecord `json:"records"`
	Activity *Activity    `json:"activity,omitempty"`
}

type LotteryRecord struct {
	Time      int64        `json:"time,omitempty"`
	Results   []LotteryHit `json:"results,omitempty"`
	CostType  int64        `json:"cost_type,omitempty"`
	CostCount int64        `json:"cost_count,omitempty"`
}

type LotteryHistoryOut struct {
	Records  []LotteryRecord `json:"records"`
	Activity *Activity       `json:"activity,omitempty"`
}

type CampCheer struct {
	CampID int64       `json:"camp_id,omitempty"`
	Cheer  int64       `json:"cheer,omitempty"`
	Tiers  []CheerTier `json:"tiers,omitempty"`
}

type CheerTier struct {
	Index   int64  `json:"index"`
	Need    int64  `json:"need,omitempty"`
	Claimed bool   `json:"claimed"`
	Rewards []Item `json:"rewards,omitempty"`
}

type CheerJoinIn struct {
	ID     int64 `json:"id"`
	CampID int64 `json:"camp_id"`
}

type CheerSubmitIn struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}

type CheerSubmitOut struct {
	Added    int64     `json:"added,omitempty"`
	Cheer    int64     `json:"cheer,omitempty"`
	Progress int64     `json:"progress,omitempty"`
	Activity *Activity `json:"activity,omitempty"`
}

type CheerClaimIn struct {
	ID   int64 `json:"id"`
	Tier int64 `json:"tier"`
}

type RecallPerson struct {
	GID       int64  `json:"gid"`
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Level     int64  `json:"level,omitempty"`
	OpenID    string `json:"open_id,omitempty"`
	Offline   int64  `json:"offline_days,omitempty"`
	LastLogin int64  `json:"last_login,omitempty"`
	RecallAt  int64  `json:"recall_at,omitempty"`
}

type RecallListOut struct {
	List []RecallPerson `json:"list"`
}

type CharityFlower struct {
	LoveID      int64         `json:"love_id,omitempty"`
	LoveCount   int64         `json:"love_count,omitempty"`
	Personal    int64         `json:"personal,omitempty"`
	Global      int64         `json:"global,omitempty"`
	MaxGlobal   int64         `json:"max_global,omitempty"`
	Share       int64         `json:"share_status,omitempty"`
	ShareItems  []Item        `json:"share_reward,omitempty"`
	Tiers       []CharityTier `json:"tiers,omitempty"`
	GlobalTiers []CharityTier `json:"global_tiers,omitempty"`
	CanDonate   bool          `json:"can_donate,omitempty"`
	Agreed      bool          `json:"agreed,omitempty"`
}

type CharityTier struct {
	Need    int64  `json:"need,omitempty"`
	Rewards []Item `json:"rewards,omitempty"`
	Reached bool   `json:"reached,omitempty"`
	Claimed bool   `json:"claimed,omitempty"`
}

type CharityClaimIn struct {
	ID    int64 `json:"id"`
	Score int64 `json:"score"`
}

type CharityAgreeIn struct {
	ID     int64 `json:"id"`
	Agreed bool  `json:"agreed"`
}

type Hunt struct {
	Treasures     []HuntTreasure `json:"treasures,omitempty"`
	DailyPool     []HuntCharm    `json:"charm_daily_pool,omitempty"`
	Equipped      []HuntCharm    `json:"charm_equipped,omitempty"`
	UnreadPlunder bool           `json:"has_unread_plundered_log,omitempty"`
}

type HuntTreasure struct {
	ID       string `json:"id,omitempty"`
	ItemID   int64  `json:"item_id,omitempty"`
	Count    int64  `json:"count,omitempty"`
	Original int64  `json:"original_count,omitempty"`
	Created  int64  `json:"create_at,omitempty"`
	Status   int64  `json:"status,omitempty"`
	StartAt  int64  `json:"start_at,omitempty"`
	EndAt    int64  `json:"end_at,omitempty"`
}

type HuntCharm struct {
	ID    int64 `json:"id"`
	Total int64 `json:"total_use_count,omitempty"`
	Used  int64 `json:"used_count,omitempty"`
}

type HuntEquipIn struct {
	ID       int64   `json:"id"`
	CharmIDs []int64 `json:"charm_ids"`
}

type HuntBattleIn struct {
	ID          int64  `json:"id"`
	DefenderGID int64  `json:"defender_gid"`
	TreasureID  string `json:"treasure_id"`
}

type HuntLogEntry struct {
	Time     int64  `json:"time,omitempty"`
	Attacker int64  `json:"attacker_gid,omitempty"`
	Name     string `json:"attacker_name,omitempty"`
	Won      bool   `json:"attacker_won,omitempty"`
	Lost     []Item `json:"lost_items,omitempty"`
	Used     []Item `json:"attacker_use_item,omitempty"`
	Injected []Item `json:"injected_items,omitempty"`
	Avatar   string `json:"attack_avatar_url,omitempty"`
}

type HuntLogOut struct {
	Logs     []HuntLogEntry `json:"logs"`
	Activity *Activity      `json:"activity,omitempty"`
}

type CharityDonateOut struct {
	Consumed int64     `json:"consumed,omitempty"`
	Added    int64     `json:"added,omitempty"`
	Personal int64     `json:"personal,omitempty"`
	Global   int64     `json:"global,omitempty"`
	Items    []Item    `json:"items,omitempty"`
	Activity *Activity `json:"activity,omitempty"`
}

type CharityXhhOut struct {
	Num        int64     `json:"num,omitempty"`
	Code       string    `json:"code,omitempty"`
	Trans      string    `json:"trans,omitempty"`
	BusinessID string    `json:"business_id,omitempty"`
	Items      []Item    `json:"items,omitempty"`
	Activity   *Activity `json:"activity,omitempty"`
}

type DropPreview struct {
	ID      int64  `json:"id,omitempty"`
	ItemID  int64  `json:"item_id,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Limit   int64  `json:"daily_limit,omitempty"`
	Dropped int64  `json:"dropped,omitempty"`
	Item    *Item  `json:"item,omitempty"`
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

type MysteryBuyOut struct {
	Items []Item        `json:"items"`
	Goods *MysteryGoods `json:"goods,omitempty"`
}

type PassLevel struct {
	Level   int64  `json:"level"`
	Free    []Item `json:"free,omitempty"`
	Premium []Item `json:"premium,omitempty"`
	Key     bool   `json:"key,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

type BattlePass struct {
	ID             int64       `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Desc           string      `json:"desc,omitempty"`
	Level          int64       `json:"level"`
	MaxLevel       int64       `json:"max_level,omitempty"`
	Exp            int64       `json:"exp,omitempty"`
	LevelExp       int64       `json:"level_exp,omitempty"`
	NextExp        int64       `json:"next_exp,omitempty"`
	Premium        bool        `json:"premium"`
	FreeClaimed    int64       `json:"free_claimed"`
	PremiumClaimed int64       `json:"premium_claimed,omitempty"`
	Price          int64       `json:"price,omitempty"`
	Levels         []PassLevel `json:"levels,omitempty"`
}

type SeasonActivity struct {
	ID   int64  `json:"id"`
	Type int64  `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	// Start 开始时间，Unix 秒。
	Start int64 `json:"start,omitempty"`
	// End 结束时间，Unix 秒。
	End int64 `json:"end,omitempty"`
}

type Season struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Phase   int64  `json:"phase,omitempty"`
	Preheat int64  `json:"preheat,omitempty"`
	// Start 赛季开始，Unix 秒。千星游记赛季号 2026072700。礼包 2026-07-29 10:00 ~ 2026-08-27 23:59，以这里为准。
	Start int64 `json:"start,omitempty"`
	// End 赛季结束，Unix 秒。本地配置没有结束日，以这里为准。
	End        int64            `json:"end,omitempty"`
	ServerTime int64            `json:"server_time,omitempty"`
	Activities []SeasonActivity `json:"activities,omitempty"`
	Next       *Season          `json:"next,omitempty"`
	Pass       BattlePass       `json:"pass"`
}

type PassClaimOut struct {
	Items    []Item     `json:"items"`
	Levels   []int64    `json:"levels,omitempty"`
	Pass     BattlePass `json:"pass,omitempty"`
	Overflow bool       `json:"overflow,omitempty"`
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
	Discount  string `json:"discount,omitempty"`
	Pic       string `json:"pic,omitempty"`
	Countdown bool   `json:"countdown,omitempty"`
}

type MallBuyIn struct {
	ID  int64 `json:"id"`
	Num int64 `json:"num"`
}

type MallProfile struct {
	ID   int64 `json:"id"`
	Type int64 `json:"type,omitempty"`
}

type MonthCard struct {
	ID            int64  `json:"id"`
	Claimable     bool   `json:"claimable"`
	Days          int64  `json:"days,omitempty"`
	TotalDays     int64  `json:"total_days,omitempty"`
	ClaimedAmount int64  `json:"claimed_amount,omitempty"`
	ExpireSeconds int64  `json:"expire_seconds,omitempty"`
	TotalCount    int64  `json:"total_count,omitempty"`
	PayID         string `json:"pay_id,omitempty"`
	Rewards       []Item `json:"rewards,omitempty"`
	PurchaseCost  *Item  `json:"purchase_cost,omitempty"`
	Claimable2    bool   `json:"claimable2,omitempty"`
}

type MonthCardClaimOut struct {
	Items []Item    `json:"items"`
	Card  MonthCard `json:"card,omitempty"`
}

type RedPacket struct {
	ID           int64  `json:"id"`
	Claimed      bool   `json:"claimed"`
	Status       int64  `json:"status,omitempty"`
	CanClaim     bool   `json:"can_claim"`
	Name         string `json:"name,omitempty"`
	Desc         string `json:"desc,omitempty"`
	Time         string `json:"time,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type RedPacketClaimOut struct {
	Status int64  `json:"status,omitempty"`
	Items  []Item `json:"items"`
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
	FromType  int64  `json:"from_type,omitempty"`
	Level     int64  `json:"level,omitempty"`
}

type VisitPopup struct {
	Unread    bool  `json:"unread"`
	LastRead  int64 `json:"last_read,omitempty"`
	NeedPopup bool  `json:"need_popup"`
}

type VisitPageIn struct {
	Page int64 `json:"page"`
}

type VisitDeleteIn struct {
	IDs []int64 `json:"ids"`
}

type Invitee struct {
	GID        int64  `json:"gid"`
	Level      int64  `json:"level,omitempty"`
	Name       string `json:"name,omitempty"`
	Avatar     string `json:"avatar_url,omitempty"`
	RegisterAt int64  `json:"register_at,omitempty"`
	InvitedAt  int64  `json:"invited_at,omitempty"`
}

type VisitFriend struct {
	GID           int64  `json:"gid"`
	Name          string `json:"name,omitempty"`
	Avatar        string `json:"avatar_url,omitempty"`
	StealCount    int64  `json:"steal_count,omitempty"`
	StealItemNum  int64  `json:"steal_item_num,omitempty"`
	HelpCount     int64  `json:"help_count,omitempty"`
	MischiefCount int64  `json:"mischief_count,omitempty"`
	Level         int64  `json:"level,omitempty"`
}

type VisitSummary struct {
	StealCount    int64         `json:"steal_count"`
	HelpCount     int64         `json:"help_count"`
	MischiefCount int64         `json:"mischief_count"`
	StealItemNum  int64         `json:"steal_item_num,omitempty"`
	Friends       []VisitFriend `json:"friends,omitempty"`
}

type RedDot struct {
	Red bool `json:"red_dot"`
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
	Layer    int64  `json:"layer,omitempty"`
	New      bool   `json:"new,omitempty"`
	Rewards  []Item `json:"rewards,omitempty"`
}

type Buff struct {
	ID    int64 `json:"id"`
	Value int64 `json:"value,omitempty"`
	Extra int64 `json:"extra,omitempty"`
}

type Album struct {
	Items     []AlbumItem `json:"items"`
	Progress  int64       `json:"progress,omitempty"`
	Level     int64       `json:"level,omitempty"`
	Rewards   []Item      `json:"rewards,omitempty"`
	Rarities  []int64     `json:"rarities,omitempty"`
	Type      int64       `json:"type,omitempty"`
	Next      int64       `json:"next,omitempty"`
	Claimable bool        `json:"claimable"`
	Claimed   bool        `json:"claimed,omitempty"`
	Buffs     []Buff      `json:"buffs,omitempty"`
	NextBuff  *Buff       `json:"next_buff,omitempty"`
}

type AlbumClaimOut struct {
	Items  []Item  `json:"items"`
	Levels []int64 `json:"levels,omitempty"`
	Level  int64   `json:"level,omitempty"`
	Next   []Item  `json:"next,omitempty"`
}

type Dog struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Protect   int64  `json:"protect,omitempty"`
	Owned     bool   `json:"owned"`
	Activated bool   `json:"activated"`
	Price     int64  `json:"price,omitempty"`
	Expire    int64  `json:"expire,omitempty"`
	LostMin   int64  `json:"lost_min,omitempty"`
	LostMax   int64  `json:"lost_max,omitempty"`
}

type DogFood struct {
	ID      int64 `json:"id"`
	Seconds int64 `json:"seconds,omitempty"`
	Count   int64 `json:"count"`
}

type DogSkill struct {
	ID    int64 `json:"id"`
	DogID int64 `json:"dog_id,omitempty"`
	Used  int64 `json:"used"`
	Max   int64 `json:"max,omitempty"`
}

type DogYard struct {
	Dogs     []Dog      `json:"dogs"`
	Deployed int64      `json:"deployed,omitempty"`
	FoodLeft int64      `json:"food_left"`
	FoodMax  int64      `json:"food_max,omitempty"`
	Foods    []DogFood  `json:"foods"`
	NewLog   bool       `json:"new_log,omitempty"`
	Pending  int64      `json:"pending,omitempty"`
	Skills   []DogSkill `json:"skills,omitempty"`
}

type FeedIn struct {
	FoodID int64 `json:"food_id"`
	Count  int64 `json:"count"`
}

type PageIn struct {
	From  int64 `json:"from"`
	Count int64 `json:"count"`
}

type ProtectLog struct {
	Time       int64  `json:"time"`
	GID        int64  `json:"gid"`
	Name       string `json:"name,omitempty"`
	Avatar     string `json:"avatar_url,omitempty"`
	Level      int64  `json:"level,omitempty"`
	Online     bool   `json:"online,omitempty"`
	LastOnline int64  `json:"last_online,omitempty"`
	Read       bool   `json:"read,omitempty"`
	Authorized int64  `json:"authorized,omitempty"`
	DogID      int64  `json:"dog_id,omitempty"`
	DogName    string `json:"dog_name,omitempty"`
	Count      int64  `json:"count,omitempty"`
	Gold       int64  `json:"gold,omitempty"`
	RecordType int64  `json:"record_type,omitempty"`
	SkillID    int64  `json:"skill_id,omitempty"`
	Skill      string `json:"skill,omitempty"`
	Trigger    int64  `json:"trigger,omitempty"`
}

type DogGiftOut struct {
	Items       []Item `json:"items"`
	Compensated []Item `json:"compensated,omitempty"`
	Claimed     int64  `json:"claimed,omitempty"`
	Pending     int64  `json:"pending,omitempty"`
}

type DogLogsOut struct {
	Logs  []ProtectLog `json:"logs"`
	Total int64        `json:"total,omitempty"`
}

type DeployOut struct {
	Deployed  int64 `json:"deployed,omitempty"`
	Previous  int64 `json:"previous,omitempty"`
	Withdrawn int64 `json:"withdrawn,omitempty"`
}

type Bulletin struct {
	ID      int64  `json:"id"`
	Title   string `json:"title,omitempty"`
	Read    bool   `json:"read"`
	Forced  bool   `json:"forced,omitempty"`
	PopType int64  `json:"pop_type,omitempty"`
	Banner  string `json:"banner,omitempty"`
}

type BulletinDetail struct {
	Title        string `json:"title,omitempty"`
	Content      string `json:"content,omitempty"`
	Start        string `json:"start,omitempty"`
	End          string `json:"end,omitempty"`
	Banner       string `json:"banner,omitempty"`
	DetailBanner string `json:"detail_banner,omitempty"`
}

type Mutant struct {
	ID         int64 `json:"id"`
	Start      int64 `json:"start,omitempty"`
	End        int64 `json:"end,omitempty"`
	ActivityID int64 `json:"activity_id,omitempty"`
	Red        bool  `json:"red_point,omitempty"`
}

type Career struct {
	GID        int64  `json:"gid,omitempty"`
	Name       string `json:"name,omitempty"`
	Level      int64  `json:"level,omitempty"`
	Exp        int64  `json:"exp,omitempty"`
	Avatar     string `json:"avatar_url,omitempty"`
	Signature  string `json:"signature,omitempty"`
	Harvested  int64  `json:"harvested"`
	Stolen     int64  `json:"stolen"`
	Remark     string `json:"remark,omitempty"`
	Gender     int64  `json:"gender,omitempty"`
	Authorized int64  `json:"authorized,omitempty"`
	Items      []Item `json:"items,omitempty"`
}

type RankIn struct {
	Type int64 `json:"type"`
	Page int64 `json:"page"`
}

type RankItem struct {
	GID    int64  `json:"gid"`
	Name   string `json:"name,omitempty"`
	Value  int64  `json:"value"`
	Rank   int64  `json:"rank"`
	Avatar string `json:"avatar_url,omitempty"`
	Level  int64  `json:"level,omitempty"`
}

type RankBoard struct {
	Items []RankItem `json:"items"`
	Total int64      `json:"total,omitempty"`
}

type TypeIn struct {
	Type int64 `json:"type"`
}

type Avatar struct {
	ID       int64 `json:"id"`
	Type     int64 `json:"type,omitempty"`
	Count    int64 `json:"count,omitempty"`
	Priority int64 `json:"priority,omitempty"`
	New      bool  `json:"new,omitempty"`
	Expire   int64 `json:"expire,omitempty"`
}

type AvatarEquipIn struct {
	ID  int64 `json:"id"`
	Off bool  `json:"off"`
}

type Skin struct {
	ID       int64 `json:"id"`
	Slot     int64 `json:"slot,omitempty"`
	Equipped bool  `json:"equipped"`
	Expire   int64 `json:"expire,omitempty"`
}

type SkinEquipIn struct {
	Current int64 `json:"current_id"`
	ID      int64 `json:"id"`
}

type DropReward struct {
	ID      int64 `json:"id"`
	Count   int64 `json:"count"`
	Chance  int64 `json:"chance,omitempty"`
	Claimed bool  `json:"claimed"`
}

type Drop struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name,omitempty"`
	Status  int64        `json:"status,omitempty"`
	Start   int64        `json:"start,omitempty"`
	End     int64        `json:"end,omitempty"`
	Dropped int64        `json:"dropped"`
	Limit   int64        `json:"limit,omitempty"`
	Rewards []DropReward `json:"rewards,omitempty"`
}

type SolarTerm struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Status  int64  `json:"status"`
	Start   int64  `json:"start,omitempty"`
	End     int64  `json:"end,omitempty"`
	Rewards []Item `json:"rewards,omitempty"`
}

type SolarClaimOut struct {
	Items []Item    `json:"items"`
	Event SolarTerm `json:"event,omitempty"`
}

type SolarBasic struct {
	ID     int64   `json:"id"`
	Season int64   `json:"season,omitempty"`
	Desc   string  `json:"desc,omitempty"`
	Events []int64 `json:"events,omitempty"`
}

type SolarOut struct {
	Terms      []SolarTerm  `json:"terms"`
	ServerTime int64        `json:"server_time,omitempty"`
	Basic      *SolarBasic  `json:"basic,omitempty"`
	Basics     []SolarBasic `json:"basics,omitempty"`
}

type Relation struct {
	Friend    bool `json:"friend"`
	Community bool `json:"community"`
	Stranger  bool `json:"stranger"`
}

type AchieveIn struct {
	Kind int64 `json:"kind"`
	ID   int64 `json:"id"`
}

type AchieveGoalIn struct {
	Kind   int64 `json:"kind"`
	ID     int64 `json:"id"`
	GoalID int64 `json:"goal_id"`
}

type AchieveGoal struct {
	ID       int64  `json:"id"`
	Cond     int64  `json:"cond,omitempty"`
	Progress int64  `json:"progress"`
	Total    int64  `json:"total,omitempty"`
	Claimed  bool   `json:"claimed"`
	Unlock   int64  `json:"unlock,omitempty"`
	Category int64  `json:"category,omitempty"`
	Exp      int64  `json:"exp,omitempty"`
	Need     []Item `json:"need,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Sort     int64  `json:"sort,omitempty"`
}

type AchieveGoalOut struct {
	Exp    int64        `json:"exp"`
	Before int64        `json:"before,omitempty"`
	After  int64        `json:"after,omitempty"`
	Scope  AchieveScope `json:"scope,omitempty"`
}

type AchieveLevel struct {
	Level   int64  `json:"level"`
	Need    int64  `json:"need,omitempty"`
	Claimed bool   `json:"claimed"`
	Rewards []Item `json:"rewards,omitempty"`
}

type AchieveScope struct {
	Kind    int64          `json:"kind"`
	ID      int64          `json:"id"`
	Level   int64          `json:"level"`
	Exp     int64          `json:"exp"`
	Next    int64          `json:"next,omitempty"`
	Claimed int64          `json:"claimed_level,omitempty"`
	Goals   []AchieveGoal  `json:"goals,omitempty"`
	Levels  []AchieveLevel `json:"levels,omitempty"`
}

type Status struct {
	LoggedIn   bool    `json:"logged_in"`
	User       User    `json:"user,omitempty"`
	ACE        ACE     `json:"ace"`
	HostGID    int64   `json:"host_gid,omitempty"`
	Heartbeats int     `json:"heartbeats,omitempty"`
	ServerMs   int64   `json:"server_ms,omitempty"`
	FirstLogin bool    `json:"first_login,omitempty"`
	Weather    Weather `json:"weather,omitempty"`
}

type ACE struct {
	Uploads   int    `json:"uploads"`
	Reports   int    `json:"status_reports"`
	Failures  int    `json:"failures"`
	LastError string `json:"last_error,omitempty"`
}

type CleanFarmEventOut struct {
	Events  []FarmSocial `json:"events,omitempty"`
	Rewards []FarmReward `json:"rewards,omitempty"`
}

type BlockAppsIn struct {
	Block bool `json:"block"`
}

type BlockAppsOut struct {
	Block bool `json:"block"`
}

type WXRecommendIn struct {
	Encrypted string `json:"encrypted_data"`
}

type WXRecommendPageIn struct {
	Offset   int64 `json:"offset"`
	PageSize int64 `json:"page_size"`
}

type WXRecommendPlayer struct {
	GID     int64  `json:"gid"`
	Name    string `json:"name,omitempty"`
	Avatar  string `json:"avatar_url,omitempty"`
	Level   int64  `json:"level,omitempty"`
	Applied bool   `json:"applied"`
}

type WXRecommendOut struct {
	Players []WXRecommendPlayer `json:"players,omitempty"`
	Total   int64               `json:"total,omitempty"`
	HasMore bool                `json:"has_more"`
}

type WXApplyResult struct {
	GID     int64 `json:"gid"`
	Success bool  `json:"success"`
	Code    int64 `json:"error_code,omitempty"`
}

type WXApplyOut struct {
	Results []WXApplyResult `json:"results,omitempty"`
}

type SkinSetIn struct {
	SkinIDs []int64 `json:"skin_ids"`
}

type SkinSetEffect struct {
	SetID  int64 `json:"set_id"`
	Type   int64 `json:"effect_type,omitempty"`
	Active bool  `json:"active"`
}

type SkinSetEffectIn struct {
	SetID   int64 `json:"set_id"`
	Type    int64 `json:"effect_type"`
	Enabled bool  `json:"enabled"`
	TypeID  int64 `json:"type_id,omitempty"`
	Param   int64 `json:"param,omitempty"`
}

type BuyPassOut struct {
	Success bool       `json:"success"`
	Items   []Item     `json:"items,omitempty"`
	Pass    BattlePass `json:"pass,omitempty"`
}

type CookiesIn struct {
	Cookies string `json:"cookies"`
}

type QQGroup struct {
	OpenID string `json:"openid,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar_url,omitempty"`
	Bound  int64  `json:"bound,omitempty"`
}

type QQAuthGroupsOut struct {
	Groups  []QQGroup `json:"groups,omitempty"`
	Cookies string    `json:"cookies,omitempty"`
}

type QQRecommendIn struct {
	Class   string `json:"class_name"`
	Session string `json:"session_info"`
	Scene   int64  `json:"scene_id"`
}

type QQRecommendGroup struct {
	OpenID  string `json:"openid,omitempty"`
	Name    string `json:"name,omitempty"`
	Avatar  string `json:"avatar_url,omitempty"`
	Auth    string `json:"join_auth,omitempty"`
	Jump    string `json:"jump_schema,omitempty"`
	Bound   int64  `json:"bound,omitempty"`
	Members int64  `json:"members,omitempty"`
}

type QQRecommendOut struct {
	Groups  []QQRecommendGroup `json:"groups,omitempty"`
	Ended   bool               `json:"ended"`
	Pos     int64              `json:"pos,omitempty"`
	Session string             `json:"session_info,omitempty"`
}

type QQBindIn struct {
	CommunityID string `json:"community_id"`
}

type QQCommunity struct {
	OpenID string `json:"openid,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar_url,omitempty"`
}

type QQBindOut struct {
	Community     QQCommunity `json:"community,omitempty"`
	BoundAt       int64       `json:"bound_at,omitempty"`
	RewardClaimed bool        `json:"reward_claimed"`
}

type QQLeaveOut struct {
	QuitLeft int64 `json:"quit_left,omitempty"`
	Cooldown int64 `json:"cooldown_until,omitempty"`
}

type QQCommunityOut struct {
	Community QQCommunity `json:"community,omitempty"`
	BoundAt   int64       `json:"bound_at,omitempty"`
	Members   int64       `json:"member_count,omitempty"`
	Friends   []Friend    `json:"friends,omitempty"`
	HasMore   bool        `json:"has_more"`
}

type QQBindInfoOut struct {
	Community     QQCommunity `json:"community,omitempty"`
	BoundAt       int64       `json:"bound_at,omitempty"`
	Cooldown      int64       `json:"cooldown_until,omitempty"`
	RewardClaimed bool        `json:"reward_claimed"`
	QuitLeft      int64       `json:"quit_left,omitempty"`
	Rewards       []Item      `json:"rewards,omitempty"`
	MaxQuit       int64       `json:"max_quit,omitempty"`
	CooldownDays  int64       `json:"cooldown_days,omitempty"`
}

type QQRevokeIn struct {
	GID         int64  `json:"gid"`
	CommunityID string `json:"community_id"`
}

type UIDIn struct {
	UID int64 `json:"uid"`
}

type GiftTokenOut struct {
	SourceType     int64  `json:"source_type,omitempty"`
	PresentOrder   string `json:"present_order_id,omitempty"`
	RedeemCode     string `json:"redeem_code,omitempty"`
	PlatformURL    string `json:"platform_url,omitempty"`
	DisplayName    string `json:"display_name,omitempty"`
	PackageInfo    string `json:"package_info,omitempty"`
	OutBillNo      string `json:"out_bill_no,omitempty"`
	TransferAmount int64  `json:"transfer_amount,omitempty"`
	MchID          string `json:"mch_id,omitempty"`
}

type GiftHistoryIn struct {
	SourceType int64 `json:"source_type"`
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
}

type GiftClaimRecord struct {
	ClaimID        string `json:"claim_id,omitempty"`
	ItemID         int64  `json:"item_id,omitempty"`
	SourceType     int64  `json:"source_type,omitempty"`
	Status         int64  `json:"status,omitempty"`
	Assigned       int64  `json:"assigned_time,omitempty"`
	Expire         int64  `json:"expire_time,omitempty"`
	Title          string `json:"title,omitempty"`
	Subtitle       string `json:"subtitle,omitempty"`
	DetailURL      string `json:"detail_url,omitempty"`
	Code           string `json:"code,omitempty"`
	PresentOrderID string `json:"present_order_id,omitempty"`
	WXActivityID   string `json:"wx_activity_id,omitempty"`
}

type GiftHistoryOut struct {
	Records []GiftClaimRecord `json:"records,omitempty"`
	Total   int64             `json:"total,omitempty"`
}

type TransferIn struct {
	OutBillNo string `json:"out_bill_no"`
}

type TransferOut struct {
	State      int64  `json:"state"`
	FailReason string `json:"fail_reason,omitempty"`
}

type FollowGiftOut struct {
	Followed bool `json:"followed"`
	Claimed  bool `json:"claimed"`
	RedDot   bool `json:"red_dot"`
}

type FollowGiftIn struct {
	Followed bool `json:"followed"`
}

type RechargeRange struct {
	Min   int64 `json:"min"`
	Max   int64 `json:"max"`
	Ratio int64 `json:"ratio"`
}

type RechargeBonusOut struct {
	Active bool            `json:"active"`
	Start  int64           `json:"start,omitempty"`
	End    int64           `json:"end,omitempty"`
	Unlock int64           `json:"unlock,omitempty"`
	Ranges []RechargeRange `json:"ranges,omitempty"`
}

type RechargeDataOut struct {
	Recharged int64 `json:"recharged,omitempty"`
	Returned  int64 `json:"returned,omitempty"`
}

type DisplayIn struct {
	Name      string `json:"name"`
	Avatar    string `json:"avatar_url"`
	Signature string `json:"signature"`
	Gender    int64  `json:"gender"`
	Remark    string `json:"remark"`
}

type DisplayOut struct {
	Name       string `json:"name,omitempty"`
	Avatar     string `json:"avatar_url,omitempty"`
	Signature  string `json:"signature,omitempty"`
	Gender     int64  `json:"gender,omitempty"`
	Remark     string `json:"remark,omitempty"`
	Authorized int64  `json:"authorized,omitempty"`
}

type SettingsKeysIn struct {
	Keys []int64 `json:"keys,omitempty"`
}

type UserSettings struct {
	DisableNudge          bool `json:"disable_nudge"`
	DisableMonthCard      bool `json:"disable_month_card"`
	DisableQQSubscribe    bool `json:"disable_qq_subscribe"`
	DisableWXRecommend    bool `json:"disable_wx_recommend"`
	DisableOfflineSummary bool `json:"disable_offline_summary"`
	AllowArkVisit         bool `json:"allow_ark_visit"`
}

type DeleteAccountIn struct {
	Name     string `json:"name"`
	CertID   string `json:"cert_id"`
	CertType int64  `json:"cert_type"`
}

type DeleteAccountOut struct {
	Success   bool   `json:"success"`
	Msg       string `json:"msg,omitempty"`
	RequestAt int64  `json:"request_time,omitempty"`
	DeleteAt  int64  `json:"delete_time,omitempty"`
}

type DecryptIn struct {
	Encrypted string `json:"encrypted_data"`
}

type QQAuthIn struct {
	Authorized bool `json:"authorized"`
}

type QQAuthOut struct {
	Authorized int64 `json:"authorized"`
}

type ReportFlowIn struct {
	OSType      int64  `json:"os_type,omitempty"`
	PlatType    int64  `json:"plat_type,omitempty"`
	OpenID      string `json:"open_id,omitempty"`
	GID         int64  `json:"gid,omitempty"`
	Name        string `json:"name,omitempty"`
	Now         int64  `json:"now,omitempty"`
	Level       int64  `json:"level,omitempty"`
	FlowType    int64  `json:"flow_type,omitempty"`
	FlowTypeStr string `json:"flow_type_str,omitempty"`
	Int1        int64  `json:"param_int1,omitempty"`
	Int2        int64  `json:"param_int2,omitempty"`
	Int3        int64  `json:"param_int3,omitempty"`
	Int4        int64  `json:"param_int4,omitempty"`
	Int5        int64  `json:"param_int5,omitempty"`
	Str6        string `json:"param_str6,omitempty"`
	Str7        string `json:"param_str7,omitempty"`
	Str8        string `json:"param_str8,omitempty"`
	Str9        string `json:"param_str9,omitempty"`
	Str10       string `json:"param_str10,omitempty"`
}

type BatchReportFlowIn struct {
	Flows []ReportFlowIn `json:"flows"`
}

type BatchReportFlowOut struct {
	Success int64 `json:"success_count"`
	Fail    int64 `json:"fail_count"`
}

type ReportUserIn struct {
	GID        int64    `json:"gid"`
	Category   int64    `json:"category"`
	Reasons    []int64  `json:"reasons,omitempty"`
	Scene      int64    `json:"scene,omitempty"`
	Desc       string   `json:"desc,omitempty"`
	Content    string   `json:"content,omitempty"`
	Pics       []string `json:"pics,omitempty"`
	Videos     []string `json:"videos,omitempty"`
	Voices     []string `json:"voices,omitempty"`
	GroupID    string   `json:"group_id,omitempty"`
	GroupName  string   `json:"group_name,omitempty"`
	BattleID   string   `json:"battle_id,omitempty"`
	BattleTime int64    `json:"battle_time,omitempty"`
	Entrance   int64    `json:"entrance,omitempty"`
	MsgBoardID string   `json:"msg_board_id,omitempty"`
}

type ReportUserOut struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}

type QQVipDailyOut struct {
	IsQQVip      bool   `json:"is_qq_vip"`
	CanClaim     bool   `json:"can_claim"`
	ClaimedToday bool   `json:"claimed_today"`
	Rewards      []Item `json:"rewards,omitempty"`
}

type QQVipClaimDailyIn struct {
	ConfigID int64 `json:"config_id"`
}

type QQVipClaimDailyOut struct {
	Rewards []Item `json:"rewards,omitempty"`
}

type QQVipRefreshOut struct {
	IsQQVip  bool  `json:"is_qq_vip"`
	VIPLevel int64 `json:"vip_level,omitempty"`
}

type QQVipConfig struct {
	Type       int64  `json:"type,omitempty"`
	Rewards    []Item `json:"rewards,omitempty"`
	Enable     bool   `json:"is_enable"`
	SeasonID   int64  `json:"season_id,omitempty"`
	ID         int64  `json:"id,omitempty"`
	Multiplier int64  `json:"reward_multiplier,omitempty"`
	Start      int64  `json:"act_start_time,omitempty"`
	End        int64  `json:"act_end_time,omitempty"`
}

type QQVipRewardsStatusOut struct {
	IsQQVip         bool          `json:"is_qq_vip"`
	CanClaim        bool          `json:"can_claim"`
	ClaimedToday    bool          `json:"claimed_today"`
	RewardsCanClaim bool          `json:"rewards_can_claim"`
	Configs         []QQVipConfig `json:"active_configs,omitempty"`
	HasRedpoint     bool          `json:"has_redpoint"`
}

type QQVipClaimRewardsIn struct {
	ConfigIDs []int64 `json:"config_ids"`
}

type QQVipClaimRewardsOut struct {
	SkinIDs  []int64 `json:"granted_skin_item_ids,omitempty"`
	FrameIDs []int64 `json:"granted_avatar_frame_item_ids,omitempty"`
	Rewards  []Item  `json:"rewards,omitempty"`
}

type MarqueeMsg struct {
	UUID         int64  `json:"uuid,omitempty"`
	ConfigID     int64  `json:"config_id,omitempty"`
	Expire       int64  `json:"expire_time,omitempty"`
	Type         int64  `json:"type,omitempty"`
	Content      string `json:"content,omitempty"`
	Priority     int64  `json:"priority,omitempty"`
	DisplayCount int64  `json:"display_count,omitempty"`
}

type MarqueeOut struct {
	Msgs []MarqueeMsg `json:"msgs,omitempty"`
}

type SystemOpenIn struct {
	SystemName int64 `json:"system_name"`
}

type SystemOpenOut struct {
	Unlocked bool `json:"unlocked"`
}

type MutantOpenInfoOut struct {
	Tips    string `json:"tips,omitempty"`
	Rewards []Item `json:"rewards,omitempty"`
}

type QQSubscribeItem struct {
	ID         int64 `json:"id"`
	Subscribed bool  `json:"subscribed"`
}

type QQSubscribeOut struct {
	Status     int64             `json:"status,omitempty"`
	Subscribed bool              `json:"subscribed"`
	Items      []QQSubscribeItem `json:"items,omitempty"`
}

type WXTemplateStatus struct {
	TemplateID string `json:"template_id"`
	Subscribed bool   `json:"subscribed"`
}

type WXSubscribeIn struct {
	Templates []WXTemplateStatus `json:"templates"`
}

type WXSubscribeOut struct {
	Templates []WXTemplateStatus `json:"templates,omitempty"`
}

type ModerateTextIn struct {
	Text   string `json:"text"`
	Reason string `json:"reason,omitempty"`
}

type ModerateTextOut struct {
	Text   string `json:"result_text,omitempty"`
	Dirty  bool   `json:"is_dirty"`
	Reason string `json:"reason,omitempty"`
}

type ModerateTextBatchIn struct {
	Items []ModerateTextIn `json:"text_items"`
}

type ModerateTextBatchOut struct {
	Items []ModerateTextOut `json:"text_items,omitempty"`
}

type ModeratePicIn struct {
	URL    string `json:"pic_url"`
	Reason string `json:"reason,omitempty"`
}

type ModeratePicOut struct {
	URL       string `json:"result_url,omitempty"`
	Dirty     bool   `json:"is_dirty"`
	DirtyType int64  `json:"dirty_type,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type ModeratePicBatchIn struct {
	Items []ModeratePicIn `json:"pic_items"`
}

type ModeratePicBatchOut struct {
	Items []ModeratePicOut `json:"pic_items,omitempty"`
}

// Session is the live game connection used by HTTP handlers.
type Session interface {
	Login(ctx context.Context, in LoginIn) (User, error)
	Info() (User, error)
	Status() (Status, error)
	Heartbeat(ctx context.Context, in RefreshIn) (HeartbeatOut, error)
	Logout(ctx context.Context) error
	Brief(ctx context.Context, in EnterIn) (User, error)
	BatchInfo(ctx context.Context, in GIDsIn) ([]User, error)
	ArkClick(ctx context.Context, in ArkIn) error
	Refresh(ctx context.Context, in RefreshIn) (LandOpOut, error)
	RefreshLands(ctx context.Context, in RefreshLandsIn) ([]Land, error)
	CleanSocial(ctx context.Context, in CleanSocialIn) (LandOpOut, error)
	PutInsects(ctx context.Context, in LandOpIn) (LandOpOut, error)
	PutWeeds(ctx context.Context, in LandOpIn) (LandOpOut, error)
	PutSocial(ctx context.Context, in PutSocialIn) (LandOpOut, error)
	CanOperate(ctx context.Context, in CanOperateIn) (CanOperateOut, error)
	Harvest(ctx context.Context, in HarvestIn) (HarvestOut, error)
	Steal(ctx context.Context, in HarvestIn) (HarvestOut, error)
	Plant(ctx context.Context, in PlantIn) (LandOpOut, error)
	Remove(ctx context.Context, in RemoveIn) (LandOpOut, error)
	Unlock(ctx context.Context, in LandIDIn) (Land, error)
	Upgrade(ctx context.Context, in LandIDIn) (Land, error)
	Farming(ctx context.Context, in LandOpIn) (LandOpOut, error)
	Water(ctx context.Context, in LandOpIn) (LandOpOut, error)
	Weed(ctx context.Context, in LandOpIn) (LandOpOut, error)
	Bug(ctx context.Context, in LandOpIn) (LandOpOut, error)
	Fertilize(ctx context.Context, in FertilizeIn) (LandOpOut, error)
	Friends(ctx context.Context) (FriendsOut, error)
	Applications(ctx context.Context) (ApplicationsOut, error)
	Accept(ctx context.Context, in GIDsIn) (AcceptOut, error)
	Reject(ctx context.Context, in GIDsIn) (RejectOut, error)
	SetTags(ctx context.Context, in TagsIn) (Friend, error)
	DeleteFriend(ctx context.Context, in EnterIn) error
	SyncFriends(ctx context.Context, in OpenIDsIn) ([]Friend, error)
	GameFriends(ctx context.Context, in GIDsIn) ([]Friend, error)
	BlockFriend(ctx context.Context, in EnterIn) (Blocked, error)
	UnblockFriend(ctx context.Context, in EnterIn) error
	BlockList(ctx context.Context) ([]Blocked, error)
	ShareKey(ctx context.Context) (string, error)
	Help(ctx context.Context, in HelpIn) error
	Enter(ctx context.Context, in EnterIn) (Visit, error)
	Leave(ctx context.Context, in EnterIn) error
	ShareCheck(ctx context.Context) (bool, error)
	ShareClaim(ctx context.Context) (ShareClaimOut, error)
	InviteInfo(ctx context.Context, in IDIn) (ShareInviteOut, error)
	InviteAward(ctx context.Context, in IDIn) (ShareAwardOut, error)
	PosterShown(ctx context.Context, in IDIn) (bool, error)
	ReportInvite(ctx context.Context, in InviteReportIn) (bool, error)
	Bag(ctx context.Context) (BagOut, error)
	Sell(ctx context.Context, in SellIn) (ItemOpOut, error)
	Use(ctx context.Context, in UseIn) (ItemOpOut, error)
	BatchUse(ctx context.Context, in BatchUseIn) (ItemOpOut, error)
	CancelNew(ctx context.Context, in IDIn) (int64, error)
	LockItems(ctx context.Context, in UIDsIn) (LockOut, error)
	UnlockItems(ctx context.Context, in UIDsIn) (LockOut, error)
	Shops(ctx context.Context) ([]Shop, error)
	Goods(ctx context.Context, in ShopIn) ([]Goods, error)
	Buy(ctx context.Context, in BuyIn) (BuyOut, error)
	AutoBuy(ctx context.Context, in AutoBuyIn) (BuyOut, error)
	Weather(ctx context.Context) (Weather, error)
	CurrentWeather(ctx context.Context) (Weather, error)
	TodayWeather(ctx context.Context) ([]Weather, error)
	Tasks(ctx context.Context) (TaskBoard, error)
	ClaimTask(ctx context.Context, in ClaimIn) (TaskClaimOut, error)
	ClaimDaily(ctx context.Context, in DailyIn) (TaskClaimOut, error)
	ReportTask(ctx context.Context, in TaskReportIn) (TaskBoard, error)
	Emails(ctx context.Context, in EmailBoxIn) ([]Email, error)
	ReadEmail(ctx context.Context, in EmailIn) (EmailDetail, error)
	ClaimEmail(ctx context.Context, in EmailIn) ([]Item, error)
	ClaimAllEmail(ctx context.Context, in EmailBoxIn) (EmailClaimOut, error)
	BatchReadEmail(ctx context.Context, in EmailIDsIn) error
	BatchDeleteEmail(ctx context.Context, in EmailIDsIn) error
	Activities(ctx context.Context) ([]Activity, error)
	ActivityGroup(ctx context.Context, in GroupIn) (ActivityGroup, error)
	SetSplashed(ctx context.Context, in IDIn) (bool, error)
	Invitees(ctx context.Context, in IDIn) ([]Invitee, error)
	Signin(ctx context.Context, in ActivityOpIn) (ActivityOpOut, error)
	ClaimProgress(ctx context.Context, in ActivityOpIn) (ActivityOpOut, error)
	ShopBuy(ctx context.Context, in ActivityShopIn) (ActivityOpOut, error)
	ShopBatchBuy(ctx context.Context, in ActivityBatchIn) (ActivityOpOut, error)
	RandBuy(ctx context.Context, in ActivityShopIn) (ActivityOpOut, error)
	RandRefresh(ctx context.Context, in IDIn) (ActivityOpOut, error)
	ClaimMega(ctx context.Context, in ActivityOpIn) (ActivityOpOut, error)
	TechSubmit(ctx context.Context, in TechIn) (ActivityOpOut, error)
	Draw(ctx context.Context, in DrawIn) (ActivityOpOut, error)
	DrawHistory(ctx context.Context, in IDIn) (DrawHistoryOut, error)
	MarkViewed(ctx context.Context, in IDIn) (ActivityOpOut, error)
	RandBatchBuy(ctx context.Context, in ActivityBatchIn) (ActivityOpOut, error)
	LotteryHistory(ctx context.Context, in IDIn) (LotteryHistoryOut, error)
	CheerJoin(ctx context.Context, in CheerJoinIn) (ActivityOpOut, error)
	CheerSubmit(ctx context.Context, in CheerSubmitIn) (CheerSubmitOut, error)
	CheerClaim(ctx context.Context, in CheerClaimIn) (ActivityOpOut, error)
	Recallable(ctx context.Context, in IDIn) (RecallListOut, error)
	Recalled(ctx context.Context, in IDIn) (RecallListOut, error)
	CharityShare(ctx context.Context, in IDIn) (ActivityOpOut, error)
	CharityDonate(ctx context.Context, in IDIn) (CharityDonateOut, error)
	CharityClaim(ctx context.Context, in CharityClaimIn) (ActivityOpOut, error)
	CharityXhh(ctx context.Context, in IDIn) (CharityXhhOut, error)
	CharityAgree(ctx context.Context, in CharityAgreeIn) (ActivityOpOut, error)
	HuntFinishCG(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntGuide(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntFeed(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntDraw(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntLog(ctx context.Context, in IDIn) (HuntLogOut, error)
	HuntClaimStory(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntClaimSeed(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntRefreshCharm(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntEquip(ctx context.Context, in HuntEquipIn) (ActivityOpOut, error)
	HuntBattle(ctx context.Context, in HuntBattleIn) (ActivityOpOut, error)
	HuntPlunderedLog(ctx context.Context, in IDIn) (HuntLogOut, error)
	HuntOpen(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntEscort(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntCompensate(ctx context.Context, in IDIn) (ActivityOpOut, error)
	HuntFriendInfo(ctx context.Context, in IDIn) (ActivityOpOut, error)
	Lottery(ctx context.Context, in LotteryIn) (LotteryOut, error)
	BrewStart(ctx context.Context, in BrewStartIn) (BrewStartOut, error)
	BrewStep(ctx context.Context, in ActivityOpIn) (BrewStepOut, error)
	BrewClaim(ctx context.Context, in BrewClaimIn) (ActivityOpOut, error)
	ClaimRecall(ctx context.Context, in ActivityOpIn) (ActivityOpOut, error)
	ClaimReturn(ctx context.Context, in ActivityOpIn) (ActivityOpOut, error)
	ClaimInvite(ctx context.Context, in InviteIn) (ActivityOpOut, error)
	ClaimNewcomer(ctx context.Context, in ActivityOpIn) (ActivityOpOut, error)
	SendGift(ctx context.Context, in GiftIn) (int64, error)
	Mystery(ctx context.Context) (MysteryShop, error)
	MysteryBuy(ctx context.Context, in IDIn) (MysteryBuyOut, error)
	MysteryAutoBuy(ctx context.Context, in MysteryAutoIn) (MysteryAutoOut, error)
	MysteryLeave(ctx context.Context) error
	Season(ctx context.Context) (Season, error)
	ClaimPass(ctx context.Context) (PassClaimOut, error)
	Mall(ctx context.Context, in MallIn) ([]Product, error)
	MallDiamonds(ctx context.Context) ([]Product, error)
	MallProfiles(ctx context.Context) ([]MallProfile, error)
	MallBuy(ctx context.Context, in MallBuyIn) (BuyOut, error)
	MonthCards(ctx context.Context) ([]MonthCard, error)
	ClaimMonthCard(ctx context.Context, in IDIn) (MonthCardClaimOut, error)
	RedPackets(ctx context.Context) ([]RedPacket, error)
	ClaimRedPacket(ctx context.Context, in IDIn) (RedPacketClaimOut, error)
	ClaimAllRedPackets(ctx context.Context) ([]Item, error)
	VisitLogs(ctx context.Context) ([]VisitLog, error)
	VisitSummary(ctx context.Context) (VisitSummary, error)
	VisitPopup(ctx context.Context) (VisitPopup, error)
	VisitPage(ctx context.Context, in VisitPageIn) ([]VisitLog, error)
	DismissVisit(ctx context.Context, in EnterIn) error
	DeleteVisit(ctx context.Context, in VisitDeleteIn) error
	Album(ctx context.Context, in AlbumIn) (Album, error)
	AlbumLevels(ctx context.Context, in AlbumIn) (AlbumLevels, error)
	ClaimAlbum(ctx context.Context, in AlbumIn) (AlbumClaimOut, error)
	MarkAlbum(ctx context.Context, in AlbumIn) error
	Dog(ctx context.Context) (DogYard, error)
	Feed(ctx context.Context, in FeedIn) (int64, error)
	ClaimDogGifts(ctx context.Context) (DogGiftOut, error)
	DogLogs(ctx context.Context, in PageIn) (DogLogsOut, error)
	DeployDog(ctx context.Context, in IDIn) (DeployOut, error)
	WithdrawDog(ctx context.Context) (DeployOut, error)
	ActivateDog(ctx context.Context, in IDIn) (Dog, error)
	BuyDog(ctx context.Context, in DogBuyIn) (DogBuyOut, error)
	Bulletins(ctx context.Context, in PageIn) ([]Bulletin, error)
	ReadBulletin(ctx context.Context, in IDIn) (BulletinDetail, error)
	Mutants(ctx context.Context) ([]Mutant, error)
	Career(ctx context.Context, in EnterIn) (Career, error)
	Ranks(ctx context.Context, in RankIn) (RankBoard, error)
	Avatars(ctx context.Context, in TypeIn) ([]Avatar, error)
	EquippedAvatars(ctx context.Context) ([]Avatar, error)
	EquipAvatar(ctx context.Context, in AvatarEquipIn) (Avatar, error)
	MarkAvatar(ctx context.Context, in IDIn) error
	Skins(ctx context.Context) ([]Skin, error)
	EquippedSkins(ctx context.Context) ([]Skin, error)
	EquipSkin(ctx context.Context, in SkinEquipIn) error
	MarkSkin(ctx context.Context, in IDIn) error
	Drops(ctx context.Context) ([]Drop, error)
	SolarTerms(ctx context.Context) (SolarOut, error)
	SolarRedDot(ctx context.Context) (RedDot, error)
	ClaimSolar(ctx context.Context, in IDIn) (SolarClaimOut, error)
	ClaimAllSolar(ctx context.Context) ([]Item, error)
	AchieveView(ctx context.Context, in AchieveIn) (AchieveScope, error)
	ClaimAchieveGoal(ctx context.Context, in AchieveGoalIn) (AchieveGoalOut, error)
	ClaimAchieveLevel(ctx context.Context, in AchieveIn) (ItemOpOut, error)
	CleanFarmEvents(ctx context.Context) (CleanFarmEventOut, error)
	BlockApplications(ctx context.Context, in BlockAppsIn) (BlockAppsOut, error)
	WXRecommend(ctx context.Context, in WXRecommendIn) (WXRecommendOut, error)
	WXRecommendPage(ctx context.Context, in WXRecommendPageIn) (WXRecommendOut, error)
	ApplyWXFriends(ctx context.Context, in GIDsIn) (WXApplyOut, error)
	EquipSkinSet(ctx context.Context, in SkinSetIn) error
	SetSkinSetEffect(ctx context.Context, in SkinSetEffectIn) error
	SkinSets(ctx context.Context) ([]SkinSetEffect, error)
	BuyPass(ctx context.Context) (BuyPassOut, error)
	MarkSeasonOpening(ctx context.Context) (bool, error)
	QQAuthGroups(ctx context.Context, in CookiesIn) (QQAuthGroupsOut, error)
	QQRecommendGroups(ctx context.Context, in QQRecommendIn) (QQRecommendOut, error)
	QQBind(ctx context.Context, in QQBindIn) (QQBindOut, error)
	QQLeave(ctx context.Context) (QQLeaveOut, error)
	QQCommunity(ctx context.Context, in PageIn) (QQCommunityOut, error)
	QQBindInfo(ctx context.Context) (QQBindInfoOut, error)
	QQClaimReward(ctx context.Context) ([]Item, error)
	QQRevokeAuth(ctx context.Context, in QQRevokeIn) error
	UseGiftToken(ctx context.Context, in UIDIn) (GiftTokenOut, error)
	GiftHistory(ctx context.Context, in GiftHistoryIn) (GiftHistoryOut, error)
	TransferStatus(ctx context.Context, in TransferIn) (TransferOut, error)
	CancelTransfer(ctx context.Context, in UIDIn) (TransferOut, error)
	FollowGiftStatus(ctx context.Context) (FollowGiftOut, error)
	SetFollowGift(ctx context.Context, in FollowGiftIn) error
	ClaimFollowGift(ctx context.Context) ([]Item, error)
	RechargeBonus(ctx context.Context) (RechargeBonusOut, error)
	RechargeBonusData(ctx context.Context) (RechargeDataOut, error)
	SetDisplay(ctx context.Context, in DisplayIn) (DisplayOut, error)
	GetSettings(ctx context.Context, in SettingsKeysIn) (UserSettings, error)
	SetSettings(ctx context.Context, in UserSettings) (UserSettings, error)
	DeleteAccount(ctx context.Context, in DeleteAccountIn) (DeleteAccountOut, error)
	DecryptOpenData(ctx context.Context, in DecryptIn) (string, error)
	SetQQRecommendAuth(ctx context.Context, in QQAuthIn) (QQAuthOut, error)
	ReportFlow(ctx context.Context, in ReportFlowIn) error
	BatchReportFlow(ctx context.Context, in BatchReportFlowIn) (BatchReportFlowOut, error)
	ReportUser(ctx context.Context, in ReportUserIn) (ReportUserOut, error)
	QQVipDailyStatus(ctx context.Context) (QQVipDailyOut, error)
	QQVipClaimDaily(ctx context.Context, in QQVipClaimDailyIn) (QQVipClaimDailyOut, error)
	QQVipRefresh(ctx context.Context) (QQVipRefreshOut, error)
	QQVipClaimRewards(ctx context.Context, in QQVipClaimRewardsIn) (QQVipClaimRewardsOut, error)
	QQVipRewardsStatus(ctx context.Context) (QQVipRewardsStatusOut, error)
	QQVipMarkRedpoint(ctx context.Context) error
	Marquee(ctx context.Context) (MarqueeOut, error)
	SystemUnlocked(ctx context.Context, in SystemOpenIn) (SystemOpenOut, error)
	MutantOpenInfo(ctx context.Context) (MutantOpenInfoOut, error)
	QQSubscribe(ctx context.Context) (QQSubscribeOut, error)
	WXSubscribe(ctx context.Context) (WXSubscribeOut, error)
	SetWXSubscribe(ctx context.Context, in WXSubscribeIn) (WXSubscribeOut, error)
	ModerateText(ctx context.Context, in ModerateTextIn) (ModerateTextOut, error)
	BatchModerateText(ctx context.Context, in ModerateTextBatchIn) (ModerateTextBatchOut, error)
	ModeratePic(ctx context.Context, in ModeratePicIn) (ModeratePicOut, error)
	BatchModeratePic(ctx context.Context, in ModeratePicBatchIn) (ModeratePicBatchOut, error)
	Close() error
}

type Factory func() Session
