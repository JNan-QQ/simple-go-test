/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package runner

import (
	"fmt"
	"gitee.com/jn-qq/simple-go-test/config"
	"gitee.com/jn-qq/simple-go-test/logger"
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
func (t *TestPackage) SelectBy(by config.By, reg string) *TestPackage {
	if reg == "" {
		return t
	}

	switch by {
	case config.ByTagName, config.ByNotTagName:
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
			if (hs && by == config.ByNotTagName) || (!hs && by == config.ByTagName) {
				t.Tests = slices.Delete(t.Tests, i, i+1)
			}
		}

	case config.ByTestName:
		for i := len(t.Tests) - 1; i >= 0; i-- {
			ts := t.Tests[i].Init()
			if hs, _ := regexp.MatchString(reg, ts.Name); !hs {
				t.Tests = slices.Delete(t.Tests, i, i+1)
			}
		}

	case config.ByPackageName:
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

	// 判断是否有测试用例
	if (t.Child == nil || len(t.Child) == 0) && (t.Tests == nil || len(t.Tests) == 0) {
		return nil
	}

	return t
}

func (t *TestPackage) Num() int {
	num := len(t.Tests)
	for _, tp := range t.Child {
		num += tp.Num()
	}
	return num
}

// Run 运行测试
func (t *TestPackage) Run() logger.ReportHtml {

	report := logger.ReportHtml{Name: t.Name}
	logger.Reporthtml = &report
	// 记录开始时间
	report.Times = time.Now().Unix()

	fmt.Println("\n进入包：" + t.Name)

	// 包结构初始化
	if t.SuitSetUp != nil {
		// 记录运行进程
		report.Step = 0
		// 实例化消息对象
		report.Setup = new(logger.Steps)
		// 异步调用
		report.Setup.Result = syncRun(t.SuitSetUp)
		// 初始化失败,停止后续运行
		if report.Setup.Result != config.Success {
			// 记录初始化失败次数
			config.SetupFailNum += 1
			fmt.Printf("%s 包初始化结果：%s\n", t.Name, colorPrinter(report.Setup.Result))
			goto PTD
		}
		fmt.Printf("%s 包初始化结果：%s\n", t.Name, colorPrinter(report.Setup.Result))
	}

	// 遍历包下测试用例
	for _, test := range t.Tests {
		// 设置初始值
		ts := test.Init()

		fmt.Println("\n\n运行测试用例：" + ts.Name)

		// 实例化测试结果对象
		cases := &logger.TestCases{
			Name:      ts.Name,
			Setup:     new(logger.Steps),
			Teardown:  new(logger.Steps),
			StepCases: new(logger.Steps),
			Times:     time.Now().Format("2006-01-02 15:04:05"),
		}
		report.Cases = append(report.Cases, cases)

		// 测试用例初始化
		report.Step = 1
		cases.Setup.Result = syncRun(test.SetUp)
		if cases.Setup.Result != config.Success {
			config.SetupFailNum += 1
			cases.Result = config.Abort
			config.AbortNum += 1
			// 初始化失败跳过测试，执行清除
			goto CTD
		}

		// 运行测试步骤
		report.Step = 2
		if len(ts.DDT) > 0 {
			// 数据驱动不为空
			// 遍历数据驱动，取出数据
			for i, ddt := range ts.DDT {
				// 将数据赋值给 Para 便于 TestStep 中用户调用
				logger.INFO(fmt.Sprintf("\n运行子测试：%d", i))
				ts.Para = ddt
				if cases.StepCases.Result = syncRun(test.TestStep); cases.StepCases.Result != config.Success {
					break
				}
			}
		} else {
			cases.StepCases.Result = syncRun(test.TestStep)
		}

		// 判断整体运行结果
		if cases.StepCases.Result == config.Success {
			cases.Result = config.Success
			config.SuccessNum += 1
		} else if cases.StepCases.Result == config.Fail {
			cases.Result = config.Fail
			config.FailNum += 1
		} else {
			cases.Result = config.Abort
			config.AbortNum += 1
		}

	CTD:
		// 测试步骤清除
		report.Step = 3
		cases.Teardown.Result = syncRun(test.TearDown)
		if cases.Teardown.Result != config.Success {
			config.TearDownFailNum += 1
		}

		fmt.Printf("测试用例 %s 运行结果：%s\n", ts.Name, colorPrinter(cases.Result))

	}

	// 递归遍历子包
	for _, child := range t.Child {
		report.Child = append(report.Child, child.Run())
	}

PTD:
	// 包结构清除
	if t.SuitTearDown != nil {
		// 重新更新，避免 递归遍历子包 数据未更新
		logger.Reporthtml = &report
		// 记录运行进程
		report.Step = 4
		// 实例化消息对象
		report.Teardown = new(logger.Steps)
		// 异步调用
		report.Teardown.Result = syncRun(t.SuitTearDown)
		if report.Teardown.Result != config.Success {
			// 记录清除失败次数
			config.TearDownFailNum += 1
		}

		fmt.Printf("%s 包清除结果：%s\n", t.Name, colorPrinter(report.Teardown.Result))
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
			if err := recover(); err != nil {
				if err == "[CHECK_POINT FAIL]" {
					r = config.Fail
				} else {
					r = config.Abort
					logger.ErrorInfo(err.(error).Error())
				}
			} else {
				r = config.Success
			}
			wg.Done()
		}()
		f()
	}()
	// 协程等待
	wg.Wait()

	return
}

// 输出结果颜色
func colorPrinter(r int) string {
	var res string
	if r == config.Success {
		res = fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 32, "PASS")
	} else if r == config.Fail {
		res = fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 31, "FAIL")
	} else if r == config.Abort {
		res = fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 33, "ABORT")
	}
	return res
}
