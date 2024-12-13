package app

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"sync"
)

type HTMLWriter struct {
	HTMLBuffer bytes.Buffer
	mutex      *sync.Mutex
}

type HeadData struct {
	Target   string
	TestDate string
	Scope    string
}

type ScanResult struct {
	ScanResultText string
	Results        []map[string]interface{}
	//PocResults     []map[string]interface{}
}

var (
	HTMLBuffer bytes.Buffer
	T          *template.Template
)

// 渲染文件结尾并整合写入文件
func WriteToHtml(path string) {
	err := T.ExecuteTemplate(&HTMLBuffer, "footer", "")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, HTMLBuffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

// 渲染结果
func RenderResult(outResult ScanResult) {
	err := T.ExecuteTemplate(&HTMLBuffer, "result", outResult)
	if err != nil {
		panic(err)
	}

}

func safeHTML(html string) template.HTML {
	return template.HTML(html)
}

// 读取模板文件并渲染header
func ReadTemplate(Data HeadData) {
	var err error
	T = template.New("template.html").Funcs(template.FuncMap{"safeHTML": safeHTML})
	T, err = T.ParseFiles("static/template.html")
	//T, err = template.ParseFiles("static/template.html")
	if err != nil {
		panic(err)
	}
	//渲染header
	err = T.ExecuteTemplate(&HTMLBuffer, "header", Data)
	if err != nil {
		panic(err)
	}

}

var FofaResults []map[string]interface{}     //全局Port_result，用于存储端口扫描结果供html保存
var PortScanResults []map[string]interface{} //全局Port_result，用于存储端口扫描结果供html保存
var UrlCheckResults []map[string]interface{} //全局Port_result，用于存储端口扫描结果供html保存
var HydraResults []map[string]interface{}    //全局Port_result，用于存储端口扫描结果供html保存
var PocResults []map[string]interface{}      //POC检测的结果

func HTMLout() {
	if Setting.OutputHtml != nil {
		resultSets := map[string]*[]map[string]interface{}{
			"fofa执行结果": &FofaResults,
			"端口扫描结果":   &PortScanResults,
			"URL检测结果":  &UrlCheckResults,
			"POC检测结果":  &PocResults,
			"弱口令检测结果":  &HydraResults,
		}
		for resultText, resultSet := range resultSets {

			if len(*resultSet) != 0 {
				data := ScanResult{
					ScanResultText: resultText,
					Results:        *resultSet,
				}
				RenderResult(data)

			}
		}
	}
}
