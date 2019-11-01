package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	MgoUri string `json:"mgo_uri"`
	Database string `json:"database"`
	Ip string 	`json:"ip"`
	Port int `json:"port"`
	TimeOut int `json:"time_out"`
	Tbk Tbk
}
type Tbk struct {
	AppId string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	AdZoneId string `json:"ad_zone_id"`
}
var Conf *Config
func Init(file string) error {
	f,err := os.Open(file)
	if err != nil {
		return nil
	}
	defer f.Close()
	conf,err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	Conf = &Config{}
	return json.Unmarshal(conf,Conf)
}