package postgresql

import "fmt"

var ClusterConf = new(ClusterConfig)

type ClusterConfig struct {
	Ip                string // if 127.0.0.1, will skip ssh
	Port              int
	Mode              string // master slave
	SSHUser           string
	SSHPassword       string
	SSHPKFile         string // ssh秘钥文件路径
	SSHPKPwd          string // ssh秘钥文件密码
	ContainerName     string // 部署的容器名字
	ContainerPort     int    // 暴露的容器端口
	ContainerPassword string // 容器的root用户登录密码
	ForceUpdate       bool   // TODO: 现在只能初始化一次 否则会卡住 是否存在 如果设置为已经存在则针对已有容器进行更新 现在只针对postgres:14的master初始化
	Version           int    // 主从的版本必须一致
	// slave
	MasterIP   string // if 127.0.0.1, will skip ssh
	MasterPort int
}

func getImageName() string {
	return fmt.Sprintf("bitnami/postgresql:%d", ClusterConf.Version)
}

func Install() error {
	if ClusterConf.Ip == "127.0.0.1" {
		return localInstall()
	} else {
		return remoteInstall()
	}
}
