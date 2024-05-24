/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package cases

import (
	"fmt"
	sgt "gitee.com/jn-qq/simple-go-test"
)

func SuiteSetUp() {
	fmt.Println("Package cases SetUp")
}

func SuiteTearDown() {
	fmt.Println("Package cases TearDown")
}

// TestNum 实现接口 simple_go_test.TestInterface
type TestNum sgt.Test

// Init 测试用例构造
func (t *TestNum) Init() *sgt.Test {
	// 设置 返回
	t.Name = "Cases001"
	t.Tags = []string{"cases", "冒烟测试", "num"}

	return (*sgt.Test)(t)
}

// SetUp 测试用例初始化
func (t *TestNum) SetUp() {
	fmt.Println("TestNum SetUp")
}

// TearDown 测试用例清除
func (t *TestNum) TearDown() {
	fmt.Println("TestNum TearDown")
}

// TestStep 测试步骤
func (t *TestNum) TestStep() {
	fmt.Println("TestNum TestStep")
}

// TestString 实现接口 simple_go_test.TestInterface
type TestString sgt.Test

func (t *TestString) Init() *sgt.Test {
	// 设置 返回
	t.Name = "Cases002"
	t.Tags = []string{"cases", "冒烟测试", "string"}
	t.DDT = []any{"1", "2", "3", "4"}

	return (*sgt.Test)(t)
}

func (t *TestString) SetUp() {
	fmt.Println("TestString SetUp")
}

func (t *TestString) TearDown() {
	fmt.Println("TestString TearDown")
}

func (t *TestString) TestStep() {
	para := t.Para
	fmt.Println("TestString TestStep By Para", para)
}
