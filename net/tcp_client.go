package net

import (
	"github.com/yhhaiua/engine/buffer"
	"net"
	"sync"
	"time"
)
var mTCPConnMap *sync.Map


type TCPClient struct {
	index 	   string
	addr       string
	listener   SocketListener
	tcpConn	   *TCPConn
}
func NewTCPClient(index string,addr string,listener SocketListener) *TCPClient{
	client:= &TCPClient{
		addr:addr,
		listener:listener,
		index:index,
	}
	return client
}

func (client *TCPClient)Start()  {
	client.Connect()
	mTCPConnMap.Store(client.index,client)
}
func (client *TCPClient)Connect()  {
	if client.tcpConn != nil{
		return
	}
	conn := client.dial()
	if conn != nil {
		gLog.Info("%s,tcp connect success:%s",client.index,client.addr)
		client.tcpConn = newTcpConn(conn,client)
		if client.tcpConn != nil{
			client.tcpConn.start()
		}
	}

}

func (client *TCPClient) dial() net.Conn {

	conn, err := net.Dial("tcp", client.addr)
	if err == nil {
		return conn
	}
	gLog.Error("%s,tcp connect err,wait again: %v",client.index,err)
	return nil
}

func (client *TCPClient)OnConnected(conn Channel){

	client.listener.OnConnected(conn)
}
func (client *TCPClient)OnDisconnected(conn Channel){

	client.listener.OnDisconnected(conn)
	client.tcpConn = nil
}
func (client *TCPClient)OnData(conn Channel,msg *buffer.ByteBuf){

	client.listener.OnData(conn,msg)
}

func run()  {
	for  {
		mTCPConnMap.Range(func(key, value interface{}) bool {
			connect, ok := value.(*TCPClient)
			if ok{
				connect.Connect()
			}
			return ok
		})
		time.Sleep(5 * time.Second)
	}

}
func init()  {
	mTCPConnMap =  new(sync.Map)
	go run()
}