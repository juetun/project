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

type PullCfgAction struct {
	arg             *common_argument.CommonArgumentStruct
	Prefix          string  `json:"prefix"`
	appList         AppList `json:"app_list"`
	DependPkgString string  `json:"depend_pkg_string"`
}

func (r *PullCfgAction) ReadData() (err error) {
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

func (r *PullCfgAction) Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (r *PullCfgAction) runItem(proPatch string) (err error) {

	if proPatch == "library" {
		log.Warnf("忽略项目目录(%s) \n", proPatch)
		return
	}

	log.Info("执行项目:", proPatch)
	goPatch := os.Getenv("GOPATH")
	pwdProject := fmt.Sprintf("%s/src/%s%s/config/apps", goPatch, r.Prefix, proPatch)
	if !r.Exists(pwdProject) {
		log.Warnf("项目目录(%s)不存在 \n", pwdProject)
		return
	}

	cmdSlice := [] utils.CmdObject{
		{Name: "git", Arg: []string{"pull"}, Dir: pwdProject,},
	}
	for _, cmd := range cmdSlice {
		if err = utils.ExeCMD(&cmd); err != nil {
			log.Error(err.Error())
			return
		}
	}

	return
}

func (r *PullCfgAction) Run() {
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

func NewPullCfgAction(arg *common_argument.CommonArgumentStruct) (res *PullCfgAction) {
	res = &PullCfgAction{arg: arg}
	if res.Prefix == "" {
		res.Prefix = "github.com/juetun/"
	}

	return
}
