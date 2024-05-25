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
	file := logger.Logger()
	defer func() {
		_ = file.Close()
	}()

	testTree = *testTree.SelectBy(selectBy, selectValue)

	config.GSTORE.SetItem("caseNum", testTree.Num())

	report := testTree.Run()

	report.Save()
}
