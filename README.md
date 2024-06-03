# simple-go-test



[TOC]



#### 介绍
go 自动化测试框架

#### 软件架构
```shell
.
├── ast             #读取更改go文件
├── config          #配置，公共常量、变量
├── logger          #日志、报告
├── runner          #按指定规则运行测试用例
├── sgt.go          #入口文件
├── statics			#静态模板文件
│   ├── demo  			#示例项目
│   └── html  			#报告模板

```
#### 安装教程
1. 保证`go version >= 1.21`
2. 安装`go install gitee.com/jn-qq/simple-go-test@latest`
3. 命令行参数：
```shell
xxx@xxx$ simple-go-test -help
Usage of simple-go-test:
  -collect
        收集测试用例
  -new string
        新建一个测试项目
  -pkg string
        包名过滤,支持正则表达式
  -report-title string
        报告标题 (default "测试报告")
  -run
        运行
  -tag string
        tag过滤,支持正则表达式
  -tagNot string
        tag反向过滤,支持正则表达式
  -test string
        用例名称过滤,支持正则表达式
  -version
        版本

```

#### 使用说明

1. 创建项目

```shell
# 进入项目要存放文件夹
simple-go-test -new 项目名称

# 开始创建项目...
# statics/demo 复制完成
# statics/demo/cases 复制完成
# statics/demo/cases/cases.go 复制完成
# statics/demo/cases/student 复制完成
# statics/demo/cases/student/student.go 复制完成
# statics/demo/cases/teacher 复制完成
# statics/demo/cases/teacher/teacher.go 复制完成
# statics/demo/go.x 复制完成
# statics/demo/main.go 复制完成
# 项目创建完成,请同步包 go mod tid
```
2.  运行
```shell
cd 项目名称

# 收集测试用例
simple-go-test -collect

# 设置过滤条件
# simple-go-test -过滤模式 过滤正则表达式
simple-go-test -pkg cases
simple-go-test -tag 冒烟测试
simple-go-test -tagNot 冒烟测试
simple-go-test -test "C00*"

# 运行
simple-go-test -run
```
**PS: 以上所有参数可以在一个命令中输入**
例：`simple-go-test -collect -pkg cases -run`

3. 测试报告：按日期时间生成在`logs`目录中，并通过默认浏览器打开

#### 

# Wiki



### 用例包结构

- 入口包应为 cases 包，包名与目录名相同
- 测试用例就是包中的一个结构体
- 包可以包含子包
- 每个包中应有组织的放入用例（数据环境）

*PS: 公共代码应放入自定义lib包中，保持测试用例简单整洁*



### 用例 结构体

- 实现接口 `runner.TestInterface`
- `组合`或设置`类型`为 `runner.Test`

```go
// TestInterface 测试用例接口
type TestInterface interface {
	Init() *Test
	SetUp()
	TearDown()
	TestStep()
}

// Test 测试用例类型
type Test struct {
	Name string
	Tags []string
	DDT  []any
	Para any
}
```

##### Example

```go
// 1.组合

// TestNum 实现接口 runner.TestInterface
type TestNum struct {
	runner.Test
    // 自定义参数
}

// Init 测试用例构造
func (t *TestNum) Init() *runner.Test {
	// 参数赋值 设置 返回
	t.Name = "C0001"
	t.Tags = []string{"cases", "冒烟测试", "num"}

	return &t.Test
}

// SetUp 测试用例初始化
func (t *TestNum) SetUp() {}

// TearDown 测试用例清除
func (t *TestNum) TearDown() {}

// TestStep 测试步骤
func (t *TestNum) TestStep() {
	// do something...
}


// 2.设置类型

// TestString 实现接口 runner.TestInterface
type TestString runner.Test

func (t *TestString) Init() *runner.Test {
	// 参数赋值 设置 返回
	t.Name = "C0002"
	t.Tags = []string{"cases", "冒烟测试", "string"}
	t.DDT = []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	return (*runner.Test)(t)
}

func (t *TestString) SetUp() {}

func (t *TestString) TearDown() {}

func (t *TestString) TestStep() {
	// do something...
}
```

**PS：结构体的方法应为指针方法，避免调用空值。组合方法可以设置自定义参数。**



##### 1. 结构体字段赋值  Init（必要）

- `Name`、`Tags`字段不能为空。
- `DDT`(数据驱动可选)，当DDT不为空时，运行时自动将DDT里的每个元素赋值给`Para`以供用户调用。

##### 2. 测试步骤 TestStep

*当数据驱动不为空时 通过 Para(类型自己断言) 获取每个数据*

`x := t.Para.(type) ...`



### 初始化 / 清除

一类、一组测试用例测试前需要执行相同操作。测试完成后恢复测试环境。

##### 单个用例

在用例结构体 中 补全对应方法

