package parser

import (
	"fmt"
	"strconv"
)

// ssr套餐信息
type SSRInfo struct {
	*Service
	Port   int64
	Passwd string
	// 可用数据总量
	TotalData string
	// 已用数据总量
	UsedData string
	// 下载用量
	Download string
	// 上传用量
	Upload string
	// 可用节点信息
	Nodes []*SSRNode
}

func NewSSRInfo(ser *Service) *SSRInfo {
	s := new(SSRInfo)
	s.Service = ser
	s.Nodes = make([]*SSRNode, 0)
	return s
}

func (s *SSRInfo) GetNodeByName(name string) *SSRNode {
	for _, node := range s.Nodes {
		if name == node.NodeName {
			return node
		}
	}

	return nil
}

func (s *SSRInfo) UsedInfo() string {
	res := ""
	res += fmt.Sprintf("服务名称：%v\n", s.Name)
	res += fmt.Sprintf("数据总量：%v\n", s.TotalData)
	res += fmt.Sprintf("已用数据：%v\n", s.UsedData)
	res += fmt.Sprintf("已下载：%v\n", s.Download)
	res += fmt.Sprintf("已上传：%v\n", s.Upload)

	return res
}

func (s *SSRInfo) String() string {
	res := "服务套餐信息\n"
	res += fmt.Sprintf("服务名称：%v\n", s.Name)
	res += fmt.Sprintf("服务费用：%v\n", s.Price)
	res += fmt.Sprintf("服务状态：%v\n", s.State)
	res += fmt.Sprintf("到期时间：%v\n", s.Expires.Format("2006-01-02"))
	res += fmt.Sprintf("数据总量：%v\n", s.TotalData)
	res += fmt.Sprintf("已用数据：%v\n", s.UsedData)
	res += fmt.Sprintf("已下载：%v\n", s.Download)
	res += fmt.Sprintf("已上传：%v\n", s.Upload)
	res += "节点信息：\n"

	for i, node := range s.Nodes {
		res += fmt.Sprintf("节点"+strconv.Itoa(i+1)+":\n%v\n", node)
	}

	return res
}
