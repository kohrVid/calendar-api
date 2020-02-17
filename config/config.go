package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() map[string]interface{} {
	switch strings.ToLower(os.Getenv("ENV")) {
	case "development":
		return getConf("development")
	case "test":
		return getConf("test")
	default:
		return getConf("development")
	}
}

func ToMapList(n interface{}) []map[string]interface{} {
	mapList := []map[string]interface{}{}

	n2, ok := n.([]interface{})
	if !ok {
		return mapList
	} else if len(n2) < 1 {
		return mapList
	}

	for _, i := range n2 {
		m := make(map[string]interface{})
		for k, v := range i.(map[interface{}]interface{}) {
			m[strings.ToLower(k.(string))] = v
		}

		mapList = append(mapList, m)
	}

	return mapList
}

func getConf(env string) map[string]interface{} {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path.Join(
		"$GOPATH",
		"src",
		"github.com",
		"kohrVid",
		"calendar-api",
		"config/",
	))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf := viper.Get(env)

	return conf.(map[string]interface{})
}
