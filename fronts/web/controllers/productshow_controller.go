package controllers

import (
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/rabbitmq"
	"Spike-Product-Demo/services"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

// 商品展示： 1. 商品页面 2. 商品秒杀功能

// controller对象：context、session、service、order
type ProductController struct {
	Ctx            iris.Context
	Session        *sessions.Session
	RabbitMq       *rabbitmq.RabbitMQ
	ProductService services.IProductService
	OrderService   services.IOrderService
}

var (
	// 保存html静态文件的目录
	htmlOutPath = "/home/yang/go/src/Spike-Product-Demo/fronts/web/htmlProductShow"
	// 模板文件目录
	templatePath = "/home/yang/go/src/Spike-Product-Demo/fronts/web/views/template"
)

// Get /detail/  获取商品详情
func (c *ProductController) GetDetail() mvc.View {
	product, err := c.ProductService.GetProductByID(6)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// Get
func (c *ProductController) GetOrder() []byte {
	// 获取URL中productID userid
	productIDString := c.Ctx.URLParam("productID")
	userIDString := c.Ctx.GetCookie("uid")
	userID, _ := strconv.ParseInt(userIDString, 10, 64)
	productID, _ := strconv.ParseInt(productIDString, 10, 64)

	// 3. 创建消息message实体
	message := datamodels.NewMessage(productID, userID)
	messageByte, err := json.Marshal(message)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	// 4. 发送到队列
	err = c.RabbitMq.PublishSimple(string(messageByte))
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}

	return []byte("true") //消息发送成功暂时返回true
	// // 获取product信息
	// product, err := c.ProductService.GetProductByID(int64(productID))
	// if err != nil {
	// 	c.Ctx.Application().Logger().Debug(err)
	// }
	// // 构建返回的订单
	// var orderID int64
	// message := "抢购失败"
	// if product.ProductNum > 0 { //可以抢购
	// 	product.ProductNum -= 1
	// 	// 更新product数据库
	// 	err := c.ProductService.UpdateProduct(product)
	// 	if err != nil {
	// 		c.Ctx.Application().Logger().Debug(err)
	// 	}
	// 	// 创建订单，更新数据库
	// 	order := &datamodels.Order{
	// 		UserID:      int64(userID),
	// 		ProductID:   int64(productID),
	// 		Orderstatus: datamodels.OrderSuccess,
	// 	}
	// 	orderID, err = c.OrderService.InsertOrder(order)
	// 	if err != nil {
	// 		c.Ctx.Application().Logger().Debug(err)
	// 	}
	// 	message = "抢购成功"
	// }
	// // 如果商品数量不足以抢购
	// return mvc.View{
	// 	Name:   "product/result.html",
	// 	Layout: "shared/productLayout.html",
	// 	Data: iris.Map{
	// 		"showMessage": message,
	// 		"orderID":     orderID,
	// 	},
	// }
}

// 访问创建静态文件
func (c *ProductController) GetGenerateHtml() {
	// 1. 获取html静态文件保存的路径
	fileName := filepath.Join(htmlOutPath, "htmlProductShow.html") //要生成的静态文件的名字
	// 2. 获取template
	contentsTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html")) // 解析模板，得到模板对象
	if err != nil {
		c.Ctx.Application().Logger().Error(err)
	}
	// 3. 获取Product
	productIDString := c.Ctx.URLParam("productID")
	productID, _ := strconv.Atoi(productIDString)
	product, err := c.ProductService.GetProductByID(int64(productID))
	// 4. 生成静态文件
	generateStaticHtml(c.Ctx, fileName, contentsTmp, product)
}
func generateStaticHtml(ctx iris.Context, fileName string, template *template.Template, product *datamodels.Product) {
	// 1. 判断file是否已经存在：如果存在则要从数据库中查询后重新生成静态文件
	if existFile(fileName) {
		// 存在则删除文件
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}
	// 2. 生成静态文件：创建文件，将template渲染后输出到文件
	//os.O_WRONLY只写打开,os.O_CREATE如果不存在则创建一个
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	// 3. 将Product渲染到文件中
	template.Execute(file, &product) //为什么要传入指针地址？因为produc是名字
}
func existFile(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
