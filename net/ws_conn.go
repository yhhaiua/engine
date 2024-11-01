package net

import (
	"github.com/gorilla/websocket"
	"github.com/yhhaiua/engine/buffer"
	"github.com/yhhaiua/engine/handler"
	"github.com/yhhaiua/engine/util"
	"sync"
)

var globalId = &util.AtomicLong{}

const WsDataLength = 100

type WSConn struct {
	closeMutex    sync.Mutex
	conn          *websocket.Conn
	receive       *buffer.ByteBuf
	listener      SocketListener
	hd            handler.Handler
	chData        chan []byte
	chStopWrite   chan struct{}
	connected     bool
	connectAtomic util.AtomicInteger
	ip            string
	closeData     bool
	id            int64
}

func (t *WSConn) Id() int64 {
	return t.id
}
func (t *WSConn) Ip() string {
	return t.ip
}

func newWSConn(conn *websocket.Conn, listener SocketListener, ip string) *WSConn {
	t := new(WSConn)
	t.conn = conn
	t.listener = listener
	t.receive = buffer.NewByteBuf()
	t.chData = make(chan []byte, WsDataLength)
	t.chStopWrite = make(chan struct{})
	t.connected = true
	hd, err := handler.NewWsDecoder()
	if err != nil {
		logger.Errorf("new WSConn err: %v", err)
		return nil
	}
	t.hd = hd
	t.ip = ip
	t.id = globalId.IncrementAndGet()
	return t
}
func (t *WSConn) start() {

	t.listener.OnConnected(t)
	go t.read()
	go t.run()
}

func (t *WSConn) read() {

	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
			t.close()
			t.listener.OnDisconnected(t)
		}
	}()
	for {
		_, b, err := t.conn.ReadMessage()
		if err != nil {
			logger.Warnf("远程连接：%s,关闭 %s", t.Ip(), err.Error())
			t.close()
			t.listener.OnDisconnected(t)
			return
		}
		length := len(b)
		if !t.hd.IsValidLength(length) {
			logger.Errorf("msg err,too large: %d,Ip:%s", length, t.Ip())
			continue
		}
		_, err = t.receive.Write(b)
		if err != nil {
			logger.Errorf("msg err: %s,Ip:%s", err.Error(), t.Ip())
			continue
		}
		msg, err := t.hd.Decode(t.receive)
		if err != nil {
			logger.Errorf("msg err: %s,Ip:%s", err.Error(), t.Ip())
			continue
		}
		if msg != nil {
			t.listener.OnData(t, msg)
		}
	}
}

func (t *WSConn) run() {
	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
			t.close()
		}
	}()

	for {
		select {
		case msg, ok := <-t.chData:
			if !ok || msg == nil {
				logger.Infof("msg close:%s", t.Ip())
				t.close()
				return
			}
			if msg != nil && t.connected {
				err := t.conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					logger.Errorf(" WriteMessage:%s", err.Error())
					t.close()
					return
				}
			}
		case <-t.chStopWrite:
			logger.Infof("chStopWrite destroy:%s", t.Ip())
			return
		}
	}

}

// Destroy 目标通知断开后销毁
func (t *WSConn) Destroy() {
	t.doDestroy()
}

// WriteAndFlush 向目标发送数据
func (t *WSConn) WriteAndFlush(msg []byte) {
	if t.closeData {
		return
	}
	t.doWrite(msg)
}

// Close 主动关闭连接（调用前先向目标发送关闭信息）
func (t *WSConn) Close() {
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()
	if t.closeData {
		return
	}
	t.chData <- nil
	t.closeData = true
}

func (t *WSConn) close() {

	if !t.connectAtomic.CompareAndSet(0, 1) {
		return
	}
	if !t.connected {
		return
	}
	t.connected = false
	_ = t.conn.Close()
}

func (t *WSConn) doDestroy() {
	t.close()
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()
	if t.closeData {
		return
	}
	close(t.chStopWrite)
	t.closeData = true
}

func (t *WSConn) doWrite(msg []byte) {
	select {
	case t.chData <- msg:
	default:
		t.full(msg)
	}
}

func (t *WSConn) full(msg []byte) {
	logger.Infof("队列已满进行等待,ip:%s,len:%d", t.Ip(), len(t.chData))
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()
	if t.closeData {
		return
	}
	t.chData <- msg
}
