package main

import (
	"flag"
	"fmt"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/manifoldco/promptui"
	"os/exec"
	"strconv"
	"strings"
)

type FlagParams struct {
	Release string `json:"release"`
}

func main() {
	// 定义几个变量，用于接收命令行的参数
	NewFlagParameter().InitFlag().
		ParamsImplement()

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
func (r *FlagParameter) ParamsImplement() {
	switch r.flagParameter.Release {
	case "develop":
		r.develop()

	}
}

func (r *FlagParameter) develop() {

	version := r.selectVersion()

	//执行当前版本要做的事
	r.cmdDevelopByVersion(version)

}

type CmdObject struct {
	Name string   `json:"name"`
	Arg  []string `json:"arg"`
}

func (r *FlagParameter) cmdDevelopByVersion(version string) {

	cmdSlice := []CmdObject{
		//{Name: "git", Arg: []string{"checkout", "master",},},
		//{Name: "git", Arg: []string{"branch", fmt.Sprintf("release/%s", version),}},
		{Name: "git", Arg: []string{"checkout", "-b", fmt.Sprintf("release/%s", version), "master"}},
		//{Name: "sed", Arg: []string{"-s", `version: "1.0.0"`, fmt.Sprintf(`version: "%s"`,version)}},
		{Name: "sed", Arg: []string{"-n", "2p", "./config/app.yaml",}},
	}
	fmt.Printf("创建分支")
	for _, item := range cmdSlice {
		err := r.exeCMD(&item)
		if err != nil {
		}
	}

}
func (r *FlagParameter) exeCMD(item *CmdObject) (err error) {
	var buf []byte
	cmdString := []interface{}{"【CMD】:", item.Name}
	for _, value := range item.Arg {
		cmdString = append(cmdString, value)
	}
	fmt.Println(cmdString...)
	cmd := exec.Command(item.Name, item.Arg...)
	buf, _ = cmd.Output()
	fmt.Printf("%s\n", buf)
	return
}

//版本选择交互
func (r *FlagParameter) selectVersion() (res string) {
	v := r.getVersion()
	versionList := r.orgVersion(v)
	fmt.Printf("Now version is :%s \r\n",v)
	prompt := promptui.Select{
		Label: fmt.Sprintf("Select Version",),
		Items: versionList,
	}
	_, vString, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("您选择的版本为：%s \n", vString)
	return vString
}

//获取当前可选择的版本
func (r *FlagParameter) orgVersion(version string) (res []string) {
	var err error
	var v int
	listVersionString := strings.Split(version, ".")
	res = make([]string, 0, len(listVersionString))
	for key, item := range listVersionString {
		listVersionString[key] = strings.TrimLeft(item, "v")
	}
	for key, item := range listVersionString {
		listVersionStringTmp := make([]string, len(listVersionString))
		copy(listVersionStringTmp, listVersionString)
		v, err = strconv.Atoi(item)
		if err != nil {
			panic(fmt.Sprintf("当前的版本号格式错误(%s;%s)", version, item))
		}
		listVersionStringTmp[key] = strconv.Itoa(v + 1)
		res = append(res, strings.Join(listVersionStringTmp, "."))
	}
	res = r.reverse(res)
	return
}

//倒序排序切片
func (r *FlagParameter) reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func (r *FlagParameter) getVersion() (res string) {
	common.PluginsApp()
	return app_obj.App.AppVersion
}

func (r *FlagParameter) InitFlag() (res *FlagParameter) {
	res = r
	// &user 就是接收命令行中输入 -u 后面的参数值，其他同理
	flag.StringVar(&r.flagParameter.Release, "release", "develop", "创建一个新的版本")

	// 解析命令行参数写入注册的flag里
	flag.Parse()
	// 输出结果
	fmt.Printf("参数为:%v", r.flagParameter)
	return
}
