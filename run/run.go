package run

import (
	"KscanPro/app"
	"KscanPro/core/cdn"
	"KscanPro/core/hydra"
	"KscanPro/core/pocScan"
	"KscanPro/core/scanner"
	"KscanPro/core/slog"
	"KscanPro/lib/color"
	"KscanPro/lib/misc"
	"KscanPro/lib/uri"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"github.com/lcvvvv/simplehttp"
	"github.com/lcvvvv/stdio/chinese"
	"github.com/mattn/go-runewidth"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Start() {
	//启用看门狗函数定时输出负载情况
	go watchDog()
	//下发扫描任务
	var wg = &sync.WaitGroup{}
	wg.Add(5)
	DomainScanner = generateDomainScanner(wg)
	IPScanner = generateIPScanner(wg)
	PortScanner = generatePortScanner(wg)
	URLScanner = generateURLScanner(wg)
	HydraScanner = generateHydraScanner(wg)

	// AQC
	// 敏感信息
	// http密码爆破  tomcat/jenkins/phpmyadmin

	//扫描器进入监听状态
	start()

	slog.Println(slog.INFO, "开始下发扫描任务")
	//开始分发扫描任务
	for _, expr := range app.Setting.Target {
		pushTarget(expr)
	}
	//slog.Println(slog.INFO, "所有扫描任务已下发完毕")
	//根据扫描情况，关闭scanner
	go stop()
	wg.Wait()
}

func pushTarget(expr string) {
	if expr == "" {
		return
	}
	if expr == "paste" || expr == "clipboard" {
		if clipboard.Unsupported == true {
			slog.Println(slog.ERROR, runtime.GOOS, "clipboard unsupported")
		}
		clipboardStr, _ := clipboard.ReadAll()
		for _, line := range strings.Split(clipboardStr, "\n") {
			line = strings.ReplaceAll(line, "\r", "")
			pushTarget(line)
		}
		return
	}
	if uri.IsIPv4(expr) {
		IPScanner.Push(net.ParseIP(expr))
		if app.Setting.Check == true {
			pushURLTarget(uri.URLParse("http://"+expr), nil)
			pushURLTarget(uri.URLParse("https://"+expr), nil)
		}
		return
	}
	if uri.IsIPv6(expr) {
		slog.Println(slog.WARN, "暂时不支持IPv6的扫描对象：", expr)
		return
	}
	if uri.IsCIDR(expr) {
		for _, ip := range uri.CIDRToIP(expr) {
			pushTarget(ip.String())
		}
		return
	}
	if uri.IsIPRanger(expr) {
		for _, ip := range uri.RangerToIP(expr) {
			pushTarget(ip.String())
		}
		return
	}
	if uri.IsDomain(expr) {
		DomainScanner.Push(expr)
		pushURLTarget(uri.URLParse("http://"+expr), nil)
		pushURLTarget(uri.URLParse("https://"+expr), nil)
		return
	}
	if uri.IsHostPath(expr) {
		pushURLTarget(uri.URLParse("http://"+expr), nil)
		pushURLTarget(uri.URLParse("https://"+expr), nil)
		if app.Setting.Check == false {
			pushTarget(uri.GetNetlocWithHostPath(expr))
		}
		return
	}
	if uri.IsNetlocPort(expr) {
		netloc, port := uri.SplitWithNetlocPort(expr)
		if uri.IsIPv4(netloc) {
			PortScanner.Push(net.ParseIP(netloc), port)
		}
		if uri.IsDomain(netloc) {
			pushURLTarget(uri.URLParse("http://"+expr), nil)
			pushURLTarget(uri.URLParse("https://"+expr), nil)
		}
		if app.Setting.Check == false {
			pushTarget(netloc)
		}
		return
	}
	if uri.IsURL(expr) {
		pushURLTarget(uri.URLParse(expr), nil)
		if app.Setting.Check == false {
			pushTarget(uri.GetNetlocWithURL(expr))
		}
		return
	}
	slog.Println(slog.WARN, "无法识别的Target字符串:", expr)
}

func pushURLTarget(URL *url.URL, response *gonmap.Response) {
	var cli *http.Client
	//判断是否初始化client
	if app.Setting.Proxy != "" || app.Setting.Timeout != 3*time.Second {
		cli = simplehttp.NewClient()
	}
	//判断是否需要设置代理
	if app.Setting.Proxy != "" {
		simplehttp.SetProxy(cli, app.Setting.Proxy)
	}
	//判断是否需要设置超时参数
	if app.Setting.Timeout != 3*time.Second {
		simplehttp.SetTimeout(cli, app.Setting.Timeout)
	}

	//判断是否存在请求修饰性参数
	if len(app.Setting.Host) == 0 && len(app.Setting.Path) == 0 {
		URLScanner.Push(URL, response, nil, cli)
		return
	}

	//如果存在，则逐一建立请求下发队列
	var reqs []*http.Request
	for _, host := range app.Setting.Host {
		req, _ := simplehttp.NewRequest(http.MethodGet, URL.String(), nil)
		req.Host = host
		reqs = append(reqs, req)
	}
	for _, path := range app.Setting.Path {
		req, _ := simplehttp.NewRequest(http.MethodGet, URL.String()+path, nil)
		reqs = append(reqs, req)
	}
	for _, req := range reqs {
		URLScanner.Push(req.URL, response, req, cli)
	}
}

var (
	DomainScanner *scanner.DomainClient
	IPScanner     *scanner.IPClient
	PortScanner   *scanner.PortClient
	URLScanner    *scanner.URLClient
	HydraScanner  *scanner.HydraClient
)

func start() {
	go DomainScanner.Start()
	go IPScanner.Start()
	go PortScanner.Start()
	go URLScanner.Start()
	go HydraScanner.Start()
	time.Sleep(time.Second * 1)
	slog.Println(slog.INFO, "Domain、IP、Port、URL、Hydra引擎已准备就绪")
}

func stop() {
	for {
		time.Sleep(time.Second)
		if DomainScanner.RunningThreads() == 0 && DomainScanner.IsDone() == false {
			DomainScanner.Stop()
			slog.Println(slog.DEBUG, "检测到所有Domian检测任务已完成，Domain扫描引擎已停止")
		}
		if IPScanner.RunningThreads() == 0 && IPScanner.IsDone() == false {
			IPScanner.Stop()
			slog.Println(slog.DEBUG, "检测到所有IP检测任务已完成，IP扫描引擎已停止")
		}
		if IPScanner.IsDone() == false {
			continue
		}
		if PortScanner.RunningThreads() == 0 && PortScanner.IsDone() == false {
			PortScanner.Stop()

			slog.Println(slog.DEBUG, "检测到所有Port检测任务已完成，Port扫描引擎已停止")
		}
		if PortScanner.IsDone() == false {
			continue
		}
		if URLScanner.RunningThreads() == 0 && URLScanner.IsDone() == false {
			URLScanner.Stop()
			slog.Println(slog.DEBUG, "检测到所有URL检测任务已完成，URL扫描引擎已停止")
		}
		if HydraScanner.RunningThreads() == 0 && HydraScanner.IsDone() == false {
			if len(app.HydraResults) != 0 {
				PrintHydraResults()
			}
			HydraScanner.Stop()
			slog.Println(slog.DEBUG, "检测到所有暴力破解任务已完成，暴力破解引擎已停止")
		}
	}
}

func generateDomainScanner(wg *sync.WaitGroup) *scanner.DomainClient {
	DomainConfig := scanner.DefaultConfig()
	DomainConfig.Threads = 10
	client := scanner.NewDomainScanner(DomainConfig)
	client.HandlerRealIP = func(domain string, ip net.IP) {
		IPScanner.Push(ip)
	}
	client.HandlerIsCDN = func(domain, CDNInfo string) {
		outputCDNRecord(domain, CDNInfo)
	}
	client.HandlerError = func(domain string, err error) {
		slog.Println(slog.DEBUG, "DomainScanner Error: ", domain, err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}

func generateIPScanner(wg *sync.WaitGroup) *scanner.IPClient {
	IPConfig := scanner.DefaultConfig()
	IPConfig.Threads = 200
	IPConfig.Timeout = 200 * time.Millisecond
	IPConfig.HostDiscoverClosed = app.Setting.ClosePing
	client := scanner.NewIPScanner(IPConfig)
	client.HandlerDie = func(addr net.IP) {
		slog.Println(slog.DEBUG, addr.String(), " is die")
	}
	client.HandlerAlive = func(addr net.IP) {
		//启用端口存活性探测任务下发器
		slog.Println(slog.DEBUG, addr.String(), " is alive")
		for _, port := range app.Setting.Port {
			PortScanner.Push(addr, port)
		}
	}
	client.HandlerError = func(addr net.IP, err error) {
		slog.Println(slog.DEBUG, "IPScanner Error: ", addr.String(), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}

func getTimeout(i int) time.Duration {
	switch {
	case i > 10000:
		return time.Millisecond * 200
	case i > 5000:
		return time.Millisecond * 300
	case i > 1000:
		return time.Millisecond * 400
	default:
		return time.Millisecond * 500
	}
}

func generatePortScanner(wg *sync.WaitGroup) *scanner.PortClient {

	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = app.Setting.Threads
	PortConfig.Timeout = getTimeout(len(app.Setting.Port))
	if app.Setting.ScanVersion == true {
		PortConfig.DeepInspection = true
	}

	client := scanner.NewPortScanner(PortConfig)
	client.HandlerClosed = func(addr net.IP, port int) {
		//nothing
	}
	client.HandlerOpen = func(addr net.IP, port int) {
		outputOpenResponse(addr, port)
	}
	client.HandlerNotMatched = func(addr net.IP, port int, response string) {
		outputUnknownResponse(addr, port, response)
	}
	client.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		URLRaw := fmt.Sprintf("%s://%s:%d", response.FingerPrint.Service, addr.String(), port)
		URL, _ := url.Parse(URLRaw)
		if appfinger.SupportCheck(URL.Scheme) == true {
			pushURLTarget(URL, response)
			return
		}
		outputNmapFinger(URL, response)
		if app.Setting.Hydra == true {
			if protocol := response.FingerPrint.Service; hydra.Ok(protocol) {
				HydraScanner.Push(addr, port, protocol)
			}
		}
	}
	client.HandlerError = func(addr net.IP, port int, err error) {
		slog.Println(slog.DEBUG, "PortScanner Error: ", fmt.Sprintf("%s:%d", addr.String(), port), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}

func generateURLScanner(wg *sync.WaitGroup) *scanner.URLClient {
	URLConfig := scanner.DefaultConfig()
	URLConfig.Threads = app.Setting.Threads/2 + 1

	client := scanner.NewURLScanner(URLConfig)
	client.HandlerMatched = func(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
		outputAppFinger(URL, banner, finger)
		// 漏洞探测逻辑入口 默认开启   --nopoc则关闭
		if app.Setting.Nopoc == false {
			url := URL.Scheme + "://" + URL.Host
			info := app.HostInfo{
				Host:  URL.Scheme,
				Ports: URL.Hostname(),
				Url:   url,
			}
			//pocScan.WebScan(url)
			pocScan.WebTitle(&info)

		}

	}
	client.HandlerError = func(url *url.URL, err error) {
		slog.Println(slog.DEBUG, "URLScanner Error: ", url.String(), " ", err)
	}
	client.Defer(func() {
		//if len(app.UrlCheckResults) != 0 || len(app.PortScanResults) != 0 {
		//	slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
		//}
		wg.Done()
	})
	return client
}

func generateHydraScanner(wg *sync.WaitGroup) *scanner.HydraClient {
	HydraConfig := scanner.DefaultConfig()
	HydraConfig.Threads = 10

	client := scanner.NewHydraScanner(HydraConfig)
	client.HandlerSuccess = func(addr net.IP, port int, protocol string, auth *hydra.Auth) {
		outputHydraSuccess(addr, port, protocol, auth)
	}
	client.HandlerError = func(addr net.IP, port int, protocol string, err error) {
		slog.Println(slog.DEBUG, fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port), err)
	}
	client.Defer(func() {
		//slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
		wg.Done()
	})
	return client
}

func outputHydraSuccess(addr net.IP, port int, protocol string, auth *hydra.Auth) {
	var target = fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
	var m = auth.Map()
	URL, _ := url.Parse(target)
	outputHandler(URL, color.Important("CrackSuccess"), m)

	keyword := "CrackSuccess"
	m["Host"] = URL.String()
	m["IP"] = URL.Hostname()
	m["Port"] = URL.Port()
	m["Service"] = URL.Scheme
	m["Keyword"] = keyword

	msi := make(map[string]interface{})
	// 遍历 strMap 中的每个元素
	for key, value := range m {
		// 将字符串值赋给对应的 interface{} 类型
		msi[key] = value
	}
	app.HydraResults = append(app.HydraResults, msi)
}

func outputNmapFinger(URL *url.URL, resp *gonmap.Response) {
	if responseFilter(resp.Raw) == true {
		return
	}
	finger := resp.FingerPrint
	m := misc.ToMap(finger)
	m["Response"] = resp.Raw
	m["IP"] = URL.Hostname()
	m["Port"] = URL.Port()
	//补充归属地信息
	if app.Setting.CloseCDN == false {
		result, _ := cdn.Find(URL.Hostname())
		m["Addr"] = result
	}
	outputHandler(URL, finger.Service, m)

}

// --check   urlcheck
func outputAppFinger(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
	if responseFilter(banner.Response, banner.Cert) == true {
		return
	}
	m := misc.ToMap(finger)
	//补充归属地信息
	if app.Setting.CloseCDN == false {
		result, _ := cdn.Find(URL.Hostname())
		m["Addr"] = result
	}
	m["Service"] = URL.Scheme
	m["FoundDomain"] = banner.FoundDomain
	m["FoundIP"] = banner.FoundIP
	m["Response"] = banner.Response
	m["Cert"] = banner.Cert
	m["Header"] = banner.Header
	m["Body"] = banner.Body
	m["ICP"] = banner.ICP
	m["FingerPrint"] = m["ProductName"]
	delete(m, "ProductName")
	//增加IP、Domain、Port字段
	m["Port"] = uri.GetURLPort(URL)
	if m["Port"] == "" {
		slog.Println(slog.WARN, "无法获取端口号：", URL)
	}
	if hostname := URL.Hostname(); uri.IsIPv4(hostname) {
		m["IP"] = hostname
	} else {
		m["Domain"] = hostname
		if v, ok := scanner.DomainDatabase.Load(hostname); ok {
			m["IP"] = v.(string)
		}
	}
	outputHandler(URL, banner.Title, m)
}

func outputCDNRecord(domain, info string) {
	if responseFilter(info) == true {
		return
	}
	//输出结果
	target := fmt.Sprintf("cdn://%s", domain)
	URL, _ := url.Parse(target)
	outputHandler(URL, "CDN资产", map[string]string{
		"CDNInfo": info,
		"Domain":  domain,
	})
}

func outputUnknownResponse(addr net.IP, port int, response string) {
	if responseFilter(response) == true {
		return
	}
	//输出结果
	target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
	URL, _ := url.Parse(target)
	outputHandler(URL, "无法识别该协议", map[string]string{
		"Response": response,
		"IP":       URL.Hostname(),
		"Port":     strconv.Itoa(port),
	})
}

func outputOpenResponse(addr net.IP, port int) {
	//输出结果
	protocol := gonmap.GuessProtocol(port)
	target := fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
	URL, _ := url.Parse(target)
	outputHandler(URL, "response is empty", map[string]string{
		"IP":   URL.Hostname(),
		"Port": strconv.Itoa(port),
	})
}

func responseFilter(strArgs ...string) bool {
	var match = app.Setting.Match
	var notMatch = app.Setting.NotMatch

	if match != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.Match) == true {
				return false
			}
		}
	}

	if notMatch != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.NotMatch) == true {
				return true
			}
		}
	}
	return false
}

