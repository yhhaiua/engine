package net

import (
	"github.com/gorilla/websocket"
	"github.com/yhhaiua/engine/buffer"
	"github.com/yhhaiua/engine/handler"
	"sync"
)

type WSConn struct {
	sync.Mutex
	conn 		*websocket.Conn
	receive 	*buffer.ByteBuf
	listener    SocketListener
	hd 			handler.Handler
	data 		chan []byte
	connected 	bool
	ip 			string
	closedata  	bool
}

func newWSConn(conn *websocket.Conn,listener SocketListener,ip string) *WSConn {
	t := new(WSConn)
	t.conn = conn
	t.listener = listener
	t.receive = buffer.NewByteBuf()
	t.data = make(chan []byte,DataLength)
	t.connected = true
	hd,err := handler.NewWsDecoder()
	if err != nil{
		gLog.Error("new WSConn err: %v",err)
		return nil
	}
	t.hd = hd
	t.ip = ip
	return t
}
func (t *WSConn) start() {

	t.listener.OnConnected(t)
	go t.read()
	go t.run()
}

func (t *WSConn)read()  {

	defer func() {
		if r := recover();r != nil{
			gLog.Error("read error : %v",r)
			t.close()
			t.listener.OnDisconnected(t)
		}
	}()
	for  {
		_, b, err := t.conn.ReadMessage()
		if err != nil{
			t.close()
			t.listener.OnDisconnected(t)
			return
		}
		t.receive.Write(b)
		msg,err := t.hd.Decode(t.receive)
		if err != nil{
			gLog.Error("msg err: %v",err)
			continue
		}
		if msg != nil{
			t.listener.OnData(t,msg)
		}
	}
}

func (t *WSConn)run()  {
	defer func() {
		if r := recover();r != nil{
			gLog.Error(" abnormal:%v",r)
			t.close()
		}
	}()

	for  {
		select {
		case msg, ok := <-t.data:
			if !ok {
				gLog.Info("run destroy:%s",t.String())
				return
			}
			if msg != nil && t.connected{
				err := t.conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil{
					gLog.Error(" WriteMessage:%v",err)
				}
			}
		}
	}
}

func (t *WSConn)WriteAndFlush(msg []byte)()  {
	t.Lock()
	defer t.Unlock()
	if t.closedata{
		return
	}
	t.data <- msg
}

func (t *WSConn)Close()  {
	t.Lock()
	defer t.Unlock()
	if t.closedata{
		return
	}
	t.closedata = true
	close(t.data)
}
func (t *WSConn)String() string  {
	return t.ip
}
func (t *WSConn)close(){

	if !t.connected{
		return
	}
	t.connected = false
	t.conn.Close()
}