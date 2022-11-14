package middleware

import "github.com/kataras/iris/v12"

// 实现进入商品详情页之前的用户是否登录验证
func AuthConProduct(c iris.Context) {
	// 没有登录则跳转到登录界面
	uid := c.GetCookie("uid")
	if uid == "" {
		c.Application().Logger().Debug("用户未登录, 必须先登录！")
		c.Redirect("/user/login")
		return
	}
	c.Application().Logger().Debug("用户" + uid + "已经登录")
	c.Next()
}
