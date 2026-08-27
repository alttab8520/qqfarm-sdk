# qqfarm-sdk

非官方 QQ 农场协议 HTTP SDK。Go 实现， `POST` 调用。

启动后打开 `/docs`，每个接口都能展开请求体和回包字段。接口是 `POST /模块/动作`，只有 `GET /YYB/Image` 例外。

回包永远是：

```json
{ "code": 0, "msg": "ok", "data": {} }
```

| `code` | 含义 |
|---|---|
| `0` | 成功，看 `data` |
| `400` | 入参不对 |
| `401` | 未登录 |
| `502` | 游戏网关失败，`msg` 是原因 |

`data` 因接口而异。常见字段：

| 字段 | 含义 |
|---|---|
| `gid` | 游戏内用户 ID |
| `host_gid` | 目标农场主人。不填或 `0` 是自己 |
| `land_ids` | 地块编号 |
| `items[].id` / `count` | 物品类型和数量。`1001` 金币，`1002` 点券，`1005` 金豆豆（雨落成诗每日限购货币） |
| `uid` | 背包格子，使用/出售/酿酒用它，不是物品类型 |
| `5001` | 天气采集瓶。雨落成诗采集消耗；每日限购和建设领取也是它 |
| `5005` | 青蛙使坏瓶。别人放在你家是农场级事件。放走 `Farm/PutSocial`（不带 `land_ids`），清走 `Farm/CleanEvents`（或 `Farm/Farming` 带 `item_ids:[5005]`） |
| `5006` | 乌云使坏瓶。别人放在地块作物上。放走 `Farm/PutSocial`（必须带地），清走 `Farm/CleanSocial`，`item_ids:[5006]`，`land_ids` 从 `Farm/Refresh` 的 `social_items` 取 |

登录后会连游戏网关，接口走真实 RPC。请求体需要官方加密运行时：默认用仓库里的 `data/tsdk-v3.9.0.wasm`，也可设置 `FARM_WASM` 指向其他路径。文件必须是 160975 字节，SHA256 `0001ab8d68ac35309dbdaa16310aeebf0fd578976a91562db12257786e3a2e54`（微信包 59，自报 `v3.9.0.1787640498`）。

登录成功后会自动跑 ACE 长会话，调用方不用自己发反作弊包。官方 `AceService` 只有 `AntiData`，包体是运行时吐出的不透明字节，组不出合法 JSON，所以没有 `/Ace/*`。进度看 `POST /System/Status` 的 `ace`。

## 范围

按模块分。活动再按玩法分。`/docs` 左侧也是这套分类。

### 系统

探活、登录状态、ACE 计数。

| 路径 | 说明 |
|---|---|
| `POST /System/Ping` | 探活 |
| `POST /System/Status` | 登录状态、ACE 计数、当前在谁家、登录时的天气 |

### ACE

网关要求每个会话带着腾讯 ACE（TSDK）运行时。`User/Login` 或 `YYB/Login` 成功后自动开始，登出时停。没有手动上报接口。

真正出网的只有 `Ace.AntiData`。调度器每 5 秒问运行时有没有待发帧：有就发出去，把下行回灌；空队列不发包。运行时里还会按固定节奏做心跳、变速检测、状态上报，这些动作本身不另开 RPC，产生的数据仍走后面的 `AntiData`。

| 节奏 | 做什么 |
|---|---|
| 每 5 秒 | 问运行时要帧，有才发 `Ace.AntiData` |
| 每 25 秒 | 运行时心跳。和 `User/Heartbeat` 不是一条 |
| 每 30 秒 | 变速检测 |
| 每 150 秒 | 状态上报 |

`System/Status` 的 `ace`：

| 字段 | 含义 |
|---|---|
| `uploads` | `AntiData` 成功次数。登录后过几轮应大于 0 |
| `status_reports` | 状态上报次数。大约 150 秒后才有 |
| `failures` | 连续失败次数 |
| `last_error` | 最近一次错误 |


### 应用宝

微信扫码拿农场小程序 `code`，再登录网关。账号存在 `yyb_data/db/yyb.db`，目录可用 `FARM_YYB_DIR` 改。代理用 `FARM_YYB_PROXY`。

| 路径 | 说明 |
|---|---|
| `POST /YYB/Accounts` | 已保存的扫码账号。带 `nickname` / `avatar` |
| `POST /YYB/QR` | 出登录二维码。回 `session_id` 和 `image`（data URL） |
| `GET /YYB/Image?session_id=` | 同一张码的 JPEG |
| `POST /YYB/Poll` | 轮询，body: `{"session_id"}`。`pending` 没扫，`scanned` 已扫未确认，`authorized` 可以 Confirm。`cancelled` / `expired` / `unknown` 是终态 |
| `POST /YYB/Confirm` | 确认并保存账号，body: `{"session_id"}`。会拉应用宝昵称和头像 URL |
| `POST /YYB/Refresh` | 刷新登录态，body 可带 `{"ref"}`。`ref` 是 id / openid / uin，空则第一条 |
| `POST /YYB/Code` | 取一次性 `code`。过期会先刷新再重试 |
| `POST /YYB/Login` | 取 code 并登录网关。之后和 `User/Login` 一样 |
| `POST /YYB/Delete` | 删本地已保存账号，body: `{"ref"}`。`ref` 必填。不是游戏删号 |
| `POST /YYB/Profile` | 再拉昵称和头像，body 可带 `{"ref"}` |
| `POST /YYB/Phone` | 取微信手机号，body 可带 `{"ref","app_id"}`。回包在 `result` |
| `POST /YYB/WXData` | 微信云函数代理，body: `{"ref","app_id","payload"}`。`payload` 必填，通常带 `api_name` |

