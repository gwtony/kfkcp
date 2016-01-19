package kfkp

import (
	"fmt"
)

func Run() {
    pconf := new(Config)
    conf, _:= pconf.ReadConf("kfkp.conf")

	if conf == nil {
        fmt.Println("no conf")
        return
    }

	fmt.Println(conf.log, conf.level)

    log := GetLogger(conf.log, conf.level)

    server, err := InitServer(conf, log)
    if err != nil {
        log.Error("init server failed")
        return
    }

    server.CoreRun()
}
