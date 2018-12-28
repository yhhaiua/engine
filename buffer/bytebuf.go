package buffer

import (
	"encoding/binary"
	"errors"
	"io"
)
var ErrTooLarge = errors.New("buffer.ByteBuf: too large")

type ByteBuf struct {
	buf       		[]byte
	readerIndex     int
	writerIndex		int
}

const MinLength  = 64
const ReadLength  = 256
//NewByteBuf 新建ByteBuf
func NewByteBuf() *ByteBuf  {
	b := new(ByteBuf)
	b.buf = b.makeSlice(MinLength)
	return b
}
func NewBuffer(buf []byte) *ByteBuf  {
	return &ByteBuf{buf: buf}
}

func (b *ByteBuf) Bytes() []byte { return b.buf[b.readerIndex:b.writerIndex] }

//ReadableBytes 可读取的字节数
func (b *ByteBuf)ReadableBytes() int  {
	return b.writerIndex - b.readerIndex;
}

//WritableBytes可写入的字节数
func (b *ByteBuf)WritableBytes() int  {
	return len(b.buf) - b.writerIndex
}

//WritableBytes可写入的字节数
func (b *ByteBuf)WritableMaxBytes() int  {
	return cap(b.buf) - b.writerIndex
}
//ReaderIndex读取开始地址
func (b *ByteBuf)ReaderIndex() int{
	return b.readerIndex
}
//ReaderToIndex readerIndex跳转
func (b *ByteBuf)ReaderToIndex(index int){
	b.readerIndex = index
}
//WritableIndex 写入开始地址
func (b *ByteBuf)WritableIndex() int{
	return b.writerIndex
}
//WriteBytes写入数据
func (b *ByteBuf)Write(p []byte) (n int, err error) {

	length := len(p)
	m,err := b.grow(length)

	if err != nil{
		return 0,err
	}
	b.writerIndex += length
	return copy(b.buf[m:], p), nil
}
func (b *ByteBuf)WriteString(p string) (n int, err error) {

	length := len(p)
	m,err := b.grow(length)

	if err != nil{
		return 0,err
	}
	b.writerIndex += length
	return copy(b.buf[m:], p), nil
}
//SkipBytes 跳过字节不读
func (b *ByteBuf)SkipBytes(length int)  {
	if length > 0 && length <= b.ReadableBytes(){
		b.readerIndex += length
	}
}
//GetUnsignedByte 获取长度
func (b *ByteBuf)GetUnsignedByte(index int) int  {
	x := b.buf[index]
	return int(x)
}

//GetUnsignedShort 获取长度
func (b *ByteBuf)GetUnsignedShort(index int,byteOrder binary.ByteOrder) int {
	x := byteOrder.Uint16(b.buf[index:index+2])
	return int(x)
}

//GetUnsignedInt 获取长度
func (b *ByteBuf)GetUnsignedInt(index int,byteOrder binary.ByteOrder) int  {
	x := byteOrder.Uint32(b.buf[index:index+4])
	return int(x)
}
//GetLong 获取长度
func (b *ByteBuf)GetLong(index int,byteOrder binary.ByteOrder) int  {
	x := byteOrder.Uint64(b.buf[index:index+8])
	return int(x)
}

func (b *ByteBuf)RetainedSlice(index,length int) *ByteBuf  {
	v := NewBuffer(b.buf[index:index+length])
	v.writerIndex = length
	return v
}
func (b *ByteBuf) tryGrowByReslice(n int) (int, bool) {
	if n <= b.WritableBytes(){
		return b.writerIndex,true
	}else if n <= b.WritableMaxBytes(){
		b.buf = b.buf[:b.writerIndex+n]
		return b.writerIndex, true
	}
	return 0, false
}

func (b *ByteBuf) grow(n int) (int,error) {
	m, ok := b.tryGrowByReslice(n)
	if ok {
		return m,nil
	}
	if b.WritableMaxBytes() + b.ReaderIndex() > 2 * n{
		copy(b.buf, b.buf[b.readerIndex:])
		b.writerIndex = b.ReadableBytes()
		b.readerIndex = 0
	}else{
		buf := b.makeSlice(2*cap(b.buf) + n)
		copy(buf, b.buf[b.readerIndex:])
		b.writerIndex = b.ReadableBytes()
		b.readerIndex = 0
		b.buf = buf
	}
	m, ok = b.tryGrowByReslice(n)
	if !ok {
		return 0,ErrTooLarge
	}
	return m,nil
}

func (b *ByteBuf)makeSlice(n int) ([]byte) {
	return make([]byte, n)
}

func (b *ByteBuf) Next(n int) []byte {
	m := b.ReadableBytes()
	if n > m {
		n = m
	}
	data := b.buf[b.readerIndex : b.readerIndex+n]
	b.readerIndex += n
	return data
}

func (b *ByteBuf) ReadFrom(r io.Reader) error {
	i,err := b.grow(ReadLength)
	if err != nil{
		return err
	}
	m, e := r.Read(b.buf[i:cap(b.buf)])
	if e != nil{
		return e
	}
	b.writerIndex += m

	return nil
}

func (b *ByteBuf) Read(p []byte) (n int, err error) {

	if b.ReadableBytes() == 0 {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.readerIndex:b.writerIndex])
	b.readerIndex += n
	return n, nil
}

func (b *ByteBuf)setData(index int,bs []byte)  {
	length := len(bs)
	if index + length > cap(b.buf){
		return
	}
	copy(b.buf[index:], bs)
}