package cli

import (
	//"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"strings"
	"time"
	//"github.com/liuchjlu/commiter/datatype"
	dockerclient "github.com/liuchjlu/commiter/dockerclient"
	etcdclient "github.com/liuchjlu/commiter/etcdclient"
	"github.com/liuchjlu/commiter/ip"
)

func Commit(etcdpath string) {

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).Format("20060102")
	log.Debugf("cli.Commit(): date=%+v\n", date)
	Container, err := dockerclient.Dockerps()
	var ContainerList []types.Container
	for _, container := range Container {
		if !strings.Contains(container.Names[0], "k8s_POD") {
			ContainerList = append(ContainerList, container)
		}
	}
	log.Infof("cli.Commit() ContainerList:%+v\n", ContainerList)
	//log.Debugf("#### Container.size:%+v, ContainerList.size:%+v\n", len(Container), len(ContainerList))
	if err != nil {
		log.Fatalf("cli.Commit():%+v\n", err)
		return
	}
	infolist, err := Getname(etcdpath)
	if err != nil {
		log.Fatalf("cli.Commit():%+v\n", err)
		return
	}
	log.Infof("cli.commit() ImportInfo:%+v\n", infolist)

	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Fatalf("cli.commit() etcdclient.NewEtcdClient faild:%+v\n", err)
	}
	var CommitCount int = 0
	for _, info := range infolist {
		for _, container := range ContainerList {
			if strings.Contains(container.Names[0], info.App+"rc-"+info.Component) {
				log.Infof("cli.commit() importinfo=:%+v\n", info)

				//get ip of container
				containerip, err := ip.GetContainerIp(container.Image, container.ID, etcdpath)
				if err != nil {
					log.Fatalf("cli.commit(). ip.getcontaineripfromassigner:%+v\n", err)
				}
				log.Debugf("cli.commit(). ip.getcontainerip:%+v\n", containerip)
				imagename := info.Repository + "/" + info.App + "-" + info.Component + "-" + containerip
				log.Infof("cli.commit() dockerclient.Dockercommit imagename:%+v\n", imagename)

				//commit repository/app-componnet-ip:date
				idresponse, err := dockerclient.Dockercommit(container.ID, imagename, date)
				if err != nil {
					log.Fatalf("cli.Commit():%+v\n", err)
				} else {
					log.Infof("cli.Commit(): Dockercommit() success,idresponse= %+v\n", idresponse)
				}

				//change the commit state of etcd
				state_to_false(etcdpath, info.App, info.Component)
				log.Infof("cli.Commit(): state_to_false()  finish")

				//set <imagename:ip>  to etcd path:/commiter/ips
				createips, err := client.CreateAbsoluteKey("/commiter/ips/"+strings.Split(imagename, "/")[1]+":"+date, containerip)
				if err != nil {
					log.Errorf("cli.commit() failed to create the <imagename:ip> in etcd:%+v\n", err)
				}
				log.Infof("cli.commit(). Set <imagename:ip> in etcd successful:%+v\n", createips)

				//Rc(app, component, containerip, containername, image, etcdpath string)
				//create podrc&upload ftpserver:/commiter/app-component
				Rc(info.App, info.Component, containerip, container.Names[0], imagename+":"+date, etcdpath)
				CommitCount += 1
			}
		}
	}
	log.Infof("cli.commit() CommitCount=%+v\n", CommitCount)
}
