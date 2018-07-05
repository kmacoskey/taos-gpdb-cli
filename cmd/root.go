package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var host string
var port string

var RootCmd = &cobra.Command{
	Use:   "taos-gpdb-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.taos.yaml)")
	RootCmd.PersistentFlags().StringVarP(&host, "API", "a", "", "Taos API to connect with")
	RootCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "Port to use when connecting to the Taos API")
}

func initConfig() {

	if host != "" && port != "" {
		viper.Set("Port", port)
		viper.Set("Host", host)
	} else {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			viper.AddConfigPath(home)
			viper.SetConfigName(".taos.yml")
		}

		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
		}
	}

}
