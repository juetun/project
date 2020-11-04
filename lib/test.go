package lib

import (
	"fmt"
	"github.com/juetun/project/utils"
	"github.com/prometheus/common/log"
)

type TestAction struct {
}

func NewTestAction() (res *TestAction) {
	return &TestAction{}
}
func (r *TestAction) Run() {
	var branchName string = "1.0.1"
	cmdSlice := []utils.CmdObject{
		{Name: "git", Arg: []string{"checkout", "-B", "master", "origin/master"}},
		{Name: "git", Arg: []string{"pull"}},
		{Name: "git", Arg: []string{"checkout", "-B", "develop", "origin/develop"}},
		{Name: "git", Arg: []string{"rebase","origin/master"}},
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
