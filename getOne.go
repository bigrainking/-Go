package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// 防止超卖

// - 预存商品数量
var productNum int64 = 10000

// - 互斥锁
var mutex sync.Mutex

// - 已经秒杀掉的商品数量sum
var sum int64 = 0

// - 获取秒杀商品功能函数：
func GetOneProduct() bool {
	//   - 加锁
	mutex.Lock()
	defer mutex.Unlock()

	//   - 判断sum是否超过限制
	if sum < productNum {
		sum += 1
		fmt.Println(sum)
		return true
	}
	return false
}

// - 秒杀商品接口
func GetProduct(rw http.ResponseWriter, req *http.Request) {
	// 获取商品数量

	if GetOneProduct() {
		rw.Write([]byte("true"))
		return
	}
	rw.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct) //绑定handler
	fmt.Println("开始运行getOne")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatal("Err:", err)
	}

}
