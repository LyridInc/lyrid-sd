package api

import (
	"github.com/gin-gonic/gin"
	"lyrid-sd/manager"
	"lyrid-sd/route"
)

func Stop(c *gin.Context) {
	mgr := manager.GetInstance()
	id := c.Param("id")
	endpoint := mgr.RouteMap[id]

	if endpoint == nil {
		c.JSON(404, "endpoint not found")
		return
	}

	endpoint.Close()
	delete(mgr.RouteMap, id)

	c.JSON(200, true)
}

func Start(c *gin.Context) {
	mgr := manager.GetInstance()
	id := c.Param("id")

	endpoint := mgr.RouteMap[id]

	if endpoint != nil {
		c.JSON(400, "unable to create endpoint")
		return
	}

	r := route.CreateNewRouter(id)
	mgr.Add(&r)
	r.Run()
}
