package geetest

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
)

//
func getFailbackPicAns(full_bg_index, img_grp_index int) int {
	full_bg_name := makeMD5([]byte(strconv.Itoa(full_bg_index)))[:9]
	bg_name := makeMD5([]byte(strconv.Itoa(img_grp_index)))[10:19]

	answer_decode := make([]byte, 9)
	for i := 0; i < 9; i++ {
		if i%2 == 0 {
			answer_decode = append(answer_decode, full_bg_name[i])
		} else if i%2 == 1 {
			answer_decode = append(answer_decode, bg_name[i])
		}
	}

	return getXPosFromStr(string(answer_decode[4:9]))
}

//
func getXPosFromStr(x_str string) int {
	if len(x_str) != 5 {
		return 0
	}

	sum, err := strconv.ParseInt(x_str, 16, 32)
	if err != nil {
		return 0
	}

	result := int(sum % 200)
	if result < 40 {
		return 40
	}

	return result
}

// 解码随机参数
func decodeResponse(challenge, str string) int {
	if len(str) > 100 {
		return 0
	}

	dict := map[int]int{
		0: 1,
		1: 2,
		2: 5,
		3: 10,
		4: 50,
	}

	count := 0
	keys := make(map[rune]int)
	exists := make([]rune, 20)

	for _, tmp := range challenge {

		if inRuneArray(tmp, exists) {
			continue
		}

		val := dict[count%5]
		exists = append(exists, tmp)
		count++
		keys[tmp] = val
	}

	res := 0
	for _, tmp := range str {
		res += keys[tmp]
	}

	return res - decodeRandBase(challenge)
}

// 输入的两位的随机数字,解码出偏移量
func decodeRandBase(challenge string) int {
	if len(challenge) < 34 {
		return 0
	}

	str := challenge[32:34]
	arr := make([]rune, 2)
	for _, char := range str {
		if char > 57 {
			arr = append(arr, char-87)
		} else {
			arr = append(arr, char-48)
		}
	}

	return int(arr[0]*36 + arr[1])
}

func inRuneArray(src rune, arr []rune) bool {
	for _, tmp := range arr {
		if tmp == src {
			return true
		}
	}

	return false
}

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
