package api

import (
	"encoding/json"
	"github.com/LyridInc/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/tkanos/gonfig"
	"io/ioutil"
	"lyrid-sd/model"
)

func GetStatus(c *gin.Context) {

}

func UpdateConfig(c *gin.Context) {
	var request model.Configuration
	if err := c.ShouldBindJSON(&request); err == nil {
		configuration := model.Configuration{}
		configuration.Bind_Address = request.Bind_Address
		configuration.Discovery_Interface = request.Discovery_Interface
		configuration.Discovery_Poll_Interval = request.Discovery_Poll_Interval
		configuration.Discovery_Port_Start = request.Discovery_Port_Start
		configuration.Max_Discovery = request.Max_Discovery
		configuration.Mngt_Port = request.Mngt_Port
		configuration.Lyrid_Key = request.Lyrid_Key
		configuration.Lyrid_Secret = request.Lyrid_Secret
		configuration.Local_Serverless_Url = request.Local_Serverless_Url
		configuration.Is_Local = request.Is_Local
		if configuration.Is_Local && len(configuration.Local_Serverless_Url) > 0 {
			sdk.GetInstance().SimulateServerless(configuration.Local_Serverless_Url)
		} else {
			sdk.GetInstance().DisableSimulate()
		}
		f, _ := json.MarshalIndent(configuration, "", " ")
		_ = ioutil.WriteFile("config/config.json", f, 0644)
		c.JSON(200, configuration)
	} else {
		c.JSON(400, err)
	}
}

func GetConfig(c *gin.Context) {
	configuration := model.Configuration{}
	err := gonfig.GetConf("config/config.json", &configuration)
	if (err == nil) {
		c.JSON(200, configuration)
	} else {
		c.JSON(400, err)
	}
}
