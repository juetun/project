package main

import (
	"flag"
	"fmt"
	"github.com/juetun/project/common_argument"
	"github.com/juetun/project/lib"
	"log"
)

type FlagParams struct {
	Release string `json:"release"`
}

func main() {
	// 定义几个变量，用于接收命令行的参数
	NewFlagParameter().InitFlag().
		Run()

	//s := fmt.Sprintf("\033[s")
	//fmt.Printf("kkkk:%s \n", s)
	//fmt.Println("1dsafasd")
	//fmt.Println("2dsafasd")
	//fmt.Println("3dsafasd")
	//fmt.Println("4dsafasd")
	//fmt.Println("5dsafasd")
	//fmt.Println("6dsafasd")
	//fmt.Println("7dsafasd")
	//fmt.Println("8dsafasd")
	//fmt.Println("9dsafasd")
	//fmt.Println("10dsafasd")
	//
	//fmt.Printf("\033[s")
	//var vString string
	//fmt.Printf("\033[%dA", 4)
	//fmt.Scanln(&vString) //读取键盘的输入，通过操作地址，赋值给x和y   阻塞式
	////fmt.Printf("\033[u")
	//fmt.Printf("\033[%dB", 4)
	//fmt.Printf("您输入的内容为:%s \n", vString)
}

type FlagParameter struct {
	flagParameter FlagParams
}

func NewFlagParameter() (res *FlagParameter) {
	return &FlagParameter{}
}

//各个参数命令的实现
func (r *FlagParameter) Run() {
	switch r.flagParameter.Release {
	case "mod":
		lib.NewGoModAction(&common_argument.CommonArgument).
			Run()
	case "develop":                  //生成开发数据
		common_argument.InitConfig() //初始化配置数据
		lib.NewDevelopAction(&common_argument.CommonArgument).
			Run()
	case "fix": //修复BUG
		lib.NewFixAction(&common_argument.CommonArgument).
			Run()
	case "test": //测试数据
		lib.NewTestAction(&common_argument.CommonArgument).
			Run()
	default:
		log.Fatalf("您输入的指令当前不支持 (%s)\n", r.flagParameter.Release)
	}
}

func (r *FlagParameter) InitFlag() (res *FlagParameter) {
	res = r
	// &user 就是接收命令行中输入 -u 后面的参数值，其他同理
	flag.StringVar(&r.flagParameter.Release, "release", "develop", "创建一个新的版本(支持如下参数):\ndevelop:\t\t生成一个新的开发分支\ntest:\t\t\t将代码发布到develop分支\n")

	// 解析命令行参数写入注册的flag里
	flag.Parse()

	// 输出结果
	fmt.Printf("参数为:%v\n", r.flagParameter)
	return
}
