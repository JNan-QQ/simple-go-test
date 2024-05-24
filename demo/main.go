package main

import (
	sgt "gitee.com/jn-qq/simple-go-test"
)

var testTree = sgt.TestPackage{}

var selectBy sgt.By = 0
var selectValue = ""

func main() {
	// 自定义日志对象
	file := sgt.Logger()
	defer func() {
		_ = file.Close()
	}()

	// 过滤
	testTree = *testTree.SelectBy(selectBy, selectValue)
	// 统计测试用例个数
	sgt.GSTORE.SetItem("caseNum", testTree.Num())

	//运行
	report := testTree.Run()

	// 保存测试报告
	report.Save()
}
