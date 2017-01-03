package cli

import (
	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/liuchjlu/commiter/etcdclient"
	"github.com/liuchjlu/commiter/ip"
)

func state_to_false(etcdpath, app, component string) {
	log.Infoln("cli.state_to_false():Start state_to_false")
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Errorf("cli.state_to_false() NewEtcdClient fail :%+v\n", err)
	}

	ipinfo, err := ip.Getip()
	if err != nil {
		log.Fatalf("cli.state_to_false():%+v\n", err)
	}

	//log.Infoln("cli.state_to_false() " + "etcd_dir:" + "/commiter/images/" + ipinfo + "/" + app + "/" + component + "/")
	etcd_dir := "/commiter/images/" + ipinfo + "/" + app + "/" + component + "/"
	log.Debugf("cli.state_to_false() "+"etcd_dir:%+v\n", etcd_dir)
	_, errr := client.CreateAbsoluteKey(etcd_dir+"state", "0")
	if errr != nil {
		log.Errorf("cli.state_to_false():%+v\n", errr)
	}

}
