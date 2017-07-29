package service

import (
	"github.com/fasthttp-contrib/sessions"
	"github.com/valyala/fasthttp"
)

// Session 用户Session
var Session sessions.Session

// InitSession 初始化Session
func InitSession(ctx *fasthttp.RequestCtx) {
	Session = sessions.StartFasthttp(ctx)
}
