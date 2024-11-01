package util

import (
	"math"
	"strconv"
)

//头段(这个是没有意义的主要是确保id整齐所有long数值型id都 10... 14个零开头) (考虑到id的可读性不采用位运算)
var HEAD_FRAGMENT = int64(math.Pow(10, 15))
//运营商段10的15次方 这样规划支持 1-999个运营商
var OPERATOR_FRAGMENT = int64(math.Pow(10, 0))
//运营商段 10的11次方 这样规划支持 1-9999台服  剩余部分为自增id 1-8亿
var SERVER_FRAGMENT = int64(math.Pow(10, 3))
//自增长部分端id 1-8亿
var AUTO_INCREMENT_FRAGMENT = int64(math.Pow(10, 7))

//BuildPrimaryKey 根据运营商、服id、字增长部分构建最终的主键值
func BuildPrimaryKey(operator int,server int,autoincrement int64)int64  {
	return HEAD_FRAGMENT + OPERATOR_FRAGMENT * int64(operator) + int64(server) * SERVER_FRAGMENT + autoincrement * AUTO_INCREMENT_FRAGMENT;
}

//ParseOperator 从指定的主键值中解析出运营商id出来
func ParseOperator(id int64)int  {
	return int((id - HEAD_FRAGMENT) % SERVER_FRAGMENT);// 头 + 自增长 + 服标识 + 运营商
}
//ParseServer 从指定的主键值中解析出服id出来
func ParseServer(id int64)int  {
	return int(((id - HEAD_FRAGMENT) % AUTO_INCREMENT_FRAGMENT) / SERVER_FRAGMENT);
}

//ParseAutoincrement 从指定的主键值中解析出子增长部分出来
func ParseAutoincrement(id int64)int64  {
	return (id - HEAD_FRAGMENT) / AUTO_INCREMENT_FRAGMENT;// 头 + 自增长 + 服标识 + 运营商
}

//ParseUnifyLogo 从指定的主键值中解析出统一标准的区服和台服
func ParseUnifyLogo(id int64)string  {
	oper := strconv.Itoa(ParseOperator(id))
	server := strconv.Itoa(ParseServer(id))
	return oper +"_"+ server;// 头 + 自增长 + 服标识 + 运营商
}