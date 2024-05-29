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
	"gitee.com/jn-qq/simple-go-test/logger"
	"gitee.com/jn-qq/simple-go-test/runner"
	"time"
)

type TestHomeWork1 runner.Test

// Init 测试用例构造
func (t TestHomeWork1) Init() runner.Test {
	// 设置 返回
	t.Name = "C0006"
	t.Tags = []string{"homework"}

	return runner.Test(t)
}

// SetUp 测试用例初始化
func (t TestHomeWork1) SetUp() {}

// TearDown 测试用例清除
func (t TestHomeWork1) TearDown() {}

// TestStep 测试步骤
func (t TestHomeWork1) TestStep() {
	logger.INFO("运行测试用例:" + t.Name)
	time.Sleep(5 * time.Second)
}
