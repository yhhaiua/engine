package net

import (
	"engine/util"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WSServer struct {
	addr        string
	listener    SocketListener
	upGrader    websocket.Upgrader
	HTTPTimeout time.Duration
	server      *http.Server
}

func NewWSServer(addr string, listener SocketListener) *WSServer {
	server := &WSServer{
		addr:        addr,
		listener:    listener,
		HTTPTimeout: 30 * time.Second,
		upGrader: websocket.Upgrader{
			HandshakeTimeout: 30 * time.Second,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}
	return server
}

func (ws *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	c, err := ws.upGrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debugf("upgrade: %s,ip:%s", err.Error(), util.GetUserIp(r))
		return
	}
	c.SetReadLimit(32767)
	wsConn := newWSConn(c, ws.listener, util.GetUserIp(r))
	if wsConn != nil {
		wsConn.start()
	}
}
func (ws *WSServer) Listen() {
	ws.server = &http.Server{
		ReadTimeout:  ws.HTTPTimeout,
		WriteTimeout: ws.HTTPTimeout,
		Addr:         ws.addr,
		Handler:      ws,
	}
	logger.Infof("websocket start monitor %s", ws.addr)
	err := ws.server.ListenAndServe()
	if err != nil {
		logger.Errorf("websocket monitor fail %s", err.Error())
	}
}
func (ws *WSServer) Close() {
	if ws.server != nil {
		_ = ws.server.Close()
	}
}
