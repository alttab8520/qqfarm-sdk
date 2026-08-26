package api

import "net/http"

const openapiYAML = `openapi: 3.0.3
info:
  title: qqfarm-sdk
  description: 非官方 QQ 农场协议 HTTP SDK。全部 POST，回包 {code,msg,data}。
  version: 0.2.0
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
  /User/GetInfo:
    post:
      tags: [User 用户]
      summary: 自己资料
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
  /Farm/Plant:
    post:
      tags: [Farm 农场]
      summary: 种植
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
