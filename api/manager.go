package api

import (
	"encoding/json"
	"github.com/LyridInc/go-sdk"
	"github.com/gin-gonic/gin"
	"log"
	"lyrid-sd/manager"
	"lyrid-sd/model"
	"lyrid-sd/utils"
	"os"
)

func GetStatus(c *gin.Context) {

}

func CheckLyridConnection(c *gin.Context) {
	config := model.GetConfig()
	if len(config.Lyrid_Key) > 0 && len(config.Lyrid_Secret) > 0 {
		user := sdk.GetInstance().GetUserProfile()
		if user != nil {
			account := sdk.GetInstance().GetAccountProfile()
			c.JSON(200, account)
		} else {
			c.JSON(200, map[string]string{"status": "OK"})
		}
	} else {
		c.JSON(200, map[string]string{"status": "ERROR"})
	}
}

func UpdateConfig(c *gin.Context) {
	var request model.Configuration
	if err := c.ShouldBindJSON(&request); err == nil {
		configuration := model.Configuration{}
		configuration.Discovery_Poll_Interval = request.Discovery_Poll_Interval
		configuration.Discovery_Port_Start = request.Discovery_Port_Start
		configuration.Scrape_Valid_Timeout = request.Scrape_Valid_Timeout
		configuration.Max_Discovery = request.Max_Discovery
		configuration.Lyrid_Key = request.Lyrid_Key
		configuration.Lyrid_Secret = request.Lyrid_Secret
		configuration.Local_Serverless_Url = request.Local_Serverless_Url
		configuration.Is_Local = request.Is_Local
		if configuration.Is_Local && len(configuration.Local_Serverless_Url) > 0 {
			sdk.GetInstance().SimulateServerless(configuration.Local_Serverless_Url)
		} else {
			sdk.GetInstance().DisableSimulate()
		}
		config := model.GetConfig()
		model.WriteConfig(configuration)
		if config.Discovery_Port_Start !=  configuration.Discovery_Port_Start {
			manager.GetInstance().ReRoute()
		}
		c.JSON(200, configuration)
	} else {
		c.JSON(400, err)
	}
}

func GetConfig(c *gin.Context) {
	configuration := model.GetConfig()
	c.JSON(200, configuration)
}

func GetExporter(c *gin.Context) {
	for _, r := range manager.GetInstance().RouteMap {
		r.SetMetricEndpoint()
	}
	c.JSON(200, manager.GetInstance().RouteMap)
}

func DeleteExporter(c *gin.Context) {
	mgr := manager.GetInstance()
	id := c.Param("id")
	exporter := mgr.RouteMap[id]

	if exporter == nil {
		c.JSON(404, "exporter not found")
		return
	}
	delete(mgr.RouteMap, id)
	exp := model.ExporterEndpoint{ID:id}
	sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "DeleteExporter", Exporter: exp}))
	exporter.Close()
}

func GetGateways(c *gin.Context) {
	response, err := sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "ListGateways"}))
	if err != nil {
		log.Println("error: ",err)
		c.JSON(404, "error on getting gateway")
	}
	var jsonresp map[string]interface{}
	json.Unmarshal([]byte(response), &jsonresp)
	if jsonresp["ReturnPayload"] != nil {
		exporters_raw := jsonresp["ReturnPayload"].([]interface{})
		c.JSON(200, exporters_raw)
	}
	//c.JSON(200, nil)
}

func DeleteGateway(c *gin.Context) {
	id := c.Param("id")
	sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "DeleteGateway", Gateway: model.Gateway{ID: id}}))
}