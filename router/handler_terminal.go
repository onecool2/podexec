package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"example/podexec/backend"

	"github.com/gorilla/websocket"
	"github.com/pressly/chi"
	ws "golang.org/x/net/websocket"
)

func terminalsHandler(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	if len(strings.TrimSpace(pid)) == 0 {
		pid = "1"
	}

	w.Write([]byte(pid))
}

func socketHandler(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	fmt.Println(" socketHandler ==========>>> pid ", pid)

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	receive := func(conn *websocket.Conn) {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("[ERROR] Receive message error:%v", err)
				return
			}

			fmt.Printf("Receive data : %s \n", string(data[:]))
		}
	}

	send := func(conn *websocket.Conn) {
		conn.WriteMessage(websocket.TextMessage, []byte("消息来自服务端"))
	}

	go receive(conn)
	go send(conn)
}

func wsHandler(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	s := ws.Server{Handler: ws.Handler(func(conn *ws.Conn) {
		req := conn.Request()
		// colsStr := r.URL.Query()["cols"][0]
		// rowsStr := r.URL.Query()["rows"][0]
		// cols, _ := strconv.ParseUint(colsStr, 10, 16)
		// rows, _ := strconv.ParseUint(rowsStr, 10, 16)
		url := req.RequestURI
		sections := strings.Split(url, "/")
		log.Println("=========>>  wsHandler 2 ", url, sections)
		// b.ExecPod("test", "nginx", "nginx", []string{"sh"}, conn, conn, conn, &term.Size{Width: uint16(cols), Height: uint16(rows)})
		b.ExecPod(sections[2], sections[3], sections[4], []string{"/bin/bash"}, conn, conn, conn, nil)
	})}
	s.ServeHTTP(w, r)
}
