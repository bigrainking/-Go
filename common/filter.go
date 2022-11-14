package common

import (
	"errors"
	"net/http"
	"strings"
)

// 拦截器对应结构体
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// 拦截器
type Filter struct {
	// 需要存储各个URL对应的拦截器
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{make(map[string]FilterHandle)}
}

// 注册拦截器
func (f *Filter) RegisteFilter(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// 获取拦截器
func (f *Filter) GetFilter(uri string) (error, FilterHandle) {
	if handler, ok := f.filterMap[uri]; !ok {
		return errors.New(uri + " 没有找到对应的拦截器！"), nil
	} else {
		return nil, handler
	}
}

// 正常注册的URL的handler
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// 执行拦截器
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, req *http.Request) {
	// 依次执行拦截器
	return func(rw http.ResponseWriter, req *http.Request) {
		for uri, handler := range f.filterMap {
			if strings.Contains(req.RequestURI, uri) {
				// if uri == req.RequestURI { //什么叫因为我们加了个一个product-id所以要扩大范围？
				// 可能是/product/order?productID=1我们需要过滤这个请求，所以拦截器需要扩大范围
				// 将请求扩大到这里
				err := handler(rw, req)
				if err != nil {
					rw.Write([]byte(err.Error())) //返回错误到writer
					return
				}
			}
			// 如果不是请求URI，则跳过 ： 为什么找到一个req.RequestURI就要跳出循环？？？
			break
		}
		// 执行正常注册的函数
		webHandle(rw, req)
	}

}
