package etcdclient

import (
	"errors"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type Etcd struct {
	client client.KeysAPI
}

type App struct {
	App       string `App`
	Component string `Component`
	Ips       []Ip   `Ips`
}

type Ip struct {
	Ip      string `Ip`
	Gateway string `Gateway`
}

var ServcieTimeout = 50 * time.Second
var BasePath = "/images"
var AppsPath = BasePath + "/apps/"
var IpsPath = BasePath + "/ips/"
var ConfigPath = BasePath + "/config"

func NewEtcdClient(etcdpath string) (*Etcd, error) {
	endpoints := strings.Split(etcdpath, ",")
	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavctxailable
		HeaderTimeoutPerRequest: ServcieTimeout,
	}
	c, err := client.New(cfg)
	kapi := client.NewKeysAPI(c)

	etcd := new(Etcd)
	etcd.client = kapi
	return etcd, err
}

func (e *Etcd) CreateKey(key, value string) (string, error) {
	// set PrevExist PrevNoExist, get more help from core/etcd/client/keys.go
	opt := &client.SetOptions{
		PrevExist: client.PrevNoExist,
	}
	log.Debugf("etcdclient.CreateKey():key=%+v,value=%+v", IpsPath+key, value)
	Response, err := e.client.Set(context.Background(), IpsPath+key, value, opt)
	if err != nil {
		return "", err
	}
	return Response.Node.Value, nil
}

func (e *Etcd) CreateAbsoluteKey(absoluteKey, value string) (string, error) {
	// this func is different with CreateKey ,this use absolute path and will update the value
	log.Debugf("etcdclient.CreateAbsoluteKey():key=%+v,value=%+v", absoluteKey, value)
	Response, err := e.client.Set(context.Background(), absoluteKey, value, nil)
	if err != nil {
		return "", err
	}
	return Response.Node.Value, nil
}

func (e *Etcd) DeleteKey(key string) error {
	return e.DeleteAbsoluteKey(IpsPath + key)
}

func (e *Etcd) DeleteAbsoluteKey(absoluteKey string) error {
	_, err := e.client.Delete(context.Background(), absoluteKey, nil)
	return err
}

func (e *Etcd) GetKey(key string) (string, error) {
	return e.GetAbsoluteKey(IpsPath + key)
}

func (e *Etcd) GetAbsoluteKey(absoluteKey string) (string, error) {
	Response, err := e.client.Get(context.Background(), absoluteKey, nil)
	if err != nil {
		return "", err
	}
	return Response.Node.Value, nil
}

func (e *Etcd) CreateAbsoluteDir(absoluteDir string) error {
	// set dir true
	opt := &client.SetOptions{
		Dir: true,
	}
	log.Debugf("etcdclient.CreateAbsoluteDir():dir=%+v", absoluteDir)
	_, err := e.client.Set(context.Background(), absoluteDir, "", opt)
	return err
}

func (e *Etcd) GetAbsoluteDir(absoluteDir string) (*client.Response, error) {
	// set recursive true
	goption := new(client.GetOptions)
	goption.Recursive = true

	Response, err := e.client.Get(context.Background(), absoluteDir, goption)
	if err != nil {
		return nil, err
	}
	return Response, nil
}

func (e *Etcd) GetAbsoluteDirIps(absoluteDir string) ([]Ip, error) {
	Response, err := e.GetAbsoluteDir(absoluteDir)
	if err != nil {
		return nil, err
	}
	// unmarshal response
	ips := make([]Ip, 0)

	for _, node := range Response.Node.Nodes {
		paths := strings.Split(node.Key, "/")
		ip := &Ip{
			Ip:      paths[len(paths)-1],
			Gateway: node.Value,
		}
		ips = append(ips, *ip)
	}
	return ips, err
}

func (e *Etcd) QueryContainerid(containerid string) (string, error) {
	Response, err := e.GetAbsoluteDir(IpsPath)
	if err != nil {
		return "", err
	}

	rep, err := regexp.Compile(containerid + ".*")
	if err != nil {
		log.Errorf("etcdclient.QueryContainerid(): regexp error, err=", err)
	}
	for _, node := range Response.Node.Nodes {
		if rep.MatchString(node.Value) {
			paths := strings.Split(node.Key, "/")
			return paths[len(paths)-1], nil
		}
	}
	return "", errors.New("no such ip")
}
