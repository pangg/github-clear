### 利用selenium+chromedriver 清空指定github账户下所有的项目仓库

> chrome和chromedriver版本需要对应  
> chrome和chromedriver版本需要对应  
> chrome和chromedriver版本需要对应  
>重要的事情说三遍......

#### 1. 当前目录chromedriver文件版本
* chrome  72.0+
* chromedriver mac版本
* 其他（linux、windows）版本下载地址：http://chromedriver.storage.googleapis.com/index.html?path=72.0.3626.7/


#### 2. 使用
2.1. 帮助说明
* 根据chrome版本下载chromedriver，放在当前drivers目录下[如果你的环境和我的一致，本步奏忽略]
* 显示帮助信息： `go run main.go` 或者 `./main`
```shell
Usage of main_go:
  -d string
    	chromedriver绝对路径 (default "./chromedriver")
  -p string
    	github登录密码
  -s int
    	chromedriver server端口 (default 9515)
  -u string
    	github登录用户名
```

2.2 开始运行
运行命令： `go run main.go -u uname -p pwd`
或者 `./mian -u uname -p pwd`

#### 3. 附录：
> chromedriver之 与 chrome 版本映射表及其下载地址:
https://blog.csdn.net/huilan_same/article/details/51896672



#### 4. 待优化
利用多线程有点bug
`limiter := make(chan int, 1)`  
1改成2 会出现重复
