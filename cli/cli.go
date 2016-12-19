package cli

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

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
		if len(os.Args) != 4 {
			importErr := "the `commit` command takes two arguments. See help"
			log.Errorln(importErr)
			return
		}
		Commit(os.Args[2], os.Args[3])
	}
	if command == "help" {
		help()
	}
	if err != nil {
		log.Errorf("cli.Run():%+v\n", err)
		return
	}
}
