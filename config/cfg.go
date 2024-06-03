/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package config

import (
	"math"
	"os"
	"strings"
)

const (
	Version = "v0.0.2"
)

var (
	Lang         = "zh"
	CasesDir     = "cases"
	PackageNames []string
	ReportName   string = "测试报告"
)

var (
	SuccessNum      int
	AbortNum        int
	FailNum         int
	AllNum          int
	SetupFailNum    int
	TearDownFailNum int
)

func GetRunPackage() string {
	file, _ := os.ReadFile("./go.mod")
	f1 := strings.Split(string(file), "\n")[0]
	return strings.Split(f1, " ")[1]
}

type By int

const (
	ByTagName By = iota
	ByNotTagName
	ByTestName
	ByPackageName
)

const (
	Fail int = iota
	Success
	Abort
)

type _GlobalStore map[string]interface{}

// SetItem 设置值
func (g _GlobalStore) SetItem(key string, value interface{}) {
	if g == nil {
		g = make(_GlobalStore)
	}
	g[key] = value
}

// GetItem 获取值，不存在返回 nil
func (g _GlobalStore) GetItem(key string) interface{} {
	if value, ok := g[key]; ok {
		return value
	} else {
		return nil
	}
}

// GetString 获取字符串, 不存在返回 ""
func (g _GlobalStore) GetString(key string) string {
	if value, ok := g[key]; ok {
		return value.(string)
	} else {
		return ""
	}
}

// GetInt 获取整数, 不存在返回 无穷小
func (g _GlobalStore) GetInt(key string) int {
	if value, ok := g[key]; ok {
		return value.(int)
	} else {
		return math.MinInt
	}
}

// GetFloat 获取浮点数，不存在返回 math.NaN
func (g _GlobalStore) GetFloat(key string) float64 {
	if value, ok := g[key]; ok {
		return value.(float64)
	} else {
		return math.NaN()
	}
}

// GetBool 获取布尔值
func (g _GlobalStore) GetBool(key string) bool {
	if value, ok := g[key]; ok {
		return value.(bool)
	} else {
		return false
	}
}

var GSTORE = new(_GlobalStore)

var (
	CasesEnd func(testName string, testResult int)
	TestEnd  func(allNum, failNum, successNum, abortNum, setupFail, teardownFail int)
)
