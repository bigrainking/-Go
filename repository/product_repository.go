// Model层
// 这是商品的数据库操作
package repository

import (
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/datamodels"
	"database/sql"
	"strconv"
)

// 增删查改数据库内容
// 链接数据库

// 商品操作数据的相关函数：增删查改商品
type IProduct interface {
	Conn() error                                      //数据库链接
	Insert(*datamodels.Product) (int64, error)        // 插入数据，返回成功操作多少行
	Delete(id int64) error                            //删除数据,返回是否删除成功
	SearchById(id int64) (*datamodels.Product, error) //查找数据，返回查找出来的结构体
	SearchAll() ([]*datamodels.Product, error)        //查找所有数据，返回多个products对应结构体
	Update(*datamodels.Product) error                 //更新某条product
	SubProductNum(productID int64) error              //商品数量-1
}

// 实现接口
// table mysqlConn仅供内部使用
type ProductManager struct {
	table     string  //表的名字
	mysqlConn *sql.DB //数据连接
}

// 创建Product操作对象
func NewProductManager(table string, mysql *sql.DB) IProduct {
	return &ProductManager{table, mysql} //因为下面都是IProductManager的指针对象对应的方法
}

// 数据库链接
func (p *ProductManager) Conn() error {
	// 数据库连接
	if p.mysqlConn == nil {
		mysqlCon, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysqlCon
	}
	// 如果没有指定数据表名称
	if p.table == "" {
		p.table = "spikeSystem.product"
	}
	return nil
}

// 增
func (p *ProductManager) Insert(product *datamodels.Product) (ret int64, err error) { //返回插入数据的ID
	// 1. 链接数据库，检查链接是否成功
	if err = p.Conn(); err != nil {
		return
	}
	// 2. 执行增语句
	sql := "INSERT " + p.table + " SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	res, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return
	}
	// 3. 处理结果
	return res.LastInsertId() //已经包含了id 和 error
}

// 删
func (p *ProductManager) Delete(id int64) error { //删除数据,返回是否删除成功
	stmt, err := p.mysqlConn.Prepare("delete from " + p.table + " where ID=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(strconv.FormatInt(id, 10))
	if err != nil {
		return err
	}
	return nil
}

// 查询一行
func (p *ProductManager) SearchById(id int64) (*datamodels.Product, error) { //查找数据，返回查找出来的结构体
	if err := p.Conn(); err != nil {
		return nil, err
	}
	// 1. 语句
	sql := "select * from " + p.table + " where ID=" + strconv.FormatInt(id, 10)
	row, err := p.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return nil, err
	}
	// 3. 将查询到的数据装进结构体
	result := common.GetResultRow(row) //转换成结构体
	if len(result) == 0 {              //如果一条数据都没有查到
		return &datamodels.Product{}, nil
	}
	product := &datamodels.Product{}
	common.DataToStructByTagSql(result, product) //将获取到的数据放入到product结构体中
	return product, nil
}

// 查询多行
func (p *ProductManager) SearchAll() (productArry []*datamodels.Product, err error) { //查找所有数据，返回多个products对应结构体
	if err = p.Conn(); err != nil {
		return
	}
	// 1. 语句数据库查询
	sql := "select * from " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return
	}
	// 3. 将查询到的数据装进结构体
	results := common.GetResultRows(rows) //转换成结构体
	if len(results) == 0 {                //如果一条数据都没有查到
		return nil, nil
	}
	// 4. 目前结果都是map[数据库列名]{值} ： 需要将其装入到对应的结构体
	for _, res := range results {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(res, product) //指针传入
		productArry = append(productArry, product)
	}
	// for _, v := range results {
	// 	product := &datamodels.Product{}
	// 	common.DataToStructByTagSql(v, product)
	// 	fmt.Println(product)
	// 	productArry = append(productArry, product)
	// }
	// 直接装进结构体
	// products := []*datamodels.Product{}

	// common.DataToStructByTagSql(result, res) //将获取到的数据放入到product结构体中
	return
}

func (p *ProductManager) Update(product *datamodels.Product) (err error) { //更新某条product
	if err = p.Conn(); err != nil {
		return
	}
	sql := "update " + p.table + " SET productName=?, productNum=?, productImage=?, productUrl=? where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ID)
	if err != nil {
		return
	}
	return nil
}
func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "Update " + p.table + " SET productNum = productNum-1 where ID=" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
