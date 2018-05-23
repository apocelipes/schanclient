package webdriver

import (
	"context"
	"errors"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"

	"schanclient/urls"
)

func dropChromeLogs(s string, v ...interface{}) {
	return
}

func NewHeadless(ctx context.Context, starturl string) (*chromedp.CDP, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("create canceled")
	default:
		run, err := runner.New(runner.Flag("headless", true),
			runner.StartURL(starturl))

		if err != nil {
			return nil, err
		}

		err = run.Start(ctx)
		if err != nil {
			return nil, err
		}

		c, err := chromedp.New(ctx, chromedp.WithRunner(run), chromedp.WithErrorf(chromedp.LogFunc(dropChromeLogs)))
		if err != nil {
			return nil, err
		}

		return c, nil
	}
}

// 获得账户登录的cookie
func GetSChannelAuth(user, passwd string) chromedp.Tasks {
	return chromedp.Tasks{ // tasks就是一系列chrome动作的组合
		// 访问URL
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

// ActionFunc
func setAllCookies(ctx context.Context, h cdp.Executor) error {
	// 获取cookie
	cookies, err := network.GetAllCookies().Do(ctx, h)
	if err != nil {
		return err
	}

	for _, v := range cookies {
		// 设置超时，TimeSinceEpoch是time.Time的自定义类型，做一下转换
		expr := cdp.TimeSinceEpoch(time.Unix(int64(v.Expires), 0))
		// 设置cookie属性，全部取自从服务器获取的cookie值
		success, err := network.SetCookie(v.Name, v.Value).
			WithExpires(&expr).
			WithDomain(v.Domain).
			WithPath(v.Path).
			WithHTTPOnly(true).
			Do(ctx, h)
		if err != nil {
			return err
		}
		if !success {
			return errors.New("could not set cookie")
		}
	}

	return nil
}
