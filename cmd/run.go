/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package main

import (
	"flag"
	"fmt"
	"gitee.com/jn-qq/simple-go-test"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type _args struct {
	version        bool
	new            string
	caseDir        string
	autoOpenReport bool
	reportTitle    string
	urlPrefix      string
	test           string
	pkg            string
	tag            string
	tagNot         string
}

var args = new(_args)

func main() {
	flag.BoolVar(&args.version, "version", false, "build version")
	flag.StringVar(&simple_go_test.Lang, "lang", "zh", "language")
	flag.StringVar(&args.new, "new", "", "create new case project")
	flag.StringVar(&simple_go_test.CasesDir, "caseDir", "cases", "指定测试目录名")
	//flag.IntVar(&simple_go_test.Logger.level, "log-level", 1, "log level")
	//flag.BoolVar(&args.autoOpenReport, "auto-open-report", false, "auto open report")
	//flag.StringVar(&args.reportTitle, "report-title", "测试报告", "report title")
	//flag.StringVar(&args.urlPrefix, "url-prefix", "http://127.0.0.1", "url prefix")
	flag.StringVar(&args.test, "test", "", "用例名称过滤")
	flag.StringVar(&args.pkg, "pkg", "", "包名过滤")
	flag.StringVar(&args.tag, "tag", "", "tag过滤")
	flag.StringVar(&args.tagNot, "tagNot", "", "tag反向过滤")
	flag.Parse()

	// 返回版本
	if args.version {
		fmt.Println(simple_go_test.Version)
		return
	}

	// 设置语言
	if !slices.Contains([]string{"zh", "en"}, simple_go_test.Lang) {
		simple_go_test.Lang = "zh"
	}

	// 新建测试项目
	if args.new != "" {
		copyDemo(args.new)
		return
	}

	// 整理过滤条件
	if args.tag != "" {
		simple_go_test.FilterBy = simple_go_test.ByTagName
		simple_go_test.FilterValue = args.tag
	} else if args.tagNot != "" {
		simple_go_test.FilterBy = simple_go_test.ByNotTagName
		simple_go_test.FilterValue = args.tagNot
	} else if args.test != "" {
		simple_go_test.FilterBy = simple_go_test.ByTestName
		simple_go_test.FilterValue = args.test
	} else {
		simple_go_test.FilterBy = simple_go_test.ByPackageName
		simple_go_test.FilterValue = simple_go_test.CasesDir
	}

	fmt.Printf("%s\n * simple-go-test %s   https://gitee.com/JNan-QQ/simple-go-test *\n%s\n",
		strings.Repeat(" *", 34), simple_go_test.Version, strings.Repeat(" *", 34))

	fmt.Println("\n\n开始格式化测试用例组织关系...")
	astPack := simple_go_test.AstFind(simple_go_test.CasesDir).ToAst()
	fmt.Println("更新 main.go 文件")
	simple_go_test.WriteToMain(astPack)
	fmt.Println("格式化完成！")

	fmt.Println("开始运行 main.go 文件")
	_ = exec.Command("go", "run", "main.go").Run()
}

func copyDemo(p string) {
	fmt.Println("开始创建项目...")
	if _, err := os.Stat(p); err == nil {
		fmt.Println("项目已存在")
		return
	} else if os.IsNotExist(err) {
		panic(err)
	}

	// 创建父目录
	if err := os.MkdirAll(p, 0755); err != nil {
		panic(err)
	}

	// 复制 demo 项目
	for _, file := range []string{
		"main.go", "go.modc",
		"cases/cases.go",
		"cases/student/student.go",
		"cases/student/homework/homework.go",
		"cases/teacher/teacher.go",
	} {
		readFile, _ := simple_go_test.FS.ReadFile(filepath.Join("demo", file))
		if file == "go.modc" {
			compile, _ := regexp.Compile("mymodulenamereplace")
			readFile = compile.ReplaceAll(readFile, []byte(p))
			compile, _ = regexp.Compile("versionreplase")
			readFile = compile.ReplaceAll(readFile, []byte(simple_go_test.Version))
			file = filepath.Join(filepath.Dir(file), "go.mod")
		}
		_ = os.MkdirAll(filepath.Dir(filepath.Join(p, file)), 0755)
		_ = os.WriteFile(filepath.Join(p, file), readFile, 0644)
		fmt.Println(file, "复制完成")
	}

	fmt.Println("项目创建完成,请同步包 go mod tid")
}
