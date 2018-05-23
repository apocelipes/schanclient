package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"schanclient/urls"
)

var (
	getDataInfo = regexp.MustCompile(`.+ \((.+(?:GB|MB|KB))\)`)
	getTotal    = regexp.MustCompile(`.+ \(流量：(.+(?:GB|MB|KB))\)`)
)

// 返回所有可用的套餐的信息
func GetService(data string) []*Service {
	res := make([]*Service, 0)

	table, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
	table.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		ser := new(Service)
		tds := s.Find("td")
		ser.Name = tds.Eq(0).Text()
		link, _ := tds.Eq(1).Find("a").Attr("href")
		ser.Link = urls.RootPath + link
		ser.Price = tds.Eq(1).Text()
		expir := tds.Eq(2).Find("span").Text()
		ser.Expires, _ = time.ParseInLocation("2006-01-02", expir, time.Local)
		ser.State = tds.Eq(3).Text()
		res = append(res, ser)
	})

	return res
}

// 获取套餐的详细使用信息
func GetSSRInfo(data string, ser *Service) *SSRInfo {
	res := NewSSRInfo(ser)

	panel, _ := goquery.NewDocumentFromReader(strings.NewReader(data))

	// 第一个table是端口和密码
	portAndPasswd := panel.Find("table").Eq(0)
	res.Port, _ = strconv.ParseInt(portAndPasswd.Find("tbody tr").Find("td").Eq(0).Text(), 10, 64)
	res.Passwd = portAndPasswd.Find("tbody tr").Find("td").Eq(1).Text()

	// 有两个header，第二个是套餐总量
	total := panel.Find("header").Eq(1).Text()
	res.TotalData = getTotal.FindStringSubmatch(total)[1]

	usage := panel.Find("#plugin-usage").Find("p")
	res.UsedData = getDataInfo.FindStringSubmatch(usage.Eq(0).Text())[1]
	res.Upload = getDataInfo.FindStringSubmatch(usage.Eq(1).Text())[1]
	res.Download = getDataInfo.FindStringSubmatch(usage.Eq(2).Text())[1]

	// 第2个table是节点信息表
	panel.Find("table").Eq(1).Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		node := new(SSRNode)
		tds := s.Children()

		node.NodeName = tds.Eq(0).Text()
		node.Type = tds.Eq(1).Text()
		node.IP = tds.Eq(2).Text()
		node.Crypto = tds.Eq(3).Text()
		node.Proto = tds.Eq(4).Text()
		node.Minx = tds.Eq(5).Text()
		node.Port = res.Port
		node.Passwd = res.Passwd

		res.Nodes = append(res.Nodes, node)
	})

	return res
}
