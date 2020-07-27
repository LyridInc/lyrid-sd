package main

import (
	"context"
	"fmt"
	"github.com/LyridInc/go-sdk"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"lyrid-sd/adapter"
	"lyrid-sd/api"
	"lyrid-sd/manager"
	lyridmodel "lyrid-sd/model"
	"os"
	"strconv"
	"time"
)

var (
	logger       log.Logger
	packetPrefix = model.MetaLabelPrefix + "lyrid_"
)

// Note: This is the struct with your implementation of the Discoverer interface (see Run function).
// Discovery retrieves target information from a Consul server and updates them via watches.
type Discovery struct {
	Address         string
	RefreshInterval int
	TagSeparator    string
	OldSourceList   map[string]bool
}

// Note: create a config struct for your custom SD type here.
type SDConfig struct {
	Address         string
	TagSeparator    string
	RefreshInterval int
}

func NewDiscovery(conf SDConfig) (*Discovery, error) {
	cd := &Discovery{
		Address:         conf.Address,
		RefreshInterval: conf.RefreshInterval,
		TagSeparator:    conf.TagSeparator,
		OldSourceList:   make(map[string]bool),
	}
	return cd, nil
}

func main() {

	godotenv.Load()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	ctx := context.Background()
	// NOTE: create an instance of your new SD implementation here.
	cfg := SDConfig{
		TagSeparator:    ",",
		Address:         "localhost",
		RefreshInterval: 30,
	}
	logger = log.NewSyncLogger(log.NewLogfmtLogger(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	disc, err := NewDiscovery(cfg)
	if err != nil {
		fmt.Println("err: ", err)
	}

	sdAdapter := adapter.NewAdapter(ctx, "lyrid_sd.json", "lyridSD", disc, logger)
	sdAdapter.Run()

	manager.GetInstance().Init()

	go manager.GetInstance().Run(context.Background())
	//for i := 9001; i <= 9005; i++ {
	//	r := route.Router{}
	//	r.Initialize(strconv.Itoa(i))
	//	go r.Run()
	//	manager.GetInstance().RouteMap[r.Port] = &r
	//}
	//g := prometheus.Gatherers{
	//	prometheus.GathererFunc(func() ([]*dto.MetricFamily, error) { return disc.GetMetricFamilies(), nil }),
	//}

	router := gin.Default()
	//	router.Use(ginprom.PromMiddleware(nil))
	//	router.GET("/metrics", ginprom.PromHandler(promhttp.HandlerFor(g, promhttp.HandlerOpts{})))
	router.GET("/status", api.GetStatus)
	router.POST("/config", api.UpdateConfig)
	router.GET("/config", api.GetConfig)
	router.Use(static.Serve("/", static.LocalFile("./web/build", true)))
	config, _ := lyridmodel.GetConfig()
	if len(config.Lyrid_Key) > 0 && len(config.Lyrid_Secret) > 0 {
		sdk.GetInstance().Initialize(config.Lyrid_Key, config.Lyrid_Secret)
		if config.Is_Local && len(config.Local_Serverless_Url) > 0 {
			sdk.GetInstance().SimulateServerless(config.Local_Serverless_Url)
		}
	}
	router.Run(":" + strconv.Itoa(config.Mngt_Port))
}

// Note: you must implement this function for your discovery implementation as part of the
// Discoverer interface. Here you should query your SD for it's list of known targets, determine
// which of those targets you care about (for example, which of Consuls known services do you want
// to scrape for metrics), and then send those targets as a target.TargetGroup to the ch channel.
func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	for c := time.Tick(time.Duration(d.RefreshInterval) * time.Second); ; {
		// Check and discover all the registered endpoints in the lyrid
		targets := make([]*targetgroup.Group, 0)

		for _, router := range manager.GetInstance().RouteMap {
			targets = append(targets, router.GetTarget())
		}

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
