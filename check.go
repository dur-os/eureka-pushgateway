package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"net/url"
	"time"
)

var logger = logging.MustGetLogger("eureka-pushgateway")

func CheckJob(host string, timeout int) {
	response, _ := http.Get(fmt.Sprintf("http://%s/metrics", host))
	var parser expfmt.TextParser
	metricFamilies, _ := parser.TextToMetricFamilies(response.Body)
	if metricFamily, ok := metricFamilies["push_time_seconds"]; ok {
		second := time.Now().Unix()
		for _, metric := range metricFamily.Metric {
			v := metric.Gauge.Value
			c := float64(second) - *v
			if c > float64(timeout) {
				var job string
				var instance string
				for _, label := range metric.Label {
					if label.Name != nil && *label.Name == "job" {
						job = *label.Value
					}
					if label.Name != nil && *label.Name == "instance" {
						instance = *label.Value
					}
				}
				deleteGroup := fmt.Sprintf("http://%s/metrics/job/%s/instance/%s", host, job, url.QueryEscape(instance))
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
