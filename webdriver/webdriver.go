package webdriver

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"

	"schanclient/urls"
)

func NewHeadless() (context.Context, context.CancelFunc) {
	opts := make([]chromedp.ExecAllocatorOption, 0)
	opts = append(opts, chromedp.ProxyServer("http://127.0.0.1:8118"))
	opts = append(opts, chromedp.Flag("headless", true))
	allocator, cancel := chromedp.NewAllocator(context.Background(), chromedp.WithExecAllocator(opts...))
	return allocator, cancel
}

// 获得账户登录的cookie
func GetSChannelAuth(user, passwd string) chromedp.Tasks {
	return chromedp.Tasks{ // tasks就是一系列chrome动作的组合
		// 访问URL
		chromedp.Navigate(urls.RootPath),
		chromedp.Navigate(urls.AuthPath),
		// 输入form的email和password
		chromedp.SendKeys("inputEmail", user, chromedp.ByID),
		chromedp.SendKeys("inputPassword", passwd, chromedp.ByID),
		// 提交表单
		chromedp.Submit("div.logincontainer form", chromedp.ByQuery),
		// 等待dologin.php完成auth并进行页面跳转
		chromedp.Sleep(3 * time.Second),
	}
}

// 获取产品列表
func GetServiceList(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 访问产品列表
		chromedp.Navigate(urls.ServiceListPath),
		// 等待直到body加载完毕
		chromedp.WaitReady("tableServicesList", chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		// 选择显示可用服务，暂不支持查看其他类型的服务
		chromedp.Click("Primary_Sidebar-My_Services_Status_Filter-Active", chromedp.ByID),
		chromedp.Sleep(2 * time.Second),
		// 获取获取产品列表HTML，由parser继续分析
		chromedp.OuterHTML("#tableServicesList_wrapper table", res, chromedp.ByQuery),
	}
}

// 获取账户界面的信息panel的HTML，后续由parser解析
func GetDataPanel(url string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("tabOverview", chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		chromedp.OuterHTML("#tabOverview div.plugin", res, chromedp.ByQuery),
	}
}
