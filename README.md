# qqfarm-sdk

非官方 QQ 农场协议 HTTP SDK。Go 实现，别人用任意语言 `POST` 调用。

启动后打开 `/docs`，接口全部是 `POST /模块/动作`，回包统一：

```json
{ "code": 0, "msg": "ok", "data": {} }
```

登录后会连游戏网关，接口走真实 RPC。请求体需要官方加密运行时：把 `tsdk-v3.9.0.wasm` 放到 `data/`，或设置 `FARM_WASM`。文件必须是 160975 字节，SHA256 `0001ab8d68ac35309dbdaa16310aeebf0fd578976a91562db12257786e3a2e54`（微信包 59，自报 `v3.9.0.1787640498`）。仓库不收录该文件。

## 范围

| 路径 | 说明 |
|---|---|
| `POST /System/Ping` | 探活 |
| `POST /System/Status` | 登录状态和 ACE 计数 |
| `POST /User/Login` | 登录，body: `{"code","open_id"}` |
| `POST /User/GetInfo` | 自己资料 |
| `POST /User/Logout` | 断开网关 |
| `POST /Farm/Refresh` | 刷新地块，body 可带 `{"host_gid"}` |
| `POST /Farm/Harvest` | 收获，body: `{"land_ids","host_gid","is_all"}` |
| `POST /Farm/Steal` | 偷菜，body 同收获，必须填好友 `host_gid` |
| `POST /Farm/Plant` | 种植，body: `{"seed_id","land_ids"}` |
| `POST /Farm/Remove` | 铲地，body: `{"land_ids"}` |
| `POST /Farm/Unlock` | 解锁地块，body: `{"land_id"}` |
| `POST /Farm/Upgrade` | 升级地块，body: `{"land_id"}` |
| `POST /Farm/Farming` | 一键浇水除草杀虫，自己地用。body: `{"land_ids","host_gid"}` |
| `POST /Farm/Water` | 浇水，body: `{"land_ids","host_gid"}` |
| `POST /Farm/Weed` | 除草，body: `{"land_ids","host_gid"}` |
| `POST /Farm/Bug` | 杀虫，body: `{"land_ids","host_gid"}` |
| `POST /Farm/Fertilize` | 施肥，body: `{"land_ids","fertilizer_id"}` |
| `POST /Friend/GetList` | 好友列表，含干旱 / 杂草 / 虫害 / 可偷数量 |
| `POST /Friend/Help` | 进好友家并浇水，body: `{"gid","land_ids"}` |
| `POST /Friend/Enter` | 进好友家，body: `{"gid"}`。回包有地块和护主犬 |
| `POST /Friend/Leave` | 离开好友家，body: `{"gid"}` |
| `POST /Friend/Applications` | 好友申请 |
| `POST /Friend/Accept` | 同意，body: `{"gids"}` |
| `POST /Friend/Reject` | 拒绝，body: `{"gids"}` |
| `POST /Friend/Delete` | 删除好友，body: `{"gid"}` |
| `POST /Share/Check` | 今日能不能领分享奖 |
| `POST /Share/Claim` | 上报分享并领奖 |
| `POST /Bag/Get` | 背包 |
| `POST /Bag/Sell` | 出售，body: `{"items":[{"id","count","uid"}]}`。`uid` 可从 `Bag/Get` 带上 |
| `POST /Bag/Use` | 使用，body: `{"id","count","uid"}` |
| `POST /Shop/List` | 商店列表 |
| `POST /Shop/Goods` | 商品，body: `{"shop_id"}`。不填则种子店 `2` |
| `POST /Shop/Buy` | 购买，body: `{"goods_id","num","price"}`。`price` 用货架上的单价 |
| `POST /Weather/Status` | 当前天气 |
| `POST /Weather/Today` | 今日天气 |
| `POST /Task/List` | 成长 / 每日任务和活跃度宝箱 |
| `POST /Task/Claim` | 领任务，body: `{"id"}` 或 `{"ids"}` |
| `POST /Task/ClaimDaily` | 领活跃度宝箱，body: `{"point_ids"}` |
| `POST /Email/List` | 邮件，body: `{"box"}`。不填为 `1`，常见还有 `2` |
| `POST /Email/Read` | 读邮件，body: `{"box","id"}` |
| `POST /Email/Claim` | 领一封，body: `{"box","id"}` |
| `POST /Email/ClaimAll` | 批量领。不填 `box` 则 `1` 和 `2` 都领 |
| `POST /Activity/List` | 活动清单。另有 `signin_*` / `shop` / `nodes` / `lottery` / `brew` / `recall` / `invite` / `gift` |
| `POST /Activity/Lottery` | 采集抽奖，body: `{"id","host_gid","free","paid"}`。自己指定好友，不自动巡访 |
| `POST /Activity/BrewStart` | 开始酿酒，body: `{"id","items":[{"uid","count"}]}`。`uid` 是背包槽位 |
| `POST /Activity/BrewStep` | 酿酒下一步，body: `{"id"}` |
| `POST /Activity/BrewClaim` | 领酿酒，body: `{"id","claim_type"}`。`1` 普通，`2` 分享。不填当 `1` |
| `POST /Activity/ClaimRecall` | 故友重逢召回奖，body: `{"id"}` |
| `POST /Activity/ClaimReturn` | 故友重逢回归礼，body: `{"id"}` |
| `POST /Activity/ClaimInvite` | 邀新红包邀请/成长档，body: `{"id","reward_type"}`。`1` 邀请，`2` 成长 |
| `POST /Activity/ClaimNewcomer` | 邀新红包新人档，body: `{"id"}` |
| `POST /Activity/SendGift` | 给指定好友送礼，body: `{"id","gid","msg_text_id"}`。不群发 |
| `POST /Activity/Signin` | 活动签到，body: `{"id","reward_id"}`。`reward_id` 可不填，从列表取 |
| `POST /Activity/ClaimProgress` | 领进度奖，body: `{"id","step"}`。`step` 一般 `0` |
| `POST /Activity/ShopBuy` | 活动商店买一件，body: `{"id","goods_id","count"}`。自己指定，不自动买 |
| `POST /Activity/ClaimMega` | 观星礼录一键领，body: `{"id"}` |
| `POST /Activity/TechSubmit` | 提交科技树节点，body: `{"id","node_id"}` |
| `POST /Mystery/Status` | 神秘商人是否在、货架 |
| `POST /Mystery/Buy` | 买一件，body: `{"id"}`。自己指定，不自动买 |
| `POST /Mystery/Leave` | 打发商人 |
| `POST /Season/Info` | 赛季和通行证 |
| `POST /Season/ClaimPass` | 领已解锁的战令奖励 |
| `POST /Mall/List` | 商城，body: `{"slot"}`。不填为 `1` |
| `POST /Mall/Buy` | 买/领商城商品，body: `{"id","num"}`。自己指定，不自动买 |
| `POST /Mall/MonthCard` | 月卡 |
| `POST /Mall/ClaimMonthCard` | 领月卡，body: `{"id"}` |
| `POST /RedPacket/List` | 今日红包 |
| `POST /RedPacket/Claim` | 领一个，body: `{"id"}` |
| `POST /RedPacket/ClaimAll` | 领完能领的 |
| `POST /Visit/Logs` | 来访记录。`action`：1浇水 2除虫 3除草 4偷菜 5放虫 6放草 |
| `POST /Album/List` | 图鉴，body: `{"type","rarity"}`。`type` 不填为普通 `1`，`2` 是超变 |
| `POST /Album/Claim` | 领图鉴奖励，body: `{"type"}` |

## 运行

默认只听本机。

```text
# 先把 tsdk-v3.9.0.wasm 放到 data/，或：
# set FARM_WASM=D:\path\to\tsdk-v3.9.0.wasm
go run .
```

文档：<http://127.0.0.1:8765/docs>

任意语言调用：

```text
curl -X POST http://127.0.0.1:8765/System/Ping
```

## 版本

当前 **1.3.0**。仓库：https://github.com/alttab8520/qqfarm-sdk

发新版：

1. 改 `internal/version/version.go` 和 `CHANGELOG.md`
2. 提交后打标签：`git tag -a v0.1.1 -m "qqfarm-sdk 0.1.1"`
3. `git push origin main --tags`

GitHub Actions 遇到 `v*` 标签会编好 Windows / Linux / macOS 二进制并开 Release。

## 说明

- 非官方，不保证能用，不保证账号安全
- 不要对公网开放
- 仓库里不要提交账号、票据、加密运行时 wasm
- 登录前必须能找到 `tsdk-v3.9.0.wasm`，否则 `/User/Login` 直接失败
- 游戏版本默认 `1.13.3.11_20260826`，可用 `FARM_GAME_VER` 覆盖
