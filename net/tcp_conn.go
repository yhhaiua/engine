package net

import (
	"engine/buffer"
	"engine/handler"
	"engine/util"
	"io"
	"net"
	"strings"
	"sync"
)

const TcpDataLength = 100

type TCPConn struct {
	closeMutex    sync.Mutex
	conn          net.Conn
	receive       *buffer.ByteBuf
	listener      SocketListener
	hd            handler.Handler
	chData        chan []byte
	connected     bool
	connectAtomic util.AtomicInteger
	closeData     bool
	id            int64
	chStopWrite   chan struct{}
}

func (t *TCPConn) Id() int64 {
	return t.id
}

func (t *TCPConn) Ip() string {
	str := t.conn.RemoteAddr().String()
	info := strings.Split(str, ":")
	return info[0]
}

func newTcpConn(conn net.Conn, listener SocketListener) *TCPConn {
	t := new(TCPConn)
	t.conn = conn
	t.listener = listener
	t.receive = buffer.NewByteBuf()
	t.chData = make(chan []byte, TcpDataLength)
	t.chStopWrite = make(chan struct{})
	t.connected = true
	hd, err := handler.NewLengthDecoder()
	if err != nil {
		logger.Errorf("new TcpConn err: %s", err.Error())
		return nil
	}
	t.hd = hd
	t.id = globalId.IncrementAndGet()
	return t
}

func (t *TCPConn) start() {

	t.listener.OnConnected(t)
	go t.read()
	go t.run()
}

func (t *TCPConn) read() {

	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
			t.close()
			t.listener.OnDisconnected(t)
		}
	}()
	for {
		err := t.receive.ReadFrom(t.conn)
		if err == io.EOF {
			logger.Infof("远程连接：%s,关闭", t.conn.RemoteAddr().String())
			t.close()
			t.listener.OnDisconnected(t)
			return
		}
		if err != nil {
			logger.Errorf("read err: %s", err.Error())
			t.close()
			t.listener.OnDisconnected(t)
			return
		}
		msg, err := t.hd.Decode(t.receive)
		if err != nil {
			logger.Errorf("msg err: %s", err.Error())
			continue
		}
		if msg != nil {
			t.listener.OnData(t, msg)
		}
	}
}

func (t *TCPConn) run() {
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
				_, err := t.conn.Write(msg)
				if err != nil {
					logger.Errorf(" Write:%s", err.Error())
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

//Destroy 目标通知断开后销毁
func (t *TCPConn) Destroy() {
	t.doDestroy()
}

//WriteAndFlush 向目标发送数据
func (t *TCPConn) WriteAndFlush(msg []byte) {
	if t.closeData {
		return
	}
	t.doWrite(msg)
}

//Close 主动关闭连接（调用前先向目标发送关闭信息）
func (t *TCPConn) Close() {
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()
	if t.closeData {
		return
	}
	t.chData <- nil
	t.closeData = true
}

func (t *TCPConn) close() {

	if !t.connectAtomic.CompareAndSet(0, 1) {
		return
	}
	if !t.connected {
		return
	}
	t.connected = false
	_ = t.conn.Close()
}

func (t *TCPConn) doDestroy() {
	t.close()
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()
	if t.closeData {
		return
	}
	close(t.chStopWrite)
	t.closeData = true
}

func (t *TCPConn) doWrite(msg []byte) {
	select {
	case t.chData <- msg:
	default:
		t.full(msg)
	}
}

func (t *TCPConn) full(msg []byte) {
	logger.Infof("队列已满进行等待,ip:%s,len:%d", t.Ip(), len(t.chData))
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()
	if t.closeData {
		return
	}
	t.chData <- msg
}
