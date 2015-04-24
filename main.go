package main

import (
	"github.com/ajaybodhe/stocks-contra/core"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

func main() {
	glog.Infoln("stocks contra begins....")
	core.Serve()
}
