package sandbox

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/pkg/configutil"

	"github.com/flyteorg/flytectl/pkg/docker"

	"github.com/docker/docker/api/types"
	"github.com/enescakir/emoji"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/k8s"
)

const (
	teardownShort = "Cleans up the sandbox environment"
	teardownLong  = `
Removes the Sandbox cluster and all the Flyte config created by 'sandbox start':
::

 flytectl sandbox teardown 
	

Usage
`
)

func teardownSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	return tearDownSandbox(ctx, cli)
}

func tearDownSandbox(ctx context.Context, cli docker.Docker) error {
	c, err := docker.GetSandbox(ctx, cli)
	if err != nil {
		return err
	}
	if c != nil {
		if err := cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			return err
		}
	}
	if err := configutil.ConfigCleanup(); err != nil {
		fmt.Printf("Config cleanup failed. Which Failed due to %v \n ", err)
	}
	if err := removeSandboxKubeContext(); err != nil {
		fmt.Printf("Kubecontext cleanup failed. Which Failed due to %v \n ", err)
	}
	fmt.Printf("%v %v Sandbox cluster is removed successfully. \n", emoji.Broom, emoji.Broom)
	return nil
}

func removeSandboxKubeContext() error {
	k8sCtxMgr := k8s.NewK8sContextManager()
	return k8sCtxMgr.RemoveContext(sandboxContextName)
}
