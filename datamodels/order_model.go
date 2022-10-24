package datamodels

// 订单对应的结构体对象
type Order struct {
	ID          int64 `json:"id" sql:"ID" form:"ID"`
	UserID      int64 `json:"UserID" sql:"userID" form:"UserID"`
	ProductID   int64 `json:"ProductID" sql:"productID" form:"ProductID"`
	Orderstatus int64 `json:"OrderStatus" sql:"orderStatus" form:"OrderStatus"` //物流状态
}

// 订单状态
const (
	OrderWait    = iota
	OrderSuccess //1：已经成功发货
	OrderFaild   //2
)
