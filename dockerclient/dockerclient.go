package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	docker_client "github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// type Docker struct {
// 	client *docker_client.Client
// }

// func NewClient() (*docker_client.Client, error) {
// 	client, err := docker_client.NewClient(docker_client.DefaultDockerHost, "", nil, nil)
// 	return client, err
// }
func Dockerps() ([]types.Container, error) {
	client, err := docker_client.NewClient(docker_client.DefaultDockerHost, "", nil, nil)
	if err != nil {
		log.Fatalf("cli.Dockerps():%+v\n", err)
		return nil, err
	}
	defer client.Close()

	list_options := &types.ContainerListOptions{}

	containers, err := client.ContainerList(context.Background(), *list_options)
	if err != nil {
		log.Fatalf("cli.Dockerps():%+v\n", err)
		return nil, err
	}
	// log.Debugf("cli.Dockerps(): containers=%+v\n", containers)
	// log.Debugf("image====================%+v\n", containers[0].Image)
	log.Infoln("dockerclient.DockerPs(): Successful get ContainerList.")
	return containers, err
}

func Dockercommit(id string, name string, tag string) (types.IDResponse, error) {
	client, err := docker_client.NewClient(docker_client.DefaultDockerHost, "", nil, nil)
	if err != nil {
		log.Fatalf("cli.Dockercommit():%+v\n", err)

	}
	defer client.Close()

	//image name & tag

	reference := name + ":" + tag
	commit_options := &types.ContainerCommitOptions{
		Reference: reference,
	}

	commitresponse, err := client.ContainerCommit(context.Background(), id, *commit_options)
	if err != nil {
		log.Fatalf("cli.Dockercommit():%+v\n", err)
		return commitresponse, err
	}
	log.Debugf("cli.Dockerpcommit(): CommitResponse=%+v\n", commitresponse)
	log.Infof("cli.dockercommit(): dockercommit() finished")
	Dockerpush(name + ":" + tag)

	return commitresponse, err
}

func Dockerpush(imagename string) {
	client, err := docker_client.NewClient(docker_client.DefaultDockerHost, "", nil, nil)

	if err != nil {
		log.Fatalf("cli.Dockerps():%+v\n", err)
		return
	}
	defer client.Close()
	//log.Debugln("#111")
	push_options := types.ImagePushOptions{
		RegistryAuth: "NotValid",
		//PrivilegeFunc: privilegeFunc,
	}

	ioreadcloser, err := client.ImagePush(context.Background(), imagename, push_options)
	// log.Debugln("#222")

	// log.Debugln("#333")
	if err != nil {
		log.Fatalf("cli.Dockerpush():%+v\n", err)
		return
	}
	log.Infof("cli.dockercommit(): dockerpush() finished")
	// log.Debugln("#444")
	ioreadcloser.Close()
	//fmt.Println(containers)
	return

}
