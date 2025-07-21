package net

import (
	"github.com/gorilla/websocket"
	"github.com/yhhaiua/engine/buffer"
	"github.com/yhhaiua/engine/handler"
	"github.com/yhhaiua/engine/util"
	"sync"
	"time"
)

var globalId = &util.AtomicLong{}

const (
	WsDataLength       = 100
	pongWait           = 2 * 60 * time.Second //等待时间
	pingPeriod         = 8 * pongWait / 10    //周期54s
	maxMsgSize   int64 = 32767 * 2            //消息最大长度
	writeWait          = 10 * time.Second     //
)

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
	//一次从管管中读取的最大长度
	t.conn.SetReadLimit(maxMsgSize)
	//连接中，每隔54秒向客户端发一次ping，客户端返回pong，所以把SetReadDeadline设为60秒，超过60秒后不允许读
	_ = t.conn.SetReadDeadline(time.Now().Add(pongWait))
	//心跳
	t.conn.SetPongHandler(func(appData string) error {
		//每次收到pong都把deadline往后推迟60秒
		_ = t.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
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
	//给前端发心跳，看前端是否还存活
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
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
				//10秒内必须把信息写给前端（写到websocket连接里去），否则就关闭连接
				_ = t.conn.SetWriteDeadline(time.Now().Add(writeWait))
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
		case <-ticker.C:
			_ = t.conn.SetWriteDeadline(time.Now().Add(writeWait))
			//心跳保持，给浏览器发一个PingMessage，等待浏览器返回PongMessage
			if err := t.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Errorf(" websocket.PingMessage:%s", err.Error())
				t.close()
				return //写websocket连接失败，说明连接出问题了，该client可以over了
			}
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
