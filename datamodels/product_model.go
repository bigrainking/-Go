package datamodels

// 创建商品对应的model
// json解析 表单解析 数据库中对应名称
type Product struct {
	ID           int64  `json:"id" sql:"ID" form:"ID"`
	ProductName  string `json:"ProductName" sql:"productName" form:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" form:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" form:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" form:"ProductUrl"`
}
