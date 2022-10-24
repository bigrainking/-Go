package repository

import (
	"Spike-Product-Demo/common"
	"Spike-Product-Demo/datamodels"
	"database/sql"
	"errors"
	"fmt"
)

type IUser interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error) //通过用户名查找用户，用于核实登录信息
	Insert(user *datamodels.User) (userId int64, err error)
}

type UserManagerRepo struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserManagerRepo(table string, mysqlConn *sql.DB) IUser {
	return &UserManagerRepo{table, mysqlConn}
}

func (u *UserManagerRepo) Conn() error {
	if u.mysqlConn == nil {
		sql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = sql
	}
	if u.table == "" {
		u.table = "spikeSystem.User"
	}
	return nil
}
func (u *UserManagerRepo) Select(userName string) (*datamodels.User, error) {
	if err := u.Conn(); err != nil {
		return nil, err //数据库链接出错
	}

	sql := "Select * from " + u.table + " where userName=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := common.GetResultRow(rows)
	if len(res) == 0 {
		return nil, errors.New("用户不存在！")
	}

	user := &datamodels.User{}
	common.DataToStructByTagSql(res, user)
	return user, nil
}
func (u *UserManagerRepo) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return //数据库链接出错
	}

	sql := "Insert " + u.table + " Set nickName=?,userName=?,hashPassword=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	res, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		fmt.Println(err)
		return
	}

	return res.LastInsertId()
}
