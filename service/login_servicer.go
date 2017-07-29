package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image-repo/models"
	"math/rand"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/valyala/fasthttp"
)

// IsLogin 判断是否登录
func IsLogin(ctx *fasthttp.RequestCtx) bool {
	user := Session.Get("user")
	if user != nil {
		return true
	}
	return false
}

// 获取密码
func GetChiper(username string) string {
	o := orm.NewOrm()
	var user *models.User
	err := o.Raw("SELECT password from user WHERE user_name=?", username).QueryRow(&user)
	if err != nil {
		return ""
	}
	return user.Password
}

// VerifyPassword 验证密码
func VerifyPassword(source string, chiper string) bool {
	if len(chiper) != 48 || Md5(chiper[:16]+source) != chiper[16:] {
		return false
	}
	return true
}

// GeneratePassword 生成密码
func GeneratePassword(str string) string {
	var salt string
	for i := 0; i < 8; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		part := fmt.Sprintf("%x", r.Intn(255))
		if len(part) < 2 {
			part = "0" + part
		}
		salt = salt + part
	}
	return salt + Md5(salt+str)
}

// Md5 md5加密
func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(string(str)))
	cipherEncode := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherEncode)
}
