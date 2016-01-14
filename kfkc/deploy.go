package kfkc

import (
	"string"
	//"fmt"
)

type Deploy struct {
	cmd			string
	ip			[]string

	prescript	string
	postscript	string

	sc			*SshContext
	hc			*HttpContext

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

	arr := strings.Split(plainMsg.Rights, "|")
	plainMsg.right = arr[0]
	plainMsg.owner = arr[1]

	return plainMsg, nil
}

func (d *Deploy) reportResult(m *Message, ip string, errMsg string) error {
	return d.hc.Report(m.Loguri, m.Flowid, m.Sid, m.Taskruntime, ip, errMsg)
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
		d.log.Error("Parse kafka message failed")
		d.reportResult(plain, "", "Parse kafka message failed")
		return err
	}

	total := len(plain.Ips)
	if total <= 0 {
		d.log.Error("No host to deploy")
		d.reportResult(plain, "", "No host to deploy")
		return nil
	}

	data, err := d.hc.GetFile(plain.Host, plain.Webfile)
	if err != nil {
		d.log.Error("Get file %s failed", plain.Webfile)
		d.reportResult(plain, "", "Get file failed")
		return err
	}

	ch := make(chan int, total)

	for _, host := range plain.Ips {
		d.log.Debug("Deploy ip: %s", host)

		go func(ip string) {
			sconn := sc.InitSshConn(ip)

			defer ch <- 1
			defer sconn.SshClose()

			/* mkdir */
			res, err := sconn.SshExec("mkdir " + plain.Dir)
			if err != nil {
				d.log.Error("Execute mkdir %s in %s failed", plain.Dir, ip)
				d.reportResult(plain, ip, "mkdir failed")
				return
			}

			/* scp */
			err = sconn.sshScp(data, plain.File, plain.Dir, plain.right)
			if err != nil {
				d.log.Error("Scp file %s to %s failed", plain.File, ip)
				d.reportResult(plain, ip, "Scp file failed")
				return
			}

			/* md5sum */
			res, err := sconn.SshExec("md5sum " + plain.Dir + "/" + plain.File)
			if err != nil {
				d.log.Error("Execute md5sum %s in %s failed", plain.File, ip)
				d.reportResult(plain, ip, "md5sum failed")
				return
			}

			/* postscript */
			res, err := sconn.SshExec(plain.Postscript)
			if err != nil {
				d.log.Error("Execute %s in %s failed", plain.Postscript, ip)
				d.reportResult(plain, ip, "Execute postscript failed")
				return
			}
		}(host)
	}

	for i := 0; i < total; i++ {
		<-ch
	}

	d.log.Info("Job done")

	return nil
}


