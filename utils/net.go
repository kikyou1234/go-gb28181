package utils

import (
	"fmt"
	"net"
)

func GetLocalIP() (string, error) {
	addrs, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range addrs {
		// 排除未启用和回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 排除 IPv6 和回环地址
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("找不到本地IP地址")
}
