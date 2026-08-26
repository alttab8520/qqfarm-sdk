# qqfarm-sdk

非官方 QQ 农场协议 HTTP SDK。Go 实现，别人用任意语言 `POST` 调用。

启动后打开 `/docs`，接口全部是 `POST /模块/动作`，回包统一：

```json
{ "code": 0, "msg": "ok", "data": {} }
```

登录后会连游戏网关，接口走真实 RPC。当前请求体加密是直通，正式服可能拒绝，下一步补加密。

## 范围

| 路径 | 说明 |
|---|---|
| `POST /System/Ping` | 探活 |
| `POST /User/Login` | 登录，body: `{"code","open_id"}` |
| `POST /User/GetInfo` | 自己资料 |
| `POST /Farm/Refresh` | 刷新地块 |
| `POST /Farm/Harvest` | 收获，body: `{"land_ids","host_gid","is_all"}` |
| `POST /Farm/Plant` | 种植，body: `{"seed_id","land_ids"}` |
| `POST /Friend/GetList` | 好友列表 |
| `POST /Friend/Help` | 帮忙浇水，body: `{"gid","land_ids"}` |

## 运行

默认只听本机。

```text
go run .
```

文档：<http://127.0.0.1:8765/docs>

任意语言调用：

```text
curl -X POST http://127.0.0.1:8765/System/Ping
```

## 版本

当前 **0.2.0**。仓库：https://github.com/alttab8520/qqfarm-sdk

发新版：

1. 改 `internal/version/version.go` 和 `CHANGELOG.md`
2. 提交后打标签：`git tag -a v0.1.1 -m "qqfarm-sdk 0.1.1"`
3. `git push origin main --tags`

GitHub Actions 遇到 `v*` 标签会编好 Windows / Linux / macOS 二进制并开 Release。

## 说明

- 非官方，不保证能用，不保证账号安全
- 不要对公网开放
- 仓库里不要提交账号、票据、更新服务器地址
