package geetest

import (
	"testing"
)

const (
	ID  string = "b46d1900d0a894591916ea94ea91bd2c"
	KEY        = "36fc3fe98530eea08dfc6ce76e3d24c4"
)

var app *App

func Test_NewApp(t *testing.T) {
	app = New(ID, KEY)
}

func Test_Register(t *testing.T) {
	user_id := "test"
	res := app.PreProcess(user_id)
	t.Log(res.Marshal())
}
