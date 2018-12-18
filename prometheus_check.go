package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func PrometheusCheckJob(host string, timeout int) {
	second := time.Now().Unix()
	response, err := http.Get(fmt.Sprintf("http://%s/api/v1/query?query=push_time_seconds&time=%d", host, second))
	if err != nil {
		return
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	jdata := new(map[string]interface{})
	err = json.Unmarshal(bytes, jdata)
	if err != nil {
		return
	}
	i, ok := (*jdata)["data"]
	if ok {
		if data, ok := i.(map[string]interface{})["result"]; ok {
			for _, v := range data.([]interface{}) {
				metric := v.(map[string]interface{})["metric"].(map[string]interface{})
				values := v.(map[string]interface{})["value"].([]interface{})
				if len(values) == 2 {
					pushTime, _ := strconv.ParseFloat(values[1].(string), 10)
					result := values[0].(float64) - pushTime - float64(timeout)
					if result > 0 {
						var job = metric["exported_job"].(string)
						var instance = metric["exported_instance"].(string)
						deleteGroup := fmt.Sprintf("http://%s/metrics/job/%s/instance/%s", metric["instance"].(string), job, url.QueryEscape(instance))
						logger.Infof("Job: %s , instance : %s now : %f, pushTime: %s", job, instance,  values[0],values[1].(string))
						request, err := http.NewRequest(http.MethodDelete, deleteGroup, nil)
						if err != nil {
							logger.Errorf("delete group URl:%s ,err: %v", deleteGroup, err)
							continue
						}
						_, err = http.DefaultClient.Do(request)
						if err != nil {
							logger.Errorf("delete group URl:%s ,err: %v", deleteGroup, err)
							continue
						}
					}
				}
			}
		}
	}

}
