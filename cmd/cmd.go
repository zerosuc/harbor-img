package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"harbor-img-clear/clean"
	"k8s.io/klog"
)

var (
	url         string
	user        string
	password    string
	projectName string
	keepNum     int
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		klog.Fatalf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}
}

func NewHarborCleanCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "harbor-img-clear",
		Short:   "harbor 仓库镜像清理",
		Long:    "harbor-img-clear 用于清理harbor的仓库中的tag，以释放存储资源",
		Example: "./harbor-img-clear  --address http://10.200.82.51  --user admin --password Harbor12345 --project appsvc  --keepNum 30",
		Run: func(cmd *cobra.Command, args []string) {
			// 从Viper获取配置
			if url == "" {
				url = viper.GetString("harbor.url")
			}
			if user == "" {
				user = viper.GetString("harbor.username")
			}
			if password == "" {
				password = viper.GetString("harbor.password")
			}
			if keepNum == 0 {
				keepNum = viper.GetInt("harbor.num")
			}
			if projectName == "" {
				projectName = viper.GetString("harbor.project")
			}

			//klog.Infof("URL: %s, User: %s, Password: %s, Project: %s, KeepNum: %d\n",
			//	url, user, password,
			//	projectName, keepNum)
			// 检查必填参数
			if url == "" || user == "" || password == "" || projectName == "" {
				cmd.Help()
				os.Exit(1)
			}

			clean.Clean(url, user, password, projectName, keepNum)
		},
	}
	// 绑定命令行参数和Viper
	rootCmd.Flags().StringVarP(&url, "address", "a", "", "例如：https://harbor.abc.com")
	rootCmd.Flags().StringVarP(&user, "user", "u", "", "用户名，例如：admin")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "密码")
	rootCmd.Flags().StringVarP(&projectName, "project", "P", "", "项目名，all表示所有项目")
	rootCmd.Flags().IntVarP(&keepNum, "keepNum", "k", 0, "保留的tag数目，例如50")

	return rootCmd
}
