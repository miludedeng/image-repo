package service

import (
	"image-repo/models"
	"time"

	"github.com/astaxie/beego/orm"
)

// 初始化数据库
func InitData() {
	user := &models.User{"xxxxx", GeneratePassword("xxxxxx"), time.Now()}
	o := orm.NewOrm()
	o.Insert(user)
}
