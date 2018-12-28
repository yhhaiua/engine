package handler

import (
	"github.com/yhhaiua/engine/buffer"
)

type Handler interface {

	Decode(in *buffer.ByteBuf)(*buffer.ByteBuf,error)

}