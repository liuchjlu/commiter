package cli

import (
	"fmt"
)

func help() {
	var helpString = `
  Usage: assigner COMMAND [args...]
  Version: 0.2.0
  Author:liuchjlu
  Email:liucaihong@iie.ac.cn

  Commands:
      import  [localyaml path] [etcd path]  [0/1/2]     Import the status of pod from yaml. 0:don't commit; 1:commit & push to registry; 2:commit & don't push to registry.
      commit  [etcd path] [images repository]           Image commit.
      help
  Use case
      imgcommit import test.yaml $ETCDPATH 1
      imgcommit commit $ETCDPATH  "192.168.11.51:5000/"
  `
	fmt.Println(helpString)
}
