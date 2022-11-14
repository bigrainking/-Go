package common

import (
	"errors"
	"net"
)

// 获取本机IP
// 返回IP地址
func GetIntranceIP() (string, error) {
	// 固定写法不需要深入了解
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		//检查Ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("获取地址异常")
}
