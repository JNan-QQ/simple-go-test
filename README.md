# simple-go-test

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
├── statics
│   ├── demo  #示例项目
│   └── html  #报告模板

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

1.  创建项目
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

3.  测试报告：按日期时间生成在`logs`目录中，并通过默认浏览器打开


#### [WIKI](https://gitee.com/jn-qq/simple-go-test/wikis)
