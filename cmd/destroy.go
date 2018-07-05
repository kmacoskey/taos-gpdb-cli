package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kmacoskey/taos-gpdb-cli/request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy clusters",
	Long:  `Destroy clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		destroyCluster(viper.GetString("destroy_id"))
	},
}

func init() {
	RootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().StringP("id", "i", "", "Destroy specific cluster with matching id")
	viper.BindPFlag("destroy_id", destroyCmd.Flags().Lookup("id"))
}

func destroyCluster(id string) {
	response, body, err := request.HttpClusterRequest("DELETE", fmt.Sprintf("cluster/%s", id), []byte(``))
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