第一次登录：出码 → 轮询到 `authorized` → 确认入库 → 取 code 登网关。`YYB/Login` 把后两步合成一步（已有账号时直接用）。

```text
curl -X POST http://127.0.0.1:8765/YYB/QR
# data.session_id 和 data.image。也可以 GET /YYB/Image?session_id=...

curl -X POST http://127.0.0.1:8765/YYB/Poll -d "{\"session_id\":\"SID\"}"
# pending / scanned / authorized。等到 authorized

curl -X POST http://127.0.0.1:8765/YYB/Confirm -d "{\"session_id\":\"SID\"}"
# data.id / openid / nickname / avatar

curl -X POST http://127.0.0.1:8765/YYB/Login -d "{\"ref\":\"1\"}"
# 或先 YYB/Code 再 User/Login
```

### 用户

登录、资料、设置、心跳、举报。

| 路径 | 说明 |
|---|---|
| `POST /User/Login` | 登录，body: `{"code","open_id"}`。`code` 可从 `YYB/Code` 拿 |
| `POST /User/GetInfo` | 自己资料 |
| `POST /User/Brief` | 查别人简要资料，body: `{"gid"}` |
| `POST /User/BatchInfo` | 按 gid 批量查，body: `{"gids"}` |
| `POST /User/SetDisplay` | 改展示资料，body: `{"name","avatar_url","signature","gender","remark"}` |
| `POST /User/Settings` | 读用户设置。body 可带 `{"keys"}`，空则全读 |
| `POST /User/SetSettings` | 写用户设置。六个开关：戳一戳 / 月卡 / QQ 订阅 / 微信推荐 / 离线摘要 / ARK 来访 |
| `POST /User/DeleteAccount` | 申请删号，body: `{"name","cert_id","cert_type"}` |
| `POST /User/DecryptOpenData` | 解密开放数据，body: `{"encrypted_data"}`。回 `json_data` |
| `POST /User/QQRecommendAuth` | QQ 好友推荐授权，body: `{"authorized"}` |
| `POST /User/ReportFlow` | 上报客户端流水。不是举报 |
| `POST /User/BatchReportFlow` | 批量上报流水，body: `{"flows"}` |
| `POST /User/Report` | 举报玩家，body: `{"gid","category","reasons"}`。不是流水 |
| `POST /User/Heartbeat` | 网关心跳，body 可带 `{"host_gid"}`。登录后也会每 25 秒自动发 |
| `POST /User/ArkClick` | ARK 点击上报，body: `{"gid","open_id","key","share_id"}`。回包为空 |
| `POST /User/Logout` | 发网关登出再断开 |

### 农场

地块、种植、务农、使坏。

| 路径 | 说明 |
|---|---|
| `POST /Farm/Refresh` | 全量刷新地块，body 可带 `{"host_gid"}`。回包 `lands` / `limits` / `events` |
| `POST /Farm/RefreshLands` | 增量刷新，body: `{"land_ids","host_gid"}`。`land_ids` 空则不发包 |
| `POST /Farm/CanOperate` | 预检，body: `{"host_gid","operation_id"}`。回包 `ok` / `steal_num`。`operation_id`：10001 务农 10004 偷菜 10005 放虫 10006 放草 |
| `POST /Farm/Plant` | 种植，body: `{"seed_id","land_ids"}`。回包 `lands` / `limits` / `drops` |
| `POST /Farm/Harvest` | 收获，body: `{"land_ids","host_gid","is_all"}`。回包 `items` / `lost` / `extra` / `lands` / `limits` / `drops` / `warnings` |
| `POST /Farm/Remove` | 铲地，body: `{"land_ids"}`。回包同种植 |
| `POST /Farm/Unlock` | 解锁地块，body: `{"land_id"}` |
| `POST /Farm/Upgrade` | 升级地块，body: `{"land_id"}` |
| `POST /Farm/Farming` | 一键浇水除草杀虫。body: `{"land_ids","host_gid","item_ids"}`。清别人放的青蛙可带 `"item_ids":[5005]`。乌云不走这里，走 `Farm/CleanSocial` |
| `POST /Farm/Water` | 浇水，body: `{"land_ids","host_gid"}`。回包 `lands` / `limits` / `drops` |
| `POST /Farm/Weed` | 除草，body: `{"land_ids","host_gid"}`。回包同浇水 |
| `POST /Farm/Bug` | 杀虫，body: `{"land_ids","host_gid"}`。回包同浇水 |
| `POST /Farm/Fertilize` | 施肥，body: `{"land_ids","fertilizer_id"}`。回包同浇水，`costs` 是肥料 |
| `POST /Farm/Steal` | 偷菜，body 同收获，必须填好友 `host_gid`。回包同收获 |
| `POST /Farm/PutInsects` | 放虫，body: `{"land_ids","host_gid"}`。必须填好友。回包同浇水 |
| `POST /Farm/PutWeeds` | 放草，body 同放虫 |
| `POST /Farm/PutSocial` | 往别人家放使坏瓶，body: `{"host_gid","item_id","land_ids"}`。青蛙 `5005` 不带地，乌云 `5006` 必须带地。回包 `lands` / `drops` / `rewards` |
| `POST /Farm/CleanSocial` | 清别人放在地块上的乌云、黄金虫。body: `{"land_ids","item_ids"}`。清乌云填 `"item_ids":[5006]`。`land_ids` 必填，从 `Refresh` 里 `social_items` 带 `item_id=5006` 的地取。青蛙不走这里 |
| `POST /Farm/CleanEvents` | 清别人放在你家的青蛙。空请求，内部走 `Farming` 带 `item_ids:[5005]`。回 `rewards`。乌云不走这里 |

