package controllers

import (
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/encrypt"
	"Spike-Product-Demo/services"
	"Spike-Product-Demo/tool"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
	Session     *sessions.Session
}

func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (u *UserController) PostRegister() {
	// 获取表单内容
	var (
		nickName = u.Ctx.FormValue("NickName")
		userName = u.Ctx.FormValue("UserName")
		passWord = u.Ctx.FormValue("Password")
	)

	// 添加到user表单
	user := &datamodels.User{
		NickName:     nickName,
		UserName:     userName,
		HashPassword: passWord,
	}
	_, err := u.UserService.AddUser(user)
	u.Ctx.Application().Logger().Debug(err) //为什么没有出错也要debug
	if err != nil {
		u.Ctx.Redirect("/user/error")
		return
	}
	u.Ctx.Redirect("/user/login")
	return
}

func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (u *UserController) PostLogin() mvc.Response {
	// 1. 获取用户表单信息
	var (
		userName = u.Ctx.FormValue("UserName")
		passWd   = u.Ctx.FormValue("Password")
	)
	user, err := u.UserService.IsPwdSucceed(userName, passWd)

	// 2.1 核验结果失败
	if err != nil {
		u.Ctx.Application().Logger().Debug(err)
		return mvc.Response{
			Path: "/user/login",
		}
	}
	// 2.2 密码输入正确:将用户ID写入到cookie中
	tool.GlobalCookie(u.Ctx, "uid", strconv.FormatInt(user.ID, 10)) //直接存入cookie
	uidByte := []byte(strconv.FormatInt(user.ID, 10))               // AES加密后存入cookie
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		u.Ctx.Application().Logger().Debug(err)
	}
	tool.GlobalCookie(u.Ctx, "signuid", uidString)
	// u.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	return mvc.Response{
		Path: "/product/",
	}
}
