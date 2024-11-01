package net

type Channel interface {
	//WriteAndFlush 向目标发送数据
	WriteAndFlush(msg []byte)
	//Close 主动关闭连接（调用前先向目标发送关闭信息）
	Close()
	//Ip 返回ip
	Ip() string
	//Destroy 目标通知断开后销毁
	Destroy()
	//Id 连接编号
	Id() int64
}
