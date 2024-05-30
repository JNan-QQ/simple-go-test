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
	"gitee.com/jn-qq/simple-go-test/logger"
	"gitee.com/jn-qq/simple-go-test/runner"
	"time"
)

func SuiteSetUp() {
	logger.INFO("进行包 student 初始化操作")
	time.Sleep(5 * time.Second)
	logger.INFO("包 student 初始化完成")
}

// TestNum1 实现接口 simple_go_test.TestInterface
type TestNum1 runner.Test

// Init 测试用例构造
func (t *TestNum1) Init() *runner.Test {
	// 设置 返回
	t.Name = "C0003"
	t.Tags = []string{"student", "冒烟测试", "num"}

	return (*runner.Test)(t)
}

// SetUp 测试用例初始化
func (t *TestNum1) SetUp() {
	logger.INFO("用例 C0003 初始化")
	time.Sleep(5 * time.Second)
}

// TearDown 测试用例清除
func (t *TestNum1) TearDown() {
	logger.INFO("用例 C0003 清除")
	time.Sleep(5 * time.Second)
}

// TestStep 测试步骤
func (t *TestNum1) TestStep() {
	logger.INFO("用例 C0003 运行")
	logger.STEP(1, "打开页面xxxx")
	time.Sleep(5 * time.Second)
	logger.CHECK_POINT("成功进入页面xxxx", false, false)
	logger.STEP(2, "关闭页面xxxx")
	time.Sleep(5 * time.Second)
}

// TestString1 实现接口 simple_go_test.TestInterface
type TestString1 struct {
	runner.Test
}

func (t *TestString1) Init() *runner.Test {
	// 设置 返回
	t.Name = "C0004"
	t.Tags = []string{"cases", "冒烟测试", "string"}
	t.DDT = []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	return &t.Test
}

func (t *TestString1) SetUp() {
	logger.INFO("用例 C0004 初始化")
	time.Sleep(5 * time.Second)
}

func (t *TestString1) TearDown() {
	logger.INFO("用例 C0004 清除")
	time.Sleep(5 * time.Second)
}

func (t *TestString1) TestStep() {
	para := t.Para
	time.Sleep(5 * time.Second)
	logger.CHECK_POINT("数据小于11", para.(int) < 11, true)
}