```go
// TestNum 实现接口 runner.TestInterface
type TestNum struct {
	runner.Test
    // 自定义参数
}


// SetUp 测试用例初始化
func (t *TestNum) SetUp() {
    // do something...
}

// TearDown 测试用例清除
func (t *TestNum) TearDown() {
    // do something...
}

...
```

 执行用例时

- 先用例结构赋值，执行 `Init`里面的代码

- 再执行 `SetUp`里面的代码
- 再执行 `TestStep`里面的代码
- 最后再执行 `TearDown`里面的代码。

PS：如果初始化失败、则用例运行异常。跳过`TestStep`直接运行清除。如果清除失败，方法算通过。

执行流程

```flow
st=>start: RUN

op=>operation: Init
op1=>operation: SetUp

cond=>condition: DDT为空

sub1=>subroutine: 遍历DDT赋值给Para运行TestStep

op2=>operation: TestStep



op3=>operation: TearDown

cond1=>condition: DDT遍历完成

e=>end: 结束框

st->op->op1->cond

cond(yes)->op2->op3->e

cond(no)->sub1->op3
```



##### 包

在改包下写入函数 `SuiteSetUp`、` SuiteTearDown` 即可

执行时

- 先执行`SuiteSetUp`里面的代码(如果有)

- 再 遍历执行包里的测试用例
- 再遍历子包、重复步骤
- 最后再执行 `SuiteTearDown`里面的代码。

执行流程：

```flow
st=>start: RUN

op=>operation: SuiteSetUp
op1=>operation: 运行本包下测试用例

cond=>condition: 子包

sub1=>subroutine: 遍历子包

op3=>operation: SuiteTearDown

cond1=>condition: DDT遍历完成

e=>end: 结束框

st->op->op1->cond

cond(no)->op3->e

cond(yes)->sub1(left)->op
```



### 数据驱动

有一批测试用例，具有 `相同的测试步骤` ，只是 `测试参数数据不同` 。将数据从测试程序中分离出来，以后增加新的测试用例，只需要修改数据。这就是数据驱动。

在用例结构体的`Init`方法中对`DDT`字段赋值即可应用数据驱动。

```go
// TestString 实现接口 runner.TestInterface
type TestString runner.Test

func (t *TestString) Init() *runner.Test {
	// 参数赋值 设置 返回
	t.Name = "C0002"
	t.Tags = []string{"cases", "冒烟测试", "string"}
	t.DDT = []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	return (*runner.Test)(t)
}

func (t *TestString) TestStep() {
	// do something...
    x := t.Para.(int)
    fmt.Println(x)
}
```

执行过程中，程序会遍历`DDT`将里面的元素赋值给` Para` 。用户可以在` TestStep` 中调用 `Para(.type)`。



### 用例 过滤

##### 标签过滤/反过滤

在用例结构体的`Init`方法中对`Tags`字符串切片赋值即为该用例标签。

```shell
# 执行包含 标签 '冒烟测试' 的用例. 
# -tag 标签正则表达式
-tag 冒烟测试  


# 执行不包含标签 '冒烟测试' 的用例.
# -tag 标签正则表达式
-tagNot 冒烟测试 
```



##### 用例名过滤

在用例结构体的`Init`方法中对`Name`字段赋值即可设置用例名。

```shell
-test 用例名正则表达式
```



##### 包名过滤

```shell
-pkg 包名正则表达式
```

**过滤包 包含子包**



### 其他功能

- `logger.STEP` 函数用来声明每个测试步骤，这样日志报告更清晰

- `logger.INFO` 函数用来打印一些信息在日志和报告中,同时输出到终端，方便出现问题时定位。

- `CHECK_POINT` 函数用来声明测试过程中的每个检查点，任何一个检查点不通过，整个测试用例就被认为不通过。缺省情况下，一个检查点不通过，后面的测试代码就不会继续执行。如果你希望 某个检查点即使不通过，后续代码仍然继续执行，可以使用参数 `allowSkip=true` 。`CHECK_POINT('即使不通过也不中止', False, allowSkip=true)`
- `GSTORE` 共享数据集合
- 每个测试用例执行结束后执行函数，在项目main.go中，对`config.CasesEnd` 赋值即可
- 项目结束后执行函数，在项目main.go中，对`config.TestEnd` 赋值即可

```
package main

import (
	...
)

var testTree = runner.TestPackage{}

var selectBy config.By = 0
var selectValue = ""

func main() {
	logger.Logger()
	
	config.CasesEnd = func(testName string, testResult int) {
		fmt.Println(testName, testResult)
	}
	
	config.TestEnd = func(allNum, failNum, successNum, abortNum, setupFail, teardownFail int) {
		fmt.Println(allNum, failNum, successNum, setupFail, teardownFail)
	}
	...

	report.Save()
}

```

