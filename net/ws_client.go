package net

import (
	"github.com/gorilla/websocket"
	"github.com/yhhaiua/engine/buffer"
	"strings"
)

type WsClient struct {
	addr     string
	listener SocketListener
	wsConn   *WSConn
}

func (w *WsClient) OnConnected(conn Channel) {
	w.listener.OnConnected(conn)
}

func (w *WsClient) OnDisconnected(conn Channel) {
	w.listener.OnDisconnected(conn)
	w.wsConn = nil
}

func (w *WsClient) OnData(conn Channel, msg *buffer.ByteBuf) {
	w.listener.OnData(conn, msg)
}

func (w *WsClient) WsConn() *WSConn {
	return w.wsConn
}

func NewWsClient(addr string, listener SocketListener) *WsClient {
	server := &WsClient{
		addr:     addr,
		listener: listener,
	}
	return server
}

func (w *WsClient) Start() {
	conn := w.dial()
	if conn != nil {
		logger.Infof("WsClient connect success:%s", w.addr)
		str := conn.RemoteAddr().String()
		info := strings.Split(str, ":")
		w.wsConn = newWSConn(conn, w, info[0])
		if w.wsConn != nil {
			w.wsConn.start()
		}
	}
}

func (w *WsClient) dial() *websocket.Conn {

	c, _, err := websocket.DefaultDialer.Dial(w.addr, nil)
	if err == nil {
		return c
	}
	logger.Errorf("WsClient connect err: %s", err.Error())
	return nil
}
