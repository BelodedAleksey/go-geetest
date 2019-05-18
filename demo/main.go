/**
 * 仅做演示，不要直接使用
 * 服务启动在 :8080
 */
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/lemon-cn/go-geetest"
	"github.com/satori/go.uuid"
)

const PORT string = ":8080"

const (
	ID  string = "b46d1900d0a894591916ea94ea91bd2c"
	KEY        = "36fc3fe98530eea08dfc6ce76e3d24c4"
)

// 对应关系为 客户端:用户
var sessions = make(map[string]*Session)
var lock = new(sync.Mutex)

type Session struct {
	UserID string // 用户ID
	Status bool   // 验证的获取状态
}

/**
 * 简单的sessions管理
 * 用于初始化 客户端与用户 对应关系
 *
 * @param  http.ResponseWriter
 * @param  *http.Request
 * @return *Session 会话数据
 */
func middleware(w http.ResponseWriter, r *http.Request) *Session {
	// 是否为新用户
	newUser := true

	var sess *Session
	sess_id := ""

	// 是否存在记录的客户端
	cookie, err := r.Cookie("sess")
	if err != nil || cookie == nil {
		// 创建新的session
		sess_id = newSN()
		// 记录到cookie
		cookie := new(http.Cookie)
		cookie.Name = "sess"
		cookie.Value = sess_id
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.MaxAge = 60 * 60
		http.SetCookie(w, cookie)
	} else {
		// 取出值
		sess_id = cookie.Value

		// 获取用户id
		var exists bool
		lock.Lock()
		sess, exists = sessions[sess_id]
		lock.Unlock()
		if exists && sess != nil && sess.UserID != "" {
			newUser = false
		}
	}

	// 新用户
	if newUser {
		// 生成唯一识别码（user_id）
		user_id := newSN()

		sess = new(Session)
		sess.UserID = user_id

		// 记录 对应关系
		lock.Lock()
		sessions[sess_id] = sess
		lock.Unlock()
	}

	return sess
}

func main() {
	app := geetest.New(ID, KEY)

	// 页面与静态资源
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 获取用户
		_ = middleware(w, r)

		http.ServeFile(w, r, "./static/login.html")
	})
	http.HandleFunc("/gt.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/gt.js")
	})

	// 接口
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		// 获取用户
		sess := middleware(w, r)

		//
		res := app.PreProcess(sess.UserID)
		sess.Status = res.Status
		w.Write(res.Marshal())
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		// 获取用户
		sess := middleware(w, r)

		challenge, validate, seccode := r.PostFormValue(geetest.FN_CHALLENGE), r.PostFormValue(geetest.FN_VALIDATE), r.PostFormValue(geetest.FN_SECCODE)

		var res bool
		if sess.Status {
			res = app.SuccessValidate(challenge, validate, seccode, sess.UserID)
		} else {
			res = app.FailValidate(challenge, validate, seccode)
		}

		// 返回
		if res {
			_, _ = w.Write([]byte(`{"status":"success"}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"fail"}`))
	})

	log.Println("服务运行中 ", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func newSN() string {
	sn, err := uuid.NewV4()
	if err != nil {
		return newSN()
	}

	return sn.String()
}
