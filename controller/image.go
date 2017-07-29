package controller

import (
	"encoding/json"
	"fmt"
	"image-repo/models"
	"image-repo/service"
	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/valyala/fasthttp"
)

// List 图片列表
func List(ctx *fasthttp.RequestCtx) {
	type ResultData struct {
		Total  int64
		Images []*models.Image
	}
	args := ctx.QueryArgs()
	callBack := string(args.Peek("callback"))
	pageIndex, err := strconv.ParseInt(string(args.Peek("pgi")), 10, 64)
	if err != nil {
		pageIndex = 0
	}
	pageSize, err := strconv.ParseInt(string(args.Peek("pgs")), 10, 64)
	if err != nil || pageSize == 0 {
		pageSize = 12
	}

	total, images := service.GetList(pageSize, pageIndex)
	data := &ResultData{Total: total, Images: images}
	jsonData, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Errorf("result parse error: %s", err)
		fmt.Println(msg)
	}
	result := callBack + "(" + string(jsonData) + ")"
	ctx.Write([]byte(result))
}

// Upload 上传图片
func Upload(ctx *fasthttp.RequestCtx) {
	//新建 需要登录
	type ResultData struct {
		Message string
	}
	if !service.IsLogin(ctx) {
		return
	}
	args := ctx.QueryArgs()
	callBack := string(args.Peek("callback"))
	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.SetStatusCode(500)
		fmt.Println("get upload file error:", err)
		return
	}
	result := service.SaveFile(f)
	data := &ResultData{}
	if result {
		data.Message = SUCCESS
	} else {
		data.Message = FAIL
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Errorf("result parse error: %s", err)
		fmt.Println(msg)
	}
	ctx.Write([]byte(callBack + "(" + string(jsonData) + ")"))

}

// Delete 删除图片
func Delete(ctx *fasthttp.RequestCtx) {
	//新建 需要登录
	type ResultData struct {
		Message string
	}
	data := &ResultData{}
	if !service.IsLogin(ctx) {
		return
	}
	args := ctx.QueryArgs()
	callBack := string(args.Peek("callback"))
	id, err := strconv.ParseInt(string(args.Peek("id")), 10, 64)
	if err != nil {
		data.Message = FAIL
	} else {
		image := &models.Image{Id: id}
		o := orm.NewOrm()
		o.Read(image)
		o.Delete(image)
		service.DeleteFile(image.Path)
	}
	data.Message = SUCCESS
	jsonData, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Errorf("result parse error: %s", err)
		fmt.Println(msg)
	}
	ctx.Write([]byte(callBack + "(" + string(jsonData) + ")"))

}
