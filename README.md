### 构想

实现一个**轻量级**扫描器，主要用于信息收集和POC探测，旨在资产发现阶段帮助获取更多有限的关键信息。


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

  https://github.com/qi4L/Qscan

  https://github.com/jjf012/gopoc (KscanPro已经接入了这个)

  https://github.com/WAY29/pocV

  https://github.com/jweny/pocassist

  先把poc框架弄出来 再去考虑指纹识别 poc识别探测
  【能不能参考gogo的思路  先去匹配框架(先指纹识别) 在根据poc去扫描探测  减少不必要的发包】

- 完善poc探测规则，先基于特征/指纹匹配poc,如果匹配不到则全部poc扫描

- 考虑是否加入finger指纹识别引擎/工具，替换原有的指纹识别机制

### References

https://github.com/chainreactors/gogo

https://github.com/niudaii/zpscan

https://github.com/SleepingBag945/dddd【如果直接接入其他项目代码 可参考这个dddd 去接入subfinder nuclei】


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
