package api

import "net/http"

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
    fetch("/openapi.json").then(function(r){ return r.json() }).then(function(spec){
      var tagOrder = (spec.tags || []).map(function(t){ return t.name })
      var pathOrder = Object.keys(spec.paths || {})
      window.ui = SwaggerUIBundle({
        spec: spec,
        dom_id: "#swagger-ui",
        tagsSorter: function(a, b){
          var ia = tagOrder.indexOf(a), ib = tagOrder.indexOf(b)
          return (ia < 0 ? 999 : ia) - (ib < 0 ? 999 : ib)
        },
        operationsSorter: function(a, b){
          return pathOrder.indexOf(a.get("path")) - pathOrder.indexOf(b.get("path"))
        }
      })
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(openAPIJSON())
}
