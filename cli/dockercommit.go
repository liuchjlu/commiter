package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"golang.org/x/net/context"
	//"strings"
)

func dockercommit(id string, name string, tag string) (types.ContainerCommitResponse, error) {
	docker_client, err := dockerclient.NewClient(dockerclient.DefaultDockerHost, "1.19", nil, nil)
	if err != nil {
		log.Fatalf("cli.Dockercommit():%+v\n", err)
	}

	//image name & tag
	reference := name + ":" + tag
	commit_options := &types.ContainerCommitOptions{
		Reference: reference,
	}

	commitresponse, err := docker_client.ContainerCommit(context.Background(), id, *commit_options)
	if err != nil {
		log.Fatalf("cli.Dockercommit():%+v\n", err)
		return commitresponse, err
	}
	log.Debugf("cli.Dockerpcommit(): CommitResponse=%+v\n", commitresponse)
	log.Infof("cli.dockercommit(): dockercommit () finished")
	return commitresponse, err

}
