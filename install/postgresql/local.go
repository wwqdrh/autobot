package postgresql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wwqdrh/autobot/internal/sshcmd/cmd"
)

// 本机安装 master与slave直接执行命令
func localInstall() error {
	if err := cmd.Cmd("docker", "ps"); err != nil {
		return fmt.Errorf("docker ps err: %w", err)
	}
	if err := cmd.Cmd("docker", "image", "inspect", getImageName()); err != nil {
		if err := cmd.Cmd("docker", "pull", "bitnami/postgresql:14"); err != nil {
			return fmt.Errorf("docker pull err: %w", err)
		}
	}

	if ClusterConf.Mode == "master" {
		cmdStr := fmt.Sprintf(`run -d --name %s -e POSTGRESQL_REPLICATION_MODE=master -e POSTGRESQL_REPLICATION_USER=replica -e POSTGRESQL_REPLICATION_PASSWORD=replica -e POSTGRESQL_PASSWORD=%s -p %d:5432 %s`, ClusterConf.ContainerName, ClusterConf.ContainerPassword, ClusterConf.ContainerPort, getImageName())
		return cmd.Cmd("docker", strings.Split(cmdStr, " ")...)
	} else if ClusterConf.Mode == "slave" {
		cmdStr := fmt.Sprintf(`run -d --name %s -e POSTGRESQL_REPLICATION_MODE=master -e POSTGRESQL_REPLICATION_USER=replica -e POSTGRESQL_REPLICATION_PASSWORD=replica -e POSTGRESQL_PASSWORD=%s -p %d:5432 %s`, ClusterConf.ContainerName, ClusterConf.ContainerPassword, ClusterConf.ContainerPort, getImageName())
		return cmd.Cmd("docker", strings.Split(cmdStr, " ")...)
	}
	return errors.New("invalid Cluster Mode")
}
