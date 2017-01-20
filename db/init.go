package db

import (
	"github.com/ajaybodhe/stocks-contra/util"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/conf"
	"database/sql"
	"os"
	"github.com/golang/glog"
)

var proddbhandle util.DB
var testdbhandle util.DB

func init() {
	//initialize production db handle
	fmt.Printf("\nconf.StocksContraConfig.DB.ConnID=%v",conf.StocksContraConfig.DB.ConnID)
	proddb, err := sql.Open("mysql", conf.StocksContraConfig.DB.ConnID + "&parseTime=True")
	if err != nil {
		glog.Errorln("error: connecting to mysql:", conf.StocksContraConfig.DB.ConnID, ":error:", err.Error())
		return
	}
	if err := proddb.Ping(); err != nil {
		glog.Fatalln("fatal: unable to connect to db:", err.Error())
		os.Exit(1)
	}
	proddbhandle.Set(proddb)
}
