package httputil

import (
	"github.com/yhhaiua/engine/util/cast"
	"github.com/yhhaiua/engine/util/treemap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Get(url string) []byte {

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.Get(url)
	if err != nil {
		logger.Errorf("Get error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Get ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}

func GetTree(url string, tree *treemap.Map) []byte {

	var b strings.Builder
	b.Grow(128)
	b.WriteString(url + "?")

	i := 0
	count := tree.Size()

	it := tree.Iterator()
	for it.Next() {
		i++
		key := it.Key()
		value := it.Value()
		b.WriteString(cast.ToString(key))
		b.WriteString("=")
		b.WriteString(cast.ToString(value))
		if i != count {
			b.WriteString("&")
		}
	}

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.Get(b.String())
	if err != nil {
		logger.Errorf("Get error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Get ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}
func GetTreeComplex(url string, tree *treemap.Map) []byte {

	var b strings.Builder
	b.Grow(128)
	b.WriteString(url + "?")

	i := 0
	count := tree.Size()

	it := tree.Iterator()
	for it.Next() {
		i++
		key := it.Key()
		value := it.Value()
		b.WriteString(cast.ToString(key))
		b.WriteString("=")
		b.WriteString(cast.ToString(value))
		if i != count {
			b.WriteString("&")
		}
	}
	return GetFormMapComplex(b.String())
}

// GetFormMapComplex 复杂请求，返回json字符串 带中文
func GetFormMapComplex(url string) []byte {

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	req.Header.Set("Content-Type", "text/html;charset=UTF-8")
	req.Header.Set("Accept", "text/html;charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Get error : %s", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Get ReadAll error : %s", err.Error())
		return nil
	}
	return body
}

func GetMap(url string, tree map[string]string) []byte {
	var b strings.Builder
	b.Grow(128)
	b.WriteString(url + "?")
	i := 0
	count := len(tree)
	for k, v := range tree {
		i++
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
		if i != count {
			b.WriteString("&")
		}
	}
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.Get(b.String())
	if err != nil {
		logger.Errorf("Get error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Get ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}

// GetMapTime http请求，设置超时响应时间
func GetMapTime(url string, tree map[string]string, second int) []byte {
	var b strings.Builder
	b.Grow(128)
	b.WriteString(url + "?")
	i := 0
	count := len(tree)
	for k, v := range tree {
		i++
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
		if i != count {
			b.WriteString("&")
		}
	}
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(second) * time.Second,
	}
	req, err := client.Get(b.String())
	if err != nil {
		logger.Errorf("Get error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Get ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}
