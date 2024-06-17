package cmd

import (
	"fmt"
	"harbor-img/version"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
	vers    bool
)

func NewHarborCmd() *cobra.Command {
	rootCmd = &cobra.Command{
		Use:     "harbor-img",
		Short:   "harbor-img 仓库镜像清理",
		Long:    "harbor-img 用于清理harbor的仓库中的tag,以释放存储资源",
		Example: "./harbor-img clear --address http://10.200.82.51  --user admin --password Harbor12345 --project appsvc  --keepNum 30",
		Run: func(cmd *cobra.Command, args []string) {
			if vers {
				fmt.Print(version.FullVersion())
			} else {
				cmd.Help()
			}
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, " harbor-img 当前版本")

	rootCmd.AddCommand(clearCommand())
	return rootCmd
}
