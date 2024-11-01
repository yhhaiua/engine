package net

import (
	"engine/buffer"
)

type SocketListener interface {
	OnConnected(conn Channel)
	OnDisconnected(conn Channel)
	OnData(conn Channel, msg *buffer.ByteBuf)
}

type SocketRunListener interface {
	SocketListener
	Run(conn Channel)
}
