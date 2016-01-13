package kfkc

import (
	//"fmt"
)

type Deploy struct {
	cmd			string
	ip			[]string

	prescript	string
	postscript	string

	sc			*SshContext

	format		*Format

	log			*Log
}

func (d *Deploy) clear() {
	d.cmd			= ""
	d.ip			= []string{}
	d.prescript		= ""
	d.postscript	= ""

}

func (d *Deploy) parse(msg []byte) (*Message, error) {
	d.clear()

	raw, err := d.format.DecodeBase64(msg)

	if err != nil {
		d.log.Error("Decode base64 failed")
		return nil, err
	}

	plainMsg := &Message{}
	err = d.format.DecodeJson(raw, &plainMsg)
	if err != nil {
		d.log.Error("Decode json failed")
		return nil, err
	}

	d.log.Debug("%+v", plainMsg)

	return plainMsg, nil
}

func InitDeploy(sc *SshContext, log *Log) (*Deploy, error) {
	d := &Deploy{}
	d.sc = sc
	d.log = log
	d.format = InitFormat()

	return d, nil
}

func (d *Deploy) RunDeploy(msg []byte) error {
	plain, err := d.parse(msg)
	if err != nil {
		return err
	}

	d.log.Debug("ips is %+v", plain.Ips)
	d.log.Debug("Host is %s", plain.Host)
	

	return nil
}


