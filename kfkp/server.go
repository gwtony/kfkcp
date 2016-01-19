package kfkp

import (
	//"reflect"
	"net/http"
	"encoding/json"
)

type Server struct {
	lport	string
	kaddr	string
	kport	string
	log		*Log
	kafka	*KafkaClient
}

func (s *Server) pubMsg(w http.ResponseWriter, req *http.Request) {
	m := &Message{}
	m.Reason = ""
	m.Retcode = 1

	defer func() {
		ret, _ := json.Marshal(m)
		w.Header().Set("Content-Type", "application/json")
		w.Write(ret)
	}()

	s.log.Debug("in pub msg")

	err := req.ParseForm()
	if err != nil {
		s.log.Error("Parse form failed", err)
		m.Reason = "Parse form failed"
		return
	}

	for k, v := range req.PostForm {
		s.log.Debug("%s %s", k, v)
	}

	tslice := req.PostForm["topic"]
	if tslice == nil {
		s.log.Error("No topic in post data")
		m.Reason = "No topic in post data"
		return
	}

	topic := tslice[0]

	mslice := req.PostForm["msg"]
	if mslice == nil {
		s.log.Error("No msg in post data")
		m.Reason = "No msg in post data"
		return
	}
	msg := mslice[0]

	err = s.kafka.sendData(topic, msg)
	if err != nil {
		s.log.Error("Send data to kafka failed")
		m.Reason = "Send data to kafka failed"
		return
	}

	m.Retcode = 0
}

func InitServer(conf *Config, log *Log) (*Server, error) {
	s := &Server{}
	s.log = log
	s.lport = conf.lport
	s.kaddr = conf.kaddr
	s.kport = conf.kport
	kafka, err := InitKafka(s.kaddr, s.kport, s.log)
	if err != nil {
		s.log.Error("Init kafka client faild")
		return nil, err
	}
	s.kafka = kafka

	return s, nil
}

func (s *Server) CoreRun() error {
	http.HandleFunc("/pub/", s.pubMsg)
	err := http.ListenAndServe(":" + s.lport, nil)
	if err != nil {
		s.log.Error("ListenAndServe: ", err)
		return err
	}

	return nil
}

