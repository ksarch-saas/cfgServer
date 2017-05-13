package utils

import (
	"fmt"
	"bytes"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const namingUri = "http://bns.noah.baidu.com"

type NamingResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetAddrFromBNS(service string) (string, error) {
	resp, err := http.Get(namingUri + "/webfoot/index.php?r=webfoot/ApiInstanceInfo&serviceName=" + service)
	if resp.StatusCode != 200 {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var namingResp NamingResponse

	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	d.Decode(&namingResp)
	if !namingResp.Success {
		return namingResp.Message, fmt.Errorf("resolve addr service failed")
	}

	serviceAddrs := []string{}

	for _, v := range namingResp.Data.([]interface{}) {
		if v.(map[string]interface{})["status"] != "0" {
			continue
		}
		serviceAddr := fmt.Sprintf("%s:%s", v.(map[string]interface{})["hostName"], v.(map[string]interface{})["port"])
		serviceAddrs = append(serviceAddrs, serviceAddr)
	}
	return strings.Join(serviceAddrs, " "), nil
}
