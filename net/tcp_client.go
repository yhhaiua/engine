package net

import (
	"github.com/yhhaiua/engine/buffer"
	"net"
	"sync"
	"time"
)

var (
	mTCPConnMap *sync.Map
	timer       *time.Ticker
	closeKey    chan string
)

type TCPClient struct {
	index    string
	addr     string
	listener SocketRunListener
	tcpConn  *TCPConn
}

func NewTCPClient(index string, addr string, listener SocketRunListener) *TCPClient {
	client := &TCPClient{
		addr:     addr,
		listener: listener,
		index:    index,
	}
	return client
}

func (client *TCPClient) Start() {
	client.Connect()
	mTCPConnMap.Store(client.index, client)
}
func (client *TCPClient) Connect() {
	if client.tcpConn != nil {
		client.listener.Run(client.tcpConn)
		return
	}
	conn := client.dial()
	if conn != nil {
		logger.Infof("%s,tcp connect success:%s", client.index, client.addr)
		client.tcpConn = newTcpConn(conn, client)
		if client.tcpConn != nil {
			client.tcpConn.start()
		}
	}

}

func (client *TCPClient) dial() net.Conn {

	conn, err := net.Dial("tcp", client.addr)
	if err == nil {
		return conn
	}
	logger.Errorf("%s,tcp connect err,wait again: %s", client.index, err.Error())
	return nil
}

func (client *TCPClient) OnConnected(conn Channel) {

	client.listener.OnConnected(conn)
}
func (client *TCPClient) OnDisconnected(conn Channel) {

	client.listener.OnDisconnected(conn)
	client.tcpConn = nil
}
func (client *TCPClient) OnData(conn Channel, msg *buffer.ByteBuf) {

	client.listener.OnData(conn, msg)
}

func (client *TCPClient) Close() {
	closeKey <- client.index
}

func run() {
	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
			go run() //出现退出异常，再次启动
		}
	}()
	for {
		select {
		//数据处理
		case <-timer.C:
			tcpConnectRun()
		case key, ok := <-closeKey:
			if ok {
				tcpConnectClose(key)
			}
		}
	}

}
func tcpConnectRun() {
	mTCPConnMap.Range(func(key, value interface{}) bool {
		connect, ok := value.(*TCPClient)
		if ok {
			connect.Connect()
		}
		return ok
	})
}
func tcpConnectClose(key string) {
	v, ok := mTCPConnMap.Load(key)
	if ok {
		connect := v.(*TCPClient)
		if connect.tcpConn != nil {
			connect.tcpConn.Close()
		}
		mTCPConnMap.Delete(key)
	}
}
func init() {
	mTCPConnMap = new(sync.Map)
	timer = time.NewTicker(5 * time.Second)
	closeKey = make(chan string)
	go run()
}
