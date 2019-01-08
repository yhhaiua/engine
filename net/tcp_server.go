package net

import (
	"github.com/yhhaiua/engine/log"
	"net"
	"time"
)

var gLog  = log.GetLogger()

type TCPServer struct {
	addr       string
	listener   SocketListener
	ln         net.Listener
}

func NewTCPServer(addr string,listener SocketListener) *TCPServer{
	server:= &TCPServer{
		addr:addr,
		listener:listener,
	}
	return server
}

func (server *TCPServer)Listen() error {

	lister, err := net.Listen("tcp", server.addr)

	if err != nil{
		gLog.Error(err)
		return err
	}
	gLog.Info("tcp success:%s",server.addr)
	server.ln = lister
	go server.run()
	return nil
}

func (server *TCPServer)run()  {

	var tempDelay time.Duration
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				gLog.Info("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			gLog.Error("accept error: %v;",err)
			return
		}
		tempDelay = 0
		tcpConn := newTcpConn(conn,server.listener)
		if tcpConn != nil{
			tcpConn.start()
		}

	}
}