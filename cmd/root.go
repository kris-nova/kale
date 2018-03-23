// Copyright © 2018 The Kubicorn Authors
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
	"fmt"
	"os"

	"github.com/kris-nova/kale/rtmp"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kale",
	Short: "OBS Bouncer for TGIK",
	Long:  `Used to launch a server that will bounce OBS streams based on arbitrary configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Starting OBS Forwarder")
		err := rtmp.ListenAndServe(o)
		if err != nil {
			logger.Critical("Epic failure: %v", err)
			os.Exit(99)
		}
		os.Exit(1)
	},
}

var o = rtmp.NewObsOptions()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 4, "Log level")

}
