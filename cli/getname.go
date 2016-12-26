package cli

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/liuchjlu/commiter/etcdclient"
	"github.com/liuchjlu/commiter/ip"
)

var App []string
var Component []string

func Getname(etcdpath string) ([]string, error) {
	var app_component []string
	var img_basename []string
	ipinfo, err := ip.Getip()
	ipstr := strings.Split(ipinfo, ".")
	//log.Infof("IPstr============%+v\n", ipstr[0]+ipstr[1]+ipstr[2]+ipstr[3])
	ip := ipstr[0] + ipstr[1] + ipstr[2] + ipstr[3]
	if err != nil {
		log.Fatalf("cli.Getname():%+v\n", err)
	}
	etcd_client, err := etcdclient.NewEtcdClient(etcdpath)

	if err != nil {
		log.Fatalf("cli.Getname():%+v\n", err)
	}
	absoluteDir := "/images/" + ipinfo
	Response, _ := etcd_client.GetAbsoluteDir(absoluteDir)
	for _, node_0 := range Response.Node.Nodes {
		log.Infoln("node_0.Key:", node_0.Key)
		log.Infoln("node_0.Nodes:", node_0.Nodes)
		for _, node_1 := range node_0.Nodes {
			log.Infoln("node_1:", node_1.Key)
			for _, node_2 := range node_1.Nodes {
				log.Infoln("node_2.Key:", node_2.Key)
				log.Infoln("node_2.Value:", node_2.Value)
				if node_2.Value == "1" {
					app_component = append(app_component, node_2.Key)
				}
			}
		}
		for i, name := range app_component {
			s1 := strings.Split(name, "/")
			img_basename = append(img_basename, s1[3]+"-"+s1[4]+"-"+ip)
			App = append(App, s1[3])
			Component = append(Component, s1[4])
			log.Debugf("cli.getname(): img_basename=%+v\n", img_basename[i])
		}
	}
	return img_basename, err
}
