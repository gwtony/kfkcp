package kfkp

import (
	//"reflect"
	"net/http"
)

type Server struct {
	lport	string
	kaddr	string
	kport	string
	log		*Log
	kafka	*KafkaClient
}

func (s *Server) pubMsg(w http.ResponseWriter, req *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte("hello, world!\n"))
	err := req.ParseForm()
	if err != nil {
		s.log.Error("parse form failed", err)
		return
	}

	//s.log.Debug(reflect.TypeOf(req.PostForm["id"]))
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
	s.lport = conf.lport
	s.kaddr = conf.kaddr
	s.kport = conf.kport
	kafka, err := InitKafka(s.kaddr, s.kport, s.log)
	if err != nil {
		s.log.Error("init kafka client faild")
		return nil, err
	}
	s.kafka = kafka

	return s, nil
}

func (s *Server) CoreRun() error {
	http.HandleFunc("/", s.pubMsg)
	err := http.ListenAndServe(":" + s.lport, nil)
	if err != nil {
		s.log.Error("ListenAndServe: ", err)
		return err
	}

	return nil
}

