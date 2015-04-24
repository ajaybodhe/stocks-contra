package conf

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pubmatic/pub-phoenix/util/conf"
)

const (
	gconfFile = "stocks-contra.conf"
)

/* CFillerConfig stores the global configuration structure for cache filler */
var StocksContraConfig struct {
	/*
		Server struct {
			AdminPort string
		}
		Redis struct {
			ServerList  string
			Port        int
			ConTO       int
			OpTO        int
			ConPoolSize int
		}
	*/
	DB struct {
		Host     string
		Port     int
		Username string
		Password string
		Protocol string
		DB       string
		ConnID   string
	}
	/*DB struct {
		host     string
		port     int
		username string
		password string
		protocol string
	}*/
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	conf.ReadConfig(gconfFile, &StocksContraConfig)
	StocksContraConfig.DB.ConnID = fmt.Sprintf("%s:%s@%s(%s:%d)/%s", StocksContraConfig.DB.Username, CFillerConfig.DB.Password, CFillerConfig.DB.Protocol,
		StocksContraConfig.DB.Host, StocksContraConfig.DB.Port, StocksContraConfig.DB.DB)
}
