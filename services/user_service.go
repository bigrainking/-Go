package services

import (
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/repository"

	"golang.org/x/crypto/bcrypt"
)

// 验证用户密码是否正确
// 用户添加

type IUserService interface {
	// 查看用户输入的名称、密码是否正确,如果正确则返回用户相关信息
	IsPwdSucceed(userName string, pwd string) (user *datamodels.User, err error)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepo repository.IUser
}

func NewUserServiceManager(repo repository.IUser) IUserService {
	return &UserService{repo}
}

func (u *UserService) IsPwdSucceed(userName string, pwd string) (user *datamodels.User, err error) {

	user, err = u.UserRepo.Select(userName)

	if err != nil {
		return
	}
	_, err = ValidatePasswd(user.HashPassword, pwd)

	if err != nil {
		return
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	// 处理密码
	hashpwd, err := GeneratePasswd(user.HashPassword)
	if err != nil {
		return
	}
	// 插入用户
	user.HashPassword = string(hashpwd)
	return u.UserRepo.Insert(user)
}

// 对比密码是否相同，一个是hash过的，一个没有
func ValidatePasswd(pwd, hashPwd string) (succeed bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(hashPwd)); err != nil {
		return
	}
	return
}

// 将明文密码hash化
func GeneratePasswd(pwd string) (hashPwd []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}
