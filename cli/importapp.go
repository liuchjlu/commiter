package cli

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/commiter/etcdclient"
)

func importapp(filepath string, etcdpath string, state string) error {
	log.Infoln("cli.importapp():Start importapp")
	// yaml unmarshal
	app, err := UnmarshalConfig(filepath)
	if err != nil {
		log.Fatalf("cli.importapp():%+v\n", err)
		return err
	}
	log.Debugf("cli.importapp(): app=%+v\n", app)

	// connect to etcd
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Fatalf("cli.importapp():%+v\n", err)
		return err
	}

	//create ip dir to etcd
	path := ""
	for _, ip := range app.Ips {
		path = etcdclient.BasePath + "/" + ip.Ip + "/"

		err = client.CreateAbsoluteDir(path)
		path = path + app.App + "/"
		err = client.CreateAbsoluteDir(path)
		path = path + app.Component + "/"
		err = client.CreateAbsoluteDir(path)
		_, err = client.CreateAbsoluteKey(path+"state", state)
		if err != nil {
			log.Errorf("cli.importapp():%+v\n", err)
		}
	}

	log.Infoln("cli.importapp(): Importapp success")
	return nil
}

func UnmarshalConfig(path string) (*etcdclient.App, error) {
	//import.yaml  path
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	app := new(etcdclient.App)     //
	err = yaml.Unmarshal(in, &app) //
	if err != nil {
		return nil, err
	}
	return app, nil
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
