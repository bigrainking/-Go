package repository

import (
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/datamodels"
	"database/sql"
	"strconv"
)

// 操作订单的增删改查

type Order interface {
	// 数据库链接
	Conn() error
	// 增删改查+查询订单相关信息
	Insert(*datamodels.Order) (id int64, err error)
	Delete(int64) error
	Update(*datamodels.Order) error
	SelectByID(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() ([]map[string]string, error) //订单号、商品名称、物流状态
}

// 创建结构体
type OrderManager struct {
	mysqlConn *sql.DB
	table     string
}

func NewOrderManagerRepo(table string, mysqlConn *sql.DB) Order {
	return &OrderManager{mysqlConn, table}
}

// 数据库链接
func (o *OrderManager) Conn() error {
	if o.mysqlConn == nil {
		conn, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = conn
	}
	if o.table == "" {
		o.table = "spikeSystem.order"
	}
	return nil
}

// 增删改查+查询订单相关信息
func (o *OrderManager) Insert(order *datamodels.Order) (id int64, err error) {
	// 先链接数据库
	if err = o.Conn(); err != nil {
		return
	}
	// 查询
	sql := "INSERT " + o.table + " SET userID=?,productID=?,orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	res, err := stmt.Exec(order.UserID, order.ProductID, order.Orderstatus)
	if err != nil {
		return
	}
	id, _ = res.LastInsertId()
	return
}
func (o *OrderManager) Delete(id int64) error {
	if err := o.Conn(); err != nil {
		return err
	}
	sql := ("Delete from " + o.table + " where ID=?")
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(strconv.FormatInt(id, 10))
	if err != nil {
		return err
	}
	return nil
}
func (o *OrderManager) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}
	sql := "Update " + o.table + " Set userID=?, productID=?, orderStatus=? where ID=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserID, order.ProductID, order.Orderstatus, order.ID)
	if err != nil {
		return err
	}
	return nil
}
func (o *OrderManager) SelectByID(id int64) (*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	// 1. 语句
	sql := "Select * from " + o.table + " where ID=" + strconv.FormatInt(id, 10)
	row, err := o.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return nil, err
	}
	// 3. 将查询到的数据装进结构体
	result := common.GetResultRow(row) //转换成结构体
	if len(result) == 0 {              //如果一条数据都没有查到
		return &datamodels.Order{}, nil
	}
	order := &datamodels.Order{}
	common.DataToStructByTagSql(result, order) //将获取到的数据放入到product结构体中
	return order, nil
}
func (o *OrderManager) SelectAll() (oderArry []*datamodels.Order, err error) {
	sql := "select * from " + o.table
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	// 3. 将查询到的数据装进结构体
	results := common.GetResultRows(rows) //转换成结构体
	if len(results) == 0 {                //如果一条数据都没有查到
		return nil, nil
	}
	// 4. 目前结果都是map[数据库列名]{值} ： 需要将其装入到对应的结构体
	for _, res := range results {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(res, order) //指针传入
		oderArry = append(oderArry, order)
	}
	return
}
func (o *OrderManager) SelectAllWithInfo() ([]map[string]string, error) { //订单号、商品名称、物流状态
	if err := o.Conn(); err != nil {
		return nil, err
	}
	// 1. 语句
	sql := "Select o.ID,p.productName,o.orderStatus From spikeSystem.order as o left join spikeSystem.product as p on o.productID=p.ID"
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	// 3. 将查询到的数据装进结构体
	results := common.GetResultRows(rows) //转换成结构体
	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}
