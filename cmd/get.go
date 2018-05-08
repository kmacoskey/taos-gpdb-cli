package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kmacoskey/taos/handlers"
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
	url := fmt.Sprintf("http://localhost:8080/cluster/%s", id)

	_, body, err := httpClusterRequest("GET", url, []byte(``))
	if err != nil {
		fmt.Println(err)
	}
	cluster_response_json := &handlers.ClusterResponse{}
	err = json.Unmarshal(body, &cluster_response_json)
	data, err := json.Marshal(cluster_response_json)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}

func getClusters() {
	url := "http://localhost:8080/clusters"

	_, body, err := httpClusterRequest("GET", url, []byte(``))
	if err != nil {
		fmt.Println(err)
	}
	cluster_response_json := &handlers.ClustersResponse{}
	err = json.Unmarshal(body, &cluster_response_json)
	data, err := json.Marshal(cluster_response_json)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}

func httpClusterRequest(request_type string, url string, body []byte) (*http.Response, []byte, error) {
	req, err := http.NewRequest(request_type, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}
	req.Close = true

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return response, nil, err
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return response, body, err
	}

	return response, body, nil
}
