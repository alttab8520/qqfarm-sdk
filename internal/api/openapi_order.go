package api

import (
	"bytes"
	"encoding/json"
	"sort"
)

type orderedObject struct {
	keys []string
	vals map[string]any
}

func (o orderedObject) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range o.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		vb, err := json.Marshal(o.vals[k])
		if err != nil {
			return nil, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		buf.Write(vb)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func tag(name, desc string) map[string]any {
	return map[string]any{"name": name, "description": desc}
}

func openAPITags() []any {
	return []any{
		tag("System 系统", "探活、登录状态、ACE 计数"),
		tag("YYB 应用宝", "微信扫码拿小程序 code，再登录网关。也可删本地账号、拉资料、取手机号"),
		tag("Resource 资源表", "官方配表。把回包里的裸 ID 翻成名字"),
		tag("User 用户", "登录、资料、设置、心跳、举报"),
		tag("Farm 农场", "地块、种植、务农、使坏"),
		tag("Weather 天气", "实时天气。雨落成诗采集要看 can_collect"),
		tag("Friend 好友", "列表、进家、申请、黑名单、微信推荐"),
		tag("Visit 来访", "来往记录和弹窗"),
		tag("Share 分享", "分享奖和邀请关系"),
		tag("Bag 背包", "格子、使用、出售"),
		tag("Shop 商店", "种子店等常驻商店"),
		tag("Mystery 神秘商人", "限时货架"),
		tag("Mall 商城", "点券、档位、月卡"),
		tag("Task 任务", "成长、每日、活跃度"),
		tag("Email 邮件", "读信、领奖"),
		tag("Activity 清单", "活动列表和开屏。每条有 start/end"),
		tag("Activity 雨落成诗", "活动号 2026070301。礼包 2026-08-26 10:00 ~ 2026-09-08 23:59。采集 / 每日限购 / 建设"),
		tag("Activity 青酿换万金", "2026080102 / 2026081202。仅活动号，开停看 List"),
		tag("Activity 鹊桥寄情", "2026081801 筑桥，2026081802 香囊。礼包 2026-08-18 10:00 ~ 2026-08-22 23:59"),
		tag("Activity 邀新领红包", "2026080701。仅活动号，开停看 List"),
		tag("Activity 故友重逢", "2026070101。仅活动号，开停看 List"),
		tag("Activity 观星礼录", "2026072701。礼包 2026-07-29 10:00 ~ 2026-08-27 23:59。赛季见 Season"),
		tag("Activity 宠物寻宝", "2026090101。仅活动号，开停看 List。cmd 901–916"),
		tag("Activity 公益小红花", "2026090901。仅活动号，开停看 List"),
		tag("Activity 阵营加油", "粽香 2026061901，足球 2026071501。仅活动号，开停看 List"),
		tag("Activity 抽奖签到店", "抽奖、随机店。开停看该条 start/end"),
		tag("Season 赛季", "千星游记 2026072700。礼包 2026-07-29 10:00 ~ 2026-08-27 23:59"),
		tag("QQGroup 群", "QQ 群授权、绑定社群"),
		tag("QQVip 会员", "QQ 会员每日礼和会员奖励"),
		tag("Marquee 跑马灯", "滚动公告"),
		tag("SystemOpen 系统解锁", "查某个系统有没有开"),
		tag("Subscribe 订阅", "QQ / 微信订阅消息"),
		tag("Moderate 审核", "文本和图片内容审核"),
		tag("Gift 礼包", "礼包码道具、领取历史、转账"),
		tag("Follow 关注有礼", "关注状态和领奖"),
		tag("Recharge 充值返利", "配置和进度。不含支付下单"),
		tag("RedPacket 红包", "每日红包"),
		tag("Album 图鉴", "普通 / 超变图鉴"),
		tag("Dog 护主犬", "狗、狗粮、护院"),
		tag("Bulletin 公告", "公告列表和正文"),
		tag("Mutant 变异", "变异活动窗口"),
		tag("Drop 掉落", "随机掉落活动"),
		tag("Solar 节气", "节气领取"),
		tag("Career 生涯", "收获 / 被偷统计"),
		tag("Rank 排行", "等级 / 金币 / 图鉴"),
		tag("Avatar 头像框", "穿戴头像框"),
		tag("Skin 皮肤", "穿戴皮肤和套装"),
		tag("Achieve 成就", "成就域和等级"),
	}
}

func orderedPaths(m map[string]any) orderedObject {
	order := []string{
		"/System/Ping", "/System/Status",
		"/YYB/Accounts", "/YYB/QR", "/YYB/Image", "/YYB/Poll", "/YYB/Confirm",
		"/YYB/Refresh", "/YYB/Code", "/YYB/Login", "/YYB/Delete", "/YYB/Profile",
		"/YYB/Phone", "/YYB/WXData",
		"/Resource/Tables", "/Resource/Lookup", "/Resource/Items", "/Resource/Refresh",
		"/User/Login", "/User/GetInfo", "/User/Brief", "/User/BatchInfo",
		"/User/SetDisplay", "/User/Settings", "/User/SetSettings",
		"/User/DeleteAccount", "/User/DecryptOpenData", "/User/QQRecommendAuth",
		"/User/ReportFlow", "/User/BatchReportFlow", "/User/Report",
		"/User/Heartbeat", "/User/ArkClick", "/User/Logout",
		"/Farm/Refresh", "/Farm/RefreshLands", "/Farm/CanOperate",
		"/Farm/Plant", "/Farm/Harvest", "/Farm/Remove", "/Farm/Unlock", "/Farm/Upgrade",
		"/Farm/Farming", "/Farm/Water", "/Farm/Weed", "/Farm/Bug", "/Farm/Fertilize",
		"/Farm/Steal", "/Farm/PutInsects", "/Farm/PutWeeds", "/Farm/PutSocial", "/Farm/CleanSocial", "/Farm/CleanEvents",
		"/Weather/Status", "/Weather/Current", "/Weather/Today",
		"/Friend/GetList", "/Friend/Enter", "/Friend/Leave", "/Friend/Help",
		"/Friend/Applications", "/Friend/Accept", "/Friend/Reject", "/Friend/Delete",
		"/Friend/SetTags", "/Friend/Sync", "/Friend/GetGameFriends",
		"/Friend/Block", "/Friend/Unblock", "/Friend/BlockList", "/Friend/BlockApplications",
		"/Friend/WXRecommend", "/Friend/WXRecommendPage", "/Friend/ApplyWX", "/Friend/ShareKey",
		"/Visit/Logs", "/Visit/Page", "/Visit/Summary", "/Visit/Popup", "/Visit/Dismiss", "/Visit/Delete",
		"/Share/Check", "/Share/Claim", "/Share/InviteInfo", "/Share/InviteAward",
		"/Share/PosterShown", "/Share/ReportInvite",
		"/Bag/Get", "/Bag/Use", "/Bag/BatchUse", "/Bag/Sell", "/Bag/Lock", "/Bag/Unlock", "/Bag/CancelNew",
		"/Shop/List", "/Shop/Goods", "/Shop/Buy", "/Shop/AutoBuy",
		"/Mystery/Status", "/Mystery/Buy", "/Mystery/AutoBuy", "/Mystery/Leave",
		"/Mall/Profiles", "/Mall/Diamonds", "/Mall/List", "/Mall/Buy", "/Mall/MonthCard", "/Mall/ClaimMonthCard",
		"/Task/List", "/Task/Claim", "/Task/ClaimDaily", "/Task/Report",
		"/Email/List", "/Email/Read", "/Email/Claim", "/Email/ClaimAll", "/Email/BatchRead", "/Email/BatchDelete",
		"/Activity/List", "/Activity/GetGroup", "/Activity/SetSplashed", "/Activity/MarkViewed",
		"/Activity/Lottery", "/Activity/LotteryHistory", "/Activity/ShopBuy", "/Activity/ShopBatchBuy", "/Activity/TechSubmit",
		"/Activity/BrewStart", "/Activity/BrewStep", "/Activity/BrewClaim", "/Activity/Signin",
		"/Activity/ClaimProgress", "/Activity/SendGift",
		"/Activity/Invitees", "/Activity/ClaimInvite", "/Activity/ClaimNewcomer",
		"/Activity/Recallable", "/Activity/Recalled", "/Activity/ClaimRecall", "/Activity/ClaimReturn",
		"/Activity/ClaimMega",
		"/Activity/HuntFinishCG", "/Activity/HuntGuide", "/Activity/HuntFeed", "/Activity/HuntDraw",
		"/Activity/HuntLog", "/Activity/HuntClaimStory", "/Activity/HuntClaimSeed",
		"/Activity/HuntRefreshCharm", "/Activity/HuntEquip", "/Activity/HuntBattle",
		"/Activity/HuntPlunderedLog", "/Activity/HuntOpen", "/Activity/HuntEscort",
		"/Activity/HuntCompensate", "/Activity/HuntFriendInfo",
		"/Activity/CharityShare", "/Activity/CharityDonate", "/Activity/CharityClaim",
		"/Activity/CharityXhh", "/Activity/CharityAgree",
		"/Activity/CheerJoin", "/Activity/CheerSubmit", "/Activity/CheerClaim",
		"/Activity/Draw", "/Activity/DrawHistory",
		"/Activity/RandBuy", "/Activity/RandRefresh", "/Activity/RandBatchBuy",
		"/Season/Info", "/Season/ClaimPass", "/Season/BuyPass", "/Season/MarkOpening",
		"/QQGroup/AuthGroups", "/QQGroup/Recommend", "/QQGroup/Bind", "/QQGroup/Leave",
		"/QQGroup/Community", "/QQGroup/BindInfo", "/QQGroup/ClaimReward", "/QQGroup/RevokeAuth",
		"/QQVip/DailyStatus", "/QQVip/ClaimDaily", "/QQVip/Refresh", "/QQVip/ClaimRewards",
		"/QQVip/RewardsStatus", "/QQVip/MarkRedpoint",
		"/Marquee/List",
		"/SystemOpen/Unlocked",
		"/Subscribe/QQ", "/Subscribe/WX", "/Subscribe/SetWX",
		"/Moderate/Text", "/Moderate/BatchText", "/Moderate/Pic", "/Moderate/BatchPic",
		"/Gift/UseToken", "/Gift/History", "/Gift/TransferStatus", "/Gift/CancelTransfer",
		"/Follow/Status", "/Follow/Set", "/Follow/Claim",
		"/Recharge/Config", "/Recharge/Data",
		"/RedPacket/List", "/RedPacket/Claim", "/RedPacket/ClaimAll",
		"/Album/List", "/Album/Levels", "/Album/Claim", "/Album/MarkViewed",
		"/Dog/Info", "/Dog/Feed", "/Dog/ClaimGifts", "/Dog/Logs",
		"/Dog/Deploy", "/Dog/Withdraw", "/Dog/Activate", "/Dog/Buy",
		"/Bulletin/List", "/Bulletin/Read",
		"/Mutant/List", "/Mutant/OpenInfo", "/Drop/List",
		"/Solar/List", "/Solar/Claim", "/Solar/ClaimAll", "/Solar/RedDot",
		"/Career/Info", "/Rank/List",
		"/Avatar/Owned", "/Avatar/Equipped", "/Avatar/Equip", "/Avatar/MarkViewed",
		"/Skin/Owned", "/Skin/Equipped", "/Skin/Equip", "/Skin/EquipSet", "/Skin/SetEffect", "/Skin/Sets", "/Skin/MarkViewed",
		"/Achieve/View", "/Achieve/ClaimGoal", "/Achieve/ClaimLevel",
	}
	seen := make(map[string]bool, len(order))
	keys := make([]string, 0, len(m))
	for _, k := range order {
		if _, ok := m[k]; ok {
			keys = append(keys, k)
			seen[k] = true
		}
	}
	var extra []string
	for k := range m {
		if !seen[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)
	return orderedObject{keys: append(keys, extra...), vals: m}
}
