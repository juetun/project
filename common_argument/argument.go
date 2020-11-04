package common_argument

import (
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
)

var CommonArgument = CommonArgumentStruct{}

func init() {
	common.PluginsApp()
	CommonArgument.Application = common.GetAppConfig()
}

type CommonArgumentStruct struct {
	*app_obj.Application
}
