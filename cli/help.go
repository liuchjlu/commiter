package cli

import (
	"fmt"
)

func help() {
	var helpString = `
  Usage: commiter COMMAND [args...]
  Version: 0.2.0
  Author:liuchjlu
  Email:liucaihong@iie.ac.cn
  Commands:
      manage  [ApiserverIp:Port] [FtpserverIp] [etcd path]  Create the config in etcd.
      import  [Localyaml Path]   [Etcd Path]   [0/1/2]      Import the status of pod from yaml. 0:don't commit; 1:commit & push to registry; 2:commit & don't push to registry.
      commit  [Etcd Path]                                   Image commit.
      help
  Use case
      commiter manage 192.168.11.52:8080 192.168.11.52 $ETCDPATH
      commiter import test.yaml $ETCDPATH 1
      commiter commit $ETCDPATH 
  `
	fmt.Println(helpString)
}
