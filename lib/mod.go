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
)

type GoModAction struct {
	arg     *common_argument.CommonArgumentStruct
	Prefix  string  `json:"prefix"`
	appList AppList `json:"app_list"`
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
	configFilePath := fmt.Sprintf("%s/appmap.yaml", dir)
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
	cmdSlice := utils.CmdObject{Name: "go", Arg: []string{"get", "-v", "github.com/juetun/base-wrapper@v0.0.199"}}
	cmdSlice.Dir = pwdProject
	if err = utils.ExeCMD(&cmdSlice); err != nil {
		log.Error(err.Error())
		return
	}
	return
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
	return
}
