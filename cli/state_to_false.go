package cli
import(
        "strings"

		"github.com/liuchjlu/pointcommit/ip"
        etcdclient "github.com/liuchjlu/pointcommit/etcdclient"
        log "github.com/Sirupsen/logrus"
)
func state_to_false(etcdpath string,name string){
        log.Infoln("cli.state_to_false():Start state_to_false")
        client, err := etcdclient.NewEtcdClient(etcdpath)
        if err != nil {
            log.Errorf("cli.state_to_false() NewEtcdClient fail :%+v\n", err)
        }
        name_info := strings.Split(name, "-")

        ipinfo,err := ip.Getip()
	if err != nil {
		log.Fatalf("cli.state_to_false():%+v\n", err)
	}
	component := strings.Join(name_info[1:],"")
	log.Infoln("etcd_dir:"+"/images/"+ipinfo+"/"+name_info[0]+"/"+component+"/")
        etcd_dir := "/images/"+ipinfo+"/"+name_info[0]+"/"+component+"/"
        log.Debugf("etcd_dir:%+v\n",etcd_dir)
        _, errr := client.CreateAbsoluteKey(etcd_dir+"state", "0")
        if errr != nil {
            log.Errorf("cli.state_to_false():%+v\n", errr)
        }

}

