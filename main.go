//author:lychao8<lychao_vip@163.com>
//date:2109-02-18
package main

import (
	"github.com/tebeka/selenium"
	"fmt"
	"github.com/tebeka/selenium/chrome"
	"time"
	"os"
	"flag"
	"runtime"
	"sync"
)

var (
	uname        = flag.String("u", "", "github登录用户名")
	pwd          = flag.String("p", "", "github登录密码")
	seleniumPath = flag.String("d", `./drivers/chromedriver`, "chromedriver路径")
	port         = flag.Int("s", 9515, "chromedriver server端口")
)

// getWD get WebDriver server
func getWD() (svc *selenium.Service, wd selenium.WebDriver, err error) {
	opts := []selenium.ServiceOption{}
	svc, err = selenium.NewChromeDriverService(*seleniumPath, *port, opts...)
	if nil != err {
		fmt.Println("start a chromedriver service falid", err.Error())
		return nil, nil, err
	}

	//注意这里，server关闭之后，chrome窗口也会关闭
	//defer svc.Stop()

	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)

	// 调起chrome浏览器
	wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", *port))
	if err != nil {
		fmt.Println("connect to the webDriver faild", err.Error())
		return nil, nil, err
	}

	return
}

func login(wd selenium.WebDriver) (flag bool, err error) {
	loginUrl := "https://github.com/login"

	wd.Get(loginUrl)
	elem, err := wd.FindElement(selenium.ByCSSSelector, "#login_field")
	if nil != err {
		return false, err
	}
	elem.Clear()
	elem.SendKeys(*uname)

	pwdelem, err := wd.FindElement(selenium.ByCSSSelector, "#password")
	if nil != err {
		return false, err
	}
	pwdelem.Clear()
	pwdelem.SendKeys(*pwd)

	formelem, err := wd.FindElement(selenium.ByTagName, "form")
	if nil != err {
		return false, err
	}
	formelem.Submit()
	return true, nil
}

func getProjectList(wd selenium.WebDriver) (list []string) {
	wd.Get(fmt.Sprintf("https://github.com/%s?tab=repositories", *uname))
	elems, _ := wd.FindElements(selenium.ByXPATH, `//*[@id="user-repositories-list"]/ul/li[*]/div[1]/div[1]/h3/a`)
	for _, aele := range elems {
		href, _ := aele.GetAttribute("href")
		list = append(list, href)
	}
	fmt.Println("url_list:", list)
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	now := time.Now()
	flag.Parse()

	if *uname == "" || *pwd == "" {
		flag.Usage()
		return
	}

	svc, wd, err := getWD()
	if err != nil {
		os.Exit(1)
	}
	defer svc.Stop()

	//关闭一个webDriver会对应关闭一个chrome窗口
	//但是不会导致seleniumServer关闭
	defer wd.Quit()

	//登录
	_, err = login(wd)
	if err != nil {
		fmt.Println("github登录失败")
		os.Exit(1)
	}

	//获取项目仓库列表
	url_list := getProjectList(wd)

	var wg sync.WaitGroup
	limiter := make(chan int, 1)
	for _, url := range url_list {
		//del(wd, url)
		wg.Add(1)
		limiter <- 1
		go func(wd selenium.WebDriver, url string) {
			defer func() {
				wg.Done()
				<-limiter
				fmt.Println("from ch", time.Now().Second())
			}()
			del(wd, url)
			return
		}(wd, url)
	}
	wg.Wait()

	wd.Get(fmt.Sprintf("https://github.com/%s?tab=repositories", *uname))
	fmt.Printf("任务完成,耗时:%#vs,3s后关闭浏览器!\n", time.Now().Sub(now).Seconds())
	time.Sleep(3 * time.Second)
}

func del(wd selenium.WebDriver, url string) {
	//当前项目url
	wd.Get(url)

	//setting导航元素
	setelem, _ := wd.FindElement(selenium.ByXPATH, `//*[@id="js-repo-pjax-container"]/div[1]/nav/a[last()]`)
	sethref, _ := setelem.GetAttribute("href")
	wd.Get(sethref)

	//project-name
	elem, _ := wd.FindElement(selenium.ByXPATH, `//*[@id="rename-field"]`)
	project_name, _ := elem.GetAttribute("value")
	fmt.Println("project_name:", project_name)

	//删除操作元素
	delelem, _ := wd.FindElement(selenium.ByXPATH, `//*[@id="options_bucket"]/div[8]/ul/li[4]/details/summary`)
	delelem.Click()

	//弹窗的input元素
	inpelem, _ := wd.FindElement(selenium.ByXPATH, `//*[@id="options_bucket"]/div[8]/ul/li[4]/details/details-dialog/div[3]/form/p/input`)
	inpelem.SendKeys(project_name)

	//确定删除元素
	butelem, _ := wd.FindElement(selenium.ByXPATH, `//*[@id="options_bucket"]/div[8]/ul/li[4]/details/details-dialog/div[3]/form/button`)
	butelem.Click()
}
