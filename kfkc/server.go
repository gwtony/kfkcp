package kfkc

import (
	"net/http"
)

type Server struct {
	kaddr	string

	zaddr	string

	sport	string
	suser	string
	skey	string
	sto		int64

	log		*Log
	kafka	*KafkaClient
	topic	string

	sc		*SshContext
}

func (s *Server) pubMsg(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		s.log.Error("parse form failed", err)
		return
	}

	for k, v := range req.PostForm {
		s.log.Debug("%s %s", k, v)
	}

	tslice := req.PostForm["topic"]
	if tslice == nil {
		s.log.Error("no topic in post data")
		return
	}
	topic := tslice[0]

	mslice := req.PostForm["msg"]
	if mslice == nil {
		s.log.Error("no msg in post data")
		return
	}
	msg := mslice[0]

	s.kafka.sendData(topic, msg)
}

func InitServer(conf *Config, log *Log) (*Server, error) {
	s := &Server{}

	s.log = log
	s.kaddr = conf.kaddr
	s.topic = conf.topic
	s.sport = conf.sport
	s.suser = conf.suser
	s.skey  = conf.skey
	s.sto	= conf.sto

	s.sc = InitSshContext(s.skey, s.user, s.to, log)

	kafka, err := InitKafka(s.kaddr, false, s.log)
	if err != nil {
		s.log.Error("init kafka client faild")
		return nil, err
	}
	kafka.s = s
	s.kafka = kafka

	return s, nil
}

func (s *Server) CoreRun() error {
	s.kafka.recvMsg("ssh", -1)
	return nil
}

