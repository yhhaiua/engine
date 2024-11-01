package handler

import (
	"github.com/yhhaiua/engine/buffer"
)

type Handler interface {
	IsValidLength(length int) bool //是否有效长度
	Decode(in *buffer.ByteBuf) (*buffer.ByteBuf, error)
}
