package model

import "github.com/tkanos/gonfig"

type Configuration struct {
	Discovery_Port_Start	int
	Max_Discovery			int
	Discovery_Poll_Interval	string
	Lyrid_Key				string
	Lyrid_Secret			string
	Local_Serverless_Url	string
	Is_Local				bool
}

func GetConfig() (Configuration, error) {
	configuration := Configuration{}
	err := gonfig.GetConf("config/config.json", &configuration)
	return configuration, err
}

