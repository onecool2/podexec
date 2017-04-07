package router

import (
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
)

var (
	upgrader = websocket.Upgrader{}

	rr = render.New(render.Options{
		Directory:  "templates",
		Delims:     render.Delims{"{[{", "}]}"},
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
	})
)
