package ast

import (
	"bytes"
	"gitee.com/jn-qq/simple-go-test/config"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// addImport 导入包
func addImport(genDecl *ast.GenDecl, packs ...string) {
	// 遍历包名
	for _, pack := range packs {
		pack = config.GetRunPackage() + "/" + pack
		hasImport := false
		// 判断是否已导入
		for _, spec := range genDecl.Specs {
			// 如果已经包含"xxx"
			if spec.(*ast.ImportSpec).Path.Value == strconv.Quote(pack) {
				hasImport = true
				break
			}
		}
		// 未导入，开始导入
		if !hasImport {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(pack),
				},
			})
		}
	}
}

// TestPackAst 测试用例语法树
type TestPackAst struct {
	name         *ast.KeyValueExpr // 包名
	suitSetup    *ast.KeyValueExpr // 包初始化
	suitTeardown *ast.KeyValueExpr // 包清除
	tests        []*ast.CallExpr   // 包测试用例
	children     []*TestPackAst    // 子包
}

// setPackName 设置包名
func (t *TestPackAst) setPackName(name string) {
	t.name = &ast.KeyValueExpr{
		Key:   ast.NewIdent("Name"),
		Value: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(name)},
	}
}

// addSuitMethod 设置初始化、清除函数
func (t *TestPackAst) addSuitMethod(packName, funName string, isSetup bool) {
	if isSetup {
		t.suitSetup = &ast.KeyValueExpr{
			Key: ast.NewIdent("SuitSetUp"),
			Value: &ast.SelectorExpr{
				X:   ast.NewIdent(packName),
				Sel: ast.NewIdent(funName),
			},
		}
	} else {
		t.suitTeardown = &ast.KeyValueExpr{
			Key: ast.NewIdent("SuitTearDown"),
			Value: &ast.SelectorExpr{
				X:   ast.NewIdent(packName),
				Sel: ast.NewIdent(funName),
			},
		}
	}
}

// addStruct 统计包下测试用例
func (t *TestPackAst) addStruct(packName, structName string) {
	if t.tests == nil {
		t.tests = make([]*ast.CallExpr, 0)
	}
	t.tests = append(t.tests, &ast.CallExpr{
		Fun: ast.NewIdent("new"),
		//Lparen:   0,
		Args: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(packName),
				Sel: ast.NewIdent(structName),
			},
		},
	})
}

// addChildPack 添加子包
func (t *TestPackAst) addChildPack(p *TestPackAst) {
	if t.children == nil {
		t.children = make([]*TestPackAst, 0)
	}
	t.children = append(t.children, p)
}

// ToAst 转换输出 AST 对象，用于替换main.go
func (t *TestPackAst) ToAst() *ast.CompositeLit {
	if t.name == nil || (t.tests == nil && t.children == nil) || (len(t.children) == 0 && len(t.tests) == 0) {
		return nil
	}
	res := &ast.CompositeLit{
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("runner"),
			Sel: ast.NewIdent("TestPackage"),
		},
		Elts: []ast.Expr{
			t.name,
		},
	}
	if t.suitSetup != nil {
		res.Elts = append(res.Elts, t.suitSetup)
	}
	if t.suitTeardown != nil {
		res.Elts = append(res.Elts, t.suitTeardown)
	}
	if t.tests != nil && len(t.tests) > 0 {
		res.Elts = append(res.Elts, &ast.KeyValueExpr{
			Key: ast.NewIdent("Tests"),
			Value: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.SelectorExpr{
						X:   ast.NewIdent("runner"),
						Sel: ast.NewIdent("TestInterface"),
					},
				},
				Elts: func() []ast.Expr {
					var ex []ast.Expr
					for _, _test := range t.tests {
						ex = append(ex, _test)
					}
					return ex
				}(),
			},
		})
	}

	if t.children != nil && len(t.children) > 0 {
		res.Elts = append(res.Elts, &ast.KeyValueExpr{
			Key: ast.NewIdent("Child"),
			Value: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.SelectorExpr{
						X:   ast.NewIdent("runner"),
						Sel: ast.NewIdent("TestPackage"),
					},
				},
				Elts: func() []ast.Expr {
					var ex []ast.Expr
					for _, c := range t.children {
						if x := c.ToAst(); x != nil {
							ex = append(ex, x)
						}
					}
					return ex
				}(),
			},
		})
	}

	return res
}

