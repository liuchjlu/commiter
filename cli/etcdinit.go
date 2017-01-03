package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/commiter/etcdclient"
)

func manage(apiserver, etcdpath string) {
	log.Infoln("cli.manage() start manage.")
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Fatalf("cli.manage():%+v\n", err)

	}
	//create dir on etcd
	client.CreateAbsoluteDir("/commiter/images")
	client.CreateAbsoluteDir("/commiter/ips")
	client.CreateAbsoluteDir("/commiter/config")

	//create config on etcd
	_, err = client.CreateAbsoluteKey("/commiter/config/apiserver", apiserver)
	if err != nil {
		log.Fatalf("cli.manage() client.createabsolutekey err:", err)
	}
	log.Infoln("cli.namage(): Manage success.")
}
