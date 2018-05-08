package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("terraform_config", "c", "", "Path to terraform configuration file")
	createCmd.Flags().StringP("timeout", "t", "", "Cluster timeout duration [h|m|s]")

	viper.BindPFlag("terraform_config", createCmd.Flags().Lookup("terraform_config"))
	viper.BindPFlag("timeout", createCmd.Flags().Lookup("timeout"))
}
