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

package main

import (
	"fmt"

	"snowflake/cmd"

	"github.com/spf13/cobra"
)

var version = "unknown"
var gitCommit = "unknown"
var buildDate = "unknown"
var goVersion = "unknown"
var env = "debug"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version:", version)
		fmt.Println("buildDate:", buildDate)
		fmt.Println("gitCommit:", gitCommit)
		fmt.Println("goVersion:", goVersion)
		if env != "prod" {
			fmt.Println("env:", env)
		}
	},
}

func init() {
	cmd.AddCommand(versionCmd)
	if env != "prod" {
		fmt.Println("env:", env)
		cmd.AddCommand(cmd.TestCmd)
	}
}

func main() {
	cmd.Execute()
}
