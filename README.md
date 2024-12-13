### 构想

在kscan基础上改进实现一个**轻量级**扫描器，主要用于信息收集和POC探测，旨在资产发现阶段帮助获取更多有限的关键信息。


### 完成项

- 初次运行时通过-v参数初始化生成conf.yml配置文件
- 加入了Ascan模块，通过爱企查查询企业信息(可手动设置控股投资比例递归到子孙企业)，也可与fofa进行联动查询，最终通过kscan完成整个扫描过程。
- 优化了fofa中获取url时泛解析问题
- 添加了敏感信息泄露模块(github gitlab codesearch email信息泄露) ~~效果非常一般(依托答辩) 已移除~~
- 添加了nopoc参数，可自定义选择是否使用Xray Poc进行探测(默认开启)

### TODO

- ascan-f 读取文件 批量扫描 尚存在问题

- 优化输出结果【还在完善中】

- 爆破模块完善【自己导入自定义字典加载】

- 完善POC模块，考虑加入nuclei POC

  参考:

  https://github.com/WAY29/pocV

  https://github.com/jweny/pocassist

  先把poc框架弄出来 再去考虑指纹识别 poc识别探测
  【能不能参考gogo的思路  先去匹配框架(先指纹识别) 在根据poc去扫描探测  减少不必要的发包】

- 完善poc探测规则，先基于特征/指纹匹配poc,如果匹配不到则全部poc扫描

- 考虑是否加入finger指纹识别引擎/工具，替换原有的指纹识别机制

### References

https://github.com/qi4L/Qscan

https://github.com/chainreactors/gogo

https://github.com/niudaii/zpscan

https://github.com/SleepingBag945/dddd【如果直接接入其他项目代码 可参考这个dddd 去接入subfinder nuclei】

### 使用参数

```

  _  __                      _____             __  __            
  | |/ /                     |  __ \           |  \/  |
  | ' / ___  ___ __ _ _ __   | |__) | __ ___   | \  / | __ ___  __
  |  < / __|/ __/ _` | '_ \  |  ___/ '__/ _ \  | |\/| |/ _` \ \/ /
  | . \\__ \ (_| (_| | | | | | |   | | | (_) | | |  | | (_| |>  < 
  |_|\_\___/\___\__,_|_| |_| |_|   |_|  \___/  |_|  |_|\__,_/_/\_\
                        KscanProMax   version: 0.3
usage: kscan [-h,--help,-v,--version,--fofa-syntax] (-t,--target,-f,--fofa,--spy]) [options] [hydra options] [fofa options] [ascan options] 
        使用前请在conf.yml下完成配置文件(初次使用可通过-v生成配置文件)
optional arguments:
  -h , --help     show this help message and exit
  -v , --version  查看版本(第一次运行会生成yaml配置文件)
  -f , --fofa     从fofa获取检测对象，需提前在conf.yml下配置:FOFA_EMAIL、FOFA_KEY 
  -t , --target   指定探测对象：
                  IP地址：114.114.114.114
                  IP地址段：114.114.114.114/24,不建议子网掩码小于12
                  IP地址段：114.114.114.114-115.115.115.115
                  URL地址：https://www.baidu.com
                  文件地址：file:/tmp/target.txt
                  剪切板: paste or clipboard
  --spy           网段探测模式，此模式下将自动探测主机可达的内网网段可接收参数为：
                  (空)、192、10、172、all、指定IP地址(将探测该IP地址B段存活网关)
options:
  --check         针对目标地址做指纹识别，仅不会进行端口探测
  --scan          将针对--fofa、--spy提供的目标对象，进行端口扫描和指纹识别
  --nopoc         在进行扫描时不开启POC探测(默认开启POC探测,目前内置xray)
  -p , --port     扫描指定端口，默认会扫描TOP400，支持：80,8080,8088-8090
  -eP, --excluded-port 跳过扫描指定的端口，支持：80,8080,8088-8090
  -o , --output   将扫描结果保存到文件
  -oJ             将扫描结果使用json格式保存到文件
  -oC             将扫描结果使用csv格式保存到文件
  -oH             将扫描结果使用html格式保存到文件
  -Pn             使用此参数后，将不会进行智能存活性探测，现在默认会开启智能存活性探测，提高效率
  -Cn             使用此参数后，控制台输出结果将不会带颜色。
  -Dn             使用此参数后，将关闭CDN识别功能
  -sV             使用此参数后，将对所有端口进行全探针探测，此参数极度影响效率，慎用！
  --top           扫描经过筛选处理的常见端口TopX，最高支持1000个，默认为TOP400
  --proxy         设置代理(socks5|socks4|https|http)://IP:Port
  --threads       线程参数,默认线程100,最大值为2048
  --path          指定请求访问的目录，只支持单个目录
  --host          指定所有请求的头部Host值
  --timeout       设置超时时间
  --encoding      设置终端输出编码，可指定为：gb2312、utf-8
  --match         对资产返回banner进行检索，剔除不存在关键字的结果记录
  --not-match     对资产返回banner进行检索，剔除存在关键字的结果记录
  --hydra         自动化爆破支持协议：ssh,rdp,ftp,smb,mysql,mssql,oracle,postgresql,mongodb,redis,默认会开启全部
hydra options:
   --hydra-user   自定义hydra爆破用户名:username or user1,user2 or file:username.txt
   --hydra-pass   自定义hydra爆破密码:password or pass1,pass2 or file:password.txt
                  若密码中存在使用逗号的情况，则使用\,进行转义，其他符号无需转义
   --hydra-update 自定义用户名、密码模式，若携带此参数，则为新增模式，会将用户名和密码补充在默认字典后面。否则将替换默认字典。
   --hydra-mod    指定自动化暴力破解模块:rdp or rdp,ssh,smb
fofa options:
   --fofa-syntax  将获取fofa搜索语法说明
   --fofa-size    将设置fofa返回条目数，默认100条
   --fofa-fix-keyword 修饰keyword，该参数中的{}最终会替换成-f参数的值
ascan options:
   --ascan        查询单个企业名称对应的网站信息【仅查询企业信息，不进行扫描】
   --ascan-s      是否递归查询对外投资[90%]以上开业公司信息
   --ascan-b      是否递归查询开业状态的子公司信息
   --ascan-f      批量查询企业名称查询并进行扫描【现在还没写导入文件的功能】
   --ascan-fofa   联动fofa批量查询企业信息

```

