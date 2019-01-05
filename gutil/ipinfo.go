package gutil

import (
	"net/http"
	"strings"
)

func GetUserIp(r *http.Request) string  {
	userip := r.Header.Get("X-Real-IP")
	if userip == ""{
		userip = r.RemoteAddr
	}
	info := strings.Split(userip,":")
	return info[0]
}
