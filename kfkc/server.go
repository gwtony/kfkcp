package kfkc

import (
	"time"
	"strings"
	"net/http"
)

type Server struct {
	kaddr	string

	zaddr	string

	sport	string
	suser	string
	skey	string
	sto		time.Duration

	log		*Log
	kafka	*KafkaClient
	topics	[]string

	sc		*SshContext
}

func (s *Server) pubMsg(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		s.log.Error("Parse form failed", err)
		return
	}

	for k, v := range req.PostForm {
		s.log.Debug("%s %s", k, v)
	}

	tslice := req.PostForm["topic"]
	if tslice == nil {
		s.log.Error("No topic in post data")
		return
	}
	topic := tslice[0]

	mslice := req.PostForm["msg"]
	if mslice == nil {
		s.log.Error("No msg in post data")
		return
	}
	msg := mslice[0]

	s.kafka.SendData(topic, msg)
}

func InitServer(conf *Config, log *Log) (*Server, error) {
	s := &Server{}

	s.log = log
	s.kaddr = conf.kaddr
	s.topics = strings.Split(conf.topic, ",")
	s.sport = conf.sport
	s.suser = conf.suser
	s.skey  = conf.skey
	s.sto	= time.Duration(conf.sto) * time.Second

	sc, err := InitSshContext(s.skey, s.suser, s.sto, log)
	if err != nil {
		s.log.Error("Init ssh context failed")
		return nil, err
	}
	s.sc = sc

	kafka, err := InitKafka(s.kaddr, false, s.log)
	if err != nil {
		s.log.Error("Init kafka client faild")
		return nil, err
	}
	kafka.s = s
	s.kafka = kafka

	return s, nil
}

func (s *Server) RecvMsg(topic string, ch chan int) {
	s.log.Debug("RecvMsg topic is %s", topic)

	s.kafka.RecvMsg(topic, -1, ch)
}

func (s *Server) CoreRun() error {
	tlen := len(s.topics)

	ch := make(chan int, tlen)

	defer close(ch)

	for _, v := range s.topics {
		go s.RecvMsg(v, ch)
	}

	for {
		select {
		case <-ch :
			s.log.Debug("Got message from channel")
			tlen--
			if tlen == 0 {
				s.log.Debug("CoreRun return")
				return nil
			}
		}
	}

	return nil
}

