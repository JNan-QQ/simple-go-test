package main

import (
	"fmt"
	"gitee.com/jn-qq/simple-go-test/config"
	"gitee.com/jn-qq/simple-go-test/logger"
	"gitee.com/jn-qq/simple-go-test/runner"
)

var testTree = runner.TestPackage{}

var selectBy config.By = 0
var selectValue = ""

func main() {
	logger.Logger()

	fmt.Println("开始过滤测试用例...")
	if tree := testTree.SelectBy(selectBy, selectValue); tree != nil {
		testTree = *tree
		config.AllNum = testTree.Num()
		fmt.Println("待执行测试用例数量：", config.AllNum)
	} else {
		fmt.Println("未有符合条件的测试用例")
		return
	}

	report := testTree.Run()

	report.Save()
}
