package api

import "net/http"

const openapiYAML = `openapi: 3.0.3
info:
  title: qqfarm-sdk
  description: 非官方 QQ 农场协议 HTTP SDK。全部 POST，回包 {code,msg,data}。
  version: 1.3.0
servers:
  - url: http://127.0.0.1:8765
paths:
  /System/Ping:
    post:
      tags: [System 系统]
      summary: 探活
      responses:
        "200": { description: OK }
  /User/Login:
    post:
      tags: [User 用户]
      summary: 登录
      responses:
        "200": { description: OK }
  /System/Status:
    post:
      tags: [System 系统]
      summary: 登录状态
      responses:
        "200": { description: OK }
  /User/GetInfo:
    post:
      tags: [User 用户]
      summary: 自己资料
      responses:
        "200": { description: OK }
  /User/Logout:
    post:
      tags: [User 用户]
      summary: 退出
      responses:
        "200": { description: OK }
  /Farm/Refresh:
    post:
      tags: [Farm 农场]
      summary: 刷新地块
      responses:
        "200": { description: OK }
  /Farm/Harvest:
    post:
      tags: [Farm 农场]
      summary: 收获
      responses:
        "200": { description: OK }
  /Farm/Steal:
    post:
      tags: [Farm 农场]
      summary: 偷菜
      responses:
        "200": { description: OK }
  /Farm/Plant:
    post:
      tags: [Farm 农场]
      summary: 种植
      responses:
        "200": { description: OK }
  /Farm/Remove:
    post:
      tags: [Farm 农场]
      summary: 铲地
      responses:
        "200": { description: OK }
  /Farm/Unlock:
    post:
      tags: [Farm 农场]
      summary: 解锁地块
      responses:
        "200": { description: OK }
  /Farm/Upgrade:
    post:
      tags: [Farm 农场]
      summary: 升级地块
      responses:
        "200": { description: OK }
  /Farm/Farming:
    post:
      tags: [Farm 农场]
      summary: 一键务农
      responses:
        "200": { description: OK }
  /Farm/Water:
    post:
      tags: [Farm 农场]
      summary: 浇水
      responses:
        "200": { description: OK }
  /Farm/Weed:
    post:
      tags: [Farm 农场]
      summary: 除草
      responses:
        "200": { description: OK }
  /Farm/Bug:
    post:
      tags: [Farm 农场]
      summary: 杀虫
      responses:
        "200": { description: OK }
  /Farm/Fertilize:
    post:
      tags: [Farm 农场]
      summary: 施肥
      responses:
        "200": { description: OK }
  /Friend/GetList:
    post:
      tags: [Friend 好友]
      summary: 好友列表
      responses:
        "200": { description: OK }
  /Friend/Help:
    post:
      tags: [Friend 好友]
      summary: 帮忙
      responses:
        "200": { description: OK }
  /Friend/Enter:
    post:
      tags: [Friend 好友]
      summary: 进好友家
      responses:
        "200": { description: OK }
  /Friend/Leave:
    post:
      tags: [Friend 好友]
      summary: 离开好友家
      responses:
        "200": { description: OK }
  /Friend/Applications:
    post:
      tags: [Friend 好友]
      summary: 好友申请
      responses:
        "200": { description: OK }
  /Friend/Accept:
    post:
      tags: [Friend 好友]
      summary: 同意申请
      responses:
        "200": { description: OK }
  /Friend/Reject:
    post:
      tags: [Friend 好友]
      summary: 拒绝申请
      responses:
        "200": { description: OK }
  /Friend/Delete:
    post:
      tags: [Friend 好友]
      summary: 删除好友
      responses:
        "200": { description: OK }
  /Share/Check:
    post:
      tags: [Share 分享]
      summary: 能否分享
      responses:
        "200": { description: OK }
  /Share/Claim:
    post:
      tags: [Share 分享]
      summary: 领分享奖
      responses:
        "200": { description: OK }
  /Bag/Get:
    post:
      tags: [Bag 背包]
      summary: 背包
      responses:
        "200": { description: OK }
  /Bag/Sell:
    post:
      tags: [Bag 背包]
      summary: 出售
      responses:
        "200": { description: OK }
  /Bag/Use:
    post:
      tags: [Bag 背包]
      summary: 使用
      responses:
        "200": { description: OK }
  /Shop/List:
    post:
      tags: [Shop 商店]
      summary: 商店列表
      responses:
        "200": { description: OK }
  /Shop/Goods:
    post:
      tags: [Shop 商店]
      summary: 商品列表
      responses:
        "200": { description: OK }
  /Shop/Buy:
    post:
      tags: [Shop 商店]
      summary: 购买
      responses:
        "200": { description: OK }
  /Weather/Status:
    post:
      tags: [Weather 天气]
      summary: 当前天气
      responses:
        "200": { description: OK }
  /Weather/Today:
    post:
      tags: [Weather 天气]
      summary: 今日天气
      responses:
        "200": { description: OK }
  /Task/List:
    post:
      tags: [Task 任务]
      summary: 任务列表
      responses:
        "200": { description: OK }
  /Task/Claim:
    post:
      tags: [Task 任务]
      summary: 领任务
      responses:
        "200": { description: OK }
  /Task/ClaimDaily:
    post:
      tags: [Task 任务]
      summary: 领活跃度宝箱
      responses:
        "200": { description: OK }
  /Email/List:
    post:
      tags: [Email 邮件]
      summary: 邮件列表
      responses:
        "200": { description: OK }
  /Email/Read:
    post:
      tags: [Email 邮件]
      summary: 读邮件
      responses:
        "200": { description: OK }
  /Email/Claim:
    post:
      tags: [Email 邮件]
      summary: 领一封
      responses:
        "200": { description: OK }
  /Email/ClaimAll:
    post:
      tags: [Email 邮件]
      summary: 批量领
      responses:
        "200": { description: OK }
  /Activity/List:
    post:
      tags: [Activity 活动]
      summary: 活动清单
      responses:
        "200": { description: OK }
  /Activity/Signin:
    post:
      tags: [Activity 活动]
      summary: 活动签到
      responses:
        "200": { description: OK }
  /Activity/ClaimProgress:
    post:
      tags: [Activity 活动]
      summary: 领进度奖
      responses:
        "200": { description: OK }
  /Activity/ShopBuy:
    post:
      tags: [Activity 活动]
      summary: 活动商店买一件
      responses:
        "200": { description: OK }
  /Activity/ClaimMega:
    post:
      tags: [Activity 活动]
      summary: 观星礼录一键领
      responses:
        "200": { description: OK }
  /Activity/TechSubmit:
    post:
      tags: [Activity 活动]
      summary: 提交科技树节点
      responses:
        "200": { description: OK }
  /Activity/Lottery:
    post:
      tags: [Activity 活动]
      summary: 采集抽奖
      responses:
        "200": { description: OK }
  /Activity/BrewStart:
    post:
      tags: [Activity 活动]
      summary: 开始酿酒
      responses:
        "200": { description: OK }
  /Activity/BrewStep:
    post:
      tags: [Activity 活动]
      summary: 酿酒下一步
      responses:
        "200": { description: OK }
  /Activity/BrewClaim:
    post:
      tags: [Activity 活动]
      summary: 领酿酒
      responses:
        "200": { description: OK }
  /Activity/ClaimRecall:
    post:
      tags: [Activity 活动]
      summary: 故友重逢召回奖
      responses:
        "200": { description: OK }
  /Activity/ClaimReturn:
    post:
      tags: [Activity 活动]
      summary: 故友重逢回归礼
      responses:
        "200": { description: OK }
  /Activity/ClaimInvite:
    post:
      tags: [Activity 活动]
      summary: 邀新红包邀请/成长档
      responses:
        "200": { description: OK }
  /Activity/ClaimNewcomer:
    post:
      tags: [Activity 活动]
      summary: 邀新红包新人档
      responses:
        "200": { description: OK }
  /Activity/SendGift:
    post:
      tags: [Activity 活动]
      summary: 给指定好友送礼
      responses:
        "200": { description: OK }
  /Mystery/Status:
    post:
      tags: [Mystery 神秘商人]
      summary: 商人状态
      responses:
        "200": { description: OK }
  /Mystery/Buy:
    post:
      tags: [Mystery 神秘商人]
      summary: 买一件
      responses:
        "200": { description: OK }
  /Mystery/Leave:
    post:
      tags: [Mystery 神秘商人]
      summary: 打发商人
      responses:
        "200": { description: OK }
  /Season/Info:
    post:
      tags: [Season 赛季]
      summary: 赛季信息
      responses:
        "200": { description: OK }
  /Season/ClaimPass:
    post:
      tags: [Season 赛季]
      summary: 领战令
      responses:
        "200": { description: OK }
  /Mall/List:
    post:
      tags: [Mall 商城]
      summary: 商城列表
      responses:
        "200": { description: OK }
  /Mall/Buy:
    post:
      tags: [Mall 商城]
      summary: 买/领
      responses:
        "200": { description: OK }
  /Mall/MonthCard:
    post:
      tags: [Mall 商城]
      summary: 月卡
      responses:
        "200": { description: OK }
  /Mall/ClaimMonthCard:
    post:
      tags: [Mall 商城]
      summary: 领月卡
      responses:
        "200": { description: OK }
  /RedPacket/List:
    post:
      tags: [RedPacket 红包]
      summary: 今日红包
      responses:
        "200": { description: OK }
  /RedPacket/Claim:
    post:
      tags: [RedPacket 红包]
      summary: 领一个
      responses:
        "200": { description: OK }
  /RedPacket/ClaimAll:
    post:
      tags: [RedPacket 红包]
      summary: 领完能领的
      responses:
        "200": { description: OK }
  /Visit/Logs:
    post:
      tags: [Visit 来访]
      summary: 来访记录
      responses:
        "200": { description: OK }
  /Album/List:
    post:
      tags: [Album 图鉴]
      summary: 图鉴
      responses:
        "200": { description: OK }
  /Album/Claim:
    post:
      tags: [Album 图鉴]
      summary: 领图鉴奖励
      responses:
        "200": { description: OK }
`

const docsHTML = `<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>qqfarm-sdk</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "/openapi.yaml",
      dom_id: "#swagger-ui"
    })
  </script>
</body>
</html>
`

func serveDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(docsHTML))
}

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	_, _ = w.Write([]byte(openapiYAML))
}
