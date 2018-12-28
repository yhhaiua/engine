package net

import (
	"github.com/yhhaiua/engine/buffer"
)

type SocketListener interface {
	OnConnected(conn Channel)
	OnDisconnected(conn Channel)
	OnData(conn Channel,msg *buffer.ByteBuf)
}
