package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/eureka"
	"github.com/hudl/fargo"
	"github.com/op/go-logging"
	"net/http"
	"strconv"
	"strings"
)

var (
	PushGatewayPort       int
	EurekaPushGatewayPort int
	HostIP                string
	EurekaUrl             string
)

func init() {
	flag.StringVar(&HostIP, "host", "127.0.0.1", "this Host IP")
	flag.StringVar(&EurekaUrl, "eureka", "127.0.0.1", "this eureka url")
	flag.IntVar(&PushGatewayPort, "port", 9091, "this PushGateway Port")
	flag.IntVar(&EurekaPushGatewayPort, "eport", 9092, "this Eureka PushGateway Port")
}

func main() {
	flag.Parse()
	logging.SetLevel(logging.INFO, "fargo")
	connection := fargo.NewConn(strings.Split(EurekaUrl, ",")...)
	logger := log.NewNopLogger()
	instance := &fargo.Instance{
		InstanceId:       HostIP + ":" + strconv.Itoa(PushGatewayPort),
		HostName:         HostIP,
		Port:             PushGatewayPort,
		PortEnabled:      true,
		App:              "prometheus-pushgateway-server",
		IPAddr:           HostIP,
		VipAddress:       "prometheus-pushgateway-server",
		SecureVipAddress: "prometheus-pushgateway-server",
		HealthCheckUrl:   fmt.Sprintf("http://%s:%d/health", HostIP, EurekaPushGatewayPort),
		StatusPageUrl:    fmt.Sprintf("http://%s:%d/status", HostIP, EurekaPushGatewayPort),
		HomePageUrl:      fmt.Sprintf("http://%s:%d/", HostIP, EurekaPushGatewayPort),
		Status:           fargo.UP,
		CountryId:        1,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
	}

	registrar := eureka.NewRegistrar(&connection, instance, logger)
	registrar.Register()
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`{"status":"UP"}`))
	})
	http.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`ok`))
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`ok`))
	})
	http.ListenAndServe(":9092", nil)
}
