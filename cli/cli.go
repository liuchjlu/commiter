package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/commiter/datatype"
	client "github.com/liuchjlu/commiter/dockerclient"
	"os"
)

var ContainerInfo []datatype.ContainerInfo
var containerinfo datatype.ContainerInfo

func initdata(etcdpath string) {
	psresponse, err := client.Dockerps()
	if err != nil {
		log.Fatalf("cli.init() failed run client.Dockerps(): %+v\n", err)

	}
	log.Debugf("cli.initdata() psresponse:%+v\n", psresponse)
	podstr, err := Getname(etcdpath)
	if err != nil {
		log.Fatalf("cli.init() failed run Getname(): %+v\n", err)
	}
	log.Debugf("cli.initdata() podstr:%+v\n", podstr)

}
func Run() {
	if len(os.Args) == 1 {
		help()
		return
	}

	var err error
	command := os.Args[1]
	log.Debugf("cli.Run(): cli args:%+v\n", os.Args)
	if command == "import" {
		if len(os.Args) != 5 {
			importErr := "the `import` command takes three arguments. See help"
			log.Errorln(importErr)
			return
		}
		filePath := os.Args[2]
		etcdPath := os.Args[3]
		state := os.Args[4]
		err = importapp(filePath, etcdPath, state)
	}
	if command == "commit" {
		if len(os.Args) != 3 {
			importErr := "the `commit` command takes one arguments. See help"
			log.Errorln(importErr)
			return
		}
		Commit(os.Args[2])
	}
	if command == "manage" {
		manage(os.Args[2], os.Args[3])
	}
	if command == "help" {
		help()
	}
	if err != nil {
		log.Errorf("cli.Run():%+v\n", err)
		return
	}
}
