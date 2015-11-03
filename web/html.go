package main

import "html/template"

var indexTemplate = template.Must(template.New("index").Parse(
	`
<html>
  <head>
  </head>
  <body>
    <h1>Services</h1>
    {{range .}}
        <h1>{{ .Name }}</h1>
    {{end}}
    <br/>
    <a hlink="/documenation">Documentation</a>
  </body>
</html>`))
