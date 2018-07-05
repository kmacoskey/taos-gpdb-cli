package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kmacoskey/taos-gpdb-cli/request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Request creation of a new Cluster",
	Long: `Request provisioning of a new Cluster with the given parameters.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		createCluster(viper.GetString("terraform_config_path"), viper.GetString("timeout"))
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("terraform_config", "c", "", "Path to terraform configuration file")
	createCmd.Flags().StringP("timeout", "t", "", "Cluster timeout duration [h|m|s]")

	viper.BindPFlag("terraform_config_path", createCmd.Flags().Lookup("terraform_config"))
	viper.BindPFlag("timeout", createCmd.Flags().Lookup("timeout"))
}

func createCluster(terraform_config_path string, timeout string) {
	terraform_config, err := ioutil.ReadFile(terraform_config_path)
	if err != nil {
		log.Fatal(err)
	}

	requestClusterStruct := request.ClusterRequest{
		TerraformConfig: string(terraform_config),
		Timeout:         timeout,
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(requestClusterStruct)

	response, body, err := request.HttpClusterRequest("PUT", "cluster", reqBodyBytes.Bytes())
	if err != nil {
		fmt.Println(err)
	}

	var cluster_response_json interface{}

	if response.StatusCode == http.StatusAccepted {
		cluster_response_json = &request.ClusterResponse{}
	} else {
		cluster_response_json = &request.ErrorResponse{}
	}

	err = json.Unmarshal(body, &cluster_response_json)
	data, err := json.Marshal(cluster_response_json)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}
