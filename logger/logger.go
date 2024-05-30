/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package logger

import (
	"fmt"
	"gitee.com/jn-qq/simple-go-test/config"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Logger 自定义log对象，写入文件
func Logger() {
	_ = os.MkdirAll("logs", os.ModePerm)
	// 日志备份
	func() {
		if _, err := os.Stat("logs/sgt.log"); err != nil {
			return
		}
		_ = os.Remove("logs/sgt.log.back1")
		_ = os.Rename("logs/sgt.log.back", "logs/sgt.log.back1")
		_ = os.Rename("logs/sgt.log", "logs/sgt.log.back")
	}()
	file, err := os.OpenFile("logs/sgt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetFlags(0)
}

var Reporthtml *ReportHtml = &ReportHtml{}

// INFO 提示信息，同时输出到终端和文件，换行输出
func INFO(msg string) {
	log.Println(msg)
	addHtml(fmt.Sprintf("<p class=\"info\">%s</p>", msg))
}

// STEP 步骤信息
func STEP(step int, msg string) {
	log.Printf("[第%02d步] %s\n", step, msg)
	addHtml(fmt.Sprintf("<p><span class=\"step\">[第%02d步]</span> %s</p>", step, msg))
}

// CHECK_POINT 检查点，如果allowRun == true,仅记录，继续运行
func CHECK_POINT(msg string, check bool, allowSkip bool) {
	if check {
		log.Printf("[CHECK_POINT PASS] %s", msg)
		addHtml(fmt.Sprintf(
			"<p class=\"check-point\"><span class=\"pass\">[CHECK_POINT PASS]</span> %s</p>",
			msg),
		)
	} else {
		log.Printf("[CHECK_POINT FAIL] %s", msg)
		addHtml(fmt.Sprintf(
			"<p class=\"check-point\"><spanspan class=\"fail\">[CHECK_POINT FAIL]</span> %s</p>",
			msg),
		)
	}
	if !check && !allowSkip {
		panic("[CHECK_POINT FAIL]")
	}
}

func ErrorInfo(msg string) {
	addHtml(fmt.Sprintf("<p class=\"fail\">%s</p>", msg))
	fmt.Println(msg)
}

func addHtml(msg string) {
	switch Reporthtml.Step {
	case 0:
		Reporthtml.Setup.addStep(msg)
	case 1:
		Reporthtml.Cases[len(Reporthtml.Cases)-1].Setup.addStep(msg)
	case 2:
		sc := Reporthtml.Cases[len(Reporthtml.Cases)-1].StepCases
		if sc == nil {
			Reporthtml.Cases[len(Reporthtml.Cases)-1].StepCases = &Steps{Html: template.HTML(msg)}
		} else {
			Reporthtml.Cases[len(Reporthtml.Cases)-1].StepCases.addStep(msg)
		}
	case 3:
		Reporthtml.Cases[len(Reporthtml.Cases)-1].Teardown.addStep(msg)
	case 4:
		Reporthtml.Teardown.addStep(msg)
	}
}

type Steps struct {
	Result int
	Html   template.HTML
	Times  string
}

func (s *Steps) time() {
	if s.Times == "" {
		s.Times = time.Now().Format("2006-01-02 15:04:05")
	}
}

func (s *Steps) addStep(h string) {
	s.time()
	s.Html += template.HTML(h)
}

func (s *Steps) setResult(r int) {
	s.time()
	s.Result = r
}

type ReportHtml struct {
	Name     string
	Setup    *Steps
	Teardown *Steps
	Cases    []*TestCases
	Child    []ReportHtml
	Step     int
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
	_, filename, _, _ := runtime.Caller(0)
	exPath := filepath.Dir(filepath.Dir(filename))
	tmpl, err := template.ParseFiles(
		filepath.Join(exPath, "statics", "html", "template.tmpl"),
		filepath.Join(exPath, "statics", "html", "js.tmpl"),
		filepath.Join(exPath, "statics", "html", "cases.tmpl"),
		filepath.Join(exPath, "statics", "html", "steps_result.tmpl"),
		filepath.Join(exPath, "statics", "html", "my.css"),
	)
	if err != nil {
		panic(err)
	}

	file, _ := os.Create(filepath.Join("logs", time.Now().Format("20060102-150405")) + ".html")
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	err = tmpl.Execute(file, html{
		ReportName:      config.ReportName,
		AllNum:          config.AllNum,
		SetupFailNum:    config.SetupFailNum,
		TeardownFailNum: config.TearDownFailNum,
		SpendTime:       r.Times,
		FailNum:         config.FailNum,
		SuccessNum:      config.SuccessNum,
		AbortNum:        config.AbortNum,
		Pack:            r,
	})
	if err != nil {
		panic(err)
	}
}

type TestCases struct {
	Name      string
	Setup     *Steps
	Teardown  *Steps
	Result    int
	StepCases *Steps
	Times     string
}
