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

1.  xxxx
2.  xxxx
3.  xxxx