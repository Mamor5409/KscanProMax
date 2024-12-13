package fofa

import (
	"fmt"
	"testing"
)

func TestDoamin(t *testing.T) {

	domain := "jnxkjk.com"
	if CheckWildcard(string(domain)) {
		println("存在泛解析")
	} else {
		println("不存在泛解析")
	}
}

func TestGetDomainDict(t *testing.T) {
	urls := []string{
		"http://www.baidu.com",
	}
	// 声明一个字典用于存放主域名
	maindomains := make(map[string]struct{})

	// 提取主域名并进行去重
	for _, url := range urls {
		// 解析URL
		parsedURL, err := parseURL(url)
		if err != nil {
			fmt.Printf("解析URL出错：%v\n", err)
			continue
		}
		println("[+] 解析的url", parsedURL.Hostname())
		// 获取主域名
		maindomain := extractMainDomain(parsedURL.Hostname())

		// 检查主域名是否已经存在于字典中
		if _, exists := maindomains[maindomain]; exists {
			fmt.Printf("已存在主域名：%s\n", maindomain)
		} else {
			// 如果不存在，则添加到字典中
			maindomains[maindomain] = struct{}{}
			fmt.Printf("已添加主域名：%s\n", maindomain)
		}
	}

	// 输出结果
	fmt.Println("唯一主域名字典:")
	for maindomain := range maindomains {
		fmt.Println(maindomain)
	}
}
