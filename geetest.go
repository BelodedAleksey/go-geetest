package geetest

import (
	"math"
	"net/url"
	"strings"

	"github.com/json-iterator/go"
	"github.com/lemon-cn/go-toolkit/random"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	FN_CHALLENGE          string = "geetest_challenge"
	FN_VALIDATE                  = "geetest_validate"
	FN_SECCODE                   = "geetest_seccode"
	GT_STATUS_SESSION_KEY        = "gt_server_status"
	API_URL                      = "http://api.geetest.com"
	REGISTER_HANDLER             = "/register.php"
	VALIDATE_HANDLER             = "/validate.php"
	GT_SDK_VERSION               = "Go.gt3-0.1.0"
)

type App struct {
	captchaID  string
	privateKey string
	res        *Response
}

type Response struct {
	Status    bool   `json:"-"`
	Success   int    `json:"success"`
	GT        string `json:"gt"`
	Challenge string `json:"challenge"`
}

func (res *Response) Marshal() []byte {
	b, _ := json.Marshal(res)
	return b
}

func (res *Response) String() string {
	return string(res.Marshal())
}

func New(id, key string) *App {
	return &App{
		id,
		key,
		nil,
	}
}

func (app *App) PreProcess(user_id string) *Response {
	addr := makeUrl(API_URL+REGISTER_HANDLER, url.Values{
		"gt":      []string{app.captchaID},
		"user_id": []string{user_id},
	})

	back := new(Response)

	res, err := httpGet(addr)
	if err != nil {
		back.Status = false
		back.Success = 0
		back.GT = app.captchaID
		back.Challenge = makeSH1(random.RandBytes(10))[:34] // 仅取34位
		return back
	}

	// 得到数据
	back.Status = true
	back.Success = 1
	back.GT = app.captchaID
	back.Challenge = makeMD5([]byte(res + app.privateKey))

	return back
}

// 正常模式获取验证结果
func (app *App) SuccessValidate(challenge, validate, seccode string, user_id string) bool {
	if !app.checkValidate(challenge, validate) {
		return false
	}

	data := url.Values{}
	data.Set("seccode", seccode)
	data.Set("sdk", GT_SDK_VERSION)
	if len(user_id) > 0 {
		data.Set("user_id", user_id)
	}

	res, err := httpPost(API_URL+VALIDATE_HANDLER, data)
	if err != nil || res == "false" {
		return false
	}

	return res == makeMD5([]byte(seccode))
}

// 宕机模式获取验证结果
func (App) FailValidate(challenge, validate, seccode string) bool {
	if len(validate) > 0 {
		arr := strings.Split(validate, "_")
		ans := decodeResponse(challenge, arr[0])
		bg_idx := decodeResponse(challenge, arr[1])
		grp_idx := decodeResponse(challenge, arr[2])
		x_pos := getFailbackPicAns(bg_idx, grp_idx)
		answer := int(math.Abs(float64(ans - x_pos)))

		return answer < 4
	}

	return false
}

func (app *App) checkValidate(challenge, validate string) bool {
	if len(validate) != 32 {
		return false
	}

	return makeMD5([]byte(app.privateKey+"geetest"+challenge)) == validate
}
