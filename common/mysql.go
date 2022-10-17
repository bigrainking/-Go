package common

// 创建数据库链接

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//创建mysql 连接
func NewMysqlConn() (db *sql.DB, err error) {
	// 用户名：密码@端口号/数据库？编码
	db, err = sql.Open("mysql", "BigRainKing:1@tcp(127.0.0.1:3306)/spikeSystem?charset=utf8")
	errlink := db.Ping()
	if errlink != nil {
		log.Printf("数据库链接失败")
		panic(errlink)
	} else {
		fmt.Println("数据库链接成功")
	}
	return
}

//获取返回值，获取一条
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	// for j := range values {
	// 	scanArgs[j] = &values[j]
	// }
	for i := 0; i < len(columns); i++ {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]string)

	for rows.Next() {
		//将行数据保存到record字典
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}

//获取所有查询结果的每一行数据的切片集合
func GetResultRows(rows *sql.Rows) (dataMaps []map[string]string) {
	// 1. 查询到的数据列名、返回值
	columns, _ := rows.Columns() //列名
	count := len(columns)
	values, valuesPoints := make([][]byte, count), make([]interface{}, count)

	// 2. 遍历Rows读取每一行
	for rows.Next() {
		// for i, v := range values { // 读取value地址到valuePoints
		// 	valuesPoints[i] = &v
		// }
		for i := 0; i < count; i++ {
			valuesPoints[i] = &values[i]
		}

		// 2.1 数据库中读取出每一行数据
		rows.Scan(valuesPoints...) //将所有内容读取进values

		// 2.2 准备接收数据的结构体Product
		row := make(map[string]string)

		// 2.3 将读取到的数据填充到product
		for i, val := range values { // val是每个列对应的值
			key := columns[i] //列名

			// 列名与值对应
			row[key] = string(val)
		}

		// 将product归到集合中
		dataMaps = append(dataMaps, row)
	}
	return
}
