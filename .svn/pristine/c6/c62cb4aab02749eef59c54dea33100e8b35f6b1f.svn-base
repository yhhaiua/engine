package util

import (
	"net/http"
	"net/url"
	"strings"
)

func GetUserIp(r *http.Request) string  {
	userIp := r.Header.Get("X-Real-IP")
	if userIp == ""{
		userIp = r.RemoteAddr
	}
	info := strings.Split(userIp,":")
	return info[0]
}
//EncodeURIComponent
func EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}