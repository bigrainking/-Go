package main

import (
	controller "Spike-Product-Demo/backends/webs/controllers"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// 1. =================框架
	app := iris.New()

	// 2. =================设置错误等级
	app.Logger().SetLevel("debug")

	// 3. =================设置模板
	// 3.1 注册动态模板
	tmplate := iris.HTML("./backends/webs/views", ".html").Layout("shared/layout.html").Reload(true)
	// template := iris.HTML("./backends/webs/views", ".html").Layout("shared/layout.html") //.Reload(true) //layout的文件地址已经在HTML中被告知了，因此只需要给出相对地址
	app.RegisterView(tmplate)

	// 3.2 设置静态模板
	app.HandleDir("/assets", "./backends/webs/assets") //请求路径 ， 实际文件夹中路径

	// 4. =================设置异常errorCode页面跳转
	app.OnAnyErrorCode(errorHandler)

	// db, _ := common.NewMysqlConn()
	ctx, cancel := context.WithCancel(context.Background()) //如何创建的context
	defer cancel()
	// 5.==================注册控制器
	product := mvc.New(app.Party("/product"))
	product.Register(ctx)                             //将context、service注册进路由组product
	product.Handle(controller.NewProductController()) //绑定对应的Controller

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
