package main

import (
	"fmt"
	"gitee.com/jn-qq/simple-go-test"
	"gitee.com/jn-qq/simple-go-test/demo/cases"
	"gitee.com/jn-qq/simple-go-test/demo/cases/student"
	"gitee.com/jn-qq/simple-go-test/demo/cases/student/homework"
	"io/fs"
)

var testTree = simple_go_test.TestPackage{
	Name:         "cases",
	Tests:        []simple_go_test.TestInterface{new(cases.TestNum), new(cases.TestString)},
	SuitSetUp:    cases.SuiteSetUp,
	SuitTearDown: cases.SuiteTearDown,
	Child: []simple_go_test.TestPackage{
		simple_go_test.TestPackage{
			Name:      "student",
			Tests:     []simple_go_test.TestInterface{new(student.TestNum1), new(student.TestString1)},
			SuitSetUp: student.SuiteSetUp,
			Child: []simple_go_test.TestPackage{
				simple_go_test.TestPackage{
					Name:  "homework",
					Tests: []simple_go_test.TestInterface{new(homework.TestHomeWork)},
				},
			},
		},
	},
}

func main() {
	//fmt.Printf("%+v", testTree.SelectBy(simple_go_test.ByTagName, "冒烟测试"))
	dir, err := fs.ReadDir(simple_go_test.FS, "demo")
	if err != nil {
		return
	}
	fmt.Println(dir)
}
