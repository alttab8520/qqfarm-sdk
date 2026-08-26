# qqfarm-sdk

非官方 QQ 农场协议 HTTP SDK。Go 实现，别人用任意语言 `POST` 调用。

启动后打开 `/docs`，接口全部是 `POST /模块/动作`，回包统一：

```json
{ "code": 0, "msg": "ok", "data": {} }
```

当前只搭好文档和路由骨架，**还没有接游戏长连接**。未实现的接口返回 `code=501`。

## 范围

| 路径 | 说明 |
|---|---|
| `POST /System/Ping` | 探活 |
| `POST /User/Login` | 登录（未接入） |
| `POST /User/GetInfo` | 自己资料（未接入） |
| `POST /Farm/Refresh` | 刷新地块（未接入） |
| `POST /Farm/Harvest` | 收获（未接入） |
| `POST /Farm/Plant` | 种植（未接入） |
| `POST /Friend/GetList` | 好友列表（未接入） |
| `POST /Friend/Help` | 帮忙（未接入） |

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

当前 **0.1.0**。仓库：https://github.com/alttab8520/qqfarm-sdk

发新版：

1. 改 `internal/version/version.go` 和 `CHANGELOG.md`
2. 提交后打标签：`git tag -a v0.1.1 -m "qqfarm-sdk 0.1.1"`
3. `git push origin main --tags`

GitHub Actions 遇到 `v*` 标签会编好 Windows / Linux / macOS 二进制并开 Release。

## 说明

- 非官方，不保证能用，不保证账号安全
- 不要对公网开放
- 仓库里不要提交账号、票据、更新服务器地址
