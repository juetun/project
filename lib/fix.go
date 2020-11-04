package lib

import "github.com/juetun/project/common_argument"

type FixAction struct {
	arg *common_argument.CommonArgumentStruct
}

func NewFixAction(args *common_argument.CommonArgumentStruct) (res *FixAction) {
	return &FixAction{
		arg: args,
	}
}

func (r *FixAction) Run() {

}
