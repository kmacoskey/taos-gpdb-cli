package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kmacoskey/taos-gpdb-cli/request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve clusters",
	Long:  `Retrieve clusters from the taos server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("all") {
			getClusters()
		} else {
			getCluster(viper.GetString("id"))
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("id", "i", "", "Get specific cluster with matching id")
	viper.BindPFlag("id", getCmd.Flags().Lookup("id"))
	getCmd.Flags().BoolP("all", "a", false, "Get all clusters")
	viper.BindPFlag("all", getCmd.Flags().Lookup("all"))
}

func getCluster(id string) {
	response, body, err := request.HttpClusterRequest("GET", fmt.Sprintf("cluster/%s", id), []byte(``))
	if err != nil {
		fmt.Println(err)
	}

	var cluster_response_json interface{}

	if response.StatusCode == http.StatusOK {
		cluster_response_json = &request.ClusterResponse{}
	} else {
		cluster_response_json = &request.ErrorResponse{}
	}

	err = json.Unmarshal(body, &cluster_response_json)
	if err != nil {
		fmt.Println(err)
	}

	data, err := json.Marshal(cluster_response_json)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}

func getClusters() {
	response, body, err := request.HttpClusterRequest("GET", "clusters", []byte(``))
	if err != nil {
		fmt.Println(err)
	}

	var cluster_response_json interface{}

	if response.StatusCode == http.StatusOK {
		cluster_response_json = &request.ClustersResponse{}
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
