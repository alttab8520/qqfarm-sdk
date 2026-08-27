package game

// 活动号不是开日。精确到几点、以及结束时间，以 Activity/List 该条 start / end（Unix 秒）为准。
// status：1 进行中，2 已结束。
const (
	// 雨落成诗。活动号 2026070301。雷雨礼包 2026-08-26 10:00 ~ 2026-09-08 23:59。
	// 精确看带 lottery / shop / nodes 那条的 start / end。client_id：采集 103，建设 104。
	ActWeather = 2026070301

	// 故友重逢。仅活动号，开停未核实。看带 recall 那条的 start / end。client_id 61。
	ActRecall = 2026070101

	// 千星游记 / 观星礼录。活动号 2026072701。千星礼包 2026-07-29 10:00 ~ 2026-08-27 23:59。
	// 精确看带 mega 那条的 start / end，赛季看 Season/Info。
	ActMega = 2026072701

	// 千星游记赛季号。礼包窗同 ActMega。精确看 Season.start / Season.end。
	ActSeasonQianxing = 2026072700

	// 青酿换万金（GotoJump 一期）。仅活动号，开停未核实。看带 brew / signin 那条的 start / end。
	// client_id：日签 71，酿酒 72。
	ActQingmei = 2026080102

	// 青酿换万金（GotoJump 二期）。仅活动号，开停未核实。看带 brew 那条的 start / end。
	ActQingmeiLater = 2026081202

	// 邀新领红包。仅活动号，开停未核实。看带 invite 那条的 start / end。
	ActInvite = 2026080701

	// 鹊桥寄情·筑桥。鹊桥礼包 2026-08-18 10:00 ~ 2026-08-22 23:59。
	// 精确看带 progress 那条的 start / end。client_id 111。
	ActQixiBridge = 2026081801

	// 鹊桥寄情·赠香囊。礼包窗同 ActQixiBridge。精确看带 gift 那条的 start / end。client_id 112。
	ActQixiGift = 2026081802

	// 宠物寻宝。仅活动号，开停未核实。看带 hunt 那条的 start / end。client_id 18。
	ActHunt = 2026090101

	// 公益小红花。仅活动号，开停未核实。看带 charity 那条的 start / end。
	ActCharity = 2026090901

	// 粽香大比拼（阵营加油）。仅活动号，开停未核实。看带 cheer 那条的 start / end。
	ActCheerZongzi = 2026061901

	// 足球狂欢季（阵营加油）。仅活动号，开停未核实。看带 cheer 那条的 start / end。
	ActCheerFootball = 2026071501
)

const (
	ClientWeatherCollect = 103
	ClientWeatherTech    = 104
	ClientRecall         = 61
	ClientQingmeiDaily   = 71
	ClientQingmeiBrew    = 72
	ClientQixiProgress   = 111
	ClientQixiGift       = 112
	ClientHunt           = 18
)
