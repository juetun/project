package lib

import (
	"fmt"
	"github.com/juetun/project/common_argument"
	"github.com/juetun/project/utils"
	"github.com/manifoldco/promptui"
	"github.com/prometheus/common/log"
	"strconv"
	"strings"
)

type DevelopAction struct {
	arg *common_argument.CommonArgumentStruct
}

func NewDevelopAction(arg *common_argument.CommonArgumentStruct) (r *DevelopAction) {
	return &DevelopAction{
		arg: arg,
	}
}

func (r *DevelopAction) Run() {
	r.develop()
}

//https://www.jianshu.com/p/34b95c5eedb6
func (r *DevelopAction) cmdDevelopByVersion(version, sourceVersion string) {

	branchName := fmt.Sprintf("release/%s", version)
	cmdSlice := []utils.CmdObject{
		{Name: "git", Arg: []string{"fetch"}},
		{Name: "git", Arg: []string{"checkout", "master"}},
		{Name: "git", Arg: []string{"branch"}},
		{Name: "git", Arg: []string{"pull"}},
		{Name: "git", Arg: []string{"checkout", "-B", branchName, "origin/master"}},
		{Name: "git", Arg: []string{"push", "--set-upstream", "origin", branchName,}},
		{Name: "git", Arg: []string{"branch"}},
		{Name: "sed", Arg: []string{"-i", ``, fmt.Sprintf(`s/"%s"/"%s"/g`, sourceVersion, version), `./config/app.yaml`}},
	}
	fmt.Printf("创建分支:\n")
	for _, item := range cmdSlice {
		if err := utils.ExeCMD(&item); err != nil {
			log.Fatalf(err.Error())
			return
		}
	}
	return
}
func (r *DevelopAction) develop() {

	version, sourceVersion := r.selectVersion()

	//执行当前版本要做的事
	r.cmdDevelopByVersion(version, sourceVersion)

}

//版本选择交互
func (r *DevelopAction) selectVersion() (res string, sourceVersion string) {
	sourceVersion = r.getVersion()
	versionList := r.orgVersion(sourceVersion)
	fmt.Printf("Now version is :%s \r\n", sourceVersion)
	prompt := promptui.Select{
		Label: fmt.Sprintf("Select Version", ),
		Items: versionList,
	}
	_, res, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("您选择的版本为：%s \n", res)
	return
}

//获取当前可选择的版本
func (r *DevelopAction) orgVersion(version string) (res []string) {

	var err error
	var v int
	listVersionString := strings.Split(version, ".")
	res = make([]string, 0, len(listVersionString)+1)
	if listVersionString[len(listVersionString)-1] == "0" {
		res = append(res, fmt.Sprintf("%s.1", version))
	}

	for key, item := range listVersionString {
		listVersionString[key] = item
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
func (r *DevelopAction) reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func (r *DevelopAction) getVersion() (res string) {
	return strings.TrimLeft(r.arg.AppVersion, "v")
}
