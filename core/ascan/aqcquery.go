package ascan

import (
	"KscanPro/app"
	"KscanPro/core/ascan/ai"
	"KscanPro/core/ascan/common"
	"KscanPro/core/ascan/common/utils"
	"KscanPro/core/ascan/runner"
	"KscanPro/core/slog"
	"fmt"
	"os"
	"testing"
)

//func main() {
//	runAQCQuery()
//}

var (
	Keyword   = ""
	Invert    = float64(0.0)
	GetBranch = false
	Target    []string
)

func Start() map[string]bool {
	slog.Println(slog.INFO, "将开始对输入的企业名称进行信息收集")
	homeSiteMaps := RunAQCQuery()

	return homeSiteMaps
}

func RunAQCQuery() map[string]bool {
	var op common.Options
	conf := common.CheckConf()
	op.KeyWord = app.Setting.Ascan
	op.IsGetBranch = app.Setting.GetBranch
	op.Invest = app.Setting.Invest
	//slog.Println(slog.INFO, "输入的Inverst值为：", op.Invest)
	op.CookieInfo = conf.Cookies
	//println("op.GetBranch  bool:", op.IsGetBranch)
	//op.Output = conf.Output
	runner.RunEnumeration(op)

	return ai.UniqueHomeSites
}

func TestrunAQCQuery(t *testing.T) {
	var op common.Options
	conf := common.CheckConf()
	if !utils.FolderExists(conf.Output) {
		os.Mkdir(conf.Output, 0777)
	}
	common.Flag(&op)
	common.Parse(&op)
	op.CookieInfo = conf.Cookies
	op.Output = conf.Output
	runner.RunEnumeration(op)

	//for domain, info := range ai.DomainMap {
	//	fmt.Printf("Domain: %s\nInfo: %v\n", domain, info)
	//}
	for domain, _ := range ai.DomainMap {
		fmt.Printf("Domain: %s\n", domain)
	}
}
