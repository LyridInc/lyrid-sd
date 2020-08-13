package route

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/LyridInc/go-sdk"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/pierrec/lz4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"io/ioutil"
	"log"
	sdmodel "lyrid-sd/model"
	"lyrid-sd/utils"
	"net/http"
	"os"
	"time"
)

var (
	packetPrefix = model.MetaLabelPrefix + "lyrid_"
)

func LabelName(postfix string) string {
	return packetPrefix + postfix
}

func ExporterStatus(c *gin.Context) {
	fmt.Println(c.Request.Host)
}

type Router struct {
	ID   string
	Port string

	MetricEndpoint string
	URL            string
	server         *http.Server
	AdditionalLabels map[string]string
}

func CreateNewRouter(port string) Router {
	return Router{
		Port: port,
	}
}

func (r *Router) Update(endpoint *sdmodel.ExporterEndpoint) {
	r.AdditionalLabels = endpoint.AdditionalLabels
}

func (r *Router) GetTarget() *targetgroup.Group {
	labels := model.LabelSet{}
	for key, value := range r.AdditionalLabels {
		labels[model.LabelName(LabelName(key))] = model.LabelValue(value)
	}
	labels[model.LabelName(LabelName("id"))] = model.LabelValue(r.ID)
	labels[model.LabelName(LabelName("port"))] = model.LabelValue(r.Port)
	return &targetgroup.Group{
		Source: fmt.Sprintf("lyrid/%s", r.ID),
		Targets: []model.LabelSet{
			model.LabelSet{
				model.AddressLabel: model.LabelValue(os.Getenv("DISCOVERY_INTERFACE") + ":" + r.Port),
			},
		},
		Labels: labels,
	}
}

func (r *Router) getMetricFamily() []*dto.MetricFamily {
	metrics := make([]*dto.MetricFamily, 0)
	exporter := sdmodel.ExporterEndpoint{ID: r.ID}
	response, _ := sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(sdmodel.LyFnInputParams{Command: "GetScrapeResult", Exporter: exporter}))
	var jsonresp map[string]*sdmodel.ScrapesEndpointResult
	json.Unmarshal([]byte(response), &jsonresp)

	if jsonresp["ReturnPayload"] != nil {
		scrapeResult := jsonresp["ReturnPayload"]
		dur, _ := time.ParseDuration(sdmodel.GetConfig().Scrape_Valid_Timeout)
		if time.Since(scrapeResult.ScrapeTime) <= dur {
			raw := []byte(scrapeResult.ScrapeResult)
			if scrapeResult.IsCompress {
				decompressed_b64bytes, _ := b64.StdEncoding.DecodeString(scrapeResult.ScrapeResult)
				r := lz4.NewReader(bytes.NewBuffer(decompressed_b64bytes))
				var err error
				raw, err = ioutil.ReadAll(r)
				if err != nil {
					fmt.Println(err)
				}
			}

			var decoded_json []*dto.MetricFamily
			json.Unmarshal(raw, &decoded_json)
			metrics = decoded_json
		}
	}

	return metrics
}

func (r *Router) GetPort() string {
	return r.Port
}

func (r *Router) SetMetricEndpoint() {
	r.MetricEndpoint = "http://" + r.server.Addr + "/metrics"
}

func (r *Router) Initialize(p string) error {

	// check port number, if used, then throw error
	r.Port = p
	return nil
}

func (r *Router) Run() {
	router := gin.Default()
	//router.Use(ginprom.PromMiddleware(nil))
	g := prometheus.Gatherers{
		prometheus.GathererFunc(func() ([]*dto.MetricFamily, error) { return r.getMetricFamily(), nil }),
	}
	router.Use(ginprom.PromMiddleware(nil))
	router.GET("/metrics", ginprom.PromHandler(promhttp.HandlerFor(g, promhttp.HandlerOpts{})))
	router.GET("/status", ExporterStatus)
	//router.GET("/metrics", ExporterStatus)
	r.server = &http.Server{
		Addr:    os.Getenv("BIND_ADDRESS") + ":" + r.Port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	/*
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shuting down server on port: " + strconv.Itoa(r.Port) )
		//router.GET("/metrics", ginprom.PromHandler(promhttp.HandlerFor(g, promhttp.HandlerOpts{})))

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		r.server.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := r.server.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}

	*/
	//endless.ListenAndServe(":" + strconv.Itoa(r.Port), router)
	//graceful.ListenAndServe(srv,10*time.Second)

	//router.Run()
	//router.Run(":" + strconv.Itoa(r.Port))

}

func (r *Router) Close() {
	log.Println("Shutting down server at: " + r.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer r.server.Close()
	defer r.server.Shutdown(ctx)
	r.server = nil
	log.Println("Server shut down at: " + r.Port)
}
