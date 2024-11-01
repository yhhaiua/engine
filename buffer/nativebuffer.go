//对包结构进行二进制解析，写入和读取规则

package buffer

import (
	"encoding/binary"
)

var byteOrder = binary.BigEndian

//数组的最大数量
const MaxArrayCount = 2000

//ReadNString 写入字符串
func (buf *ByteBuf)ReadNString() string  {
	length := buf.ReadNInt16()
	return  string(buf.Next(length))
}

//ReadNShort 读取int16
func (buf *ByteBuf)ReadNInt8() int8 {
	var v int8
	binary.Read(buf, byteOrder,&v)
	return v
}
//ReadNInt16 读取int16 返回int
func (buf *ByteBuf)ReadNInt16() int {
	var v int16
	binary.Read(buf, byteOrder,&v)
	return int(v)
}

//ReadNShort 读取int16
func (buf *ByteBuf)ReadNShort() int16 {
	var v int16
	binary.Read(buf, byteOrder,&v)
	return v
}

//ReadNInt32 读取int32
func (buf *ByteBuf)ReadNInt32() int32 {
	var v int32
	binary.Read(buf, byteOrder,&v)
	return v
}

//ReadNInt64 读取int64
func (buf *ByteBuf)ReadNInt64() int64  {
	var v int64
	binary.Read(buf, byteOrder,&v)
	return v
}

//ReadNBytes 读取bytes
func (buf *ByteBuf)ReadNBytes() []byte {
	len := buf.ReadNInt16()
	if len == 0{
		return nil
	}
	result := buf.Next(len)
	return result
}
//ReadNInt16Array 读取int16数组
func (buf *ByteBuf)ReadNInt16Array() []int16 {
	len := buf.ReadNInt16()
	if len == 0 || len > MaxArrayCount {
		return nil
	}
	result := make([]int16,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNShort()
	}
	return result
}

//ReadNInt32Array 读取int32数组
func (buf *ByteBuf)ReadNInt32Array() []int32 {
	len := buf.ReadNInt16()
	if len == 0 || len > MaxArrayCount {
		return nil
	}
	result := make([]int32,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNInt32()
	}
	return result
}
//ReadNInt64Array 读取int64数组
func (buf *ByteBuf)ReadNInt64Array() []int64 {
	len := buf.ReadNInt16()
	if len == 0 || len > MaxArrayCount {
		return nil
	}
	result := make([]int64,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNInt64()
	}
	return result
}
//ReadNStringArray 读取字符串数组
func (buf *ByteBuf)ReadNStringArray() []string {
	len := buf.ReadNInt16()
	if len == 0 || len > MaxArrayCount {
		return nil
	}
	result := make([]string,len)
	for i:= 0;i < len;i++{
		result[i] = buf.ReadNString()
	}
	return result
}

//WriteNString 写入字符串
func (buf *ByteBuf)WriteNString(v string)  {
	buf.WriteNInt16(len(v))
	_, _ = buf.WriteString(v)
}
//
func (buf *ByteBuf)WriteNInt8(v int8)  {
	_ = binary.Write(buf, byteOrder, v)
}

//WriteNInt16 写入int16
func (buf *ByteBuf)WriteNInt16(v int)  {
	_ = binary.Write(buf, byteOrder, int16(v))
}

//WriteNShort 写入int16
func (buf *ByteBuf)WriteNShort(v int16)  {
	_ = binary.Write(buf, byteOrder, v)
}

//WriteNInt32 写入int32
func (buf *ByteBuf)WriteNInt32(v int32)  {
	_ = binary.Write(buf, byteOrder, v)
}

//WriteNInt64 写入int64
func (buf *ByteBuf)WriteNInt64(v int64)  {
	_ = binary.Write(buf, byteOrder, v)
}

//WriteNBytes 写字节数组
func (buf *ByteBuf)WriteNBytes(v []byte){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	_, _ = buf.Write(v)
}
//WriteNInt16Array 写入int16数组
func (buf *ByteBuf)WriteNInt16Array(v []int16){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNShort(v[i])
	}
}

//WriteNInt32Array 写入int32数组
func (buf *ByteBuf)WriteNInt32Array(v []int32){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNInt32(v[i])
	}
}

//WriteNInt64Array 写入int64数组
func (buf *ByteBuf)WriteNInt64Array(v []int64){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNInt64(v[i])
	}
}

//WriteNStringArray 写入字符串数组
func (buf *ByteBuf)WriteNStringArray(v []string){
	if v == nil{
		buf.WriteNInt16(0)
		return
	}
	buf.WriteNInt16(len(v))
	for i:= 0;i < len(v);i++{
		buf.WriteNString(v[i])
	}
}
//SetInt 从特点位置传入值
func (buf *ByteBuf)SetInt(index int,value int)  {
	var bs [4]byte
	byteOrder.PutUint32(bs[:], uint32(value))
	buf.setData(index,bs[:])
}