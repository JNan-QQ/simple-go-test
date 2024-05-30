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
	"gitee.com/jn-qq/simple-go-test/logger"
	"gitee.com/jn-qq/simple-go-test/runner"
	"time"
)

func SuiteSetUp() {
	logger.INFO("进行包 cases 初始化操作")
	time.Sleep(5 * time.Second)
	logger.INFO("包 cases 初始化完成")
}

func SuiteTearDown() {
	logger.INFO("进行包 cases 清除操作")
	time.Sleep(5 * time.Second)
	logger.INFO("包 cases 清除完成")
}

// TestNum 实现接口 simple_go_test.TestInterface
type TestNum struct {
	runner.Test
}

// Init 测试用例构造
func (t *TestNum) Init() *runner.Test {
	// 设置 返回
	t.Name = "C0001"
	t.Tags = []string{"cases", "冒烟测试", "num"}

	return &t.Test
}

// SetUp 测试用例初始化
func (t *TestNum) SetUp() {
	logger.INFO("用例 C0001 初始化")
	time.Sleep(5 * time.Second)
}

// TearDown 测试用例清除
func (t *TestNum) TearDown() {
	logger.INFO("用例 C0001 清除")
	time.Sleep(5 * time.Second)
}

// TestStep 测试步骤
func (t *TestNum) TestStep() {
	logger.INFO("用例 C0001 运行")
	logger.STEP(1, "打开页面xxxx")
	time.Sleep(5 * time.Second)
	logger.CHECK_POINT("成功进入页面xxxx", true, false)
	logger.STEP(2, "关闭页面xxxx")
	time.Sleep(5 * time.Second)
}

// TestString 实现接口 simple_go_test.TestInterface
type TestString runner.Test

func (t *TestString) Init() *runner.Test {
	// 设置 返回
	t.Name = "C0002"
	t.Tags = []string{"cases", "冒烟测试", "string"}
	t.DDT = []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	return (*runner.Test)(t)
}

func (t *TestString) SetUp() {
	logger.INFO("用例 C0002 初始化")
	time.Sleep(5 * time.Second)
}

func (t *TestString) TearDown() {
	logger.INFO("用例 C0002 清除")
	time.Sleep(5 * time.Second)
}

func (t *TestString) TestStep() {
	para := t.Para
	time.Sleep(5 * time.Second)
	logger.CHECK_POINT("数据小于6", para.(int) < 6, true)
}
