package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func HttpClusterRequest(request_type string, route string, body []byte) (*http.Response, []byte, error) {
	base_url := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("port"))
	url := fmt.Sprintf("%s/%s", base_url, route)

	req, err := http.NewRequest(request_type, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
