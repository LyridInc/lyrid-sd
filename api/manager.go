package api

import (
	"encoding/json"
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
