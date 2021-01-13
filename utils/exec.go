package utils

import (
	"bytes"
	"fmt"
	"github.com/prometheus/common/log"
	"os/exec"
)

type ErrorCallBack func(e error, errString string) (err error)
type CmdObject struct {
	Name        string        `json:"name"`
	Arg         []string      `json:"arg"`
	ErrCallBack ErrorCallBack `json:"err_call_back"`
}


func ExeCMD(item *CmdObject) (err error) {

	var buf []byte
	var stderr bytes.Buffer
	cmdString := []interface{}{"【CMD】:", item.Name}
	for _, value := range item.Arg {
		cmdString = append(cmdString, value)
	}
	fmt.Println(cmdString...)

	cmd := exec.Command(item.Name, item.Arg...)
	cmd.Stderr = &stderr

	
	if buf, err = cmd.Output();err == nil {
		fmt.Printf("%s\n", buf)
		return
	}

	if item.ErrCallBack == nil {
		log.Error("错误信息输出:", err.Error())
		log.Error(stderr.String())
		return
	}

	log.Infof("err:%s \n%s \n", err.Error(), stderr.String())
	item.ErrCallBack(err, stderr.String())
	return
}
