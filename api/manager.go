package api

import (
	"encoding/json"
	"github.com/LyridInc/go-sdk"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"lyrid-sd/manager"
	"lyrid-sd/model"
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
		f, _ := json.MarshalIndent(configuration, "", " ")
		_ = ioutil.WriteFile(os.Getenv("CONFIG_DIR") + "/config.json", f, 0644)
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
