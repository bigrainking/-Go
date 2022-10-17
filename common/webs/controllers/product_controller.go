package controller

import (
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/repository"
	service "Spike-Product-Demo/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// =====================创建Controller对象(类似于MVC里面的engine)
// struct中的对象都在main中注册好，类似Engine中的Context都在main中注册好
type ProductController struct {
	Ctx            iris.Context            //这个Controller的上下文
	ProductService service.IProductService // Service逻辑处理
}

func NewProductController() *ProductController {
	db, _ := common.NewMysqlConn()
	repositoryProduct := repository.NewProductManager("product", db)
	return &ProductController{ProductService: service.NewIRoductSeviceManager(repositoryProduct)}
}

// ----------尝试用NewController的方式创建，并初始化struct----

// ===================相关匹配方法

// 1. 显示所有商品 Get /product/all 显示商品详情页面：获取所有商品
// 自动匹配：用Get方法访问/product/all会调用该方法
func (p *ProductController) GetAll() mvc.View {
	// 调用Service去获取数据
	// productArray, _ := p.ProductService.GetAllProduct()
	// 渲染view
	productArray, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray": productArray,
		},
	}
}

// 2. 商品管理页面 GET /product/manager：展示多条商品，页面有修改and删除操作按钮
func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// 3. 修改商品功能 POST /product/update ： 提交修改商品表单
func (p *ProductController) PostUpdate() {
	// 1. 处理提交的表单，解析表单，将表单中的数据填充到product结构体里面
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()                                       //解析上传的表单
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "form"}) // 通过form解析传入的表单
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil { //将表单中的内容填充到product
		p.Ctx.Application().Logger().Debug(err) //debug级别的error
	}
	// 2. 更新数据库中的商品信息
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	// 3. 更新完毕后跳转到指定页面
	p.Ctx.Application().Logger().Info("成功修改商品, ID为", product.ID)
	p.Ctx.Redirect("/product/all")
}

// 4. 删除商品 GET localhost:8080/product/delete
func (p *ProductController) GetDelete() {
	// 1. 获取商品ID
	idString := p.Ctx.URLParam("id")              //json中是id
	id, err := strconv.ParseInt(idString, 10, 16) //字符串转换：字符串，字符串进制，返回结果bit的大小
	if err != nil {
		p.Ctx.Application().Logger().Debug("GetDelete：数字转换错误", err)
	}
	// 2. Service删除商品
	err = p.ProductService.DeleteProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug("删除商品失败，ID为：%i", id, err)
	}

	// 3. 跳转页面
	p.Ctx.Application().Logger().Info("成功删除商品, ID为", id)
	p.Ctx.Redirect("/product/all")
}

// 5. 添加商品页面 GET /product/add ：add.html
func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add.html",
	}
}

// 6. 添加商品按钮 POST /product/add：提交表单
func (p *ProductController) PostAdd() {
	// 1.获取HTML页面填写的表单信息 2.解析表单填充到数据库 3.跳转到展示所有商品页面
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()                                         // 获取表单并解析
	dec := common.NewDecoder((&common.DecoderOptions{TagName: "form"})) // 通过Tag将表单内容填充到product
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Application().Logger().Info("成功添加商品, ID为", product.ID)
	p.Ctx.Redirect("/product/all")
}
