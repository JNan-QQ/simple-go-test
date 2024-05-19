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
	"embed"
	"os"
	"strings"
)

const (
	Version = "v1.0.0"
)

var (
	Lang         = "zh"
	CasesDir     = "cases"
	PackageNames []string
	FilterBy     By
	FilterValue  string
)

//go:embed demo/*
var FS embed.FS

func getRunPackage() string {
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