var (
	//disableKey       = []string{"MatchRegexString", "Service", "ProbeName", "Response", "Cert", "Header", "Body", "IP"}
	disableKey       = []string{"MatchRegexString", "ProbeName", "Response", "Cert", "Header", "Body", "IP"}
	ImportantKey     = []string{"ProductName", "DeviceType"}
	VaryImportantKey = []string{"Hostname", "FingerPrint", "ICP"}
)

func getHTTPDigest(s string) string {
	var length = 24
	var digestBuf []rune
	_, body := simplehttp.SplitHeaderAndBody(s)
	body = chinese.ToUTF8(body)
	for _, r := range []rune(body) {
		buf := []byte(string(r))
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if len(buf) > 1 {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(body, length-len(digestBuf))
}

func getRawDigest(s string) string {
	var length = 24
	if len(s) < length {
		return s
	}
	var digestBuf []rune
	for _, r := range []rune(s) {
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if 0x20 <= r && r <= 0x7E {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(s, length-len(digestBuf))
}

//	func outputHandler(URL *url.URL, keyword string, m map[string]string) {
//		m = misc.FixMap(m)
//		if respRaw := m["Response"]; respRaw != "" {
//			if m["Service"] == "http" || m["Service"] == "https" {
//				m["Digest"] = strconv.Quote(getHTTPDigest(respRaw))
//			} else {
//				m["Digest"] = strconv.Quote(getRawDigest(respRaw))
//			}
//		}
//		m["Length"] = strconv.Itoa(len(m["Response"]))
//		sourceMap := misc.CloneMap(m)
//		for _, keyword := range disableKey {
//			delete(m, keyword)
//		}
//		for key, value := range m {
//			if key == "FingerPrint" {
//				continue
//			}
//			m[key] = misc.StrRandomCut(value, 24)
//		}
//		fingerPrint := color.StrMapRandomColor(m, true, ImportantKey, VaryImportantKey)
//		fingerPrint = misc.FixLine(fingerPrint)
//		format := "%-30v %-" + strconv.Itoa(misc.AutoWidth(color.Clear(keyword), 26+color.Count(keyword))) + "v %s"
//		printStr := fmt.Sprintf(format, URL.String(), keyword, fingerPrint)
//		slog.Println(slog.DATA, printStr)
//
//		if jw := app.Setting.OutputJson; jw != nil {
//			sourceMap["URL"] = URL.String()
//			sourceMap["Keyword"] = keyword
//			jw.Push(sourceMap)
//		}
//		if cw := app.Setting.OutputCSV; cw != nil {
//			sourceMap["URL"] = URL.String()
//			sourceMap["Keyword"] = keyword
//			delete(sourceMap, "Header")
//			delete(sourceMap, "Cert")
//			delete(sourceMap, "Response")
//			delete(sourceMap, "Body")
//			sourceMap["Digest"] = strconv.Quote(sourceMap["Digest"])
//			for key, value := range sourceMap {
//				sourceMap[key] = chinese.ToUTF8(value)
//			}
//			cw.Push(sourceMap)
//		}
//	}
//
//
//

// 用于跟踪是否已经打印过表头
var isPortHeaderPrinted = false

func outputHandler(URL *url.URL, keyword string, m map[string]string) {
	m = misc.FixMap(m)

	if respRaw := m["Response"]; respRaw != "" {
		if m["Service"] == "http" || m["Service"] == "https" {
			m["Digest"] = strconv.Quote(getHTTPDigest(respRaw))
		} else {
			m["Digest"] = strconv.Quote(getRawDigest(respRaw))
		}
	}
	m["Length"] = strconv.Itoa(len(m["Response"]))
	sourceMap := misc.CloneMap(m)
	for _, keyword := range disableKey {
		delete(m, keyword)
	}
	for key, value := range m {
		if key == "FingerPrint" {
			m[key] = misc.FixLine(value)
			continue
		}
		m[key] = misc.StrRandomCut(value, 24)
	}
	if m["Service"] == "" {
		m["Service"] = URL.Scheme
	}
	m["ICP"] = runewidth.FillRight(m["ICP"], 19)
	m["ICP"] = misc.FixLine(m["ICP"])
	if m["FingerPrint"] == "" {
		if m["ProductName"] != "" && m["Version"] != "" {
			m["FingerPrint"] = fmt.Sprintf(m["ProductName"] + " " + m["Version"])
		} else if m["ProductName"] != "" {
			m["FingerPrint"] = m["ProductName"]
		} else {
			m["FingerPrint"] = m["Version"]
		}
	}
	m["IP"] = URL.Hostname()
	m["Keyword"] = keyword
	if !app.Args.Check {
		if !isPortHeaderPrinted {
			slog.Println(slog.INFO, "开始PortScan")
			slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
			slog.Println(slog.DATA, fmt.Sprintf("\r| %-15s | %-5s | %-15s | %-7s | %-40s | %-60s ", "IP", "Port", "Service", "OS", "Keyword", "Finger"))
			slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
			isPortHeaderPrinted = true
		}
		m["Host"] = URL.String()
		msi := make(map[string]interface{})
		// 遍历 strMap 中的每个元素
		for key, value := range m {
			// 将字符串值赋给对应的 interface{} 类型
			msi[key] = value
		}
		app.PortScanResults = append(app.PortScanResults, msi)
		slog.Println(slog.DATA, fmt.Sprintf("\r| %-15s | %-5s | %-15s | %-7s | %-"+strconv.Itoa(misc.AutoWidth(m["Keyword"], 40))+"v | %-60s  ", URL.Hostname(), m["Port"], m["Service"], m["OperatingSystem"], m["Keyword"], m["FingerPrint"]))
	} else if !app.Args.Scan {
		if !isPortHeaderPrinted {
			slog.Println(slog.INFO, "开始URLCheck")
			slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
			slog.Println(slog.DATA, fmt.Sprintf("\r| %-45s | %-5s | %-6s | %-19s | %-60s ", "URL", "Port", "Length", "ICP", "Finger"))
			slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
			isPortHeaderPrinted = true
		}
		m["Host"] = URL.String()
		msi := make(map[string]interface{})
		// 遍历 strMap 中的每个元素
		for key, value := range m {
			// 将字符串值赋给对应的 interface{} 类型
			msi[key] = value
		}
		app.UrlCheckResults = append(app.UrlCheckResults, msi)
		slog.Println(slog.DATA, fmt.Sprintf("\r| %-45s | %-5s | %-6s | "+m["ICP"]+" | %-60s ", URL.String(), m["Port"], m["Length"], m["FingerPrint"]))
	}

	//if app.Args.Hydra {
	//	if !isHydraHeaderPrinted {
	//		slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
	//		//fmt.Printf("| %-15s | %-4s | %-10s | %-10s \n", "IP", "port", "Service", "Addr")
	//		slog.Println(slog.DATA, fmt.Sprintf("\r| %-15s | %-5s | %-15s | %-7s | %-30s | %-60s ", "IP", "Port", "Service", "Username", "Password", "keyword"))
	//		slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
	//		isHydraHeaderPrinted = true
	//	}
	//	app.HydraResults = append(app.HydraResults, m)
	//	slog.Println(slog.DATA, fmt.Sprintf("\r| %-15s | %-5s | %-15s | %-7s | %-30s | %-60s ", URL.Hostname(), m["Port"], m["Service"], m["Username"], m["Password"], keyword))
	//}

	if jw := app.Setting.OutputJson; jw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		jw.Push(sourceMap)
	}
	if cw := app.Setting.OutputCSV; cw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		delete(sourceMap, "Header")
		delete(sourceMap, "Cert")
		delete(sourceMap, "Response")
		delete(sourceMap, "Body")
		sourceMap["Digest"] = strconv.Quote(sourceMap["Digest"])
		for key, value := range sourceMap {
			sourceMap[key] = chinese.ToUTF8(value)
		}
		cw.Push(sourceMap)
	}
}

func PrintPocResults() {
	slog.Println(slog.DATA, fmt.Sprintf("+------------------------------------------------------------------------------------------------+"))
	slog.Println(slog.DATA, fmt.Sprintf("| %-45v | %-10v | %-20v ", "Host", "VulnLevel", "POC Name"))
	slog.Println(slog.DATA, fmt.Sprintf("+------------------------------------------------------------------------------------------------+"))
	for _, result := range app.PocResults {
		//red := colorR.BgRed.Render
		format := "| %-45v | %-10v | %-20v "
		printStr := fmt.Sprintf(format, result["Host"], result["VulnLevel"], result["PocName"])
		slog.Println(slog.DATA, printStr)
	}
	if len(app.PocResults) != 0 {
		slog.Println(slog.DATA, fmt.Sprintf("+------------------------------------------------------------------------------------------------+"))
	}
}

func PrintHydraResults() {
	slog.Println(slog.INFO, "开始Hydra")
	slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
	slog.Println(slog.DATA, fmt.Sprintf("\r| %-15s | %-5s | %-15s | %-15s | %-30s | %-60s ", "IP", "Port", "Service", "Username", "Password", "Keyword"))
	slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
	for _, result := range app.HydraResults {
		slog.Println(slog.DATA, fmt.Sprintf("\r| %-15s | %-5s | %-15s | %-15s | %-30s | %-60s ", result["IP"], result["Port"], result["Service"], result["Username"], result["Password"], result["Keyword"]))
	}
	if len(app.HydraResults) != 0 {
		slog.Println(slog.DATA, fmt.Sprintf("+-------------------------------------------------------------------------------------------------------------------------------------+"))
	}
}

func watchDog() {
	for {
		time.Sleep(time.Second * 1)
		var (
			nDomain = DomainScanner.RunningThreads()
			nIP     = IPScanner.RunningThreads()
			nPort   = PortScanner.RunningThreads()
			nURL    = URLScanner.RunningThreads()
			nHydra  = HydraScanner.RunningThreads()
		)
		if time.Now().Unix()%180 == 0 {
			warn := fmt.Sprintf("当前存活协程数：Domain：%d 个，IP：%d 个，Port：%d 个，URL：%d 个，Hydra：%d 个", nDomain, nIP, nPort, nURL, nHydra)
			slog.Println(slog.WARN, warn)
		}
	}
}
