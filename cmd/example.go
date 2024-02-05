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
	"snowflake/comm"
	"snowflake/config"
	"snowflake/debug"
	"snowflake/example"
	"snowflake/fd"
	"snowflake/log"

	"github.com/spf13/cobra"
)

// exampleCmd .
var cityCmd = &cobra.Command{
	Use:   "example",
	Short: "run example",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("event:%s msg:example_exception %v caller_stack:%s",
					comm.EventWarn, r, comm.GetCallerStackLog())
			}
			log.Close()
		}()

		log.Init(config.GetViper().GetBool("example.debug"), config.GetViper().GetString("example.logout"))
		fd.IncreaseFDLimit()

		server := example.NewExampleServer()
		go comm.GoControlListen(server.(comm.ControlService), config.GetViper().GetString("example.controller"))
		go debug.GoPprofListen(config.GetViper().GetString("example.pprof"))
		comm.Run(server)
	},
}

func init() {
	rootCmd.AddCommand(cityCmd)
}
