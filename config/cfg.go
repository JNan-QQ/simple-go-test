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
	"os"
	"strings"
)

const (
	Version = "v0.0.1"
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

func (g _GlobalStore) SetItem(key string, value interface{}) {
	if g == nil {
		g = make(_GlobalStore)
	}
	g[key] = value
}

func (g _GlobalStore) GetItem(key string) interface{} {
	if value, ok := g[key]; ok {
		return value
	} else {
		return nil
	}
}

var GSTORE = new(_GlobalStore)
