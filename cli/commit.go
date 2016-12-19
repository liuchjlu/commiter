package cli

import (
	//        "fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

func Commit(etcdpath string, repository string) {
	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).Format("20060102")
	log.Debugf("cli.Commit(): date=%+v\n", date)
	Container, err := dockerps()
	if err != nil {
		log.Fatalf("cli.Commit():%+v\n", err)
		return
	}

	Name, err := Getname(etcdpath)
	if err != nil {
		log.Fatalf("cli.Commit():%+v\n", err)
		return
	}
	for _, name := range Name {
		for _, container := range Container {
			if strings.Contains(container.Names[0], name) {
				dockercommit(container.ID, repository+name, date)
				log.Infof("cli.Commit(): Commit()  finish")
				state_to_false(etcdpath, name)
				log.Infof("cli.Commit(): state_to_false()  finish")
			}
		}
	}
}
