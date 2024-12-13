package leakscan

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type EmailInfo struct {
	SubDomain string   `json:"subDomain"`
	Emails    []string `json:"Emails"`
}

type AllEmailInfo struct {
	CompanyName string      `json:"companyName"`
	AllEmails   []EmailInfo `json:"AllEmails"`
}

func SearchSkymem(domain string) (emailResult EmailInfo) {
	emailResult.SubDomain = domain
	//var mailResult []string
	baseurl := "https://www.skymem.info"
	url := baseurl + "/srch?q=" + domain + "&ss=srch"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("[-] 创建http请求错误，原因:", err)
		return
	}
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[-] 发送请求错误，原因:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		// 读取响应体内容
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err2)
		}

		htmlData := string(bodyBytes)
		if strings.Contains(htmlData, "did not match any documents") {
			fmt.Println("[-] 未找到该域名相关邮箱:(")
			return
		} else {
			// 判断页数 是否存在第二页
			if strings.Contains(htmlData, "More emails for") {
				//fmt.Println("[+] 存在多页邮箱地址，由于网站限制最多查看25条该域名相关邮箱！")
				// //div/div/a/@href
				root, _ := htmlquery.Parse(strings.NewReader(htmlData))
				hrefs := htmlquery.Find(root, "//div/div/a/@href")
				nextpageurl := htmlquery.InnerText(hrefs[0])
				//   /domain/5662c76012ad5b179450b607?p=2

				url2 := baseurl + nextpageurl
				//fmt.Println(url2)
				//url2inhtml := "skymem.info" + nextpageurl
				re := regexp.MustCompile(`(\?p=)\d+`)
				//url3inhtml := re.ReplaceAllString(url2inhtml, "${1}3")
				nextpageurl2 := re.ReplaceAllString(nextpageurl, "${1}3")
				url3 := baseurl + nextpageurl2

				//fmt.Println(url3inhtml)

				req, err2 := http.NewRequest("GET", url2, nil)
				if err2 != nil {
					fmt.Println("[-] 创建http请求错误，原因:", err2)
					return
				}
				// 发送请求
				client := &http.Client{}
				resp2, err2 := client.Do(req)
				if err2 != nil {
					fmt.Println("[-] 发送请求错误，原因:", err2)
					return
				}
				defer resp2.Body.Close()

				if resp2.StatusCode == 200 {
					// 读取响应体内容
					bodyBytes2, err2 := ioutil.ReadAll(resp2.Body)
					if err != nil {
						panic(err2)
					}
					//fmt.Println("[+] 正在爬取第二页邮箱")
					htmlData2 := string(bodyBytes2)
					//fmt.Println("[dbg]", htmlData2)
					root2, _ := htmlquery.Parse(strings.NewReader(htmlData2))
					emails := htmlquery.Find(root2, "//tr/td/a/text()")

					for _, email := range emails {
						emailResult.Emails = append(emailResult.Emails, htmlquery.InnerText(email))
						//fmt.Println(htmlquery.InnerText(email))
					}
					//fmt.Println("[+] 第二页收集完毕")
					//fmt.Println("[+] 判断是否存在第三页")
					//println(url3inhtml)
					if strings.Contains(htmlData2, nextpageurl2) {
						//fmt.Println("[+] 存在第三页内容 准备爬取")

						req, err := http.NewRequest("GET", url3, nil)
						if err != nil {
							fmt.Println("[-] 创建http请求错误，原因:", err)
							return
						}
						// 发送请求
						client := &http.Client{}
						resp3, err := client.Do(req)
						if err != nil {
							fmt.Println("[-] 发送请求错误，原因:", err)
							return
						}
						defer resp3.Body.Close()
						if resp.StatusCode == 200 {
							// 读取响应体内容
							bodyBytes3, err2 := ioutil.ReadAll(resp3.Body)
							if err != nil {
								panic(err2)
							}

							//fmt.Println("[+] 正在爬取第三页邮箱")
							htmlData3 := string(bodyBytes3)
							root3, _ := htmlquery.Parse(strings.NewReader(htmlData3))
							emails := htmlquery.Find(root3, "//tr/td/a/text()")
							//fmt.Println(emails)
							for _, email := range emails {
								emailResult.Emails = append(emailResult.Emails, htmlquery.InnerText(email))
								//fmt.Println(htmlquery.InnerText(email))
							}

						} else {
							fmt.Println("[-] 第三页请求失败 StatusCode: ", resp3.StatusCode)
							return
						}

					} else {
						//fmt.Println("[+] 该域名相关邮箱收集完毕")
						return
					}

				} else {
					fmt.Println("[-] 请求第二页出现错误!")
					return
				}

			}

			htmlData := string(bodyBytes)
			root, _ := htmlquery.Parse(strings.NewReader(htmlData))
			emails := htmlquery.Find(root, "//tr/td/a/text()")
			//fmt.Println("[+] 第一页")
			for _, email := range emails {
				emailResult.Emails = append(emailResult.Emails, htmlquery.InnerText(email))
				//mailResult = append(mailResult, htmlquery.InnerText(email))
				//fmt.Println(htmlquery.InnerText(email))
			}
			//fmt.Println("[+] 第一页 完毕")
		}

	} else {
		fmt.Println("[-] Request skymem Error! Status Code：", resp.StatusCode)
		return
	}

	//fmt.Println(mailResult)
	//fmt.Println(len(mailResult))
	//jsonData, err := json.Marshal(mailResult)
	//fmt.Println(string(jsonData))

	return
}
