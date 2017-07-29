package service

import (
	"bytes"
	"fmt"
	"image-repo/models"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

// StorePath 图片存储路径
var StorePath string

// GetList 获取图片列表
func GetList(pageSize int64, pageIndex int64) (total int64, images []*models.Image) {
	start := pageSize * pageIndex
	o := orm.NewOrm()
	o.Raw("SELECT count(1) FROM image").QueryRow(&total)
	num, err := o.Raw("SELECT * FROM image order by id desc limit ?,?", start, pageSize).QueryRows(&images)
	if err == nil {
		fmt.Println("image nums: ", num)
	} else {
		fmt.Println(err)
	}
	return total, images
}

// SaveFile 保存文件
func SaveFile(f *multipart.FileHeader) bool {
	oldName := f.Filename
	suffix := oldName[strings.LastIndex(oldName, "."):]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	part := fmt.Sprintf("%x", r.Intn(9999))
	for len(part) < 4 {
		part = "0" + part
	}
	filePreName := fmt.Sprintf("%x", time.Now().UnixNano()) + "_" + part
	fileName := filePreName + suffix
	fh, err := f.Open()
	if err != nil {
		fmt.Println("open upload file error:", err)
		return false
	}
	defer fh.Close() // 记得要关

	// 打开保存文件句柄
	fp, err := os.OpenFile(StorePath+string(os.PathSeparator)+fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		msg := fmt.Errorf("open saving file error:%s", err)
		fmt.Println(msg)
		return false
	}
	defer fp.Close() // 记得要关

	if _, err = io.Copy(fp, fh); err != nil {
		msg := fmt.Errorf("save upload file error:%s", err)
		fmt.Println(msg)
		return false
	}
	if err != nil {
		msg := fmt.Errorf("id error :%s", err)
		fmt.Println(msg)
		return false
	}
	ScaleImage(filePreName, suffix)
	image := &models.Image{Path: fileName, AddTime: time.Now()}
	o := orm.NewOrm()
	o.Insert(image)
	return true
}

// DeleteFile 删除文件
func DeleteFile(fileName string) {
	suffix := fileName[strings.LastIndex(fileName, "."):]
	name := fileName[:strings.LastIndex(fileName, ".")]
	thumbnail := name + "_thumbnail" + suffix
	os.Remove(StorePath + fileName)
	os.Remove(StorePath + thumbnail)
}

// ScaleImage 缩放图片
func ScaleImage(fileName string, suffix string) {
	source := fmt.Sprintf("%s%s%s", StorePath, fileName, suffix)
	dest := fmt.Sprintf("%s%s%s%s", StorePath, fileName, "_thumbnail", suffix)
	cmd := exec.Command("gm", "convert", source,
		"-resize", "350x200", "+profile", "\"*\"", dest)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	fmt.Println(cmd.Args)
	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	fmt.Println(out.String())
}
