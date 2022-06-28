package common_argument

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/plugins"
)

var CommonArgument = CommonArgumentStruct{}

func InitConfig() {
	var preLoad=&app_start.PluginsOperate{}
	plugins.PluginsApp(preLoad)
	CommonArgument.Application = common.GetAppConfig()
}

type CommonArgumentStruct struct {
	*app_obj.Application
}
