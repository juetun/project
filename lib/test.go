package lib

import (
	"fmt"
	"github.com/juetun/project/common_argument"
	"github.com/juetun/project/utils"
	"github.com/prometheus/common/log"
	"strings"
)

type TestAction struct {
	arg *common_argument.CommonArgumentStruct
}

func NewTestAction(arg *common_argument.CommonArgumentStruct) (res *TestAction) {
	return &TestAction{
		arg: arg,
	}
}
func (r *TestAction) Run() {
	var branchName = strings.TrimLeft(r.arg.AppVersion,"v")
	cmdSlice := []utils.CmdObject{
		{Name: "git", Arg: []string{"checkout", "-B", "master", "origin/master"}},
		{Name: "git", Arg: []string{"pull"}},
		{Name: "git", Arg: []string{"checkout", "-B", "develop", "origin/develop"}},
		{Name: "git", Arg: []string{"rebase", "origin/master"},},
		{Name: "git", Arg: []string{"branch"}},
		{Name: "git", Arg: []string{"pull", "origin", fmt.Sprintf("release/%s", branchName)}},
		{Name: "git", Arg: []string{"push", "--set-upstream", "origin", "develop",}},
	}
	fmt.Printf("合并推送分支:\n")
	for _, item := range cmdSlice {
		err := utils.ExeCMD(&item)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
