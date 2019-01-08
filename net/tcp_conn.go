package net

import (
	"github.com/yhhaiua/engine/buffer"
	"github.com/yhhaiua/engine/handler"
	"net"
	//"sync"
)

const DataLength  = 32

type TCPConn struct {
	conn 		net.Conn
	receive 	*buffer.ByteBuf
	listener    SocketListener
	hd 			handler.Handler
	data 		chan []byte
	connected 	bool
	closedata  	bool
}

func newTcpConn(conn net.Conn,listener SocketListener) *TCPConn {
	t := new(TCPConn)
	t.conn = conn
	t.listener = listener
	t.receive = buffer.NewByteBuf()
	t.data = make(chan []byte,DataLength)
	t.connected = true
	hd,err := handler.NewLengthDecoder()
	if err != nil{
		gLog.Error("new TcpConn err: %v",err)
		return nil
	}
	t.hd = hd
	return t
}

func (t *TCPConn) start() {

	t.listener.OnConnected(t)
	go t.read()
	go t.run()
}

func (t *TCPConn)read()  {

	defer func() {
		if r := recover();r != nil{
			gLog.Error("read error : %v",r)
			t.close()
			t.listener.OnDisconnected(t)
		}
	}()
	for  {
		err := t.receive.ReadFrom(t.conn)
		if err != nil{
			t.close()
			t.listener.OnDisconnected(t)
			return
		}
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

func (t *TCPConn)run()  {
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
				t.conn.Write(msg)
			}
		}
	}
}

func (t *TCPConn)WriteAndFlush(msg []byte)()  {
	if t.closedata{
		return
	}
	defer func() {
		if r := recover();r != nil{
			gLog.Error(" WriteAndFlush:%v",r)
		}
	}()
	t.data <- msg
}

func (t *TCPConn)Close()  {
	if t.closedata{
		return
	}
	t.closedata = true
	close(t.data)
}
func (t *TCPConn)String() string  {
	return t.conn.RemoteAddr().String()
}
func (t *TCPConn)close(){
	//t.Lock()
	//defer t.Unlock()
	if !t.connected{
		return
	}
	t.connected = false
	t.conn.Close()
}