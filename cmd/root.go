// Copyright © 2021 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"snowflake/comm"
	"snowflake/config"
	"snowflake/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugRun bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snowflake",
	Short: "snowflake",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// AddCommand .
func AddCommand(cmds ...*cobra.Command) {
	rootCmd.AddCommand(cmds...)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./main.cfg)")
	rootCmd.PersistentFlags().BoolVarP(&debugRun, "debug", "d", false, "debug run mode")

	go func() { // 直到进程结束
		for {
			time.Sleep(time.Minute * 5)
			log.Info(comm.GetProcessStatusLog())
		}
	}()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("json")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("main.cfg")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())

		v := viper.New()
		data, err := os.ReadFile(viper.ConfigFileUsed())
		if err != nil {
			fmt.Println("read config file err:", err)
			os.Exit(0)
		}
		v.SetConfigType("json")
		v.ReadConfig(bytes.NewBuffer(data))
		v.Set("debug", debugRun)
		config.ReplaceViper(v) // 全新只读实例

		// 监控配置更新
		// 配置文件保存操作会触发回调
		// 注意：ubuntu下软链接的配置文件无法自动触发热更新
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Infof("Config file changed: %s", e.Name)
			// OnConfigChange可能有个bug, 修改配置文件后, 会触发两次。
			// 2次调用对本程序运行无影响，不特殊处理
			config.ReloadConfig("json")
		})
		viper.WatchConfig()
	} else {
		var skip bool
		if len(os.Args) >= 2 {
			if os.Args[1] == "version" || os.Args[1] == "help" {
				skip = true
			}
		}
		if !skip {
			fmt.Println("config file err:", err)
			os.Exit(1)
		}
	}
}
