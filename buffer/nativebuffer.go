package buffer

import (
	"encoding/binary"
)

type NativeBuffer struct {
	*ByteBuf
	byteOrder binary.ByteOrder
}
//Wrap 新建
func Wrap(b *ByteBuf) *NativeBuffer  {
	buf := new(NativeBuffer)
	buf.ByteBuf = b
	buf.byteOrder = binary.BigEndian
	return buf
}
//NewNativeBuffer 新建
func NewNativeBuffer() *NativeBuffer  {
	buf := new(NativeBuffer)
	buf.byteOrder = binary.BigEndian
	buf.ByteBuf = NewByteBuf()
	return buf
}
//ReadNString 写入字符串
func (buf *NativeBuffer)ReadNString() string  {
	length := buf.ReadNInt16()
	return  string(buf.Next(length))
}

//ReadNShort 读取int16
func (buf *NativeBuffer)ReadNInt8() int8 {
	var v int8
	binary.Read(buf, buf.byteOrder,&v)
	return v
}
//ReadNInt16 读取int16 返回int
func (buf *NativeBuffer)ReadNInt16() int {
	var v int16
	binary.Read(buf, buf.byteOrder,&v)
	return int(v)
}

//ReadNShort 读取int16
func (buf *NativeBuffer)ReadNShort() int16 {
	var v int16
	binary.Read(buf, buf.byteOrder,&v)
	return v
}

//ReadNInt32 读取int32
func (buf *NativeBuffer)ReadNInt32() int32 {
	var v int32
	binary.Read(buf, buf.byteOrder,&v)
	return v
}

//ReadNInt64 读取int64
func (buf *NativeBuffer)ReadNInt64() int64  {
	var v int64
	binary.Read(buf, buf.byteOrder,&v)
	return v
}

func (buf *NativeBuffer)ReadNBytes() []byte {
	len := buf.ReadNInt16()
	if len == 0{
		return nil
	}
	result := buf.Next(len)
	return result
}

func (buf *NativeBuffer)ReadNInt16Array() []int16 {
	len := buf.ReadNInt16()
	if len == 0 || len > 1000{
		return nil
	}
	result := make([]int16,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNShort()
	}
	return result
}

func (buf *NativeBuffer)ReadNInt32Array() []int32 {
	len := buf.ReadNInt16()
	if len == 0 || len > 1000{
		return nil
	}
	result := make([]int32,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNInt32()
	}
	return result
}
func (buf *NativeBuffer)ReadNInt64Array() []int64 {
	len := buf.ReadNInt16()
	if len == 0 || len > 1000{
		return nil
	}
	result := make([]int64,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNInt64()
	}
	return result
}
func (buf *NativeBuffer)ReadNStringArray() []string {
	len := buf.ReadNInt16()
	if len == 0 || len > 1000{
		return nil
	}
	result := make([]string,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNString()
	}
	return result
}

//WriteNString 写入字符串
func (buf *NativeBuffer)WriteNString(v string)  {
	buf.WriteNInt16(len(v))
	buf.WriteString(v)
}
//
func (buf *NativeBuffer)WriteNInt8(v int8)  {
	binary.Write(buf, buf.byteOrder, v)
}

//WriteNInt16 写入int16
func (buf *NativeBuffer)WriteNInt16(v int)  {
	binary.Write(buf, buf.byteOrder, int16(v))
}

//WriteNShort 写入int16
func (buf *NativeBuffer)WriteNShort(v int16)  {
	binary.Write(buf, buf.byteOrder, v)
}

//WriteNInt32 写入int32
func (buf *NativeBuffer)WriteNInt32(v int32)  {
	binary.Write(buf, buf.byteOrder, v)
}

//WriteNInt64 写入int64
func (buf *NativeBuffer)WriteNInt64(v int64)  {
	binary.Write(buf, buf.byteOrder, v)
}
func (buf *NativeBuffer)WriteNBytes(v []byte){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	buf.Write(v)
}
func (buf *NativeBuffer)WriteNInt16Array(v []int16){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNShort(v[i])
	}
}
func (buf *NativeBuffer)WriteNInt32Array(v []int32){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNInt32(v[i])
	}
}

func (buf *NativeBuffer)WriteNInt64Array(v []int64){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNInt64(v[i])
	}
}
func (buf *NativeBuffer)WriteNStringArray(v []string){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNString(v[i])
	}
}

func (buf *NativeBuffer)SetInt(index int,value int)  {
	var bs [4]byte
	buf.byteOrder.PutUint32(bs[:], uint32(value))
	buf.setData(index,bs[:])
}