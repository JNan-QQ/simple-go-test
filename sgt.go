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
	"embed"
	"flag"
	"fmt"
	"gitee.com/jn-qq/simple-go-test/ast"
	"gitee.com/jn-qq/simple-go-test/config"
	"io/fs"
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
	run            bool
}

var args = new(_args)

//go:embed statics/demo/*
var FS embed.FS

func main() {
	flag.BoolVar(&args.version, "version", false, "build version")
	flag.StringVar(&config.Lang, "lang", "zh", "language")
	flag.StringVar(&args.new, "new", "", "create new case project")
	flag.StringVar(&config.CasesDir, "caseDir", "cases", "指定测试目录名")
	//flag.IntVar(&config.Logger.level, "log-level", 1, "log level")
	//flag.BoolVar(&args.autoOpenReport, "auto-open-report", false, "auto open report")
	//flag.StringVar(&args.reportTitle, "report-title", "测试报告", "report title")
	//flag.StringVar(&args.urlPrefix, "url-prefix", "http://127.0.0.1", "url prefix")
	flag.StringVar(&args.test, "test", "", "用例名称过滤")
	flag.StringVar(&args.pkg, "pkg", "", "包名过滤")
	flag.StringVar(&args.tag, "tag", "", "tag过滤")
	flag.StringVar(&args.tagNot, "tagNot", "", "tag反向过滤")
	flag.BoolVar(&args.run, "run", false, "运行")
	flag.Parse()

	// 返回版本
	if args.version {
		fmt.Println(config.Version)
		return
	}

	// 设置语言
	if !slices.Contains([]string{"zh", "en"}, config.Lang) {
		config.Lang = "zh"
	}

	// 新建测试项目
	if args.new != "" {
		copyDemo(args.new)
		return
	}

	// 整理过滤条件
	if args.tag != "" {
		config.FilterBy = config.ByTagName
		config.FilterValue = args.tag
	} else if args.tagNot != "" {
		config.FilterBy = config.ByNotTagName
		config.FilterValue = args.tagNot
	} else if args.test != "" {
		config.FilterBy = config.ByTestName
		config.FilterValue = args.test
	} else if args.pkg != "" {
		config.FilterBy = config.ByPackageName
		config.FilterValue = config.CasesDir
	} else {
		config.FilterBy = config.ByTagName
		config.FilterValue = ""
	}

	fmt.Println("\n\n开始格式化测试用例组织关系...")
	astPack := ast.Find(config.CasesDir).ToAst()
	fmt.Println("更新 main.go 文件")
	ast.WriteToMain(astPack)
	fmt.Println("格式化完成！")

	if args.run {
		fmt.Printf("%s\n * simple-go-test %s   https://gitee.com/JNan-QQ/simple-go-test *\n%s\n",
			strings.Repeat(" *", 34), config.Version, strings.Repeat(" *", 34))

		fmt.Println("开始运行 main.go 文件")
		_ = exec.Command("go", "run", "main.go").Run()
	}
}

func copyDemo(p string) {
	fmt.Println("开始创建项目...")
	abs, _ := filepath.Abs(p)
	if _, err := os.Stat(abs); err == nil {
		fmt.Println("项目已存在")
		return
	}

	// 复制 demo 项目
	if err := fs.WalkDir(FS, "demo", func(path string, d fs.DirEntry, err error) error {
		_path := strings.Replace(path, "demo", p, 1)
		if d.IsDir() {
			// 创建目录
			if err := os.MkdirAll(_path, 0755); err != nil {
				return err
			}
		} else {
			readFile, err := fs.ReadFile(FS, path)
			if err != nil {
				return err
			}
			if d.Name() == "go.x" || d.Name() == "main.go" {
				compile, _ := regexp.Compile("packname")
				readFile = compile.ReplaceAll(readFile, []byte(p))
				compile, _ = regexp.Compile("configVersion")
				readFile = compile.ReplaceAll(readFile, []byte(config.Version))
				_path = strings.Replace(_path, "go.x", "go.mod", 1)
			}
			if err = os.WriteFile(_path, readFile, 0644); err != nil {
				return err
			}
		}
		fmt.Println(path, "复制完成")
		return nil
	}); err != nil {
		return
	}

	fmt.Println("项目创建完成,请同步包 go mod tid")
}
