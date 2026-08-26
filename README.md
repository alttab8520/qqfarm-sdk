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
| `POST /Farm/Refresh` | 刷新地块 |
| `POST /Farm/Harvest` | 收获，body: `{"land_ids","host_gid","is_all"}`。偷菜把 `host_gid` 填好友 |
| `POST /Farm/Plant` | 种植，body: `{"seed_id","land_ids"}` |
| `POST /Farm/Water` | 浇水，body: `{"land_ids","host_gid"}` |
| `POST /Farm/Weed` | 除草，body: `{"land_ids","host_gid"}` |
| `POST /Farm/Bug` | 杀虫，body: `{"land_ids","host_gid"}` |
| `POST /Farm/Fertilize` | 施肥，body: `{"land_ids","fertilizer_id"}` |
| `POST /Friend/GetList` | 好友列表 |
| `POST /Friend/Help` | 进好友家并浇水，body: `{"gid","land_ids"}` |

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

当前 **0.5.0**。仓库：https://github.com/alttab8520/qqfarm-sdk

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
