package geetest

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"strings"
)

// 拼接Get请求参数
func makeUrl(link string, params url.Values) string {
	if params == nil {
		return link
	}
	if strings.Contains(link, "?") {
		if link[len(link)-1:] == "?" { // 最后一个就是？
			link += params.Encode() // 直接拼接
		} else {
			link += "&" + params.Encode()
		}
	} else {
		link += "?" + params.Encode()
	}
	return link
}

func makeMD5(p []byte) string {
	hash := md5.New()
	hash.Write(p)
	return hex.EncodeToString(hash.Sum(nil))
}

func makeSH1(p []byte) string {
	hash := sha1.New()
	hash.Write(p)
	return hex.EncodeToString(hash.Sum(nil))
}
