package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chromedp/chromedp"

	"schanclient/config"
	"schanclient/parser"
	_ "schanclient/pyclient"
	"schanclient/ssr"
	"schanclient/webdriver"
)

func main() {
	// 命令行参数
	showInfo := flag.Bool("show-info", false, "show client's infomation.")
	showUsed := flag.Bool("show-used", false, "show user's data use.")
	showConf := flag.Bool("show-conf", false, "show current ssrnode config.")
	setNode := flag.String("set-node", "", "set SSR node to ssr_config_file.")
	clientFlag := flag.String("d", "", "start, stop or restart the ssr client.")
	flag.Parse()

	// 信号处理
	ctxt, AllocatorCancel := webdriver.NewHeadless()
	defer AllocatorCancel()
	c, cancel := chromedp.NewContext(ctxt)
	defer cancel()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	go func(allocatorCancel, cancel context.CancelFunc) {
		select {
		case <-sig:
			cancel()
			allocatorCancel()
		}
	}(AllocatorCancel, cancel)

	// 获取配置
	conf := new(config.UserConfig)
	err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("load config error: %v\n", err)
	}

	// 重定向log
	logpath, err := conf.LogFile.AbsPath()
	if err != nil {
		log.Fatalf("load LogFile path error: %v\n", err)
	}
	l, _ := os.OpenFile(logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	log.SetOutput(l)

	// 显示ssr node配置
	if *showConf {
		showNodeConf(conf)
		return
	}

	// 启动ssr
	if *clientFlag != "" {
		var ssrl ssr.SSRLauncher
		ssrl = ssr.NewLauncher("python", conf)
		if ssrl == nil {
			log.Fatalln("create ssrclient error.")
		}

		switch *clientFlag {
		case "start":
			err := ssrl.Start()
			if err != nil {
				log.Fatalf("start ssrclient error: %v\n", err)
			}
			fmt.Println("已开始")
		case "stop":
			err := ssrl.Stop()
			if err != nil {
				log.Fatalf("stop ssrclient error: %v\n", err)
			}
			fmt.Println("已停止")
		case "restart":
			err := ssrl.Restart()
			if err != nil {
				log.Fatalf("restart ssrclient error: %v\n", err)
			}
			fmt.Println("已重启")
		default:
			log.Fatalln("Unknow -d flags.")
		}

		// -d flag can't use with other flags
		return
	}

	// 查询信息
	if *showInfo || *showUsed || *setNode != "" {
		// 登录schannel
		err = chromedp.Run(ctxt, webdriver.GetSChannelAuth(conf.UserName, conf.Passwd))
		if err != nil {
			log.Fatalln(err.Error() + "GetSChannelAuth")
		}

		// 关闭chrome
		defer func() {
			cancel()
			AllocatorCancel()
			err = chromedp.FromContext(c).Browser.Shutdown()
			if err != nil {
				log.Fatalln(err.Error() + "shutdown")
			}
		}()
	} else {
		flag.Usage()
		return
	}

	if *showInfo {
		fmt.Println("show info")
		showServiceInfo(c)
		fmt.Println("info down")
		return
	}

	if *showUsed {
		showServiceUsed(c)
		return
	}

	if *setNode != "" {
		ssrconfpath, err := conf.SSRConfigPath.AbsPath()
		if err != nil {
			log.Fatalln(err)
		}

		setSSRNode(c, *setNode, ssrconfpath)
		return
	}
}

func showServiceInfo(ctxt context.Context) {
	// run task list
	var res string
	err := chromedp.Run(ctxt, webdriver.GetServiceList(&res))
	if err != nil {
		log.Fatalln(err.Error() + "GetServiceList")
	}

	arr := parser.GetService(res)
	for _, v := range arr {
		var data string
		err = chromedp.Run(ctxt, webdriver.GetDataPanel(v.Link, &data))
		if err != nil {
			log.Fatalln(err.Error() + "GetDataPanel")
		}

		info := parser.GetSSRInfo(data, v)
		fmt.Println(info)
	}
}

func showServiceUsed(ctxt context.Context) {
	var res string
	err := chromedp.Run(ctxt, webdriver.GetServiceList(&res))
	if err != nil {
		log.Fatalln(err.Error() + "GetServiceList")
	}

	arr := parser.GetService(res)
	for _, v := range arr {
		var data string
		err = chromedp.Run(ctxt, webdriver.GetDataPanel(v.Link, &data))
		if err != nil {
			log.Fatalln(err.Error() + "GetDataPanel")
		}

		info := parser.GetSSRInfo(data, v)
		fmt.Println(info.UsedInfo())
	}
}

func setSSRNode(ctxt context.Context, nodename, path string) {
	var res string
	err := chromedp.Run(ctxt, webdriver.GetServiceList(&res))
	if err != nil {
		log.Fatalln(err.Error() + "GetServiceList")
	}

	ser := parser.GetService(res)[0]
	var data string
	err = chromedp.Run(ctxt, webdriver.GetDataPanel(ser.Link, &data))
	if err != nil {
		log.Fatalln(err.Error() + "GetDataPanel")
	}
	// 获得套餐信息
	info := parser.GetSSRInfo(data, ser)

	// 取得节点
	node := info.GetNodeByName(nodename)
	if node == nil {
		log.Fatalln("node: " + nodename + " is not exists.")
	}

	err = node.Store(path)
	if err != nil {
		log.Fatalf("store node error: %v\n", err)
	}
	fmt.Printf("设置成功。\n节点名称：%v\n文件位置：%v\n", nodename, path)
}

func showNodeConf(conf *config.UserConfig) {
	path, err := conf.SSRConfigPath.AbsPath()
	if err != nil {
		log.Fatalf("show ssr node conf get conf path: %v\n", err)
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalln("show ssr node conf:", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("show ssr node conf read configure:", err)
	}

	fmt.Println(string(data))
}
