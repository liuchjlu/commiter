package cli

import (
	//        "fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	client "github.com/liuchjlu/commiter/dockerclient"
)

func Commit(etcdpath string, repository string) {

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).Format("20060102")
	log.Debugf("cli.Commit(): date=%+v\n", date)
	Container, err := client.Dockerps()
	if err != nil {
		log.Fatalf("cli.Commit():%+v\n", err)
		return
	}

	Name, err := Getname(etcdpath)
	log.Infof("cli.commit() Getname:%+v\n", Name)
	if err != nil {
		log.Fatalf("cli.Commit():%+v\n", err)
		return
	}
	for _, name := range Name {
		for _, container := range Container {
			if strings.Contains(container.Names[0], name) {
				log.Infof("cli.commit() name=:%+v\n", name)
				idresponse, err := client.Dockercommit(container.ID, repository+name, date)
				if err != nil {
					log.Fatalf("cli.Commit():%+v\n", err)
				} else {
					log.Infof("cli.Commit(): Dockercommit() success,idresponse= %+v\n", idresponse)
				}
				state_to_false(etcdpath, name)
				log.Infof("cli.Commit(): state_to_false()  finish")
			}
		}
	}
}