### 天气

实时天气。雨落成诗采集要看 `can_collect`。

| 路径 | 说明 |
|---|---|
| `POST /Weather/Status` | 农场实时天气。`can_collect=true` 才能做雨落成诗采集 |
| `POST /Weather/Current` | 此刻天气，带名字和起止时间 |
| `POST /Weather/Today` | 今日天气 |

### 好友

列表、进家、申请、黑名单、微信推荐。

| 路径 | 说明 |
|---|---|
| `POST /Friend/GetList` | 好友列表，含干旱 / 杂草 / 虫害 / 可偷数量，以及 `application_count` / 黑名单 / 新好友 / 关注 |
| `POST /Friend/Enter` | 进好友家，body: `{"gid"}`。进别人家带 `reason=2`。回地块、护主犬、操作次数、天气、在不在家 |
| `POST /Friend/Leave` | 离开好友家，body: `{"gid"}` |
| `POST /Friend/Help` | 进好友家浇水后离开，body: `{"gid","land_ids"}`。进家带 `reason=2` |
| `POST /Friend/Applications` | 好友申请。回 `applications` / `blocked` |
| `POST /Friend/Accept` | 同意，body: `{"gids"}`。回新好友和成功/失败/已满数 |
| `POST /Friend/Reject` | 拒绝，body: `{"gids"}`。回拒绝成功数 |
| `POST /Friend/Delete` | 删除好友，body: `{"gid"}` |
| `POST /Friend/SetTags` | 设置标签。协议请求无字段，body 可带但不会发出去 |
| `POST /Friend/Sync` | 同步微信好友，body: `{"open_ids"}` |
| `POST /Friend/GetGameFriends` | 按 gid 批量查，body: `{"gids"}` |
| `POST /Friend/Block` | 拉黑，body: `{"gid"}` |
| `POST /Friend/Unblock` | 取消拉黑，body: `{"gid"}` |
| `POST /Friend/BlockList` | 黑名单 |
| `POST /Friend/BlockApplications` | 拒收好友申请，body: `{"block"}`。不是拉黑 |
| `POST /Friend/WXRecommend` | 微信推荐好友。body 可带 `{"encrypted_data"}` |
| `POST /Friend/WXRecommendPage` | 微信推荐分页，body: `{"offset","page_size"}`。`page_size` 不填当 `20` |
| `POST /Friend/ApplyWX` | 微信推荐批量加，body: `{"gids"}`。不是申请加好友 |
| `POST /Friend/ShareKey` | 自己的分享 key |

### 来访

来往记录和弹窗。

| 路径 | 说明 |
|---|---|
| `POST /Visit/Logs` | 来访记录。`action`：1浇水 2除虫 3除草 4偷菜 5放虫 6放草。`from_type` 是来源类型 |
| `POST /Visit/Page` | 分页来访，body: `{"page"}`。不填为 `1` |
| `POST /Visit/Summary` | 来往摘要：总偷菜 / 帮忙 / 捣乱和按好友明细 |
| `POST /Visit/Popup` | 来访弹窗：未读 / 要不要弹 |
| `POST /Visit/Dismiss` | 关掉弹窗，body: `{"gid"}` |
| `POST /Visit/Delete` | 删记录，body: `{"ids"}` |

### 分享

分享奖和邀请关系。

| 路径 | 说明 |
|---|---|
| `POST /Share/Check` | 今日能不能领分享奖 |
| `POST /Share/Claim` | 上报分享并领奖。回包 `success` / `has_reward` / `items` |
| `POST /Share/InviteInfo` | 邀请进度。body 可带 `{"id"}`，`id` 是 share_cfg_id |
| `POST /Share/InviteAward` | 领邀请奖，body: `{"id"}`。`id` 是 share_cfg_id。回 `info` / `awards` / `awarded` |
| `POST /Share/PosterShown` | 标记海报已展示，body: `{"id"}`。`id` 是活动 ID |
| `POST /Share/ReportInvite` | 上报邀请关系，body: `{"open_id","share_key"}` |

### 背包

格子、使用、出售。

| 路径 | 说明 |
|---|---|
| `POST /Bag/Get` | 背包。格子带 `uid` / `expire` / `locked` / 卖价货币，另有 `max` / `used` |
| `POST /Bag/Use` | 使用，body: `{"id","count","uid"}`。回包 `used` 消耗、`items` 获得、`compensated` 补偿 |
| `POST /Bag/BatchUse` | 批量使用，body: `{"items":[{"id","count","uid"}]}`。回包同使用 |
| `POST /Bag/Sell` | 出售，body: `{"items":[{"id","count","uid"}]}`。`uid` 可从 `Bag/Get` 带上。回包 `used` 卖出、`items` 货币 |
| `POST /Bag/Lock` | 锁格子，body: `{"uids"}` |
| `POST /Bag/Unlock` | 解锁格子，body: `{"uids"}` |
| `POST /Bag/CancelNew` | 清新标记，body: `{"id"}` |

### 商店

常驻商店。

