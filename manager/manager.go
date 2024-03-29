package manager

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LyridInc/go-sdk"
	sdkModel "github.com/LyridInc/go-sdk/model"
	"github.com/go-kit/kit/log/level"
	"lyrid-sd/logger"
	"lyrid-sd/model"
	"lyrid-sd/route"
)

type NodeManager struct {
	StartPort         int
	NextPortAvailable int
	RouteMap          map[string]model.Router
	Apps              []*sdkModel.App
}

type customSD struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
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
	config := model.GetConfig()
	if config.Discovery_Max_Port_Used > config.Discovery_Port_Start {
		manager.StartPort = config.Discovery_Max_Port_Used
	} else {
		manager.StartPort = config.Discovery_Port_Start
	}
	manager.NextPortAvailable = manager.StartPort
	manager.Apps = sdk.GetInstance().GetApps()
}

func (manager *NodeManager) ReRoute() {
	// Close created route
	level.Info(logger.GetInstance().Logger).Log("Message", "Re reoute")
	for _, r := range manager.RouteMap {
		r.Close()
		r = nil
	}
	manager.RouteMap = make(map[string]model.Router)
	config := model.GetConfig()
	manager.StartPort = config.Discovery_Port_Start
	manager.NextPortAvailable = manager.StartPort
}
func getUsedPort(id string, sd []customSD) int {
	for _, item := range sd {
		if item.Labels[route.LabelName("id")] == id && len(item.Labels[route.LabelName("port")]) > 0 {
			port, _ := strconv.Atoi(item.Labels[route.LabelName("port")])
			return port
		}
	}
	return 0
}

func isReserved(p int, sd []customSD) bool {
	for _, item := range sd {
		if len(item.Labels[route.LabelName("port")]) > 0 {
			port, _ := strconv.Atoi(item.Labels[route.LabelName("port")])
			if port == p {
				return true
			}
		}
	}
	return false
}

func (manager *NodeManager) Run(ctx context.Context) {
	config := model.GetConfig()
	duration, _ := time.ParseDuration(config.Discovery_Poll_Interval)
	for c := time.Tick(duration); ; {
		var sd []customSD
		if len(manager.RouteMap) == 0 {
			// first run, check old used ports on config file
			jsonFile, err := os.Open(os.Getenv("CONFIG_DIR") + "/lyrid_sd.json")
			defer jsonFile.Close()
			if err == nil {
				byteValue, _ := ioutil.ReadAll(jsonFile)
				_ = json.Unmarshal([]byte(byteValue), &sd)
			}
		}
		list := manager.GetExporterList()
		config := model.GetConfig()
		maxPortused := config.Discovery_Max_Port_Used
		for _, endpoint := range list {
			if manager.RouteMap[endpoint.ID] == nil {
				// route to this id doesn't exist
				level.Info(logger.GetInstance().Logger).Log("Message", "Route to ID doesn't exist", "EndpointID", endpoint.ID)
				r := route.Router{ID: endpoint.ID, URL: endpoint.URL, AdditionalLabels: endpoint.AdditionalLabels}
				if sd == nil {
					r.Initialize(strconv.Itoa(manager.NextPortAvailable))
					manager.NextPortAvailable++
					maxPortused = manager.NextPortAvailable
				} else {
					port := getUsedPort(endpoint.ID, sd)
					if port != 0 {
						r.Initialize(strconv.Itoa(port))
					} else {
						for ok := true; ok; ok = isReserved(port, sd) {
							port = manager.NextPortAvailable
							manager.NextPortAvailable++
						}
						r.Initialize(strconv.Itoa(port))
					}
					if port > maxPortused {
						maxPortused = port
					} else {
						maxPortused = manager.NextPortAvailable
					}
				}
				go r.Run()
				manager.RouteMap[endpoint.ID] = &r
				// notify Discovery Engine to create target over in the in json file
			} else {
				// update labels
				manager.RouteMap[endpoint.ID].Update(&endpoint)
			}
		}
		config = model.GetConfig()
		if config.Discovery_Max_Port_Used < maxPortused {
			config.Discovery_Max_Port_Used = maxPortused
			model.WriteConfig(config)
		}
		select {
		case <-c:
			continue
		case <-ctx.Done():
			return

		}
	}
}

func (manager *NodeManager) ExecuteFunction(body string) ([]byte, error) {
	response, err := sdk.GetInstance().ExecuteFunctionByName(model.GetConfig().Noc_App_Name, os.Getenv("NOC_MODULE_NAME"), os.Getenv("NOC_TAG"), os.Getenv("NOC_FUNCTION_NAME"), body)
	level.Debug(logger.GetInstance().Logger).Log("Response", response)
	return response, err
}

func (manager *NodeManager) ExecuteFunctionWithURIAndMethod(method string, uri string, body string) ([]byte, error) {
	for _, app := range manager.Apps {
		if strings.Contains(strings.ToLower(app.Name), strings.ToLower(os.Getenv("NOC_APP_NAME"))) {
			level.Debug(logger.GetInstance().Logger).Log("App name", app.Name)
			response, err := sdk.GetInstance().ExecuteApp(app.Name, os.Getenv("NOC_MODULE_NAME"), os.Getenv("NOC_TAG"), os.Getenv("NOC_FUNCTION_NAME"), uri, method, body)
			level.Debug(logger.GetInstance().Logger).Log("Response", response)
			return response, err
		}
	}
	return nil, errors.New("unable to execute")
}

func (manager *NodeManager) Add(r model.Router) {
	manager.RouteMap[r.GetPort()] = r
}

func (manager *NodeManager) GetExporterList() []model.ExporterEndpoint {
	exporter_list := make([]model.ExporterEndpoint, 0)
	response, err := manager.ExecuteFunctionWithURIAndMethod("GET", "/api/exporters", "")
	if err != nil {
		level.Error(logger.GetInstance().Logger).Log("Error", err)
	}
	var jsonresp []model.ExporterEndpoint
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err == nil {
		exporter_list = jsonresp
	}
	return exporter_list
}

func (manager *NodeManager) GetMetricsFromEndpoint(id string) {
	//Result []*dto.MetricFamily
}
