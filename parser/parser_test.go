package parser

import "testing"

func TestGetTotal(t *testing.T) {
	data := `使用报表 (流量：50GB)`
	res := getTotal.FindStringSubmatch(data)
	if res[1] != "50GB" {
		t.Error("regexp getTotal has some problem.")
	}
}

func TestGetDataInfo(t *testing.T) {
	data1 := `已使用 (16.14GB)`
	data2 := `上传 (14.66MB)`
	data3 := `下载 (16.12GB)`

	if getDataInfo.FindStringSubmatch(data1)[1] != "16.14GB" {
		t.Error("regexp getDataInfo has some problem on getting used.")
	}

	if getDataInfo.FindStringSubmatch(data2)[1] != "14.66MB" {
		t.Error("regexp getDataInfo has some problem on getting upload.")
	}

	if getDataInfo.FindStringSubmatch(data3)[1] != "16.12GB" {
		t.Error("regexp getDataInfo has some problem on getting download.")
	}
}
