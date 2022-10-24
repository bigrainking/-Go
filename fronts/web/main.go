package main

import (
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/fronts/web/controllers"
	"Spike-Product-Demo/repository"
	"Spike-Product-Demo/services"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func main() {
	// 1. =================框架
	app := iris.New()

	// 2. =================设置错误等级
	app.Logger().SetLevel("debug")

	// 3. =================设置模板
	// 3.1 注册动态模板
	tmplate := iris.HTML("./views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	// 3.2 设置静态模板
	app.HandleDir("/public", "./public") //请求路径，实际文件夹中路径

	// 4. =================设置异常errorCode页面跳转
	app.OnAnyErrorCode(errorHandler)

	// 5 注册控制器
	initController(app)

	// 6. =====================启动
	app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(),
	)
}

func errorHandler(ctx iris.Context) {
	ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问页面出错"))
	ctx.ViewLayout("")
	ctx.View("shared/error.html") //将数据渲染到页面中:注册模板时已经告知了模板的位置
}

func initController(app *iris.Application) {
	// 5.=================== 注册控制器
	db, err := common.NewMysqlConn()
	if err != nil {
		app.Logger().Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// session
	sess := sessions.New(sessions.Config{
		Cookie: "AdminCookie",
	})
	//  product控制器
	repoUser := repository.NewUserManagerRepo("spikeSystem.user", db)
	serviceUser := services.NewUserServiceManager(repoUser)
	user := mvc.New(app.Party("/user"))
	user.Register(ctx, serviceUser, sess.Start)
	user.Handle(new(controllers.UserController))
}
