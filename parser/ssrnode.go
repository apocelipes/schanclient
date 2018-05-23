package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// ssr节点信息
type SSRNode struct {
	// 节点名字
	NodeName string `json:"-"`
	// 节点类型
	Type string `json:"-"`
	// 节点IP地址
	IP     string `json:"server"`
	Port   int64  `json:"server_port"`
	Passwd string `json:"password"`
	// 加密算法
	Crypto string `json:"method"`
	// 连接协议
	Proto string `json:"protocol"`
	// 混淆算法
	Minx string `json:"obfs"`
}

func (s *SSRNode) Store(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func (s *SSRNode) String() string {
	format := "Name:%-10s IP:%-15s Port:%-5v Pswd:%-15s Crypto:%-11s Protocol:%-7s Obfs:%-6s"
	res := fmt.Sprintf(format, s.NodeName, s.IP, s.Port, s.Passwd, s.Crypto, s.Proto, s.Minx)
	return res
}
