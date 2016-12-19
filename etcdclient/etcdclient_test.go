package etcdclient

import (
	"fmt"
	"testing"
)

var etcdpath = "http://192.168.11.51:2379"
var key = "192.168.11.52"

func TestCreateKey(t *testing.T) {
	// test connect
	client, err := NewEtcdClient(etcdpath)
	fmt.Printf("client=%+v,err:=%+v\n", client, err)
	// test create key
	containerid, err := client.CreateKey(key, "abcdef")
	fmt.Printf("containerid=%+v,err:=%+v\n", containerid, err)
	// test create same key
	containerid, err = client.CreateKey(key, "123456")
	fmt.Printf("containerid=%+v,err:=%+v\n", containerid, err)
}

func TestDeleteKey(t *testing.T) {
	// test connect
	client, err := NewEtcdClient(etcdpath)
	fmt.Printf("client=%+v,err:=%+v\n", client, err)
	// test delete key
	err = client.DeleteKey(key)
	fmt.Printf("err:=%+v\n", err)
}

func TestGetKey(t *testing.T) {
	// test connect
	client, err := NewEtcdClient(etcdpath)
	fmt.Printf("client=%+v,err:=%+v\n", client, err)
	// test create key
	containerid, err := client.GetKey(key)
	fmt.Printf("containerid=%+v,err:=%+v\n", containerid, err)
}

func TestCreateDir(t *testing.T) {
	// test connect
	client, err := NewEtcdClient(etcdpath)
	fmt.Printf("client=%+v,err:=%+v\n", client, err)
	// test create dir
	client.CreateDir(AppsPath)
	client.CreateDir(IpsPath)
}

func TestGetDir(t *testing.T) {
	// test connect
	client, err := NewEtcdClient(etcdpath)
	fmt.Printf("client=%+v,err:=%+v\n", client, err)
	// test get all key from dir
	containerids, err := client.GetDir(IpsPath)
	fmt.Printf("containerids=%+v,err:=%+v\n", containerids, err)
}
