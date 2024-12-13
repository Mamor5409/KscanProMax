package leakscan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type SearchcodeResp struct {
	Matchterm    string `json:"matchterm"`
	Previouspage int    `json:"previouspage"`
	Searchterm   string `json:"searchterm"`
	Query        string `json:"query"`
	Total        int    `json:"total"`
	Page         int    `json:"page"`
	Nextpage     int    `json:"nextpage"`
	Results      []struct {
		Repo       string   `json:"repo"`
		Language   string   `json:"language"` // 语言 根据后缀去判断
		Linescount int      `json:"linescount"`
		Location   string   `json:"location"`
		Name       string   `json:"name"`
		Url        string   `json:"url"` // searchcode中文件链接
		Md5Hash    string   `json:"md5hash"`
		Lines      struct { // 关键词出现的相关内容
			Field1   string `json:"4,omitempty"`
			Field2   string `json:"5,omitempty"`
			Field3   string `json:"6,omitempty"`
			Field4   string `json:"256,omitempty"`
			Field5   string `json:"257,omitempty"`
			Field6   string `json:"258,omitempty"`
			Field7   string `json:"77,omitempty"`
			Field8   string `json:"78,omitempty"`
			Field9   string `json:"79,omitempty"`
			Field10  string `json:"14,omitempty"`
			Field11  string `json:"15,omitempty"`
			Field12  string `json:"16,omitempty"`
			Field13  string `json:"48,omitempty"`
			Field14  string `json:"49,omitempty"`
			Field15  string `json:"50,omitempty"`
			Field16  string `json:"104,omitempty"`
			Field17  string `json:"105,omitempty"`
			Field18  string `json:"106,omitempty"`
			Field19  string `json:"41,omitempty"`
			Field20  string `json:"42,omitempty"`
			Field21  string `json:"43,omitempty"`
			Field22  string `json:"55,omitempty"`
			Field23  string `json:"56,omitempty"`
			Field24  string `json:"57,omitempty"`
			Field25  string `json:"53,omitempty"`
			Field26  string `json:"54,omitempty"`
			Field27  string `json:"68,omitempty"`
			Field28  string `json:"69,omitempty"`
			Field29  string `json:"70,omitempty"`
			Field30  string `json:"26,omitempty"`
			Field31  string `json:"27,omitempty"`
			Field32  string `json:"28,omitempty"`
			Field33  string `json:"94,omitempty"`
			Field34  string `json:"95,omitempty"`
			Field35  string `json:"96,omitempty"`
			Field36  string `json:"7,omitempty"`
			Field37  string `json:"8,omitempty"`
			Field38  string `json:"245,omitempty"`
			Field39  string `json:"246,omitempty"`
			Field40  string `json:"247,omitempty"`
			Field41  string `json:"52,omitempty"`
			Field42  string `json:"171,omitempty"`
			Field43  string `json:"172,omitempty"`
			Field44  string `json:"173,omitempty"`
			Field45  string `json:"37,omitempty"`
			Field46  string `json:"38,omitempty"`
			Field47  string `json:"39,omitempty"`
			Field48  string `json:"10,omitempty"`
			Field49  string `json:"11,omitempty"`
			Field50  string `json:"12,omitempty"`
			Field51  string `json:"33,omitempty"`
			Field52  string `json:"34,omitempty"`
			Field53  string `json:"35,omitempty"`
			Field54  string `json:"45,omitempty"`
			Field55  string `json:"46,omitempty"`
			Field56  string `json:"47,omitempty"`
			Field57  string `json:"64,omitempty"`
			Field58  string `json:"65,omitempty"`
			Field59  string `json:"66,omitempty"`
			Field60  string `json:"91,omitempty"`
			Field61  string `json:"92,omitempty"`
			Field62  string `json:"93,omitempty"`
			Field63  string `json:"102,omitempty"`
			Field64  string `json:"103,omitempty"`
			Field65  string `json:"114,omitempty"`
			Field66  string `json:"115,omitempty"`
			Field67  string `json:"116,omitempty"`
			Field68  string `json:"124,omitempty"`
			Field69  string `json:"125,omitempty"`
			Field70  string `json:"126,omitempty"`
			Field71  string `json:"127,omitempty"`
			Field72  string `json:"128,omitempty"`
			Field73  string `json:"129,omitempty"`
			Field74  string `json:"73,omitempty"`
			Field75  string `json:"74,omitempty"`
			Field76  string `json:"75,omitempty"`
			Field77  string `json:"9,omitempty"`
			Field78  string `json:"25,omitempty"`
			Field79  string `json:"153,omitempty"`
			Field80  string `json:"154,omitempty"`
			Field81  string `json:"155,omitempty"`
			Field82  string `json:"109,omitempty"`
			Field83  string `json:"110,omitempty"`
			Field84  string `json:"111,omitempty"`
			Field85  string `json:"44,omitempty"`
			Field86  string `json:"63,omitempty"`
			Field87  string `json:"76,omitempty"`
			Field88  string `json:"13,omitempty"`
			Field89  string `json:"23,omitempty"`
			Field90  string `json:"24,omitempty"`
			Field91  string `json:"36,omitempty"`
			Field92  string `json:"80,omitempty"`
			Field93  string `json:"81,omitempty"`
			Field94  string `json:"51,omitempty"`
			Field95  string `json:"3,omitempty"`
			Field96  string `json:"17,omitempty"`
			Field97  string `json:"18,omitempty"`
			Field98  string `json:"19,omitempty"`
			Field99  string `json:"31,omitempty"`
			Field100 string `json:"32,omitempty"`
			Field101 string `json:"2,omitempty"`
			Field102 string `json:"20,omitempty"`
			Field103 string `json:"21,omitempty"`
			Field104 string `json:"146,omitempty"`
			Field105 string `json:"147,omitempty"`
			Field106 string `json:"148,omitempty"`
			Field107 string `json:"179,omitempty"`
			Field108 string `json:"180,omitempty"`
			Field109 string `json:"181,omitempty"`
			Field110 string `json:"72,omitempty"`
			Field111 string `json:"71,omitempty"`
			Field112 string `json:"143,omitempty"`
			Field113 string `json:"144,omitempty"`
			Field114 string `json:"145,omitempty"`
			Field115 string `json:"40,omitempty"`
			Field116 string `json:"101,omitempty"`
			Field117 string `json:"112,omitempty"`
			Field118 string `json:"108,omitempty"`
			Field119 string `json:"113,omitempty"`
			Field120 string `json:"117,omitempty"`
			Field121 string `json:"118,omitempty"`
			Field122 string `json:"119,omitempty"`
			Field123 string `json:"120,omitempty"`
			Field124 string `json:"30,omitempty"`
			Field125 string `json:"100,omitempty"`
			Field126 string `json:"98,omitempty"`
			Field127 string `json:"99,omitempty"`
			Field128 string `json:"82,omitempty"`
			Field129 string `json:"83,omitempty"`
			Field130 string `json:"84,omitempty"`
			Field131 string `json:"58,omitempty"`
			Field132 string `json:"59,omitempty"`
			Field133 string `json:"141,omitempty"`
			Field134 string `json:"142,omitempty"`
			Field135 string `json:"138,omitempty"`
			Field136 string `json:"139,omitempty"`
			Field137 string `json:"140,omitempty"`
			Field138 string `json:"223,omitempty"`
			Field139 string `json:"224,omitempty"`
			Field140 string `json:"225,omitempty"`
			Field141 string `json:"107,omitempty"`
			Field142 string `json:"67,omitempty"`
			Field143 string `json:"87,omitempty"`
			Field144 string `json:"88,omitempty"`
			Field145 string `json:"89,omitempty"`
			Field146 string `json:"1,omitempty"`
			Field147 string `json:"220,omitempty"`
			Field148 string `json:"221,omitempty"`
			Field149 string `json:"222,omitempty"`
			Field150 string `json:"206,omitempty"`
			Field151 string `json:"207,omitempty"`
			Field152 string `json:"208,omitempty"`
			Field153 string `json:"209,omitempty"`
			Field154 string `json:"210,omitempty"`
			Field155 string `json:"211,omitempty"`
			Field156 string `json:"212,omitempty"`
			Field157 string `json:"213,omitempty"`
			Field158 string `json:"214,omitempty"`
			Field159 string `json:"29,omitempty"`
			Field160 string `json:"199,omitempty"`
			Field161 string `json:"200,omitempty"`
			Field162 string `json:"201,omitempty"`
			Field163 string `json:"134,omitempty"`
			Field164 string `json:"135,omitempty"`
			Field165 string `json:"136,omitempty"`
			Field166 string `json:"187,omitempty"`
			Field167 string `json:"188,omitempty"`
			Field168 string `json:"189,omitempty"`
			Field169 string `json:"137,omitempty"`
			Field170 string `json:"97,omitempty"`
			Field171 string `json:"235,omitempty"`
			Field172 string `json:"236,omitempty"`
			Field173 string `json:"237,omitempty"`
			Field174 string `json:"248,omitempty"`
			Field175 string `json:"249,omitempty"`
			Field176 string `json:"85,omitempty"`
			Field177 string `json:"86,omitempty"`
			Field178 string `json:"22,omitempty"`
			Field179 string `json:"60,omitempty"`
			Field180 string `json:"61,omitempty"`
		} `json:"lines"`
		Id       int    `json:"id"`
		Filename string `json:"filename"` //文件名 +后缀
	} `json:"results"`
	LanguageFilters []struct {
		Count    int    `json:"count"`
		Language string `json:"language"`
		Id       int    `json:"id"`
	} `json:"language_filters"`
	SourceFilters []struct {
		Count  int    `json:"count"`
		Source string `json:"source"`
		Id     int    `json:"id"`
	} `json:"source_filters"`
}

