package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"blog/models"
	"blog/libs"
)

type LoginController struct {
	BaseController
}

//登录 TODO:XSRF过滤
func (self *LoginController) LoginIn() {
	//if self.userId > 0 {
	//	self.redirect(beego.URLFor("HomeController.Index"))
	//}
	beego.ReadFromRequest(&self.Controller)

	if self.QisPost() {
		username := strings.TrimSpace(self.GetString("username"))
		password := strings.TrimSpace(self.GetString("password"))

		if username != "" && password != "" {
			user, err := models.AdminGetByName(username)
			fmt.Println(user)
			flash := beego.NewFlash()
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == -1 {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIp = self.QgetClientIp()
				user.LastLogin = time.Now().Unix()
				user.Update()
				authkey := libs.Md5([]byte(self.QgetClientIp() + "|" + user.Password + user.Salt))
				self.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)

				self.redirect(beego.URLFor("HomeController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&self.Controller)
			self.redirect(beego.URLFor("LoginController.LoginIn"))
		}
	}
	self.TplName = "login/login.html"
}

//登出
func (self *LoginController) LoginOut() {
	self.Ctx.SetCookie("auth", "")
	self.redirect(beego.URLFor("LoginController.LoginIn"))
}

func (self *LoginController) NoAuth() {
	self.Ctx.WriteString("没有权限")
}