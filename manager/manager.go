package manager

import (
	"context"
	"encoding/json"
	"github.com/LyridInc/go-sdk"
	"log"
	"lyrid-sd/model"
	"lyrid-sd/route"
	"lyrid-sd/utils"
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
	config, _ := model.GetConfig()
	manager.StartPort = config.Discovery_Port_Start
	manager.NextPortAvailable = manager.StartPort
}

func (manager *NodeManager) ReRoute() {
	// Close created route
	log.Println("Re route")
	for _, r  := range manager.RouteMap {
		r.Close()
		r = nil
	}
	manager.RouteMap = make(map[string]model.Router)
	config, _ := model.GetConfig()
	manager.StartPort = config.Discovery_Port_Start
	manager.NextPortAvailable = manager.StartPort
}

func (manager *NodeManager) Run(ctx context.Context) {
	config, _ := model.GetConfig()
	duration, _ := time.ParseDuration(config.Discovery_Poll_Interval)
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
	exporter_list := make([]model.ExporterEndpoint, 0)
	response, _ := sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "ListExporter"}))
	log.Println("response: ",string(response))
	var jsonresp map[string]interface{}
	json.Unmarshal([]byte(response), &jsonresp)
	/*
	exporter_list := make([]model.ExporterEndpoint, 0)
	url := "http://localhost:8080"

	request := make(map[string]interface{})
	request["Command"] = "ListExporter"

	jsonreq, _ := json.Marshal(request)
	fmt.Println()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonreq))
	req.Header.Add("content-type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return exporter_list
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var jsonresp map[string]interface{}

	json.Unmarshal(body, &jsonresp)

	 */

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