### 常用命令

1.企业信息查询后再fofa查询资产，并进行端口扫描

`go run kscan.go  --ascan-fofa  山东XXXX有限公司 --scan  `

```
+-------------------------------------------------------------------------------------------------------------------------------------+
| IP              | Port  | Service         | OS      | Keyword                                  | Finger                                                       
+-------------------------------------------------------------------------------------------------------------------------------------+
| 1xx.xx.xx.xxx| 22    | ssh             |         | ssh                                      | OpenSSH 7.4                                                   
| 115.236.121.126 | 25    | smtp            |         | smtp                                     |                                                               
| 112.230.203.162 | 23    | telnet          |         | telnet                                   |                                                               
| 115.236.121.126 | 110   | pop3            |         | pop3                                     |                                                               
| 115.236.121.126 | 143   | imap            |         | imap                                     |                                                               
| 1xx.xx.xx.xxx| 443   | https           |         | 测试测试测试                       | JQuery;HTML5;xinjingxiang-system;nginx/1.24.0;nginx           
| 1xx.xx.xx.xxx| 3306  | mysql           |         | mysql                                    | MySQL                                                         
| 1xx.xx.xx.xxx| 80    | http            |         | 测试测试测试                       | JQuery;HTML5;xinjingxiang-system;nginx/1.24.0;nginx           
| 1xx.xx.xx.xxx | 443   | https           |         | WebUserLogin                             | PasswordField;TOPSEC                                          
| 1xx.xx.xx.xxx| 8080  | http            |         | 测试测试测试                       | JQuery;HTML5;xinjingxiang-system                              
[+]2024/05/31 15:22:50 程序执行总时长为：[28.703585004s]

```

2.企业信息查询后再fofa查询资产，并进行url存活探测

` go run kscan.go  --ascan-fofa  山东XXXX有限公司 --check --oH 2.html `

```
[+]2024/05/31 16:01:44 开始URLCheck
+-------------------------------------------------------------------------------------------------------------------------------------+
| URL                                           | Port  | Length | ICP                 | Finger                                                       
+-------------------------------------------------------------------------------------------------------------------------------------+
| https://coims.xxxxx.com:4433                | 4433  | 5802   |  | 若依(RuoYi)-管理系统;Bootstrap;zepto.JS;ruoyi-System;Jiusi-OA;PasswordField;JQuery;HTML5;ruoyi-system;若依ruoyi管理系统 
| https://1xx.xx.xx.xxx                      | 443   | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| http://pop.xxxxx.com                        | 80    | 339    |  |                                                              
| http://1xx.xx.xx.xxx                       | 80    | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| http://smtp.xxxxx.com                       | 80    | 339    |  |                                                              
| http://www.xxxxx.com                        | 80    | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| http://imap.xxxxx.com                       | 80    | 339    |  |                                                              
| https://xxxxx.com                           | 443   | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| https://www.xxxxx.com                       | 443   | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| https://lnruyvqdth.xxxxx.cn                 | 443   | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| https://httpsmail.xxxxx.cn                  | 443   | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| http://bb.ccc.dddd.xxxxx.cn                 | 80    | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| http://what.website.xxxxx.cn                | 80    | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| http://ccc.dddd.xxxxx.cn                    | 80    | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system  

```

