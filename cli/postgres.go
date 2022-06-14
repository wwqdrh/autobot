// mode: default master slave
package cli

import (
	"github.com/wwqdrh/autobot/install/postgresql"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/logger"
)

var (
	// flag
	PostgresCmd = &cobra.Command{
		Use:          "postgres",
		Short:        "install postgres",
		Example:      "...",
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := postgresql.Install()
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			return err
		},
	}
)

func init() {
	// ssh
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.Ip, "ip", "", "目标地址")
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.SSHUser, "username", "", "登录用户名")
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.SSHPassword, "password", "", "登录密码")
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.SSHPKFile, "pk", GetUserHomeDir()+"/.ssh/id_rsa", "private key for ssh")
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.SSHPKPwd, "pk-passwd", "", "private key password for ssh")
	// basic
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.Mode, "mode", "default", "部署的应用类型: default、master、slave")
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.ContainerName, "name", "", "部署的容器名")
	PostgresCmd.Flags().IntVar(&postgresql.ClusterConf.ContainerPort, "dbport", 5432, "容器的端口")
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.ContainerPassword, "dbpassword", "", "数据库的密码")
	PostgresCmd.Flags().BoolVar(&postgresql.ClusterConf.ForceUpdate, "update", false, "如果指定容器已经存在是否在已有容器上进行master初始化")
	PostgresCmd.Flags().IntVar(&postgresql.ClusterConf.Version, "version", 14, "主从数据库版本，必须一致")
	// slave db
	PostgresCmd.Flags().StringVar(&postgresql.ClusterConf.MasterIP, "master_ip", "", "主数据库所在的ip")
	PostgresCmd.Flags().IntVar(&postgresql.ClusterConf.MasterPort, "master_port", -1, "主数据库所在的端口")
}
