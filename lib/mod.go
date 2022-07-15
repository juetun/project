package lib

import (
	"fmt"
	"github.com/juetun/project/common_argument"
	"github.com/juetun/project/utils"
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type GoModAction struct {
	arg             *common_argument.CommonArgumentStruct
	Prefix          string  `json:"prefix"`
	appList         AppList `json:"app_list"`
	DependPkgString string  `json:"depend_pkg_string"`
}

type AppList struct {
	Apps []string `json:"apps" yaml:"apps"`
}

func (r *GoModAction) ReadData() (err error) {
	var yamlFile []byte
	var dir string

	if dir, err = os.Getwd(); err != nil {
		return
	}
	configFilePath := fmt.Sprintf("%s/config/apps/appcommon.yaml", dir)
	log.Infof("读取配置文件:%s", configFilePath)
	if yamlFile, err = ioutil.ReadFile(configFilePath); err != nil {
		log.Fatalf("yamlFile.Get err #%v \n", err)
	}
	var appList AppList
	if err = yaml.Unmarshal(yamlFile, &appList); err != nil {
		log.Fatalf("Load  appMap config failure(%#v) \n", err)
	}
	r.appList = appList
	return
}

func (r *GoModAction) runItem(proPatch string) (err error) {
	log.Info("执行项目:", proPatch)
	goPatch := os.Getenv("GOPATH")
	pwdProject := fmt.Sprintf("%s/src/%s%s", goPatch, r.Prefix, proPatch)
	if !r.Exists(pwdProject) {
		log.Warnf("目录(%s)不存在", pwdProject)
		return
	}
	cmdSlice := [] utils.CmdObject{
		{Name: "go", Arg: []string{"get", "-v", r.DependPkgString}, Dir: pwdProject,},
		{Name: "git", Arg: []string{"commit", "-am", r.DependPkgString}, Dir: pwdProject,},
	}
	for _, cmd := range cmdSlice {
		if err = utils.ExeCMD(&cmd); err != nil {
			log.Error(err.Error())
			return
		}
	}
	cmdSlice = [] utils.CmdObject{
		{Name: "git", Arg: []string{"push", "origin"}, Dir: pwdProject,},
	}
	time.Sleep(500*time.Millisecond)
	for _, cmd := range cmdSlice {
		if err = utils.ExeCMD(&cmd); err != nil {
			log.Error(err.Error())
			return
		}
	}
	return
}

func (r *GoModAction) Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func (r *GoModAction) Run() {
	var err error
	if err = r.ReadData(); err != nil {
		return
	}

	fmt.Printf("开始拉取最新项目:\n")
	//
	//for _, projectPatch := range r.appList.Apps {
	//	r.runItem(projectPatch)
	//}
	var syncG sync.WaitGroup
	var i int
	var have bool
	for _, projectPatch := range r.appList.Apps {
		i++
		syncG.Add(1)
		have = true
		go func(proPatch string) {
			defer syncG.Done()
			r.runItem(proPatch)
		}(projectPatch)
		if i%5 == 0 {
			syncG.Wait()
			have = false
		}
	}
	if have {
		log.Info("结束操作")
		syncG.Wait()
	}
	log.Info("操作完成")
	return
}

func NewGoModAction(arg *common_argument.CommonArgumentStruct) (res *GoModAction) {
	res = &GoModAction{arg: arg}
	if res.Prefix == "" {
		res.Prefix = "github.com/juetun/"
	}
	if res.DependPkgString == "" {
		res.DependPkgString = "github.com/juetun/base-wrapper@latest"
	}
	return
}