| 路径 | 说明 |
|---|---|
| `POST /Shop/List` | 商店列表 |
| `POST /Shop/Goods` | 商品，body: `{"shop_id"}`。不填则种子店 `2` |
| `POST /Shop/Buy` | 购买，body: `{"goods_id","num","price"}`。`price` 用货架上的单价。回包另有更新后的 `goods` |
| `POST /Shop/AutoBuy` | 按物品查货架再买，body: `{"item_id","num","shop_id"}`。不填 `shop_id` 则种子店。超限购会拒绝或截到剩余次数 |

### 神秘商人

限时货架。

| 路径 | 说明 |
|---|---|
| `POST /Mystery/Status` | 神秘商人是否在、货架 |
| `POST /Mystery/Buy` | 买一件，body: `{"id"}`。自己指定。回 `items` 和更新后的 `goods` |
| `POST /Mystery/AutoBuy` | 买当前货架未购商品。body 可带 `{"currency"}` 过滤。失败项在 `failed`，不会半路停 |
| `POST /Mystery/Leave` | 打发商人 |

### 商城

点券、档位、月卡。

| 路径 | 说明 |
|---|---|
| `POST /Mall/Profiles` | 商城档位 |
| `POST /Mall/Diamonds` | 点券货架 |
| `POST /Mall/List` | 商城，body: `{"slot"}`。不填为 `1` |
| `POST /Mall/Buy` | 买/领商城商品，body: `{"id","num"}`。自己指定。回包 `success` / `items`，不是商店那种 `cost` |
| `POST /Mall/MonthCard` | 月卡。带回剩余天数、过期秒数、购买消耗 |
| `POST /Mall/ClaimMonthCard` | 领月卡，body: `{"id"}`。回奖励和更新后的月卡 |

### 任务

成长、每日、活跃度。

| 路径 | 说明 |
|---|---|
| `POST /Task/List` | 成长 / 每日任务和活跃度宝箱 |
| `POST /Task/Claim` | 领任务，body: `{"id"}` 或 `{"ids"}`。回 `items` / `compensated` / `board` |
| `POST /Task/ClaimDaily` | 领活跃度宝箱，body: `{"point_ids"}` |
| `POST /Task/Report` | 上报任务进度，body: `{"id","progress"}` |

### 邮件

读信、领奖。

| 路径 | 说明 |
|---|---|
| `POST /Email/List` | 邮件，body: `{"box"}`。不填为 `1`，常见还有 `2`。列表 `read` 是已读，已领看详情 |
| `POST /Email/Read` | 读邮件，body: `{"box","id"}` |
| `POST /Email/Claim` | 领一封，body: `{"box","id"}` |
| `POST /Email/ClaimAll` | 批量领。不填 `box` 则 `1` 和 `2` 都领。可带 `id` 领单封。没领成的在 `unclaimed` |
| `POST /Email/BatchRead` | 批量已读，body: `{"box","ids"}` |
| `POST /Email/BatchDelete` | 批量删除，body: `{"box","ids"}` |

### 活动 · 清单

先打 List。每条有 `start` / `end`。档期见下表。

| 路径 | 说明 |
|---|---|
| `POST /Activity/List` | 活动清单。每条有 `start` / `end`（Unix 秒）和 `status`（`1` 进行中，`2` 已结束）。活动号不是开日。礼包核实了雨落成诗、鹊桥、千星，其余看该条 `start` / `end`。档期见下表 |
| `POST /Activity/GetGroup` | 活动组详情，body: `{"id"}` 或 `{"group_id"}`。List 里没有每日限购货架或建设节点时打这个。窗口看组的 `start` / `end` |
| `POST /Activity/SetSplashed` | 标记活动开屏已看，body: `{"id"}`。窗口看该条 `start` / `end` |
| `POST /Activity/MarkViewed` | 标记活动已看，body: `{"id"}`。开停看该条 `start` / `end` |

### 活动 · 雨落成诗

活动号 `2026070301`。礼包 2026-08-26 10:00 ~ 2026-09-08 23:59。精确看 List。采集 / 每日限购 / 建设。

| 路径 | 说明 |
|---|---|
| `POST /Activity/Lottery` | 雨落成诗采集。body: `{"id","host_gid","free","paid"}`。消耗天气采集瓶 `5001`。`host_gid` 必须是 `can_collect` 的农场，自己指定，不自动巡访。建设/领取/每日购买不是这个接口。活动号 `2026070301`。礼包 2026-08-26 10:00 ~ 2026-09-08 23:59，精确看带 `lottery` 那条的 `start` / `end` |
| `POST /Activity/LotteryHistory` | 雨落成诗采集历史，body: `{"id"}`。活动号 `2026070301`。礼包 2026-08-26 10:00 ~ 2026-09-08 23:59，精确看带 `lottery` 那条的 `start` / `end` |
| `POST /Activity/ShopBuy` | 活动商店买一件。雨落成诗每日限购：金豆豆 `1005` 买采集瓶 `5001`，`goods_id` 从同组 `shop` 取。回包 `items` / `costs` / `activity`。活动号 `2026070301`。礼包 2026-08-26 10:00 ~ 2026-09-08 23:59，精确看带 `shop` 那条的 `start` / `end` |
| `POST /Activity/ShopBatchBuy` | 活动商店批量买，body: `{"id","items":[{"goods_id","count"}]}`。开停看带 `shop` 那条的 `start` / `end` |
| `POST /Activity/TechSubmit` | 雨落成诗建设/领取。body: `{"id","node_id"}`。`id` 是带 `nodes` 的那条活动，不是采集那条。`status` 2 或 3 可提交。回 `items` / `unlocked` / `activity`。活动号 `2026070301`。礼包 2026-08-26 10:00 ~ 2026-09-08 23:59，精确看带 `nodes` 那条的 `start` / `end` |

