package cli

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/liuchjlu/commiter/etcdclient"
	"github.com/liuchjlu/commiter/ip"
)

type ImportInfo struct {
	App        string
	Component  string
	Repository string
}

func Getname(etcdpath string) ([]ImportInfo, error) {
	var InfoOne ImportInfo
	var ImportInfoList []ImportInfo
	var flag bool = false
	var str []string
	// var img_basename []string
	ipinfo, err := ip.Getip()
	// ipstr := strings.Split(ipinfo, ".")
	// //log.Infof("IPstr============%+v\n", ipstr[0]+ipstr[1]+ipstr[2]+ipstr[3])
	// ip := ipstr[0] + ipstr[1] + ipstr[2] + ipstr[3]
	if err != nil {
		log.Fatalf("cli.Getname():%+v\n", err)
	}
	etcd_client, err := etcdclient.NewEtcdClient(etcdpath)

	if err != nil {
		log.Fatalf("cli.Getname():%+v\n", err)
	}
	absoluteDir := "/commiter/images/" + ipinfo
	Response, _ := etcd_client.GetAbsoluteDir(absoluteDir)
	for _, node_0 := range Response.Node.Nodes {
		// log.Infoln("node_0.Key:", node_0.Key)
		// log.Infoln("node_0.Nodes:", node_0.Nodes)
		for _, node_1 := range node_0.Nodes {
			// log.Infoln("node_1:", node_1.Key)
			// log.Infof("node_1.Nodes:%+v\n", node_1.Nodes)
			for _, node_2 := range node_1.Nodes {
				// log.Infoln("node_2.Key:", node_2.Key)
				// log.Infoln("node_2.Value:", node_2.Value)
				Key := strings.Split(node_2.Key, "/")[6]
				if Key == "state" {
					if node_2.Value == "1" {
						flag = true
						str = strings.Split(node_2.Key, "/")
						InfoOne.App = str[4]
						InfoOne.Component = str[5]
						//app_component = append(app_component, node_2.Key)
						log.Debugf("cli.Getname() str:%+v\n", str)
					}
				} else if Key == "repository" {
					InfoOne.Repository = node_2.Value
					log.Debugf("cli.Getname() Repository:%+v\n", node_2.Value)
				}
			}
			if flag {
				ImportInfoList = append(ImportInfoList, InfoOne)
				log.Infof("cli.Getname() infoOne:%+v\n", InfoOne)
				log.Infof("cli.Getname() ImportInfoList:%+v\n", ImportInfoList)
				flag = false
			}
		}
	}
	return ImportInfoList, err
}
