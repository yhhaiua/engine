package buffer

type CPacket struct {
	Code int
}

type PacketCons interface {
	Write(buf *ByteBuf)
	Read(buf *ByteBuf)
	Copy() PacketCons
	CodeId() int
	Module() string
	GetStage() int
}

//type PacketCons interface {
//	Code() int
//	Module() int
//}
