package net

import (
	"github.com/gorilla/websocket"
	"github.com/yhhaiua/engine/gutil"
	"net/http"
	"time"
)

type WSServer struct {
	addr     string
	listener SocketListener
	upgrader  websocket.Upgrader
	HTTPTimeout time.Duration
}

func NewWSServer(addr string,listener SocketListener) *WSServer{
	server:= &WSServer{
		addr:addr,
		listener:listener,
		HTTPTimeout : 30 * time.Second,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: 30 * time.Second,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}
	return server
}

func (ws *WSServer)ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		gLog.Error("upgrade: %v", err)
		return
	}
	c.SetReadLimit(32767)
	wsConn := newWSConn(c,ws.listener,gutil.GetUserIp(r))
	if wsConn != nil{
		wsConn.start()
	}
}
func (ws *WSServer)Listen() error {
	srv := &http.Server{
		ReadTimeout: ws.HTTPTimeout,
		WriteTimeout: ws.HTTPTimeout,
		Addr:ws.addr,
		Handler : ws,
	}
	gLog.Info("websocket start monitor :%s",ws.addr)
	err := srv.ListenAndServe()
	if err != nil {
		gLog.Error("websocket monitor fail %s", err)
		return err
	}

	return nil
}