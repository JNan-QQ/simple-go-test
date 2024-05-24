/**
  Copyright (c) [2024] [JiangNan]
  [simple-go-test] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

package test

import (
	"fmt"
	sgt "gitee.com/jn-qq/simple-go-test/logger"
	"html/template"
	"net/http"
	"testing"
	"time"
)

func TestHtml(t *testing.T) {
	type html struct {
		ReportName      string
		AllNum          int
		SetupFailNum    int
		TeardownFailNum int
		SpendTime       int
		FailNum         int
		SuccessNum      int
		AbortNum        int
		Pack            sgt.ReportHtml
	}
	tmplDemo := func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(
			"D:\\Documents\\Go\\src\\simple-go-test\\html\\template.tmpl",
			"D:\\Documents\\Go\\src\\simple-go-test\\html\\js.tmpl",
			"D:\\Documents\\Go\\src\\simple-go-test\\html\\cases.tmpl",
			"D:\\Documents\\Go\\src\\simple-go-test\\html\\steps_result.tmpl",
			"D:\\Documents\\Go\\src\\simple-go-test\\html\\my.css",
		)
		if err != nil {
			fmt.Println(err)
		}

		err = tmpl.Execute(w, html{
			ReportName:      "第一次测试报告",
			AllNum:          500,
			SetupFailNum:    20,
			TeardownFailNum: 10,
			SpendTime:       5631,
			FailNum:         10,
			SuccessNum:      400,
			AbortNum:        90,
			Pack: sgt.ReportHtml{
				Name: "cases",
				Times: []time.Time{
					time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Setup: &sgt.Steps{
					Result: 1,
					Html:   "<p>初始化</p>",
					Times:  time.Now().Format("2006-01-02 15:04:05"),
				},
				Teardown: &sgt.Steps{
					Result: 1,
					Html:   "<p>清除</p>",
					Times:  time.Now().Format("2006-01-02 15:04:05"),
				},
				Cases: []sgt.TestCases{
					{
						Name: "C0001",
						Setup: sgt.Steps{
							Result: 1,
							Html:   "<p>初始化</p>",
						},
						Times:    time.Now().Format("2006-01-02 15:04:05"),
						Teardown: sgt.Steps{},
						Result:   0,
						Ddts:     nil,
					},
					{
						Name:  "C0002",
						Times: time.Now().Format("2006-01-02 15:04:05"),
						Setup: sgt.Steps{
							Result: 1,
							Html:   "<p>初始化</p>",
						},
						Teardown: sgt.Steps{
							Result: 1,
							Html:   "<p>清除</p>",
						},
						Result: 1,
						Ddts: []sgt.Steps{
							{
								Result: 1,
								Html:   "<p>do somethings ...</p>",
							},
							{
								Result: 1,
								Html:   "<p>do somethings ...</p>",
							},
						},
					},
				},
				Child: nil,
			},
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	server := http.Server{
		Addr: "127.0.0.1:80",
	}
	http.HandleFunc("/", tmplDemo)
	server.ListenAndServe()
}