// Find 遍历包下的go文件
func Find(root string) *TestPackAst {
	astPack := new(TestPackAst)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		// 统一路径分隔符
		path = strings.ReplaceAll(path, "\\", "/")
		// 格式化包路径
		packPath := strings.TrimPrefix(path, strings.TrimSuffix(root, config.CasesDir))
		if d.IsDir() {
			// 子包深度
			deep := len(strings.Split(packPath, "/"))

			// 包名
			config.PackageNames = append(config.PackageNames, packPath)

			var astP = astPack
			for i := 1; i < deep; i++ {
				if astP.children == nil {
					astP.children = make([]*TestPackAst, 0)
				}
				astP.children = append(astP.children, new(TestPackAst))
				astP = astP.children[len(astP.children)-1]
			}
		} else {
			// 文件目录深度
			deep := len(strings.Split(filepath.Dir(packPath), "/"))
			var astP = astPack
			for i := 1; i < deep; i++ {
				astP = astP.children[len(astP.children)-1]
			}
			readAst(path, astP)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return astPack
}

// 读取文件语法树
func readAst(path string, t *TestPackAst) {
	file, _ := filepath.Abs(path)
	fileSet := token.NewFileSet()
	parseFile, _ := parser.ParseFile(fileSet, file, nil, parser.ParseComments)

	var packName string
	ast.Inspect(parseFile, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.File:
			packName = n.Name.Name
			t.setPackName(packName)
		case *ast.FuncDecl:
			if n.Recv != nil {
				return true
			}
			if n.Name.Name == "SuiteSetUp" {
				t.addSuitMethod(packName, n.Name.Name, true)
			} else if n.Name.Name == "SuiteTearDown" {
				t.addSuitMethod(packName, n.Name.Name, false)
			}
		case *ast.TypeSpec:
			switch v := n.Type.(type) {
			case *ast.StructType:
				for _, field := range v.Fields.List {
					vv := field.Type.(*ast.SelectorExpr)
					if vv.X.(*ast.Ident).Name == "runner" && vv.Sel.Name == "Test" {
						t.addStruct(packName, n.Name.Name)
						break
					}
				}
			case *ast.SelectorExpr:
				if v.X.(*ast.Ident).Name == "runner" && v.Sel.Name == "Test" {
					t.addStruct(packName, n.Name.Name)
				}
			}
		}
		return true
	})
}

// WriteToMain 更新import TestObject 保存到main.go
func WriteToMain(v *ast.CompositeLit, selects ...any) {
	file, err := filepath.Abs("main.go")
	if err != nil {
		panic(err)
	}
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.GenDecl:
			if v != nil {
				if n.Tok == token.IMPORT {
					// 查找有没有import context包
					addImport(n, config.PackageNames...)
				} else if n.Tok == token.VAR {
					x := n.Specs[0].(*ast.ValueSpec)
					if x.Names[0].Name == "testTree" && v != nil {
						x.Values[0] = v
					}
				}
			} else if len(selects) == 2 {
				if n.Tok == token.VAR {
					x := n.Specs[0].(*ast.ValueSpec)
					if x.Names[0].Name == "selectBy" && selects != nil {
						x.Values[0].(*ast.BasicLit).Value = strconv.Itoa(selects[0].(int))
					} else if x.Names[0].Name == "selectValue" {
						x.Values[0].(*ast.BasicLit).Value = strconv.Quote(selects[1].(string))
					}
				}

			}
		}
		return true
	})

	var output []byte
	buffer := bytes.NewBuffer(output)
	_ = format.Node(buffer, fileSet, f)

	fss, _ := os.OpenFile("main.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer func() {
		_ = fss.Close()
	}()

	n, _ := fss.Seek(0, os.SEEK_END)
	_, err = fss.WriteAt(buffer.Bytes(), n)
}
