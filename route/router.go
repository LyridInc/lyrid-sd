package route

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func ExporterStatus(c *gin.Context) {
	fmt.Println(c.Request.Host)
}

func Hello(c *gin.Context) {
	c.String(200, "Hello from: "+c.Request.Host)
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

func (r *Router) getMetricFamily(p int) {

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
	router.GET("/hello", Hello)
	router.GET("/status", ExporterStatus)

	r.server = &http.Server{
		Addr:    ":" + r.Port,
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
