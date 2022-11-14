// 逻辑处理层
package services

import (
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/repository"
)

// 处理来自Controller的请求，调用repository操作数据库，根据Controller的要求
// 将repository中的内容处理后返回给Controller

type IProductService interface {
	// 增删查改
	// 获取产品
	GetProductByID(ID int64) (product *datamodels.Product, err error)
	// 获取所有产品
	GetAllProduct() (allProducts []*datamodels.Product, err error)
	// 增删改
	InsertProduct(product *datamodels.Product) (ID int64, err error)
	DeleteProductByID(ID int64) error
	UpdateProduct(product *datamodels.Product) error
	SubProductNum(productID int64) error //商品数量-1
}

type ProductServiceManager struct {
	repositoryProduct repository.IProduct
}

func NewIPoductSeviceManager(repositoryProduct repository.IProduct) IProductService {
	return &ProductServiceManager{repositoryProduct}
}

// 实现接口
// 下面都是直接返回repository操作数据的结果
func (pService *ProductServiceManager) InsertProduct(product *datamodels.Product) (int64, error) {
	return pService.repositoryProduct.Insert(product)
}
func (pService *ProductServiceManager) DeleteProductByID(id int64) error {
	return pService.repositoryProduct.Delete(id)
}
func (pService *ProductServiceManager) UpdateProduct(product *datamodels.Product) error {
	return pService.repositoryProduct.Update(product)
}
func (pService *ProductServiceManager) GetProductByID(id int64) (*datamodels.Product, error) {
	return pService.repositoryProduct.SearchById(id)
}
func (pService *ProductServiceManager) GetAllProduct() ([]*datamodels.Product, error) {
	return pService.repositoryProduct.SearchAll()
}
func (pService *ProductServiceManager) SubProductNum(productID int64) error {
	return pService.repositoryProduct.SubProductNum(productID)
}
