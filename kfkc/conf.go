package kfkc

import (
	goconf "github.com/msbranco/goconfig"
	"path/filepath"
	//"time"
)

const DEFAULT_CONF_PATH		= "../conf"
const DEFAULT_CONF			= "kfkc.conf"
const DEFAULT_SSH_TIMEOUT	= 5

type Config struct {
	log		string	/* log file */
	level	string	/* log level */

	topic	string	/* topic */

	kaddr	string	/* kafka addr */

	sport	string	/* ssh port */
	suser	string	/* ssh user */
	skey	string	/* ssh key */
	sto		int64	/* ssh timeout */
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
	conf.topic, _		= c.GetString("default", "topic")
	conf.kaddr, _		= c.GetString("default", "kafka_addr")
	conf.sport, _		= c.GetString("default", "ssh_port")
	conf.suser, _		= c.GetString("default", "ssh_user")
	conf.skey, _		= c.GetString("default", "ssh_key")
	conf.sto, err		= c.GetInt64("default", "ssh_timeout")
	if err != nil {
		conf.sto = DEFAULT_SSH_TIMEOUT
	}

	return conf, nil
}

