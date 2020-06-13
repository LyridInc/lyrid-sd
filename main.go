package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"io/ioutil"
	"lyrid-sd/adapter"
	"lyrid-sd/api"
	"lyrid-sd/manager"
	"lyrid-sd/route"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	logger       log.Logger
	packetPrefix = model.MetaLabelPrefix + "lyrid_"
)

// Note: create a config struct for your custom SD type here.
type sdConfig struct {
	Address         string
	TagSeparator    string
	RefreshInterval int
}

// Note: This is the struct with your implementation of the Discoverer interface (see Run function).
// Discovery retrieves target information from a Consul server and updates them via watches.
type discovery struct {
	address         string
	refreshInterval int
	tagSeparator    string
	logger          log.Logger
	oldSourceList   map[string]bool
}

func newDiscovery(conf sdConfig) (*discovery, error) {
	cd := &discovery{
		address:         conf.Address,
		refreshInterval: conf.RefreshInterval,
		tagSeparator:    conf.TagSeparator,
		logger:          logger,
		oldSourceList:   make(map[string]bool),
	}
	return cd, nil
}

func labelName(postfix string) string {
	return packetPrefix + postfix
}

func (d *discovery) createTarget() *targetgroup.Group {
	return &targetgroup.Group{
		Source: fmt.Sprintf("packet/%s", "2921a651-1188-49ff-b8cc-1e3cf6a25588"),
		Targets: []model.LabelSet{
			model.LabelSet{
				model.AddressLabel: model.LabelValue("10.1.17.49:9182"),
			},
			model.LabelSet{
				model.AddressLabel: model.LabelValue("10.1.17.53:9182"),
			},
		},
		Labels: model.LabelSet{
			model.LabelName(labelName("tags")): model.LabelValue("sample_tag"),
		},
	}
}

// Note: you must implement this function for your discovery implementation as part of the
// Discoverer interface. Here you should query your SD for it's list of known targets, determine
// which of those targets you care about (for example, which of Consuls known services do you want
// to scrape for metrics), and then send those targets as a target.TargetGroup to the ch channel.
func (d *discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	for c := time.Tick(time.Duration(d.refreshInterval) * time.Second); ; {
		logger.Log("Hello")

		// Check and discover all the registered endpoints in the lyrid

		targets := make([]*targetgroup.Group, 0)
		targets = append(targets, d.createTarget())
		ch <- targets
		// Wait for ticker or exit when ctx is closed.
		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (d *discovery) GetMetricFamilies() []*dto.MetricFamily {

	// get metric from http endpoint and then pass it back
	var metrics []*dto.MetricFamily
	response, err := http.Get("http://localhost:8081/endpoints/scrape/e678e31f-7829-44b4-9f64-c9c65c0f0050")

	if err != nil {
		fmt.Printf("The HTTP request to Lyrid Server failed. %s\n", err)
		os.Exit(1)
	} else {
		databyte, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(databyte, &metrics)

	}

	return metrics
}

func main() {

	manager.GetInstance().Init()

	ctx := context.Background()
	// NOTE: create an instance of your new SD implementation here.
	cfg := sdConfig{
		TagSeparator:    ",",
		Address:         "localhost",
		RefreshInterval: 30,
	}
	logger = log.NewSyncLogger(log.NewLogfmtLogger(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	disc, err := newDiscovery(cfg)
	if err != nil {
		fmt.Println("err: ", err)
	}

	sdAdapter := adapter.NewAdapter(ctx, "custom_sd.json", "exampleSD", disc, logger)
	sdAdapter.Run()

	for i := 9001; i <= 9005; i++ {
		r := route.Router{}
		r.Initialize(strconv.Itoa(i))
		go r.Run()
		manager.GetInstance().RouteMap[r.Port] = &r
	}
	//g := prometheus.Gatherers{
	//	prometheus.GathererFunc(func() ([]*dto.MetricFamily, error) { return disc.GetMetricFamilies(), nil }),
	//}

	router := gin.Default()
	//	router.Use(ginprom.PromMiddleware(nil))
	//	router.GET("/metrics", ginprom.PromHandler(promhttp.HandlerFor(g, promhttp.HandlerOpts{})))
	router.POST("/stop/:id", api.Stop)
	router.POST("/start/:id", api.Start)
	go router.Run(":9000")

	<-ctx.Done()
}
