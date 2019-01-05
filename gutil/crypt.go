package gutil

import (
	"crypto/md5"
	"encoding/hex"
)

//GetMD5 获取MD5
func GetMD5(src string) string  {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := md5Ctx.Sum(nil)
	sign := hex.EncodeToString(cipherStr)
	return sign
}
