package buffer

import (
	"bytes"
	"encoding/binary"
)

type NativeBuffer struct {
	*bytes.Buffer
}

//WriteNString 写入字符串
func (buf *NativeBuffer)ReadNString() string  {
	length := buf.ReadNInt16()
	return  string(buf.Next(length))
}
//ReadNInt16 读取int16 返回int
func (buf *NativeBuffer)ReadNInt16() int {
	var v int16
	binary.Read(buf, binary.BigEndian,&v)
	return int(v)
}

//ReadNShort 读取int16
func (buf *NativeBuffer)ReadNShort() int16 {
	var v int16
	binary.Read(buf, binary.BigEndian,&v)
	return v
}

//ReadNInt32 读取int32
func (buf *NativeBuffer)ReadNInt32() int32 {
	var v int32
	binary.Read(buf, binary.BigEndian,&v)
	return v
}

//ReadNInt64 读取int64
func (buf *NativeBuffer)ReadNInt64() int64  {
	var v int64
	binary.Read(buf, binary.BigEndian,&v)
	return v
}

//WriteNString 写入字符串
func (buf *NativeBuffer)WriteNString(v string)  {
	buf.WriteNInt16(len(v))
	buf.WriteString(v)
}

//WriteNInt16 写入int16
func (buf *NativeBuffer)WriteNInt16(v int)  {
	binary.Write(buf, binary.BigEndian, int16(v))
}

//WriteNShort 写入int16
func (buf *NativeBuffer)WriteNShort(v int16)  {
	binary.Write(buf, binary.BigEndian, v)
}

//WriteNInt32 写入int32
func (buf *NativeBuffer)WriteNInt32(v int32)  {
	binary.Write(buf, binary.BigEndian, v)
}

//WriteNInt64 写入int64
func (buf *NativeBuffer)WriteNInt64(v int64)  {
	binary.Write(buf, binary.BigEndian, v)
}