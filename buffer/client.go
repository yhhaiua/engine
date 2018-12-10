package buffer

import (
	"bytes"
	"golang.org/x/net/websocket"
)

//NewClientByte 获取前端发送字符
func newClientByte(cmd CmdCode)[]byte  {
	buf := new(NativeBuffer)
	buf.Buffer = new(bytes.Buffer)
	var init int32
	var init2 int64
	buf.WriteNInt32(init)
	buf.WriteNInt64(init2)
	cmd.Write(buf)
	return buf.Bytes()
}

func ClientSend(ws *websocket.Conn,cmd CmdCode) error  {
	data := newClientByte(cmd)
	return websocket.Message.Send(ws, data)
}