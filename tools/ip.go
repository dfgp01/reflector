package tools

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// 返回本地IP地址
func GetLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ip := ipnet.IP.To4()
			if IsPrivateIPv4(ip) {
				return ip, nil
			}
		}
	}
	return nil, errors.New("no private ip address")
}

// 是否是私有IP（内网地址）
func IsPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || (ip[0] == 172 && ip[1] >= 16 && ip[1] < 32) || (ip[0] == 192 && ip[1] == 168))
}

// 获取外网IP（未测试）
func GetOutBoundIP() {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip := strings.Split(localAddr.String(), ":")[0]
	fmt.Println(ip, err)
}
