package main

import (
	"gitee.com/jn-qq/simple-go-test/config"
	"gitee.com/jn-qq/simple-go-test/logger"
	"gitee.com/jn-qq/simple-go-test/runner"
)

var testTree = runner.TestPackage{}

var selectBy config.By = 0
var selectValue = ""

func main() {
	// 自定义日志对象
	file := logger.Logger()
	defer func() {
		_ = file.Close()
	}()

	// 过滤
	testTree = *testTree.SelectBy(selectBy, selectValue)
	// 统计测试用例个数
	config.GSTORE.SetItem("caseNum", testTree.Num())

	//运行
	report := testTree.Run()

	// 保存测试报告
	report.Save()
}
