/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package simple_go_test

import (
	"regexp"
	"slices"
	"sync"
	"time"
)

type TestInterface interface {
	Init() *Test
	SetUp()
	TearDown()
	TestStep()
}

type Test struct {
	Name string
	Tags []string
	DDT  []any
	Para any
}

// TestPackage 测试结构整合
type TestPackage struct {
	Name         string
	Tests        []TestInterface
	SuitSetUp    func()
	SuitTearDown func()
	Child        []TestPackage
}

// SelectBy 根据正则表达式过滤测试用例
func (t *TestPackage) SelectBy(by By, reg string) *TestPackage {
	if reg == "" {
		return t
	}
	switch by {
	case ByTagName, ByNotTagName:
		// 遍历测试用例tag
		for i := len(t.Tests) - 1; i >= 0; i-- {
			ts := t.Tests[i].Init()
			var hs bool
			for _, tag := range ts.Tags {
				hs, _ = regexp.MatchString(reg, tag)
				if hs {
					break
				}
			}
			if (hs && by == ByNotTagName) || (!hs && by == ByTagName) {
				t.Tests = slices.Delete(t.Tests, i, i+1)
			}
		}

	case ByTestName:
		for i := len(t.Tests) - 1; i >= 0; i-- {
			ts := t.Tests[i].Init()
			hs, _ := regexp.MatchString(reg, ts.Name)
			if !hs {
				t.Tests = slices.Delete(t.Tests, i, i+1)
			}
		}

	case ByPackageName:
		hs, _ := regexp.MatchString(reg, t.Name)
		if !hs {
			t.Tests = nil
		} else {
			return t
		}

	default:
		return t
	}

	// 递归遍历子包
	for i := len(t.Child) - 1; i >= 0; i-- {
		child := t.Child[i].SelectBy(by, reg)
		if child != nil {
			t.Child[i] = *child
		} else {
			t.Child = slices.Delete(t.Child, i, i+1)
		}
	}

	if (t.Child == nil || len(t.Child) == 0) && (len(t.Tests) == 0 || t.Tests == nil) {
		return nil
	} else {
		return t
	}
}

func (t *TestPackage) Num() int {
	allNum = len(t.Tests)
	for _, tp := range t.Child {
		allNum += tp.Num()
	}
	return allNum
}

// Run 运行测试
func (t *TestPackage) Run() ReportHtml {

	report := ReportHtml{Name: t.Name}
	reportHtml = &report
	// 记录开始时间
	report.Times = time.Now().Unix()

	// 包结构初始化
	if t.SuitSetUp != nil {
		// 记录运行进程
		report.step = 0
		// 实例化消息对象
		report.Setup = new(Steps)
		// time
		report.Setup.Times = time.Now().Format("2006-01-02 15:04:05")
		// 异步调用
		report.Setup.Result = syncRun(t.SuitSetUp)
		// 初始化失败,停止后续运行
		if report.Setup.Result != success {
			// 记录初始化失败次数
			setupFailNum += 1
			return report
		}
	}
	// 包结构清除
	defer func() {
		if t.SuitTearDown != nil {
			// 记录运行进程
			report.step = 4
			// 实例化消息对象
			report.Teardown = new(Steps)
			report.Teardown.Times = time.Now().Format("2006-01-02 15:04:05")
			// 异步调用
			report.Teardown.Result = syncRun(t.SuitTearDown)
			if report.Teardown.Result != success {
				// 记录清除失败次数
				tearDownFailNum += 1
			}
		}
	}()

	// 遍历包下测试用例
	for _, test := range t.Tests {
		// 设置初始值
		ts := test.Init()

		// 实例化测试结果对象
		cases := &TestCases{Name: ts.Name}
		report.Cases = append(report.Cases, *cases)
		cases.Times = time.Now().Format("2006-01-02 15:04:05")

		// 测试用例初始化
		report.step = 1
		cases.Setup.Result = syncRun(test.TestStep)
		if cases.Setup.Result != success {
			setupFailNum += 1
			// 初始化失败跳过测试，执行清除
			goto casesTeardown
		}

		// 运行测试步骤
		report.step = 2
		if len(ts.DDT) > 0 {
			// 数据驱动不为空
			var hasFailed bool
			// 遍历数据驱动，取出数据
			for _, ddt := range ts.DDT {
				report.Cases[len(report.Cases)-1].Ddts = append(report.Cases[len(report.Cases)-1].Ddts, Steps{})
				// 将数据赋值给 Para 便于 TestStep 中用户调用
				ts.Para = ddt
				if syncRun(test.TestStep) != success {
					hasFailed = true
				}
			}
			// 判断整体运行结果
			if hasFailed {
				failNum += 1
				cases.Result = fail
			} else {
				successNum += 1
				cases.Result = success
			}

		} else {
			cases.Result = syncRun(test.TestStep)
			if cases.Result == success {
				successNum += 1
			} else if cases.Result == fail {
				failNum += 1
			} else {
				abortNum += 1
			}
		}

	casesTeardown:
		// 测试步骤清除
		report.step = 3
		cases.Teardown.Result = syncRun(test.TearDown)
		if cases.Teardown.Result != success {
			tearDownFailNum += 1
		}
	}

	// 递归遍历子包
	for _, child := range t.Child {
		report.Child = append(report.Child, child.Run())
	}

	// 记录结束时间
	report.Times = time.Now().Unix() - report.Times

	return report
}

// 协程运行函数，捕获异常
func syncRun(f func()) (r int) {
	// 协程对象
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// 捕获异常
		defer func() {
			if err := recover(); err == nil {
				r = success
			} else if err == "" {
				r = abort
			} else {
				r = fail
			}
			wg.Done()
		}()
		f()
	}()
	// 协程等待
	wg.Wait()

	return
}
