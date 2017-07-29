package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"

	"image-repo/models"
	"image-repo/service"
	"image-repo/controller"
)

const (
	DB_DRIVER_NAME = "sqlite3"
	DB_DATABASE    = "default"
	DB_FILE        = "data/gallery.db"
)

func InitConn() {
	orm.RegisterDriver(DB_DRIVER_NAME, orm.DRSqlite)
	orm.RegisterDataBase(DB_DATABASE, DB_DRIVER_NAME, DB_FILE)
	orm.RegisterModel(new(models.User), new(models.Image))
	orm.Debug = true
	err := orm.RunSyncdb(DB_DATABASE, false, true) //建表
	if err != nil {
		fmt.Errorf("db register failed : %s", err)
	}
}

func main() {
	port := flag.String("p", "3000", "Server port")
	appName := flag.String("n", "", "AppName")
	storePath := flag.String("P", "/data", "Store Path")
	flag.Parse()
	service.StorePath = *storePath
	InitConn() //初始化数据库链接
	service.InitData()

	m := func(ctx *fasthttp.RequestCtx) {
		service.InitSession(ctx)
		ctx.SetContentType("application/json")
		path := string(ctx.Path())
		if *appName != "" {
			path = path[len(*appName)+1:]
		}
		switch path {
		case "/list":
			controller.List(ctx)
		case "/upload":
			controller.Upload(ctx)
		case "/delete":
			controller.Delete(ctx)
		case "/login":
			controller.Login(ctx)
		case "/login_info":
			controller.LoginInfo(ctx)
		case "/logout":
			controller.Logout(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	fmt.Printf("Store path is :%s\n", service.StorePath)
	fmt.Printf("Server name is :%s, port is :%s\n", *appName, *port)
	fasthttp.ListenAndServe(":"+*port, m)
}
