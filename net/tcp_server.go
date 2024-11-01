package net

import (
	"github.com/yhhaiua/engine/log"
	"net"
	"time"
)

var logger = log.GetLogger()

type TCPServer struct {
	addr     string
	listener SocketListener
	ln       net.Listener
}

func NewTCPServer(addr string, listener SocketListener) *TCPServer {
	server := &TCPServer{
		addr:     addr,
		listener: listener,
	}
	return server
}

func (server *TCPServer) Listen() {

	lister, err := net.Listen("tcp", server.addr)

	if err != nil {
		logger.Errorf("Listen error:%s", err.Error())
		return
	}
	logger.Infof("tcp success:%s", server.addr)
	server.ln = lister
	server.run()
}

func (server *TCPServer) run() {

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
				logger.Infof("accept error: %s; retrying in %v", err.Error(), tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			logger.Errorf("accept error: %s;", err.Error())
			return
		}
		tempDelay = 0
		tcpConn := newTcpConn(conn, server.listener)
		if tcpConn != nil {
			tcpConn.start()
		}

	}
}
func (server *TCPServer) Close() {
	if server.ln != nil {
		_ = server.ln.Close()
	}
}
