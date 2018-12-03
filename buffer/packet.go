package buffer

type CPacket struct {
	Code int
}

type CmdCode interface {
	Write(buf *NativeBuffer)
	Read(buf *NativeBuffer)
}