package kfkc

import (
	"fmt"
)

func Run() {
    cconf := new(Config)
    conf, _:= cconf.ReadConf("kfkc.conf")

	if conf == nil {
        fmt.Println("no conf")
        return
    }

	//if conf.daemon {
	//	fmt.Println("To be daemon")
	//	err := Daemon(1)
	//	if err != nil {
	//		fmt.Println("Daemonize failed")
	//		return
	//	}
	//}
	//fmt.Println("Not to be daemon")

    log := GetLogger(conf.log, conf.level)

    server, err := InitServer(conf, log)
    if err != nil {
        log.Error("init server failed")
        return
    }

    server.CoreRun()
}
