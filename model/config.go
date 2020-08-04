package model

import (
	"encoding/json"
	"github.com/tkanos/gonfig"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Discovery_Port_Start	int
	Max_Discovery			int
	Discovery_Poll_Interval	string
	Lyrid_Key				string
	Lyrid_Secret			string
	Local_Serverless_Url	string
	Is_Local				bool
}

func GetConfig() (Configuration) {
	filePath := os.Getenv("CONFIG_DIR")+"/config.json"
	configuration := Configuration{}
	err := gonfig.GetConf(filePath, &configuration)
	if err != nil {
		configuration = Configuration{
			Discovery_Port_Start:    9001,
			Max_Discovery:           1024,
			Discovery_Poll_Interval: "15s",
			Lyrid_Key:               "",
			Lyrid_Secret:            "",
			Local_Serverless_Url:    "http://localhost:8080",
			Is_Local:                true,
		}
		f, _ := json.MarshalIndent(configuration, "", " ")
		_ = os.Mkdir(os.Getenv("CONFIG_DIR"), 0755)

		//file, er := os.Create(filePath)
		//file.Close()
		ioutil.WriteFile(filePath, f, 0644)
	}
	return configuration
}

