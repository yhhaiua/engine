package net

import (
	"github.com/yhhaiua/engine/buffer"
)

type TcpSession struct {
	channel Channel
}

func (t *TcpSession) Channel() Channel {
	return t.channel
}

func NewTcpSession(channel Channel) *TcpSession {
	return &TcpSession{channel: channel}
}
func (t *TcpSession) Post(cmd buffer.PacketCons) {
	msg := buffer.NewByteBuf()
	msg.WriteNInt32(0)
	cmd.Write(msg)
	msg.SetInt(0, msg.ReadableBytes()-4)
	t.channel.WriteAndFlush(msg.Bytes())
}

type TcpSInterFace interface {
	Post(cmd buffer.PacketCons)
}
