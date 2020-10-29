package main

import (
	"flag"
	"fmt"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/project/utils/bar"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FlagParams struct {
	Release string `json:"release"`
}

func main() {
	// 定义几个变量，用于接收命令行的参数
	NewFlagParameter().InitFlag().
		ParamsImplement()
	s := fmt.Sprintf("\033[s")
	fmt.Printf("kkkk:%s \n", s)
	fmt.Println("1dsafasd")
	fmt.Println("2dsafasd")
	fmt.Println("3dsafasd")
	fmt.Println("4dsafasd")
	fmt.Println("5dsafasd")
	fmt.Println("6dsafasd")
	fmt.Println("7dsafasd")
	fmt.Println("8dsafasd")
	fmt.Println("9dsafasd")
	fmt.Println("10dsafasd")

	fmt.Printf("\033[s")
	var vString string
	fmt.Printf("\033[%dA", 4)
	fmt.Scanln(&vString) //读取键盘的输入，通过操作地址，赋值给x和y   阻塞式
	//fmt.Printf("\033[u")
	fmt.Printf("\033[%dB", 4)
	fmt.Printf("您输入的内容为:%s \n", vString)
	//Bar()
}
func Bar() {

	pgbar := bar.New("多线程进度条")
	bar.Println("进度条1")
	b := pgbar.NewBar("1st", 20000)
	b.SetSpeedSection(900, 100)
	//bar.Println("独立进度条")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20000; i++ {
			b.Add()
			time.Sleep(time.Second / 2000)
		}
	}()
	wg.Wait()
}

type FlagParameter struct {
	flagParameter FlagParams
}

func NewFlagParameter() (res *FlagParameter) {
	return &FlagParameter{}
}

//各个参数命令的实现
func (r *FlagParameter) ParamsImplement() {
	var err error
	switch r.flagParameter.Release {
	case "develop":
		version := r.develop()
		err= exec.Command("git","checkout master").Run()
		if err != nil {
			panic(err.Error())
		}

		err = exec.Command("git", "branch", fmt.Sprintf("release/%s", version)).Run()
		if err != nil {
			panic(err.Error())
		}
	}
}

func (r *FlagParameter) develop() string {
	version := r.getVersion()
	versionList := r.orgVersion(version)
	fmt.Printf("您可以选择如下版本的：\n")
	defaultValue := ""
	var v string
	for _, value := range versionList {

		if defaultValue == "" {
			defaultValue = value
			v = fmt.Sprintf("v%s", value)
			v = fmt.Sprintf("》%s", v)
		} else {
			v = fmt.Sprintf("v%s", value)
		}
		fmt.Printf("%s\n", v)
	}

	var vString string
	fmt.Println("请选择版本：")
	fmt.Scanln(&vString) //读取键盘的输入，通过操作地址，赋值给x和y   阻塞式
	if vString == "" {
		vString = defaultValue
	}
	fmt.Printf("您选择的版本为：%s \n", vString)
	return vString
	//var reader = bufio.NewReader(os.Stdin) //os.Stdin标椎输入（键盘）
	//s1, _ := reader.ReadString('\n')       //读到\n结束
	//fmt.Println("读到的数据：", s1)
}

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
