package postgresql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wwqdrh/autobot/install/base"
	"github.com/wwqdrh/autobot/internal/sshcmd/sshutil"
)

var sshClient *sshutil.SSH

func remoteInstall() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	sshClient = &sshutil.SSH{
		User:       ClusterConf.SSHUser,
		Password:   ClusterConf.SSHPassword,
		PkFile:     ClusterConf.SSHPKFile,
		PkPassword: ClusterConf.SSHPKPwd,
	}

	// docker环境检查
	result := string(sshClient.Cmd(ClusterConf.Ip, "docker --version"))
	if strings.Index(result, "Docker version") != 0 {
		return errors.New("目标机器未安装docker环境")
	}

	if ClusterConf.Mode == "master" {
		return remoteMaster()
	} else if ClusterConf.Mode == "slave" {
		return remoteSlave()
	} else {
		return errors.New("invalid mode")
	}
}

func remoteMaster() error {
	res := sshClient.CmdToString(ClusterConf.Ip, "docker ps -a | grep "+ClusterConf.ContainerName, " ")
	if res != "" && !ClusterConf.ForceUpdate {
		return errors.New("容器已经存在")
	}

	if err := sshClient.CmdAsync(ClusterConf.Ip, "docker pull "+getImageName()); err != nil {
		return err
	}
	if res == "" {
		return sshClient.CmdAsync(ClusterConf.Ip, fmt.Sprintf(`docker run -d --name %s \
		-e POSTGRESQL_REPLICATION_MODE=master \
		-e POSTGRESQL_REPLICATION_USER=replica \
		-e POSTGRESQL_REPLICATION_PASSWORD=replica \
		-e POSTGRESQL_PASSWORD=%s \
		-p %d:5432 \
		%s`, ClusterConf.ContainerName, ClusterConf.ContainerPassword, ClusterConf.ContainerPort, getImageName()))
	} else {
		if err := base.Install(sshClient, ClusterConf.Ip); err != nil {
			return err
		}
		if err := sshClient.CmdAsync(ClusterConf.Ip, fmt.Sprintf("docker exec -it %s /bin/bash -c 'rm -rf /usr/local/autobot && mkdir -p /usr/local/autobot'", ClusterConf.ContainerName)); err != nil {
			return err
		}
		if err := sshClient.CmdAsync(ClusterConf.Ip, fmt.Sprintf("docker cp /usr/local/autobot/base %s:/usr/local/autobot/base", ClusterConf.ContainerName)); err != nil {
			return err
		}
		if err := sshClient.CmdAsync(ClusterConf.Ip, fmt.Sprintf("docker cp /usr/local/autobot/postgresql %s:/usr/local/autobot/postgresql", ClusterConf.ContainerName)); err != nil {
			return err
		}
		if err := sshClient.CmdAsync(ClusterConf.Ip, fmt.Sprintf("docker exec -it %s /bin/bash -c 'bash -x /usr/local/autobot/postgresql/master.sh'", ClusterConf.ContainerName)); err != nil {
			return err
		}
		return nil
	}
}

func remoteSlave() error {
	if ClusterConf.MasterIP == "" || ClusterConf.MasterPort == 0 {
		return errors.New("未设置master的ip以及端口")
	}

	if err := sshClient.CmdAsync(ClusterConf.Ip, "docker pull "+getImageName()); err != nil {
		return err
	}
	return sshClient.CmdAsync(ClusterConf.Ip, fmt.Sprintf(`docker run -d --name %s \
	-e POSTGRESQL_REPLICATION_MODE=slave \
	-e POSTGRESQL_MASTER_HOST=%s \
	-e POSTGRESQL_MASTER_PORT_NUMBER=%d \
	-e POSTGRESQL_REPLICATION_USER=replica \
	-e POSTGRESQL_REPLICATION_PASSWORD=replica \
	-e POSTGRESQL_PASSWORD=%s \
	-p %d:5432 \
	%s`, ClusterConf.ContainerName, ClusterConf.MasterIP, ClusterConf.MasterPort, ClusterConf.ContainerPassword, ClusterConf.ContainerPort, getImageName()))
}