type ezoutputResult2 struct {
	SourceCome         string `json:"sourceCome"`
	RepositoryFullName string `json:"repositoryFullName"`
	//FileUrl            string `json:"fileUrl"`
}

func SearchBySearchcodeApi(strCode string) []ezoutputResult2 {
	var searchcodeResp SearchcodeResp
	var ezoutputresult2 ezoutputResult2
	var ezoutputresult2arr []ezoutputResult2
	apiurl := "https://searchcode.com/api/codesearch_I/?q=" + strCode + "&p=1&per_page=200" //per_page 参数没用

	req, err := http.NewRequest("GET", apiurl, nil)

	if err != nil {
		return nil
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[-]发送请求错误:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		htmlData := string(body)
		//fmt.Println(htmlData)

		err = json.Unmarshal([]byte(htmlData), &searchcodeResp)
		if err != nil {
			fmt.Println("[-] Error:", err)
			return nil
		}

		if searchcodeResp.Total == 0 {
			fmt.Println("[-] 通过Searchcode 未找到该关键词相关代码信息 :(")
			return nil
		} else {
			//fmt.Printf("[+] 找到 %d 条\n", searchcodeResp.Total)
			for _, result := range searchcodeResp.Results {
				//fmt.Printf("[+] ----------第 %d 条 ---------\n", i+1)
				//fmt.Println("[+] Repo:", result.Repo)
				//fmt.Println("[+] Language:", result.Language)
				//fmt.Println("[+] Name:", result.Name)
				//fmt.Println("[+] FileName:", result.Filename)
				//fmt.Println("[+] FileUrl:", result.Url)
				//fmt.Println("[+] FileContent:", result.Lines)

				//ezoutputresult2.FileUrl = result.Url
				ezoutputresult2.RepositoryFullName = result.Repo
				ezoutputresult2.SourceCome = getSourceCome(result.Repo)
				ezoutputresult2arr = append(ezoutputresult2arr, ezoutputresult2)
			}

		}
	} else {
		fmt.Println("[-] Error , Status Code:", resp.StatusCode)
	}
	return ezoutputresult2arr
}

func GetJsonBySearchCode(strcode string) string {
	res := SearchBySearchcodeApi(strcode)
	jsonBody, err := json.Marshal(res)
	if err != nil {
		fmt.Println("json解析错误:", err)
		return ""
	}
	return string(jsonBody)
}

func getSourceCome(url string) string {
	re := regexp.MustCompile(`https?://([\w-]+\.)*([\w-]*)`)
	subs := re.FindStringSubmatch(url)

	if len(subs) > 2 {
		return strings.TrimSuffix(subs[1], ".")
	}

	return ""
}
