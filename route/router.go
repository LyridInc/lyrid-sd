package route

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"io/ioutil"
	"log"
	lyridmodel "lyrid-sd/model"
	"net/http"
	"time"
)

var (
	packetPrefix = model.MetaLabelPrefix + "lyrid_"
)

func labelName(postfix string) string {
	return packetPrefix + postfix
}

func ExporterStatus(c *gin.Context) {
	fmt.Println(c.Request.Host)
}

type Router struct {
	ID   string
	Port string

	metricendpoint string

	server *http.Server
}

func CreateNewRouter(port string) Router {
	return Router{
		Port: port,
	}
}

func (r *Router) GetTarget() *targetgroup.Group {
	config, _ := lyridmodel.GetConfig()
	return &targetgroup.Group{
		Source: fmt.Sprintf("lyrid/%s", r.ID),
		Targets: []model.LabelSet{
			model.LabelSet{
				model.AddressLabel: model.LabelValue(config.Discovery_Interface + ":" + r.Port),
			},
		},
		Labels: model.LabelSet{
			model.LabelName(labelName("tags")): model.LabelValue("sample_tag"),
			model.LabelName(labelName("id")):   model.LabelValue(r.ID),
		},
	}
}

func (r *Router) getMetricFamily() []*dto.MetricFamily {
	metrics := make([]*dto.MetricFamily, 0)
	// todo: Change to lyrid-sdk later
	url := "http://localhost:8080"

	request := make(map[string]interface{})
	request["Command"] = "GetScrapeResult"
	exporter := make(map[string]interface{})
	exporter["ID"] = r.ID
	request["Exporter"] = exporter

	jsonreq, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonreq))
	req.Header.Add("content-type", "application/json")

	response, _ := http.DefaultClient.Do(req)

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var jsonresp map[string]interface{}

	json.Unmarshal(body, &jsonresp)

	if jsonresp["ReturnPayload"] != nil {
		raw := jsonresp["ReturnPayload"].(string)
		var decoded_json []*dto.MetricFamily
		json.Unmarshal([]byte(raw), &decoded_json)
		metrics = decoded_json
	}

	return metrics
}

func (r *Router) GetPort() string {
	return r.Port
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
	config, _ := lyridmodel.GetConfig()
	r.server = &http.Server{
		Addr:    config.Bind_Address + ":" + r.Port,
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
