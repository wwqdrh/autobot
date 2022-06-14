package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "自动化部署工具",
	Short:             "自动化部署工具",
	SilenceUsage:      true,
	Long:              `自动化部署工具`,
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	os.Setenv("PATH", os.Getenv("PATH")+":/usr/local/bin")

	rootCmd.AddCommand(PostgresCmd)
}

func GetUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Help()
	}
}
