package ip

import (
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/liuchjlu/commiter/etcdclient"
)

var (
	BasePath string = "/commiter/ips/"
)

//get ip of Physical machine
func Getip() (string, error) {
	addrsbr0, err := net.InterfaceByName("br0")
	if err != nil {
		log.Fatalf("ip.Getip():%+v\n", err)
		return "", err
	}
	ip_br0, err := addrsbr0.Addrs()
	ip_info := strings.Split(ip_br0[0].String(), "/")
	if err != nil {
		log.Fatalf("ip.Getip():%+v\n", err)
		return ip_info[0], err
	}
	log.Debugf("ip.Getip(): ip=%+v\n", ip_info[0])
	return ip_info[0], nil
}
func GetContainerIp(imagename, containerid, etcdpath string) (string, error) {
	log.Infoln("ip.GetContainerIp(): Strat getting container ip from commiter")
	commiterip, err := Getcontaineripfromcommiter(strings.Split(imagename, "/")[1], etcdpath)
	if commiterip == "" {
		log.Infof("ip.GetContainerIp(). Commiter do not have the ip of this container. Container id:%+v Container image:%+v\n", containerid, imagename)
	} else {
		log.Infof("ip.GetContainerIp(). Successful get container ip. Container ip:%+v Container id:%+v Container image:%+v\n", commiterip, containerid, imagename)
		return commiterip, err
	}
	log.Infoln("Ip.GetContainerIp(): Start getting container ip from assigner")
	assignerip, err := Getcontaineripfromassigner(containerid, etcdpath)
	return assignerip, err
}

//get ip of virtual machine
func Getcontaineripfromassigner(containerid, etcdpath string) (string, error) {
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Errorf("cli.query():%+v\n", err)
		return "", err
	}

	//query from etcd
	ip, err := client.QueryContainerid(containerid)
	if err != nil {
		log.Fatalf("ip.getcontainerip() Failed to  get container ip from assigner: %+v\n", err)
	} else {
		log.Infof("ip.getcontainerip() successful get ip from assigner : %+v\n", ip)
	}

	return ip, err
}
func Getcontaineripfromcommiter(imagename, etcdpath string) (string, error) {
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Errorf("cli.query():%+v\n", err)
		return "", err
	}
	ip, err := client.GetAbsoluteKey(BasePath + imagename)
	//Response, err := client.GetAbsoluteDir("/images/204/databus")
	if err != nil {
		log.Debugf("ip.Getcontaineripfromcommiter()  Failed to  get container ip from commiter: %+v\n", err)
		return "", nil
	}
	if ip != "" {
		return ip, nil
	}

	return "", nil
}

// func (e *Etcd) QueryContainerid(containerid string) (string, error) {
// 	Response, err := e.GetAbsoluteDir(IpsPath)
// 	if err != nil {
// 		return "", err
// 	}

// 	rep, err := regexp.Compile(containerid + ".*")
// 	if err != nil {
// 		log.Errorf("etcdclient.QueryContainerid(): regexp error, err=", err)
// 	}
// 	for _, node := range Response.Node.Nodes {
// 		if rep.MatchString(node.Value) {
// 			paths := strings.Split(node.Key, "/")
// 			return paths[len(paths)-1], nil
// 		}
// 	}
// 	return "", errors.New("no such ip")
// }
