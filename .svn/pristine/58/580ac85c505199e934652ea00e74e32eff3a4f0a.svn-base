package httputil

import (
	"bytes"
	"encoding/json"
	"engine/jsonx"
	"engine/log"
	"engine/util/cast"
	"engine/util/treemap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var logger = log.GetLogger()

func Post(url string, body io.Reader) []byte {

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}

func PostForm(url string, data url.Values) []byte {

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.PostForm(url, data)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}

func PostFormTreeJson(url string, tree *treemap.Map) []byte {
	m := MapToUrlValue3(tree)
	data, err := json.Marshal(m)
	if err != nil {
		logger.Errorf("%s", err.Error())
		return nil
	}
	return PostJson(url, data)
}
func PostFormTree(url string, tree *treemap.Map) []byte {
	u := MapToUrlValue1(tree)

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.PostForm(url, u)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}
func PostJson(url string, data []byte) []byte {
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}

func PostFormJson(url string, tree interface{}) []byte {
	data, err := json.Marshal(tree)
	if err != nil {
		logger.Errorf("%s", err.Error())
		return nil
	}
	return PostJson(url, data)
}
func PostFormMap(url string, tree map[string]string) []byte {
	u := MapToUrlValue2(tree)
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := client.PostForm(url, u)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return bodys
}

func PostFormMapComplexHeader(url string, header, tree map[string]string) []byte {
	u := MapToUrlValue2(tree)
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(u.Encode()))
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return body
}
func PostFormMapjsonHeader(url string, header, tree map[string]string) []byte {
	data := jsonx.Marshal(tree)
	//u := MapToUrlValue2(tree)
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return body
}
func PostFormMapjsonHeader2(url string, header map[string]string, data []byte) []byte {
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return body
}

// PostFormMapComplex 复杂请求，返回json字符串 带中文
func PostFormMapComplex(url string, tree map[string]string) []byte {
	u := MapToUrlValue2(tree)
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(u.Encode()))
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Accept", "application/json;charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Post error : %s", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s", err.Error())
		return nil
	}
	return body
}

// MapToUrlValue1 生成url
func MapToUrlValue1(tree *treemap.Map) url.Values {
	it := tree.Iterator()
	u := make(url.Values)
	for it.Next() {
		key := it.Key()
		value := it.Value()
		u.Set(cast.ToString(key), cast.ToString(value))
	}
	return u
}

// MapToUrlValue2 生成url
func MapToUrlValue2(tree map[string]string) url.Values {
	u := make(url.Values)
	for key, value := range tree {
		u.Set(key, value)
	}
	return u
}

func MapToUrlValue3(tree *treemap.Map) map[string]interface{} {
	it := tree.Iterator()
	u := make(map[string]interface{})
	for it.Next() {
		key := it.Key()
		value := it.Value()
		u[cast.ToString(key)] = value
	}
	return u
}
