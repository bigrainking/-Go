package datamodels

type User struct {
	ID           int64  `json:"id" form:"ID" sql:"ID"`
	NickName     string `json:"NickName" form:"NickName" sql:"nickName"`
	UserName     string `json:"UserName" form:"UserName" sql:"userName"`
	HashPassword string `json:"-" form:"HashPassword" sql:"hashPassword"`
}
