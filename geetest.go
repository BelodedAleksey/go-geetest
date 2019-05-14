package geetest

import (
	"bytes"
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
	GT_SDK_VERSION               = "Go_3.0.0"
)

type App struct {
	captchaID  string
	privateKey string
	res        *Response
}

type Response struct {
	Success   int    `json:"success"`
	GT        string `json:"gt"`
	Challenge string `json:"challenge"`
}

func (res *Response) Marshal() string {
	b, _ := json.Marshal(res)
	return string(b)
}

func New(id, key stirng) *App {
	return &App{
		id,
		key,
	}
}

func (app *App) PreProcess(user_id string) *Response {
	addr := makeUrl(API_URL+REGISTER_HANDLER, url.Values{
		"gt":      app.captchaID,
		"user_id": user_id,
	})

	back := new(Response)

	res, err := httpGet(addr)
	if err != nil {
		back.Success = 0
		back.GT = app.captchaID
		back.Challenge = makeSH1(random.RandBytes(10))[:34] // 仅取34位
		return back
	}

	// 得到数据
	back.Success = 1
	back.GT = app.captchaID
	back.Challenge = makeMD5([]byte(res + app.privateKey))

	return back
}

// 正常模式获取验证结果
func (app *App) SuccessValidate(challenge, validate, seccode string, user_id string) int {
	if !app.checkValidate(challenge, validate) {
		return 0
	}

	data := url.Values{
		"seccode": seccode,
		"sdk":     GT_SDK_VERSION,
	}
	if len(user_id) > 0 {
		data.Set("user_id", user_id)
	}

	res, err := httpPost(API_URL+VALIDATE_HANDLER, data)
	if err != nil {
		return 0
	}

	if res == makeMD5([]byte(seccode)) {
		return 1
	}

	return 0
}

// /**
//  * 宕机模式获取验证结果
//  *
//  * @param $challenge
//  * @param $validate
//  * @param $seccode
//  * @return int
//  */
// public function fail_validate($challenge, $validate, $seccode) {
//     if ($validate) {
//         $value   = explode("_", $validate);
//         $ans     = $this->decode_response($challenge, $value['0']);
//         $bg_idx  = $this->decode_response($challenge, $value['1']);
//         $grp_idx = $this->decode_response($challenge, $value['2']);
//         $x_pos   = $this->get_failback_pic_ans($bg_idx, $grp_idx);
//         $answer  = abs($ans - $x_pos);
//         if ($answer < 4) {
//             return 1;
//         } else {
//             return 0;
//         }
//     } else {
//         return 0;
//     }
// }
// 宕机模式获取验证结果
func (app *App) FailValidate(challenge, validate, seccode string) int {
	if len(validate) > 0 {
		value := strings.Split(validate, "_")

		// $value   = explode("_", $validate);
		// $ans     = $this->decode_response($challenge, $value['0']);
		// $bg_idx  = $this->decode_response($challenge, $value['1']);
		// $grp_idx = $this->decode_response($challenge, $value['2']);
		// $x_pos   = $this->get_failback_pic_ans($bg_idx, $grp_idx);
		// $answer  = abs($ans - $x_pos);
		// if ($answer < 4) {
		//     return 1;
		// } else {
		//     return 0;
		// }
	}

	return 0
}

func (app *App) checkValidate(challenge, validate string) bool {
	if len(validate) != 32 {
		return false
	}

	return makeMD5([]byte(app.privateKey+"geetest"+challenge)) == validate
}