### 活动 · 青酿换万金

活动号 `2026080102` / `2026081202`。仅活动号，开停未核实。酿酒、日签。

| 路径 | 说明 |
|---|---|
| `POST /Activity/BrewStart` | 开始酿酒，body: `{"id","items":[{"uid","count"}]}`。`uid` 是背包槽位。青酿换万金 `2026080102` / `2026081202`，开停未核实，看带 `brew` 那条的 `start` / `end` |
| `POST /Activity/BrewStep` | 酿酒下一步，body: `{"id"}`。档期同 `BrewStart` |
| `POST /Activity/BrewClaim` | 领酿酒，body: `{"id","claim_type"}`。`1` 普通，`2` 分享。不填当 `1`。档期同 `BrewStart` |
| `POST /Activity/Signin` | 活动签到，body: `{"id","reward_id"}`。`reward_id` 可不填，从列表取。青酿日签走 `2026080102`，开停看该条 `start` / `end` |

### 活动 · 鹊桥寄情

筑桥 `2026081801`，香囊 `2026081802`。礼包 2026-08-18 10:00 ~ 2026-08-22 23:59。精确看 List。

| 路径 | 说明 |
|---|---|
| `POST /Activity/ClaimProgress` | 领进度奖，body: `{"id","step"}`。`step` 一般 `0`。鹊桥寄情筑桥 `2026081801`。礼包 2026-08-18 10:00 ~ 2026-08-22 23:59，精确看带 `progress` 那条的 `start` / `end` |
| `POST /Activity/SendGift` | 给指定好友送礼，body: `{"id","gid","msg_text_id"}`。不群发。鹊桥寄情香囊 `2026081802`。礼包窗同上，精确看带 `gift` 那条的 `start` / `end` |

### 活动 · 邀新领红包

活动号 `2026080701`。仅活动号，开停未核实。

| 路径 | 说明 |
|---|---|
| `POST /Activity/Invitees` | 邀新领红包被邀请名单，body: `{"id"}`。活动号 `2026080701`，开停看带 `invite` 那条的 `start` / `end` |
| `POST /Activity/ClaimInvite` | 邀新领红包邀请/成长档，body: `{"id","reward_type"}`。`1` 邀请，`2` 成长。活动号 `2026080701`，开停看带 `invite` 那条的 `start` / `end` |
| `POST /Activity/ClaimNewcomer` | 邀新领红包新人档，body: `{"id"}`。活动号 `2026080701`，开停看带 `invite` 那条的 `start` / `end` |

### 活动 · 故友重逢

活动号 `2026070101`。仅活动号，开停未核实。

| 路径 | 说明 |
|---|---|
| `POST /Activity/Recallable` | 故友重逢可召回名单，body: `{"id"}`。活动号 `2026070101`，开停看带 `recall` 那条的 `start` / `end` |
| `POST /Activity/Recalled` | 故友重逢已召回名单，body: `{"id"}`。活动号 `2026070101`，开停看带 `recall` 那条的 `start` / `end` |
| `POST /Activity/ClaimRecall` | 故友重逢召回奖，body: `{"id"}`。活动号 `2026070101`，开停看带 `recall` 那条的 `start` / `end` |
| `POST /Activity/ClaimReturn` | 故友重逢回归礼，body: `{"id"}`。活动号 `2026070101`，开停看带 `recall` 那条的 `start` / `end` |

### 活动 · 观星礼录

活动号 `2026072701`，赛季号 `2026072700`。礼包 2026-07-29 10:00 ~ 2026-08-27 23:59。精确看 List / `Season/Info`。

| 路径 | 说明 |
|---|---|
| `POST /Activity/ClaimMega` | 观星礼录一键领，body: `{"id"}`。千星游记 `2026072701`。礼包 2026-07-29 10:00 ~ 2026-08-27 23:59，精确看带 `mega` 那条的 `start` / `end` |

### 活动 · 宠物寻宝

活动号 `2026090101`。仅活动号，开停未核实。cmd `901–916`。

