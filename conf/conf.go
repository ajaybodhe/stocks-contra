package conf

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"

	"code.google.com/p/gcfg"
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

const allowAllFilesCommand = "allowAllFiles=true"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	ReadConfig(gconfFile, &StocksContraConfig)
	StocksContraConfig.DB.ConnID = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?%s", StocksContraConfig.DB.Username, StocksContraConfig.DB.Password, StocksContraConfig.DB.Protocol,
		StocksContraConfig.DB.Host, StocksContraConfig.DB.Port, StocksContraConfig.DB.DB, allowAllFilesCommand)
}

/*ReadConfig - reads the flags for --conf and if its found reads file and sets configuration into out. If --conf is not provided, then defaultPath is used. */
func ReadConfig(defaultPath string, out interface{}) {
	confFile := flag.String("conf", defaultPath, "Configuration file path")
	flag.Parse()
	glog.Info("conffile:", *confFile)
	err := gcfg.ReadFileInto(out, *confFile)
	if err != nil {
		glog.Fatal("error: util.conf.init:", err.Error())
	}
	glog.Info(os.Stdout, "boot.util.conf.init.success:\n***************Configuration:***************\n%+v\n*****************END****************\n", out)
}
