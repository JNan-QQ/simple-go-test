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
	"bufio"
	"cmp"
	"embed"
	"flag"
	"fmt"
	"gitee.com/jn-qq/simple-go-test/ast"
	"gitee.com/jn-qq/simple-go-test/config"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
)

type _args struct {
	version        bool
	new            string
	autoOpenReport bool
	test           string
	pkg            string
	tag            string
	tagNot         string
	run            bool
	collect        bool
}

var args = new(_args)

//go:embed statics/demo
var FS embed.FS

func main() {
	flag.BoolVar(&args.version, "version", false, "版本")
	//flag.StringVar(&config.Lang, "lang", "zh", "language")
	flag.StringVar(&args.new, "new", "", "新建一个测试项目")
	flag.BoolVar(&args.collect, "collect", false, "收集测试用例")
	flag.BoolVar(&args.autoOpenReport, "auto-open-report", true, "自动打开报告")
	flag.StringVar(&config.ReportName, "report-title", "测试报告", "报告标题")
	flag.StringVar(&args.test, "test", "", "用例名称过滤,支持正则表达式")
	flag.StringVar(&args.pkg, "pkg", "", "包名过滤,支持正则表达式")
	flag.StringVar(&args.tag, "tag", "", "tag过滤,支持正则表达式")
	flag.StringVar(&args.tagNot, "tagNot", "", "tag反向过滤,支持正则表达式")
	flag.BoolVar(&args.run, "run", false, "运行")
	flag.Parse()

	// 返回版本
	if args.version {
		fmt.Println(config.Version)
		return
	}

	// 设置语言
	//if !slices.Contains([]string{"zh", "en"}, config.Lang) {
	//	config.Lang = "zh"
	//}

	// 新建测试项目
	if args.new != "" {
		copyDemo(args.new)
		return
	}

	if args.collect {
		fmt.Println("\n开始格式化测试用例组织关系...")
		astPack := ast.Find(config.CasesDir).ToAst()
		fmt.Println("更新 main.go 文件")
		ast.WriteToMain(astPack)
		fmt.Println("格式化完成！")
	}

	// 整理过滤条件
	if args.tag != "" {
		ast.WriteToMain(nil, int(config.ByTagName), args.tag)
		fmt.Println("过滤设置成功")
	} else if args.tagNot != "" {
		ast.WriteToMain(nil, int(config.ByNotTagName), args.tagNot)
		fmt.Println("过滤设置成功")
	} else if args.test != "" {
		ast.WriteToMain(nil, int(config.ByTestName), args.test)
		fmt.Println("过滤设置成功")
	} else if args.pkg != "" {
		ast.WriteToMain(nil, int(config.ByPackageName), args.pkg)
		fmt.Println("过滤设置成功")
	}

	if args.run {
		fmt.Printf("%s\n * simple-go-test %s   https://gitee.com/JNan-QQ/simple-go-test *\n%s\n",
			strings.Repeat(" *", 34), config.Version, strings.Repeat(" *", 34))

		fmt.Println("开始运行 main.go 文件")
		execCommand()
	}

	if args.autoOpenReport {
		var cmd string
		var cArgs []string
		switch runtime.GOOS {
		case "linux":
			cmd = "xdg-open"
		case "windows":
			cmd = "cmd"
			cArgs = []string{"/c", "start"}
		case "darwin":
			cmd = "open"
		default:
			panic("未知运行平台，无法自动打开！")
		}
		// 获取文件路径
		glob, _ := filepath.Glob("logs/*.html")
		if len(glob) == 0 {
			panic("未生成报告")
		}
		slices.SortFunc(glob, func(a, b string) int {
			return cmp.Compare(b, a)
		})
		rp, _ := filepath.Abs(glob[0])
		cArgs = append(cArgs, rp)
		_ = exec.Command(cmd, cArgs...).Start()
	}
}

// 执行main文件，实时获取输出
func execCommand() {
	cmd := exec.Command("go", "run", "main.go")

	// 结果输出到管道
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	printer := func(rb io.Reader, w *sync.WaitGroup) {
		defer w.Done()
		reader := bufio.NewReader(rb)
		for {
			read, _, err := reader.ReadLine()
			if err != nil || err == io.EOF {
				return
			}
			fmt.Println(string(read))
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go printer(stdout, &wg)
	go printer(stderr, &wg)

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

// 复制示例项目
func copyDemo(p string) {
	fmt.Println("开始创建项目...")
	abs, _ := filepath.Abs(p)
	if _, err := os.Stat(abs); err == nil {
		fmt.Println("项目已存在")
		return
	}

	// 复制 demo 项目
	if err := fs.WalkDir(FS, "statics/demo", func(path string, d fs.DirEntry, err error) error {
		_path := strings.Replace(path, "statics/demo", p, 1)
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
				compile, _ = regexp.Compile("sgtVersion")
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
