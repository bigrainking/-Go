package main

import (
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/fronts/web/controllers"
	"Spike-Product-Demo/fronts/web/middleware"
	"Spike-Product-Demo/rabbitmq"
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
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	// 3.2 设置静态模板
	app.HandleDir("/public", "./web/public") //请求路径，实际文件夹中路径
	app.HandleDir("/html", "./web/htmlProductShow")
	// 4. =================设置异常errorCode页面跳转
	app.OnAnyErrorCode(errorHandler)

	// 5 注册控制器
	initController(app)

	// 6. =====================启动
	app.Run(
		// iris.Addr(":8081"),
		iris.Addr(":8081"),
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

	// 后面cookie优化，不需要session了
	// session
	// // session中保存一个cookie的名字
	// sess := sessions.New(sessions.Config{
	// 	Cookie: "AdminCookie",
	// })
	//  user控制器
	repoUser := repository.NewUserManagerRepo("spikeSystem.user", db)
	serviceUser := services.NewUserServiceManager(repoUser)
	user := mvc.New(app.Party("/user"))
	user.Register(ctx, serviceUser /*,sess.Start*/)
	user.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewRabbitMQsimple("spikeProduct")
	// productshow控制器
	repoProduct := repository.NewProductManager("spikeSystem.product", db)
	serviceProduct := services.NewIPoductSeviceManager(repoProduct)
	repoOrder := repository.NewOrderManagerRepo("spikeSystem.order", db)
	serviceOrder := services.NewOrderServiceManager(repoOrder)
	product := mvc.New(app.Party("/product"))
	product.Router.Use(middleware.AuthConProduct)
	product.Register(ctx, serviceProduct, serviceOrder, rabbitmq) //, sess.Start)
	product.Handle(new(controllers.ProductController))
}