| 路径 | 说明 |
|---|---|
| `POST /Activity/HuntFinishCG` | 宠物寻宝结束 CG。body: `{"id"}`。cmd `901`。活动号 `2026090101`，开停看带 `hunt` 那条的 `start` / `end` |
| `POST /Activity/HuntGuide` | 宠物寻宝引导。body: `{"id"}`。cmd `902`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntFeed` | 宠物寻宝喂养。body: `{"id"}`。cmd `903`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntDraw` | 宠物寻宝抽奖。body: `{"id"}`。cmd `904`。不是 `Activity/Draw`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntLog` | 宠物寻宝日志。body: `{"id"}`。cmd `905`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntClaimStory` | 宠物寻宝领剧情。body: `{"id"}`。cmd `906`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntClaimSeed` | 宠物寻宝领种子。body: `{"id"}`。cmd `907`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntRefreshCharm` | 宠物寻宝刷新护符池。body: `{"id"}`。cmd `908`。消耗看配置，点券 `1002`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntEquip` | 宠物寻宝装备护符。body: `{"id","charm_ids"}`。cmd `909`。护符 `101`–`105`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntBattle` | 宠物寻宝开战/掠夺。body: `{"id","defender_gid","treasure_id"}`。cmd `910`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntPlunderedLog` | 宠物寻宝被掠日志。body: `{"id"}`。cmd `911`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntOpen` | 宠物寻宝开宝。body: `{"id"}`。cmd `913`。有可领宝藏时打这个。档期同 `HuntFinishCG` |
| `POST /Activity/HuntEscort` | 宠物寻宝护送。body: `{"id"}`。cmd `914`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntCompensate` | 宠物寻宝领掠夺补偿。body: `{"id"}`。cmd `915`。档期同 `HuntFinishCG` |
| `POST /Activity/HuntFriendInfo` | 宠物寻宝查好友活动。body: `{"id"}`。cmd `916`。档期同 `HuntFinishCG` |

### 活动 · 公益小红花

活动号 `2026090901`。仅活动号，开停未核实。

| 路径 | 说明 |
|---|---|
| `POST /Activity/CharityShare` | 公益花领分享奖，body: `{"id"}`。活动号 `2026090901`，开停看带 `charity` 那条的 `start` / `end` |
| `POST /Activity/CharityDonate` | 公益花捐出全部爱心，body: `{"id"}`。档期同 `CharityShare` |
| `POST /Activity/CharityClaim` | 公益花领个人档，body: `{"id","score"}`。`score` 是档位所需个人积分。档期同 `CharityShare` |
| `POST /Activity/CharityXhh` | 公益花领星辉回，body: `{"id"}`。档期同 `CharityShare` |
| `POST /Activity/CharityAgree` | 公益花合规同意，body: `{"id","agreed"}`。档期同 `CharityShare` |

### 活动 · 阵营加油

粽香 `2026061901`，足球 `2026071501`。仅活动号，开停未核实。

| 路径 | 说明 |
|---|---|
| `POST /Activity/CheerJoin` | 加入加油阵营，body: `{"id","camp_id"}`。粽香 `2026061901`，足球 `2026071501`。仅活动号，开停看带 `cheer` 那条的 `start` / `end` |
| `POST /Activity/CheerSubmit` | 提交加油，body: `{"id","count"}`。档期同 `CheerJoin` |
| `POST /Activity/CheerClaim` | 领加油档，body: `{"id","tier"}`。档期同 `CheerJoin` |

### 活动 · 抽奖 / 随机店

开停看该条 `start` / `end`。

| 路径 | 说明 |
|---|---|
| `POST /Activity/Draw` | 活动抽奖，body: `{"id","count"}`。`count` 不填当 `1`。列表看 `draw`。开停看该条 `start` / `end` |
| `POST /Activity/DrawHistory` | 活动抽奖历史，body: `{"id"}`。开停看该条 `start` / `end` |
| `POST /Activity/RandBuy` | 随机商店买一件，body: `{"id","goods_id","count"}`。开停看带 `rand_shop` 那条的 `start` / `end` |
| `POST /Activity/RandRefresh` | 刷新随机商店，body: `{"id"}`。`costs` 是扣掉的，新货在 `activity.rand_shop`。开停看带 `rand_shop` 那条的 `start` / `end` |
| `POST /Activity/RandBatchBuy` | 随机商店批量买，body: `{"id","items":[{"goods_id","count"}]}`。开停看带 `rand_shop` 那条的 `start` / `end` |

### 赛季

千星游记赛季号 `2026072700`。礼包 2026-07-29 10:00 ~ 2026-08-27 23:59。精确看 `Season/Info`。

| 路径 | 说明 |
|---|---|
| `POST /Season/Info` | 赛季和通行证。带回预热、活动简介、战令各档。千星游记赛季号 `2026072700`。礼包 2026-07-29 10:00 ~ 2026-08-27 23:59，精确看 `Season.start` / `Season.end` |
| `POST /Season/ClaimPass` | 领已解锁的战令奖励。回 `items` / 刚领等级 / 更新后的战令 |
| `POST /Season/BuyPass` | 买高级战令。空请求。回 `success` / `items` / `pass` |
| `POST /Season/MarkOpening` | 标记赛季开幕已看 |

### 红包

每日红包。

| 路径 | 说明 |
|---|---|
| `POST /RedPacket/List` | 今日红包。`can_claim` 是今日未领且 `status` 非 0 |
| `POST /RedPacket/Claim` | 领一个，body: `{"id"}`。回 `status` 和这一件奖励 |
| `POST /RedPacket/ClaimAll` | 领完能领的 |

### 图鉴

普通 / 超变图鉴。

| 路径 | 说明 |
|---|---|
| `POST /Album/List` | 图鉴，body: `{"type","rarity"}`。`type` 不填为普通 `1`，`2` 是超变。带回等级 / 总奖励 / 加成 |
| `POST /Album/Levels` | 图鉴等级列表，body: `{"type"}`。不填为普通 `1` |
| `POST /Album/Claim` | 领图鉴奖励，body: `{"type"}`。回刚领的等级和下一级预览 |
| `POST /Album/MarkViewed` | 清图鉴新解锁标记，body: `{"type"}`。不填为普通 `1` |

### 护主犬

狗、狗粮、护院。

| 路径 | 说明 |
|---|---|
| `POST /Dog/Info` | 护主犬、狗粮剩余秒、粮库存、待领礼物 |
| `POST /Dog/Feed` | 喂狗粮，body: `{"food_id","count"}`。`food_id` 从 `Dog/Info` 取。不自动补 |
| `POST /Dog/ClaimGifts` | 领技能礼物。回本次数量、剩余待领、补偿 |
| `POST /Dog/Logs` | 护院日志，body: `{"from","count"}`。`count` 不填为 `50`。回 `total` |
| `POST /Dog/Deploy` | 出战，body: `{"id"}` |
| `POST /Dog/Withdraw` | 收回出战犬。回 `withdrawn` 是被收回的狗 |
| `POST /Dog/Activate` | 激活已有的狗，body: `{"id"}` |
| `POST /Dog/Buy` | 买并激活，body: `{"id","price"}`。`price` 不填从 `Dog/Info` 取 |

### 公告

公告列表和正文。

| 路径 | 说明 |
|---|---|
| `POST /Bulletin/List` | 公告，body: `{"from","count"}`。`count` 不填为 `20` |
| `POST /Bulletin/Read` | 公告正文，body: `{"id"}` |

### 变异 / 掉落 / 节气

限时窗口看各条自己的时间。

| 路径 | 说明 |
|---|---|
| `POST /Mutant/List` | 变异宝典。条目是活动窗口，不是果实图鉴。开停看每条自己的窗口 |
| `POST /Mutant/OpenInfo` | 变异系统开放奖励。提示和物品 |
| `POST /Drop/List` | 随机掉落活动，只读。开停看每条自己的窗口，或 `Activity/List` 的 `drops` |
| `POST /Solar/List` | 节气。`status`：1未到 2可领 3已领 4过期未领 5过期已领。带回服务端时间和赛季配置 |
| `POST /Solar/Claim` | 领一个节气，body: `{"id"}`。回更新后的节气 |
| `POST /Solar/ClaimAll` | 领完当前可领的节气 |
| `POST /Solar/RedDot` | 节气红点 |

### 生涯 / 排行

| 路径 | 说明 |
|---|---|
| `POST /Career/Info` | 生涯收获 / 被偷统计。请求无字段，只查当前登录号 |
| `POST /Rank/List` | 排行榜，body: `{"type","page"}`。都不填则 `type=1`、`page=1`。`type`：1等级 2金币 3图鉴 |

### 装扮

头像框和皮肤。

| 路径 | 说明 |
|---|---|
| `POST /Avatar/Owned` | 已有头像框，body: `{"type"}` 可筛选 |
| `POST /Avatar/Equipped` | 当前头像框 |
| `POST /Avatar/Equip` | 穿/脱，body: `{"id","off"}`。`off=true` 卸下 |
| `POST /Avatar/MarkViewed` | 标记头像框已查看，body: `{"id"}` |
| `POST /Skin/Owned` | 已有皮肤 |
| `POST /Skin/Equipped` | 当前皮肤 |
| `POST /Skin/Equip` | 换皮肤，body: `{"current_id","id"}`。`id=0` 卸下 |
| `POST /Skin/EquipSet` | 穿套装，body: `{"skin_ids"}` |
| `POST /Skin/SetEffect` | 开关套装特效，body: `{"set_id","effect_type","enabled"}` |
| `POST /Skin/Sets` | 当前套装特效 |
| `POST /Skin/MarkViewed` | 标记皮肤已查看，body: `{"id"}` |

### QQ 群

授权、绑定社群。会员每日礼看 QQ 会员。

| 路径 | 说明 |
|---|---|
| `POST /QQGroup/AuthGroups` | 已授权群。body 可带 `{"cookies"}` |
| `POST /QQGroup/Recommend` | 推荐群，body: `{"class_name","session_info","scene_id"}` |
| `POST /QQGroup/Bind` | 绑定社群，body: `{"community_id"}` |
| `POST /QQGroup/Leave` | 退出社群。回剩余可退次数和冷却 |
| `POST /QQGroup/Community` | 社群信息。`from` 当页码，不填当 `1` |
| `POST /QQGroup/BindInfo` | 绑定详情和奖励预览 |
| `POST /QQGroup/ClaimReward` | 领社群绑定奖 |
| `POST /QQGroup/RevokeAuth` | 撤销授权，body: `{"gid","community_id"}` |

### QQ 会员

每日礼和会员奖励。不是 QQ 群。

| 路径 | 说明 |
|---|---|
| `POST /QQVip/DailyStatus` | 是不是会员、今天能不能领、奖励预览 |
| `POST /QQVip/ClaimDaily` | 领每日礼。body 可带 `{"config_id"}` |
| `POST /QQVip/Refresh` | 刷新会员信息。回 `is_qq_vip` / `vip_level` |
| `POST /QQVip/RewardsStatus` | 会员奖励档、红点、能不能领 |
| `POST /QQVip/ClaimRewards` | 领会员奖励，body: `{"config_ids"}`。回皮肤 / 头像框 ID 和 `rewards` |
| `POST /QQVip/MarkRedpoint` | 清会员红点 |

### 跑马灯

| 路径 | 说明 |
|---|---|
| `POST /Marquee/List` | 当前滚动消息。每条有 `uuid` / `content` / `expire_time` / `priority` |

### 系统解锁

| 路径 | 说明 |
|---|---|
| `POST /SystemOpen/Unlocked` | 查某个系统有没有开，body: `{"system_name"}`。`system_name` 是官方枚举编号 |
| `POST /Mutant/OpenInfo` | 变异系统开放奖励。走 `Mutant.GetSystemOpenInfo`，不是上面这条解锁查询 |

### 订阅

QQ 订阅和微信订阅消息。不是支付订阅。

| 路径 | 说明 |
|---|---|
| `POST /Subscribe/QQ` | QQ 订阅状态。回包可能是标量 `status`，也可能是 `items` |
| `POST /Subscribe/WX` | 微信订阅模板状态 |
| `POST /Subscribe/SetWX` | 写微信订阅，body: `{"templates":[{"template_id","subscribed"}]}` |

### 审核

文本和图片内容审核。

| 路径 | 说明 |
|---|---|
| `POST /Moderate/Text` | 审一条文本，body: `{"text","reason"}`。回 `result_text` / `is_dirty` |
| `POST /Moderate/BatchText` | 批量审文本，body: `{"text_items"}` |
| `POST /Moderate/Pic` | 审一张图，body: `{"pic_url","reason"}`。回 `result_url` / `is_dirty` / `dirty_type` |
| `POST /Moderate/BatchPic` | 批量审图，body: `{"pic_items"}` |

### 礼包

礼包码道具、领取历史、转账。`uid` 是背包格子，不是兑换码。

| 路径 | 说明 |
|---|---|
| `POST /Gift/UseToken` | 使用礼包码道具，body: `{"uid"}`。兑换码在回包 `redeem_code` |
| `POST /Gift/History` | 领取历史，body: `{"source_type","page","page_size"}` |
| `POST /Gift/TransferStatus` | 查转账，body: `{"out_bill_no"}` |
| `POST /Gift/CancelTransfer` | 取消红包转账，body: `{"uid"}` |

### 关注有礼

| 路径 | 说明 |
|---|---|
| `POST /Follow/Status` | 已关注 / 已领 / 红点 |
| `POST /Follow/Set` | 标记已关注，body: `{"followed"}` |
| `POST /Follow/Claim` | 领关注礼 |

### 充值返利

只读配置和进度。不含支付下单。

| 路径 | 说明 |
|---|---|
| `POST /Recharge/Config` | 是否开、窗口、档位比例 |
| `POST /Recharge/Data` | 已充 / 已返 |

### 成就

| 路径 | 说明 |
|---|---|
| `POST /Achieve/View` | 成就域，body: `{"kind","id"}` |
| `POST /Achieve/ClaimGoal` | 领成就目标，body: `{"kind","id","goal_id"}`。回域经验和更新后的域，不是物品 |
| `POST /Achieve/ClaimLevel` | 领成就域等级，body: `{"kind","id"}`。回奖励和补偿 |

## 活动档期

活动号不是开日。精确到几点、以及结束，一律看 `Activity/List` 该条 `start` / `end`（Unix 秒）。`status`：`1` 进行中，`2` 已结束。

礼包商店只核实了三档：雨落成诗、鹊桥寄情、千星游记。其余没有礼包窗，编号前 8 位不能当开日。

| 活动 | 活动号 | 开始 | 结束 | 依据 | 列表怎么认 |
|---|---|---|---|---|---|
| 粽香大比拼 | `2026061901` | 该条 `start` | 该条 `end` | 仅活动号，未核实 | `cheer` |
| 故友重逢 | `2026070101` | 该条 `start` | 该条 `end` | 仅活动号，未核实 | `recall` |
| 足球狂欢季 | `2026071501` | 该条 `start` | 该条 `end` | 仅活动号，未核实 | `cheer` |
| 千星游记 / 观星礼录 | `2026072701` | 2026-07-29 10:00 | 2026-08-27 23:59 | 千星礼包；活动号是 7-27 | `mega`；赛季 `2026072700` 看 `Season/Info` |
| 青酿换万金 | `2026080102` / `2026081202` | 该条 `start` | 该条 `end` | 仅活动号，未核实 | `brew`；日签也可能在 `signin` |
| 邀新领红包 | `2026080701` | 该条 `start` | 该条 `end` | 仅活动号，未核实 | `invite` |
| 鹊桥寄情·筑桥 | `2026081801` | 2026-08-18 10:00 | 2026-08-22 23:59 | 鹊桥寄情礼包 | `progress` |
| 鹊桥寄情·赠香囊 | `2026081802` | 2026-08-18 10:00 | 2026-08-22 23:59 | 鹊桥寄情礼包 | `gift` |
| 雨落成诗 | `2026070301` | 2026-08-26 10:00 | 2026-09-08 23:59 | 雷雨每日/礼包 | `lottery` / `shop` / `nodes` |
| 宠物寻宝 | `2026090101` | 该条 `start` | 该条 `end` | 仅活动号，未核实 | `hunt` |
| 公益小红花 | `2026090901` | 该条 `start` | 该条 `end` | 仅活动号，未核实 | `charity` |
| 抽奖 / 签到 / 随机店 | 当期活动号 | 该条 `start` | 该条 `end` | 看列表 | `draw` / `signin` / `rand_shop` |

## 运行

默认只听本机。

```text
go run .
# 默认用仓库里的 data/tsdk-v3.9.0.wasm
# set FARM_WASM=D:\path\to\tsdk-v3.9.0.wasm
```

文档：<http://127.0.0.1:8765/docs>

任意语言调用：

```text
curl -X POST http://127.0.0.1:8765/System/Ping
```

## 版本

当前 **1.23.0**。仓库：https://github.com/alttab8520/qqfarm-sdk

发新版：

1. 改 `internal/version/version.go` 和 `CHANGELOG.md`
2. 提交后打标签：`git tag -a v1.23.0 -m "qqfarm-sdk 1.23.0"`
3. `git push origin main --tags`

GitHub Actions 遇到 `v*` 标签会编好 Windows / Linux / macOS 二进制并开 Release。

## 说明

- 非官方，不保证能用，不保证账号安全
- 不要对公网开放
- 仓库里不要提交账号、票据
- 登录前必须能找到 `tsdk-v3.9.0.wasm`，否则 `/User/Login` 直接失败
- ACE 登录后自动跑，不要自己发 `Ace.AntiData`。看 `System/Status` 的 `ace`
- 游戏版本默认 `1.13.3.11_20260826`，可用 `FARM_GAME_VER` 覆盖
