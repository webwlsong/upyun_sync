package controllers

import (
	"fmt"
	"reflect"
	"upyun_sync/models"
	"upyun_sync/upyun"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Index() {

	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Ctx.Redirect(302, "/login")
	}
	fmt.Println(reflect.TypeOf(userinfo))
	user := new(models.Space)
	//Magic code:: reflect to get all field
	user.Name = reflect.ValueOf(userinfo).Elem().Field(0).String()
	user.UserName = reflect.ValueOf(userinfo).Elem().Field(1).String()
	user.PassWord = reflect.ValueOf(userinfo).Elem().Field(2).String()
	u := upyun.NewUpYun(user.Name, user.UserName, user.PassWord)
	u.Debug = true
	u.SetApiDomain(upyun.EdAuto)
	v, err := u.GetBucketUsage()
	if err != nil {
		fmt.Println(err)
		return
	}
	dirs, err := u.ReadDir("/")
	fmt.Printf("ReadDir: %v\n", err)
	for i, d := range dirs {
		fmt.Printf("\t%d: %v\n", i, d)
	}

	this.Data["UsedSize"] = v / 1024 / 1024
	files, err := u.ReadDir("/")
	if err != nil {
		fmt.Println(err)
		return
	}
	this.Data["Files"] = files
	this.TplNames = "index.html"
}

func (this *MainController) Login() {
	userinfo := this.GetSession("userinfo")
	fmt.Println(userinfo)
	if userinfo != nil {
		this.Ctx.Redirect(302, "/index")
	}
	this.TplNames = "login.html"
	spacename := this.GetString("spacename")
	username := this.GetString("username")
	fmt.Println(username)
	if username != "" && spacename != "" {
		password := this.GetString("password")
		user := new(models.Space)
		user.Name = spacename
		user.UserName = username
		user.PassWord = password
		this.SetSession("userinfo", user)
		this.Ctx.Redirect(302, "/index")
		return
	}
	return
}
func (this *MainController) Logout() {
	this.DelSession("userinfo")
	this.Ctx.Redirect(302, "/login")
}
