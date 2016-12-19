package ip

import (
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func Getip()(string ,error) {
    addrsbr0, err := net.InterfaceByName("br0")
    if err != nil {
      log.Fatalf("ip.Getip():%+v\n", err)
      return "",err
    }
    ip_br0,err := addrsbr0.Addrs()
    ip_info := strings.Split(ip_br0[0].String(), "/")
    if err != nil {
      log.Fatalf("ip.Getip():%+v\n", err)
      return ip_info[0],err
    }
    log.Debugf("ip.Getip(): ip=%+v\n",ip_info[0])
    return ip_info[0],nil
}
