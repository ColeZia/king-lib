/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos/v2"
	"github.com/spf13/cobra"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/service/appinfo"
	"gl.king.im/king-lib/framework/service/cmd/migrate"
)

var (
	confpath    string
	showVersion bool
	srvInfo     *ServiceInfo
)

type Cmd struct {
	Cmd  *cobra.Command
	Conf interface{}
}

type ServiceInfo struct {
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	Id               string
	BootConf         interface{}
	MigrateBc        interface{}
	MigrateRunFun    migrate.MigrateRunFun
	SeedRunFun       migrate.SeedRunFun
	Cmds             []*Cmd
	InitAppWrapper   func(klog.Logger) (*kratos.App, func(), error)
	InitAppWrapperV2 func(*service.ServiceBootData) (*kratos.App, func(), error)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "boss-service",
	Short: "BOSS微服务命令行工具",
	Long:  `BOSS微服务命令行工具，包括微服务启动、数据迁移等命令`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(srvInfo.Version)
			return
		}
		runService()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(si *ServiceInfo) {
	srvInfo = si

	appinfo.AppInfoIns = appinfo.AppInfo{
		Id:        srvInfo.Id,
		Framework: "kratos",
		Name:      srvInfo.Name,
		Version:   srvInfo.Version,
	}

	service.AppInfoIns = service.AppInfo{
		Id:        srvInfo.Id,
		Framework: "kratos",
		Name:      srvInfo.Name,
		Version:   srvInfo.Version,
	}

	//加载配置文件--在这里执行的话，flag还没有被scan
	//scanCnf()

	migrate.InitCmd(rootCmd, srvInfo.MigrateRunFun, srvInfo.SeedRunFun, srvInfo.MigrateBc)
	for _, v := range si.Cmds {
		rootCmd.AddCommand(v.Cmd)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.boss-cmd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.PersistentFlags().StringVar(&confpath, "conf", "", "配置文件")
	rootCmd.PersistentFlags().StringVarP(&confpath, "conf", "c", "", "配置文件")
	//rootCmd.PersistentFlags().BoolVar(&showVersion, "version", false, "显示版本号")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "显示版本号")
}
