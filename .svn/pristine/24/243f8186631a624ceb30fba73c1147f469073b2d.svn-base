//对websocket接到的包进行解析
//websocket数据不需要处理粘包和半包

package handler

import (
	"encoding/binary"
	"engine/buffer"
	"errors"
	"strconv"
)

type WsDecoder struct {
	byteOrder      binary.ByteOrder
	maxFrameLength int
}

func NewWsDecoder() (Handler, error) {
	handler := new(WsDecoder)
	err := handler.Init(binary.BigEndian, 32767)
	return handler, err
}

func (decoder *WsDecoder) Init(byteOrder binary.ByteOrder, maxFrameLength int) error {

	if byteOrder == nil {
		return errors.New("byteOrder nil")
	}
	if maxFrameLength <= 0 {
		return errors.New("maxFrameLength must be a positive integer: " + strconv.Itoa(maxFrameLength))
	}

	decoder.maxFrameLength = maxFrameLength
	decoder.byteOrder = byteOrder
	return nil
}

//IsValidLength 判断是否是有效长度
func (decoder *WsDecoder) IsValidLength(length int) bool {
	if length <= 0 || length > decoder.maxFrameLength {
		return false
	}
	return true
}
func (decoder *WsDecoder) Decode(in *buffer.ByteBuf) (*buffer.ByteBuf, error) {

	frameLength := in.ReadableBytes()
	if frameLength > decoder.maxFrameLength {
		in.SkipBytes(in.ReadableBytes())
		return nil, errors.New("Adjusted frame length (" + strconv.Itoa(frameLength) + ") is more " +
			"than maxFrameLength: " + strconv.Itoa(decoder.maxFrameLength))
	}
	readerIndex := in.ReaderIndex()
	actualFrameLength := frameLength
	frame := in.RetainedSlice(readerIndex, actualFrameLength)
	in.ReaderToIndex(readerIndex + actualFrameLength)
	return frame, nil
}
