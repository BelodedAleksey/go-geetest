package geetest

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"time"
	// "net/http/cookiejar"
)

var cli = new(http.Client)

func init() {
	// jar, _ := cookiejar.New(nil)
	// cli.Jar = jar

	// 统一超时 时间
	cli.Timeout = time.Second * 30
}

func httpGet(addr string) (string, error) {
	var err error

	// 请求对象
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return "", err
	}

	// 请求对象设置
	// ......

	// 发送请求 并返回
	var res *http.Response
	res, err = cli.Do(req)
	if err != nil {
		return "", err
	}

	// 读取内容
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func httpPost(addr string, val url.Values) (string, error) {
	var err error

	// post数据
	postParams := strings.NewReader("")
	if val != nil {
		postParams = strings.NewReader(val.Encode())
	}

	// 请求对象
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, addr, postParams)
	if err != nil {
		return "", err
	}

	// 请求对象设置
	// ......

	// 请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求 并返回
	var res *http.Response
	res, err = cli.Do(req)
	if err != nil {
		return "", err
	}

	// 读取内容
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
