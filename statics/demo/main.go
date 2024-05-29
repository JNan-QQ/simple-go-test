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
	logger.Logger()

	testTree = *testTree.SelectBy(selectBy, selectValue)

	config.AllNum = testTree.Num()

	report := testTree.Run()

	report.Save()
}
