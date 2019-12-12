/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/lbernardo/lambda-local/controller"
	"github.com/spf13/cobra"
)

var YamlConfig string

// endpointsCmd represents the endpoints command
var endpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Short: "List endpoints API Gateway",
	Long:  `List endpoints API Gateway`,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteCmdEndpoints(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(endpointsCmd)
	endpointsCmd.PersistentFlags().StringVar(&YamlConfig, "yaml", "serverless.yml", "File yaml serverless.yml [default serverless.yml]")
}

func ExecuteCmdEndpoints(cmd *cobra.Command, args []string) {
	endController := controller.NewEndpointsController(YamlConfig, "127.0.0.1", "3000")
	endController.ListEndpoints()
}
