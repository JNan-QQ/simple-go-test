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
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Logger 自定义log对象，写入文件
func Logger() *os.File {
	_ = os.MkdirAll("logs", os.ModePerm)
	// 日志备份
	func() {
		if _, err := os.Stat("logs/sgt.log"); err == nil {
			return
		}
		_ = os.Remove("logs/sgt.log.back1")
		_ = os.Rename("logs/sgt.log.back", "logs/sgt.log.back1")
		_ = os.Rename("logs/sgt.log", "logs/sgt.log.back")
	}()
	file, err := os.OpenFile("logs/run.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	l := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, file), nil))
	slog.SetDefault(l)

	return file
}

var reportHtml *ReportHtml

// INFO 提示信息，同时输出到终端和文件
func INFO(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	slog.Info(msg)
	addHtml(fmt.Sprintf("<p class=\"s%d info\">%s</p>", reportHtml.step, msg))
}

// STEP 步骤信息
func STEP(step int, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	slog.Info(fmt.Sprintf("[第%02d步] %s", step, msg))
	addHtml(fmt.Sprintf("<p class=\"s%d step\"><span>[第%02d步]</span> %s</p>", reportHtml.step, step, msg))
}

// CHECK_POINT 检查点，如果allowRun!=nil,仅记录，继续运行
func CHECK_POINT(msg string, check bool, allowRun ...bool) {
	var htmlMsg string
	if check {
		slog.Info(fmt.Sprintf("[CHECK_POINT PASS] %s", msg))
		htmlMsg = fmt.Sprintf("<p class=\"s%d check-point success\"><span>[CHECK_POINT PASS]</span> %s</p>", reportHtml.step, msg)
	} else {
		slog.Info(fmt.Sprintf("[CHECK_POINT FAIL] %s", msg))
		htmlMsg = fmt.Sprintf("<p class=\"s%d check-point fail\"><span>[CHECK_POINT FAIL]</span> %s</p>", reportHtml.step, msg)
	}
	addHtml(htmlMsg)
	if !check && allowRun == nil {
		panic("[CHECK_POINT FAIL]")
	}
}

func addHtml(msg string) {
	switch reportHtml.step {
	case 0:
		if reportHtml.Setup == nil {
			reportHtml.Setup = new(Steps)
		}
		reportHtml.Setup.addStep(msg)
	case 1:
		reportHtml.Cases[len(reportHtml.Cases)-1].Setup.addStep(msg)
	case 2:
		ddts := reportHtml.Cases[len(reportHtml.Cases)-1].Ddts
		if ddts == nil {
			reportHtml.Cases[len(reportHtml.Cases)-1].Ddts = append(
				reportHtml.Cases[len(reportHtml.Cases)-1].Ddts,
				Steps{Html: template.HTML(msg)},
			)
		} else {
			reportHtml.Cases[len(reportHtml.Cases)-1].Ddts[len(ddts)-1].addStep(msg)
		}
	case 3:
		reportHtml.Cases[len(reportHtml.Cases)-1].Setup.addStep(msg)
	case 4:
		if reportHtml.Teardown == nil {
			reportHtml.Teardown = new(Steps)
		}
		reportHtml.Teardown.addStep(msg)
	}
}

type Steps struct {
	Result int
	Html   template.HTML
	Times  string
}

func (s Steps) addStep(h string) {
	s.Html += template.HTML(h)
}

func (s Steps) setResult(r int) {
	s.Result = r
}

type ReportHtml struct {
	Name     string
	Setup    *Steps
	Teardown *Steps
	Cases    []TestCases
	Child    []ReportHtml
	step     int
	Times    int64
}

func (r ReportHtml) Save() {
	type html struct {
		ReportName      string
		AllNum          int
		SetupFailNum    int
		TeardownFailNum int
		SpendTime       int64
		FailNum         int
		SuccessNum      int
		AbortNum        int
		Pack            ReportHtml
	}
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	tmpl, err := template.ParseFiles(
		filepath.Join(exPath, "html", "template.tmpl"),
		filepath.Join(exPath, "html", "js.tmpl"),
		filepath.Join(exPath, "html", "cases.tmpl"),
		filepath.Join(exPath, "html", "steps_result.tmpl"),
		filepath.Join(exPath, "html", "my.css"),
	)
	if err != nil {
		panic(err)
	}

	file, _ := os.Create(filepath.Join("logs", time.Now().Format("20060102-150405")) + ".html")
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	err = tmpl.Execute(file, html{
		ReportName:      ReportName,
		AllNum:          allNum,
		SetupFailNum:    setupFailNum,
		TeardownFailNum: tearDownFailNum,
		SpendTime:       r.Times,
		FailNum:         failNum,
		SuccessNum:      successNum,
		AbortNum:        abortNum,
		Pack:            r,
	})
	if err != nil {
		panic(err)
	}
}

type TestCases struct {
	Name     string
	Setup    Steps
	Teardown Steps
	Result   int
	Ddts     []Steps
	Times    string
}
