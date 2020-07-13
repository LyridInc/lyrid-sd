package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lyrid-sd/model"
	"lyrid-sd/route"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type NodeManager struct {
	StartPort         int
	NextPortAvailable int
	RouteMap          map[string]model.Router
}

var instance *NodeManager
var once sync.Once

func GetInstance() *NodeManager {
	once.Do(func() {
		instance = &NodeManager{}
	})
	return instance
}

func (manager *NodeManager) Init() {
	manager.RouteMap = make(map[string]model.Router)
	manager.StartPort, _ = strconv.Atoi(os.Getenv("DISCOVERY_PORT_START"))
	manager.NextPortAvailable = manager.StartPort
}

func (manager *NodeManager) Run(ctx context.Context) {
	duration, _ := time.ParseDuration(os.Getenv("DISCOVERY_POLL_INTERVAL"))
	for c := time.Tick(duration); ; {

		log.Println("Polling Lyrid For new Info")

		list := manager.GetExporterList()

		for _, endpoint := range list {
			if manager.RouteMap[endpoint.ID] == nil {
				// route to this id doesn't exist
				log.Println("Route to ID doesn't exist: ", endpoint.ID)

				r := route.Router{ID: endpoint.ID}
				r.Initialize(strconv.Itoa(manager.NextPortAvailable))
				manager.NextPortAvailable++
				go r.Run()
				manager.RouteMap[endpoint.ID] = &r
				// notify Discovery Engine to create target over in the in json file

			}
		}

		select {
		case <-c:
			continue
		case <-ctx.Done():
			return

		}
	}
}

func (manager *NodeManager) Add(r model.Router) {
	manager.RouteMap[r.GetPort()] = r
}

func (manager *NodeManager) GetExporterList() []model.ExporterEndpoint {
	// todo: Change to lyrid-sdk later

	exporter_list := make([]model.ExporterEndpoint, 0)
	url := "http://localhost:8080"

	request := make(map[string]interface{})
	request["Command"] = "ListExporter"

	jsonreq, _ := json.Marshal(request)
	fmt.Println()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonreq))
	req.Header.Add("content-type", "application/json")
	response, err := http.DefaultClient.Do(req)
	log.Println(response);
	if err != nil {
		return exporter_list
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var jsonresp map[string]interface{}

	json.Unmarshal(body, &jsonresp)

	if jsonresp["ReturnPayload"] != nil {
		exporters_raw := jsonresp["ReturnPayload"].([]interface{})
		for _, raw := range exporters_raw {
			raw_iface := raw.(map[string]interface{})

			exporter := model.ExporterEndpoint{
				ID:           raw_iface["ID"].(string),
				Gateway:      raw_iface["Gateway"].(string),
				URL:          raw_iface["URL"].(string),
				ExporterType: raw_iface["ExporterType"].(string),
			}
			exporter_list = append(exporter_list, exporter)
		}
	}

	return exporter_list
}

func (manager *NodeManager) GetMetricsFromEndpoint(id string) {
	//Result []*dto.MetricFamily
}
