//仿照java netty下 LengthFieldBasedFrameDecoder 类解析器构造解析

package handler

import (
	"encoding/binary"
	"errors"
	"github.com/yhhaiua/engine/buffer"
	"math"
	"strconv"
)

type LengthDecoder struct {
	byteOrder              binary.ByteOrder
	maxFrameLength         int
	lengthFieldOffset      int
	lengthFieldLength      int
	lengthFieldEndOffset   int
	lengthAdjustment       int
	initialBytesToStrip    int
	discardingTooLongFrame bool
	tooFrameLength         int
	bytesToDiscard         int
}

//NewLengthDecoder 构建解码
func NewLengthDecoder() (Handler, error) {
	handler := new(LengthDecoder)
	err := handler.Init(binary.BigEndian, 327670, 0,
		4, 0, 4)
	return handler, err
}

//NewLengthDecoderClient 构建解码
func NewLengthDecoderClient(length int) (Handler, error) {
	handler := new(LengthDecoder)
	err := handler.Init(binary.BigEndian, length, 0,
		4, 0, 4)
	return handler, err
}

// Init 初始化构建解码器
func (decoder *LengthDecoder) Init(byteOrder binary.ByteOrder, maxFrameLength, lengthFieldOffset,
	lengthFieldLength, lengthAdjustment, initialBytesToStrip int) error {

	if byteOrder == nil {
		return errors.New("byteOrder nil")
	}
	if maxFrameLength <= 0 {
		return errors.New("maxFrameLength must be a positive integer: " + strconv.Itoa(maxFrameLength))
	}

	if lengthFieldOffset < 0 {
		return errors.New("lengthFieldOffset must be a non-negative integer: " + strconv.Itoa(lengthFieldOffset))
	}

	if initialBytesToStrip < 0 {
		return errors.New("initialBytesToStrip must be a non-negative integer: " + strconv.Itoa(initialBytesToStrip))
	}

	if lengthFieldOffset > maxFrameLength-lengthFieldLength {
		return errors.New("maxFrameLength (" + strconv.Itoa(maxFrameLength) + ") " +
			"must be equal to or greater than " +
			"lengthFieldOffset (" + strconv.Itoa(lengthFieldOffset) + ") + " +
			"lengthFieldLength (" + strconv.Itoa(lengthFieldLength) + ").")
	}

	decoder.maxFrameLength = maxFrameLength
	decoder.lengthFieldOffset = lengthFieldOffset
	decoder.lengthFieldLength = lengthFieldLength
	decoder.lengthAdjustment = lengthAdjustment
	decoder.lengthFieldEndOffset = lengthFieldOffset + lengthFieldLength
	decoder.initialBytesToStrip = initialBytesToStrip
	decoder.byteOrder = byteOrder
	return nil
}

//IsValidLength 判断是否是有效长度
func (decoder *LengthDecoder) IsValidLength(length int) bool {
	if length <= 0 || length > decoder.maxFrameLength {
		return false
	}
	return true
}

//Decode 对数据进行解码
func (decoder *LengthDecoder) Decode(in *buffer.ByteBuf) (*buffer.ByteBuf, error) {

	if decoder.discardingTooLongFrame {
		bytesToDiscard := decoder.bytesToDiscard
		localBytesToDiscard := int(math.Min(float64(bytesToDiscard), float64(in.ReadableBytes())))
		in.SkipBytes(localBytesToDiscard)
		bytesToDiscard -= localBytesToDiscard
		decoder.bytesToDiscard = bytesToDiscard
		decoder.failIfNecessary()
	}

	if in.ReadableBytes() < decoder.lengthFieldEndOffset {
		return nil, nil
	}
	actualLengthFieldOffset := in.ReaderIndex() + decoder.lengthFieldOffset

	frameLength, err := decoder.getUnadjustedFrameLength(in, actualLengthFieldOffset, decoder.lengthFieldLength)
	if err != nil {
		in.SkipBytes(decoder.lengthFieldEndOffset)
		return nil, err
	}
	if frameLength < 0 {
		in.SkipBytes(decoder.lengthFieldEndOffset)
		return nil, errors.New("negative pre-adjustment length field: " + strconv.Itoa(frameLength))
	}
	frameLength += decoder.lengthAdjustment + decoder.lengthFieldEndOffset

	if frameLength < decoder.lengthFieldEndOffset {
		in.SkipBytes(decoder.lengthFieldEndOffset)
		return nil, errors.New("Adjusted frame length (" + strconv.Itoa(frameLength) + ") is less " +
			"than lengthFieldEndOffset: " + strconv.Itoa(decoder.lengthFieldEndOffset))
	}
	if frameLength > decoder.maxFrameLength {
		discard := frameLength - in.ReadableBytes()
		decoder.tooFrameLength = frameLength
		if discard < 0 {
			in.SkipBytes(frameLength)
		} else {
			decoder.discardingTooLongFrame = true
			decoder.bytesToDiscard = discard
			in.SkipBytes(in.ReadableBytes())
		}
		decoder.failIfNecessary()
		return nil, errors.New("Adjusted frame length (" + strconv.Itoa(frameLength) + ") is more " +
			"than maxFrameLength: " + strconv.Itoa(decoder.maxFrameLength))
	}

	// never overflows because it's less than maxFrameLength
	if in.ReadableBytes() < frameLength {
		return nil, nil
	}
	if decoder.initialBytesToStrip > frameLength {
		in.SkipBytes(frameLength)
		return nil, errors.New("Adjusted frame length (" + strconv.Itoa(frameLength) + ") is less " +
			"than initialBytesToStrip: " + strconv.Itoa(decoder.initialBytesToStrip))
	}
	in.SkipBytes(decoder.initialBytesToStrip)

	readerIndex := in.ReaderIndex()
	actualFrameLength := frameLength - decoder.initialBytesToStrip
	frame := in.RetainedSlice(readerIndex, actualFrameLength)
	in.ReaderToIndex(readerIndex + actualFrameLength)
	return frame, nil
}

func (decoder *LengthDecoder) getUnadjustedFrameLength(buf *buffer.ByteBuf, offset, length int) (int, error) {
	var frameLength int
	switch length {
	case 1:
		frameLength = buf.GetUnsignedByte(offset)
	case 2:
		frameLength = buf.GetUnsignedShort(offset, decoder.byteOrder)
	case 4:
		frameLength = buf.GetUnsignedInt(offset, decoder.byteOrder)
	case 8:
		frameLength = buf.GetLong(offset, decoder.byteOrder)
	default:
		return 0, errors.New("unsupported lengthFieldLength: " + strconv.Itoa(decoder.lengthFieldLength) + " (expected: 1, 2, 4, or 8)")
	}
	return frameLength, nil
}

func (decoder *LengthDecoder) failIfNecessary() {
	if decoder.bytesToDiscard == 0 {
		decoder.tooFrameLength = 0
		decoder.discardingTooLongFrame = false
	}
}
