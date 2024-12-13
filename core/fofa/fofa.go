package fofa

import (
	"KscanPro/app"
	"KscanPro/core/slog"
	"KscanPro/lib/fofa"
	"KscanPro/lib/misc"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var this *fofa.Client
var keywordSlice []string
var results []fofa.Result

func Init(email, key string) {
	//设置日志输出器
	fofa.SetLogger(slog.Debug())
	//初始化fofa模块
	this = fofa.New(email, key)
	this.SetSize(app.Setting.FofaSize)
	//获取所有关键字
	keywordSlice = makeKeywordSlice()

}

func Run() {
	//对每个关键字进行查询
	for _, keyword := range keywordSlice {
		slog.Printf(slog.WARN, "本次搜索关键字为：%v", keyword)
		_, r := this.Search(keyword)
		//size, r := this.Search(keyword)
		displayResponse(r)
		//slog.Printf(slog.INFO, "本次搜索，返回结果总条数为：%d，此次返回条数为：%d  [已过滤掉泛解析Domain的结果]", size, len(r)-tmpnum)
		results = append(results, r...)
	}
}

func makeKeywordSlice() []string {
	var keywordSlice []string
	if app.Setting.FofaFixKeyword == "" {
		keywordSlice = app.Setting.Fofa
	} else {
		for _, keyword := range app.Setting.Fofa {
			keyword = strings.ReplaceAll(app.Setting.FofaFixKeyword, "{}", keyword)
			keywordSlice = append(keywordSlice, keyword)
		}
	}
	return keywordSlice
}

func GetUrlTarget() []string {
	var strSlice []string
	for _, r := range results {
		Fix(&r)
		strSlice = append(strSlice, r.Host)
	}
	strSlice = misc.RemoveDuplicateElement(strSlice)
	return strSlice
}

func GetHostTarget() []string {
	var strSlice []string
	for _, r := range results {
		strSlice = append(strSlice, r.Ip)
	}
	strSlice = misc.RemoveDuplicateElement(strSlice)
	return strSlice
}

// 记录泛解析domain相关url数
var tmpnum int = 0

// 声明一个全局字典用于存放主域名
var maindomains = make(map[string]bool)
var wildcardmaindomains = make(map[string]bool)
var isFoFaHeaderPrinted = false

func displayResponse(r []fofa.Result) {
	for _, row := range r {
		Fix(&row)
		m := row.Map()

		domain := m["Domain"]

		// 检查主域名是否已经存在于泛解析主域名的字典中  如果存在 则跳过当前url
		if _, exists := maindomains[domain]; exists {
			if _, exists := wildcardmaindomains[domain]; exists {
				//tmpnum++
				continue
			}
		} else {
			maindomains[domain] = true
			if CheckWildcard(domain) {
				wildcardmaindomains[domain] = true
				//tmpnum++
				continue
			}
		}

		m = misc.FixMap(m)

		// 检查 Ip 是否为 0.0.0.0，若是则跳过该条记录
		if strings.Contains(m["Ip"], "0.0.0.0") || m["Protocol"] == "" {
			continue // 跳过当前条目
		}

		if m["Banner"] != "" {
			//	m["Banner"] = misc.FixLine(m["Banner"])
			//	m["Banner"] = misc.StrRandomCut(m["Banner"], 20)
			m["Banner"] = strings.ReplaceAll(m["Banner"], "\n", "<br>")
		}
		if m["Header"] != "" {
			m["Header"] = strings.ReplaceAll(m["Header"], "\n", "<br>")
		}

		if !isFoFaHeaderPrinted {
			slog.Println(slog.DATA, fmt.Sprintf("+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+"))
			slog.Println(slog.DATA, fmt.Sprintf("\r| %-45v | %-30v | %-15v | %-20v | %-10v | %-25v | %v ", "Host", "Title", "IP", "Domain", "Protocol", "Addr", "Server"))
			slog.Println(slog.DATA, fmt.Sprintf("+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+"))
			isFoFaHeaderPrinted = true
		}
		line := fmt.Sprintf("\r| %-45v | %-"+strconv.Itoa(misc.AutoWidth(row.Title, 30))+"v | %-15v | %-"+strconv.Itoa(misc.AutoWidth(row.Domain, 20))+"v | %-"+strconv.Itoa(misc.AutoWidth(row.Protocol, 10))+"v | %-"+strconv.Itoa(misc.AutoWidth(fmt.Sprintf(m["Country"]+" "+row.Province+" "+row.City), 25))+"v | %v ", row.Host, row.Title, row.Ip, row.Domain, m["Protocol"], m["Addr"], m["Server"])
		slog.Println(slog.DATA, line)
		msi := make(map[string]interface{})
		// 遍历 strMap 中的每个元素
		for key, value := range m {
			// 将字符串值赋给对应的 interface{} 类型
			msi[key] = value
		}
		app.FofaResults = append(app.FofaResults, msi)
	}
}

func Fix(r *fofa.Result) {
	//修复title
	if r.Title == "" && r.Protocol != "" {
		r.Title = strings.ToUpper(r.Protocol)
	}
	r.Title = misc.FixLine(r.Title)
	//修改host
	if r.Host == "" {
		r.Host = r.Ip
	}

	if regexp.MustCompile("\\w+://.*").MatchString(r.Host) == false {
		if r.Host == "" {
			r.Protocol = "http"
		}
		r.Host = r.Protocol + "://" + r.Host
	}
}

// 泛解析判断:代码摘自zpscan
func CheckWildcard(domain string) (ok bool) {
	for i := 0; i < 2; i++ {
		_, err := net.LookupHost(uuid.NewV4().String() + "." + domain)
		if err == nil {
			return true
		}
	}
	return false
}

// 解析URL
func parseURL(urlString string) (*url.URL, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// 提取主域名
func extractMainDomain(hostname string) string {
	parts := strings.Split(hostname, ".")
	var mainDomain string
	if len(parts) >= 3 {
		mainDomain = parts[len(parts)-2] + "." + parts[len(parts)-1]
	} else {
		mainDomain = hostname
	}
	return mainDomain
}
