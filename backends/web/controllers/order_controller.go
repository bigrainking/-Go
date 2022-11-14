package controllers

import (
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// ====== 注册一个Controller的控制器
type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService //需要外部可使用，大写
}

// ====== Controller相关方法

// 显示所有订单 Get /order/all
func (o *OrderController) GetAll() mvc.View {
	orderArr, err := o.OrderService.GetAllWithInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArr,
		},
	}
}

// 添加商品页面 GET /order/add ：add.html
func (o *OrderController) GetAdd() mvc.View {
	return mvc.View{
		Name: "order/add.html",
	}
}

// 添加商品按钮 POST /order/add：提交表单
func (o *OrderController) PostAdd() {
	// 1.获取HTML页面填写的表单信息 2.解析表单填充到数据库 3.跳转到展示所有商品页面
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()                                         // 获取表单并解析
	dec := common.NewDecoder((&common.DecoderOptions{TagName: "form"})) // 通过Tag将表单内容填充到product
	if err := dec.Decode(o.Ctx.Request().Form, order); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	id, err := o.OrderService.InsertOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	} else {
		o.Ctx.Application().Logger().Info("成功添加订单, ID为", id)
	}
	o.Ctx.Redirect("/order/all")
}

// 修改商品页面 GET /order/manager：展示需要修改的商品的信息
func (o *OrderController) GetManager() mvc.View {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	order, err := o.OrderService.GetOrderByID(id)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "order/manager.html",
		Data: iris.Map{
			"order": order,
		},
	}
}

// 修改商品表单提交 POST /order/update ： 提交修改商品表单
func (o *OrderController) PostUpdate() {
	// 1. 处理提交的表单，解析表单，将表单中的数据填充到product结构体里面
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()                                       //解析上传的表单
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "form"}) // 通过form解析传入的表单
	if err := dec.Decode(o.Ctx.Request().Form, order); err != nil {   //将表单中的内容填充到product
		o.Ctx.Application().Logger().Debug(err) //debug级别的error
	}
	// 2. 更新数据库中的商品信息
	err := o.OrderService.UpdateOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	// 3. 更新完毕后跳转到指定页面
	o.Ctx.Application().Logger().Info("成功修改订单, ID为", order.ID) //通过解析form表单此处已经获得了ID
	o.Ctx.Redirect("/order/all")
}

// 删除商品 GET localhost:8080/order/delete
func (o *OrderController) GetDelete() {
	// 1. 获取商品ID
	idString := o.Ctx.URLParam("id")              //json中是id
	id, err := strconv.ParseInt(idString, 10, 16) //字符串转换：字符串，字符串进制，返回结果bit的大小
	if err != nil {
		o.Ctx.Application().Logger().Debug("GetDelete：数字转换错误", err)
	}
	// 2. Service删除商品
	err = o.OrderService.DeleteOrder(id)
	if err != nil {
		o.Ctx.Application().Logger().Debug("删除商品失败，ID为：%i", id, err)
	} else {
		o.Ctx.Application().Logger().Info("成功删除订单, ID为", id)
	}

	// 3. 跳转页面
	o.Ctx.Redirect("/order/all")
}
