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
	"reflect"
	"regexp"
	"slices"
)

type TestInterface interface {
	Init() Test
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
	switch by {
	case ByTagName, ByNotTagName:
		// 遍历测试用例tag
		for i := len(t.Tests) - 1; i >= 0; i-- {
			ts := t.Tests[i]
			ts.Init()
			v := reflect.ValueOf(ts).Elem()
			tags := v.FieldByName("Tags")

			var hs bool
			for j := 0; j < tags.Len(); j++ {
				hs, _ = regexp.MatchString(reg, tags.Index(j).String())
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
			ts := t.Tests[i]
			ts.Init()
			hs, _ := regexp.MatchString(reg, reflect.ValueOf(ts).Elem().FieldByName("Name").String())
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

// Run 运行测试
func (t *TestPackage) Run() {
	if t.SuitSetUp != nil {
		t.SuitSetUp()
	}
	defer func() {
		if t.SuitTearDown != nil {
			t.SuitTearDown()
		}
	}()
	for _, test := range t.Tests {
		test.Init()
		test.SetUp()
		v := reflect.ValueOf(test).Elem()
		ddt := v.FieldByName("DDT")
		if !ddt.IsValid() {
			for i := 0; i < ddt.Len(); i++ {
				v.FieldByName("Para").Set(ddt.Index(i))
				test.TestStep()
			}
		} else {
			test.TestStep()
		}
		test.TearDown()
	}

	for _, child := range t.Child {
		child.Run()
	}
}
