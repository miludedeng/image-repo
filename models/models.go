package models


import "time"

type Image struct {
	Id        int64
	Path      string
	AddTime   time.Time
}

type User struct {
	UserName string  `orm:"pk"`
	Password string
	Addtime  time.Time
}