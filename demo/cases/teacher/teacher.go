/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package teacher

import (
	"fmt"
	"gitee.com/jn-qq/simple-go-test"
)

type TestHomeWork1 simple_go_test.Test

// Init 测试用例构造
func (t *TestHomeWork1) Init() simple_go_test.Test {
	// 设置 返回
	t.Name = "Cases006"
	t.Tags = []string{"homework"}

	return simple_go_test.Test(*t)
}

// SetUp 测试用例初始化
func (t *TestHomeWork1) SetUp() {}

// TearDown 测试用例清除
func (t *TestHomeWork1) TearDown() {}

// TestStep 测试步骤
func (t *TestHomeWork1) TestStep() {
	fmt.Println("TestHomeWork1 TestStep")
}
