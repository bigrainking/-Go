package main

import (
	"Spike-Product-Demo/backends/web/controllers"
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/repository"
	"Spike-Product-Demo/services"
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
	tmplate := iris.HTML("./backends/web/views", ".html").Layout("shared/layout.html").Reload(true)
	// template := iris.HTML("./backends/webs/views", ".html").Layout("shared/layout.html") //.Reload(true) //layout的文件地址已经在HTML中被告知了，因此只需要给出相对地址
	app.RegisterView(tmplate)

	// 3.2 设置静态模板
	app.HandleDir("/assets", "./backends/web/assets") //请求路径 ， 实际文件夹中路径

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

	//  product控制器
	repoProduct := repository.NewProductManager("product", db)
	serviceProduct := services.NewIPoductSeviceManager(repoProduct)
	product := mvc.New(app.Party("/product"))
	product.Register(ctx, serviceProduct)
	product.Handle(new(controllers.ProductController))

	// order控制器
	repoOrder := repository.NewOrderManagerRepository("order", db)
	serviceOrder := services.NewOrderServiceManager(repoOrder)
	order := mvc.New(app.Party("/order"))
	order.Register(ctx, serviceOrder)
	order.Handle(new(controllers.OrderController))
}