3.对ip资产进行扫描 默认为端口扫描

`go run kscan.go  -t 1xx.xx.xx.xx`

```
[+]2024/05/31 15:27:17 开始PortScan
+-------------------------------------------------------------------------------------------------------------------------------------+
| IP              | Port  | Service         | OS      | Keyword                                  | Finger                                                       
+-------------------------------------------------------------------------------------------------------------------------------------+
| 1xx.xx.xx.xxx| 22    | ssh             |         | ssh                                      | OpenSSH 7.4                                                   
| 1xx.xx.xx.xxx| 3306  | mysql           |         | mysql                                    | MySQL                                                         
| 1xx.xx.xx.xxx| 443   | https           |         | 测试测试测试                       | JQuery;HTML5;xinjingxiang-system;nginx/1.24.0;nginx           
| 1xx.xx.xx.xxx| 80    | http            |         | 测试测试测试                       | JQuery;HTML5;xinjingxiang-system;nginx/1.24.0;nginx           
| 1xx.xx.xx.xxx| 8080  | http            |         | 测试测试测试                       | JQuery;HTML5;xinjingxiang-system                       
```


` go run kscan.go  -t 1xx.xx.xx.xxx --check`

``` 
[+]2024/05/31 15:27:58 开始URLCheck
+-------------------------------------------------------------------------------------------------------------------------------------+
| URL                                           | Port  | Length | ICP                 | Finger                                                       
+-------------------------------------------------------------------------------------------------------------------------------------+
| https://1xx.xx.xx.xxx                      | 443   | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
| http://1xx.xx.xx.xxx                       | 80    | 62394  | ICP测试测试测试 | JQuery;HTML5;xinjingxiang-system                             
[+]2024/05/31 15:28:11 程序执行总时长为：[15.38529118s]

```

4.通过fofa进行信息收集，同上可进行端口扫描或url探活
`go run kscan.go  -f 1xx.xx.xx.xxx`

```
+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Host                                          | Title                          | IP              | Domain               | Protocol   | Addr                      | Server 
+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| http://1xx.xx.xx.xxx                       | 测试测试测试             | 1xx.xx.xx.xxx|                      | http       |                           | nginx/1.24.0 
| http://1xx.xx.xx.xxx                       | HTTP                           | 1xx.xx.xx.xxx|                      | http       |                           | nginx/1.24.0 
| https://1xx.xx.xx.xxx                      | HTTPS                          | 1xx.xx.xx.xxx|                      | https      |                           | nginx/1.24.0 
| http://www.xxxxx.com                        | 测试测试测试             | 1xx.xx.xx.xxx| xxxxx.com          | http       |                           | nginx/1.24.0 
| https://1xx.xx.xx.xxx                      | 测试测试测试             | 1xx.xx.xx.xxx|                      | https      |                           | nginx/1.24.0 
| https://xxxxx.com                           | 测试测试测试             | 1xx.xx.xx.xxx| xxxxx.com          | https      |                           | nginx/1.24.0 
| https://www.xxxxx.com                       | 测试测试测试             | 1xx.xx.xx.xxx| xxxxx.com          | https      |                           | nginx/1.24.0 
| http://xxxxx.com                            | 测试测试测试             | 1xx.xx.xx.xxx| xxxxx.com          | http       |                           | nginx/1.24.0 
| https://ctf.xxxxx.com                       | 测试测试测试             | 1xx.xx.xx.xxx| xxxxx.com          | https      |                           | nginx/1.24.0 
| http://ctf.xxxxx.com                        | 301MovedPermanently            | 1xx.xx.xx.xxx| xxxxx.com          | http       |                           | nginx/1.24.0 
| ssh://1xx.xx.xx.xx:22                      | SSH                            | 1xx.xx.xx.xxx|                      | ssh        |                           |  
| mysql://1xx.xx.xx.xx:3306                  | MYSQL                          | 1xx.xx.xx.xxx|                      | mysql      |                           |  
[*]2024/05/31 15:55:36 可以使用--check参数对fofa扫描结果进行存活性及指纹探测，也可以使用--scan参数对fofa扫描结果进行端口扫描

```
