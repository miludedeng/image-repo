package controller

import (
	"encoding/json"
	"fmt"
	"image-repo/models"
	"image-repo/service"

	"github.com/valyala/fasthttp"
)

// Login 登录
func Login(ctx *fasthttp.RequestCtx) {
	type LoginResult struct {
		Result string
	}
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")
	ctx.Response.Header.Set("Access-Control-Max-Age", "1000")
	args := ctx.PostArgs()
	username := string(args.Peek("username"))
	password := string(args.Peek("password"))
	callBack := string(args.Peek("callback"))
	if username == "" {
		args = ctx.QueryArgs()
		username = string(args.Peek("username"))
		password = string(args.Peek("password"))
		callBack = string(args.Peek("callback"))
	}
	isPass := service.VerifyPassword(password, service.GetChiper(username))
	var res string
	if isPass {
		fmt.Printf("%s login success!\n", username)
		fmt.Printf("Session %v\n", service.Session)
		service.Session.Set("user", &models.User{
			UserName: username,
			Password: password,
		})
		res = SUCCESS
	} else {
		fmt.Printf("%s login failed!\n", username)
		res = FAIL
	}
	data := &LoginResult{res}
	jsonData, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Errorf("result parse error: %s", err)
		fmt.Println(msg)
	}

	result := callBack + "(" + string(jsonData) + ")"
	ctx.Write([]byte(result))
}

// LoginInfo 登录信息
func LoginInfo(ctx *fasthttp.RequestCtx) {
	type LoginInfo struct {
		IsLogin  bool
		UserName string
	}
	args := ctx.QueryArgs()
	callBack := string(args.Peek("callback"))
	var user *models.User
	userInterface := service.Session.Get("user")
	switch userInterface.(type) {
	case *models.User:
		user = userInterface.(*models.User)
	}
	var data *LoginInfo
	if user != nil {
		data = &LoginInfo{true, user.UserName}
	} else {
		data = &LoginInfo{false, ""}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Errorf("result parse error: %s", err)
		fmt.Println(msg)
	}
	result := callBack + "(" + string(jsonData) + ")"
	ctx.Write([]byte(result))
}

// Logout 退出
func Logout(ctx *fasthttp.RequestCtx) {
	type LogoutResult struct {
		Status string
	}
	args := ctx.QueryArgs()
	callBack := string(args.Peek("callback"))
	service.Session.Set("user", nil)
	var data *LogoutResult
	data = &LogoutResult{"success"}

	jsonData, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Errorf("result parse error: %s", err)
		fmt.Println(msg)
	}
	result := callBack + "(" + string(jsonData) + ")"
	ctx.Write([]byte(result))
}
