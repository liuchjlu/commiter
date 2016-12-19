package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"golang.org/x/net/context"
	//"strings"
)

func dockerpush() ([]types.Container, error) {
	client, err := dockerclient.NewClient(dockerclient.DefaultDockerHost, "1.19", nil, nil)
	if err != nil {
		log.Fatalf("cli.Dockerps():%+v\n", err)
		return nil, err
	}
	list_options := &types.ContainerListOptions{}
	containers, err := client.ContainerList(context.Background(), *list_options)
	if err != nil {
		log.Fatalf("cli.Dockerps():%+v\n", err)
		return nil, err
	}
	log.Debugf("cli.Dockerps(): containers=%+v\n", containers)
	//fmt.Println(containers)
	return containers, err

}
