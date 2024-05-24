/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package student

import (
	"fmt"
	sgt "gitee.com/jn-qq/simple-go-test"
)

func SuiteSetUp() {
	fmt.Println("Package student SetUp")
}

// TestNum1 实现接口 simple_go_test.TestInterface
type TestNum1 sgt.Test

// Init 测试用例构造
func (t *TestNum1) Init() *sgt.Test {
	// 设置 返回
	t.Name = "Cases003"
	t.Tags = []string{"student", "冒烟测试", "num"}

	return (*sgt.Test)(t)
}

// SetUp 测试用例初始化
func (t *TestNum1) SetUp() {
	fmt.Println("TestNum1 SetUp")
}

// TearDown 测试用例清除
func (t *TestNum1) TearDown() {
	fmt.Println("TestNum1 TearDown")
}

// TestStep 测试步骤
func (t *TestNum1) TestStep() {
	fmt.Println("TestNum1 TestStep")
}

// TestString1 实现接口 simple_go_test.TestInterface
type TestString1 sgt.Test

func (t *TestString1) Init() *sgt.Test {
	// 设置 返回
	t.Name = "Cases004"
	t.Tags = []string{"student", "冒烟测试", "string"}
	t.DDT = []any{"1", "2", "3", "4"}

	return (*sgt.Test)(t)
}

func (t *TestString1) SetUp() {
	fmt.Println("TestString1 SetUp")
}

func (t *TestString1) TearDown() {
	fmt.Println("TestString1 TearDown")
}

func (t *TestString1) TestStep() {
	para := t.Para
	fmt.Println("TestString1 TestStep By Para", para)
}
