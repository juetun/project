package utils

import (
	"bytes"
	"fmt"
	"github.com/prometheus/common/log"
	"os/exec"
)
type CmdObject struct {
	Name string   `json:"name"`
	Arg  []string `json:"arg"`
}
func  ExeCMD(item *CmdObject) (err error) {
	var buf []byte
	var stderr bytes.Buffer
	cmdString := []interface{}{"【CMD】:", item.Name}
	for _, value := range item.Arg {
		cmdString = append(cmdString, value)
	}
	fmt.Println(cmdString...)
	cmd := exec.Command(item.Name, item.Arg...)
	cmd.Stderr = &stderr
	buf, err = cmd.Output()
	if err != nil {
		log.Error("错误信息输出", err.Error())
		log.Error(stderr.String())
	}

	fmt.Printf("%s\n", buf)
	return
}

