package cmd

import (
	"github.com/lbernardo/lambda-local/controller"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start local functions lambda",
	Long:  `Start local functions lambda`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		ExecuteCmdStart(cmd, args)
	},
}

var Port string
var Host string
var Yaml string
var Volume string
var Network string
var Environment string

func init() {
	startCmd.PersistentFlags().StringVar(&Port, "port", "3000", "port usage [default 3000]")
	startCmd.PersistentFlags().StringVar(&Host, "host", "0.0.0.0", "host usage [default 0.0.0.0]")
	startCmd.PersistentFlags().StringVar(&Yaml, "yaml", "serverless.yml", "File yaml serverless.yml [default serverless.yml]")
	startCmd.PersistentFlags().StringVar(&Environment, "env", "", "File for using environment variables other than serverless. Can replace serverless variables")
	startCmd.Flags().StringVar(&Volume, "volume", "", "Docker volume mount execution [ex: --volume $PWD]")
	startCmd.Flags().StringVar(&Network, "network", "", "Set network name usage")
	startCmd.MarkFlagRequired("volume")
}

func ExecuteCmdStart(cmd *cobra.Command, args []string) {

	se := controller.Server{
		Host:            Host,
		Port:            Port,
		Yaml:            Yaml,
		Volume:          Volume,
		Network:         Network,
		EnvironmentFile: Environment,
	}
	se.StartServer()
}
