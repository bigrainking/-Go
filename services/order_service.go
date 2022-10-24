package services

import (
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/repository"
)

// 订单的逻辑函数

type OrderService interface {
	// 增删改查
	InsertOrder(order *datamodels.Order) (id int64, err error)
	DeleteOrder(id int64) error
	UpdateOrder(order *datamodels.Order) error
	GetOrderByID(id int64) (*datamodels.Order, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllWithInfo() ([]map[string]string, error) //获取订单号\商品名\发货状态
}

type OrderManager struct {
	repoOrder repository.Order
}

func NewOrderServiceManager(repo repository.Order) OrderService {
	return &OrderManager{repo}
}

func (o *OrderManager) InsertOrder(order *datamodels.Order) (id int64, err error) {
	return o.repoOrder.Insert(order)
}
func (o *OrderManager) DeleteOrder(id int64) error {
	return o.repoOrder.Delete(id)
}
func (o *OrderManager) UpdateOrder(order *datamodels.Order) error {
	return o.repoOrder.Update(order)
}
func (o *OrderManager) GetOrderByID(id int64) (*datamodels.Order, error) {
	return o.repoOrder.SelectByID(id)
}
func (o *OrderManager) GetAllOrder() ([]*datamodels.Order, error) {
	return o.repoOrder.SelectAll()
}
func (o *OrderManager) GetAllWithInfo() ([]map[string]string, error) { //获取订单号\商品名\发货状态
	return o.repoOrder.SelectAllWithInfo()
}
