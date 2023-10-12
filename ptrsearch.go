package main

import (
        "fmt"
        "net"
        "regexp"
)

func main() {
        // 定义要查询的IP地址范围
        //cidrs := []string{"4.0.0.0/9", "8.0.0.0/9"}
        cidrs := []string{"8.0.0.0/9"}
        ch := make(chan struct{}, 1000)
        // 创建PTR查询
        ptrQuery := "CHINA" // 不用在查询中加上"."

        // 遍历每个CIDR范围
        for _, cidr := range cidrs {
                _, ipNet, err := net.ParseCIDR(cidr)
                if err != nil {
                        fmt.Println("无效的CIDR范围:", cidr)
                        continue
                }

                // 遍历IP地址范围
                for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
                        ch <- struct{}{}
                        go func() {

                                hosts, err := net.LookupAddr(ip.String())
                                if err != nil {
                                        //continue // 跳过查询失败的IP地址
                                }

                                // 遍历查询结果，查找包含"china"的主机名
                                for _, host := range hosts {
                                        //fmt.Printf("IP地址: %s, 主机名: %s\n", ip.String(), host)
                                        if contains(host, ptrQuery) {
                                                fmt.Printf("IP地址: %s, 主机名: %s\n", ip.String(), host)
                                        }

                                }
                                <-ch

                        }()

                        // 执行PTR查询

                }
        }
}

// inc 增加IP地址
func inc(ip net.IP) {
        for j := len(ip) - 1; j >= 0; j-- {
                ip[j]++
                if ip[j] > 0 {
                        break
                }
        }
}

// containsChina 检查主机名是否包含"china"
func contains(host, query string) bool {
        regex := regexp.MustCompile("(?i)" + regexp.QuoteMeta(query))
        return regex.MatchString(host)
}
