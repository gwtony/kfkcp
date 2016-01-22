package kfkp

import (
	goconf "github.com/msbranco/goconfig"
	"path/filepath"
	//"time"
)

const DEFAULT_CONF_PATH = "../conf"
const DEFAULT_CONF		= "kfkp.conf"

type Config struct {
	log		string		/* log file */
	level	string		/* log level */
	//timeout time.Duration
	lport	string		/* listen port */
	kaddr	string		/* kafka addr */
}

func (conf *Config) ReadConf(file string) (*Config, error) {
	if file == "" {
		file = filepath.Join(DEFAULT_CONF_PATH, DEFAULT_CONF)
	}

	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		return nil, err
	}

	//TODO: check
	conf.log, _			= c.GetString("default", "log")
	conf.level, _		= c.GetString("default", "level")
	//conf.timeout, _		= c.GetInt64("default", "timeout")
	//conf.timeout		= time.Duration(timeout) * time.Millisecond
	conf.lport, _		= c.GetString("default", "listen_port")
	conf.kaddr, _		= c.GetString("default", "kafka_addr")

	return conf, nil
}

