package kfkc

import (
	"strings"
	"errors"
	//"fmt"
)

type Deploy struct {
	cmd			string
	ip			[]string

	prescript	string
	postscript	string

	sc			*SshContext
	hc			*HttpClient

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

	dmsg := strings.Split(string(msg), "\"")
	if len(dmsg) < 3 {
		d.log.Error("Message format error, len is %d", len(dmsg))
		return nil, errors.New("Message format error")
	}

	d.log.Debug("Splited msg is %s", dmsg[1])

	raw, err := d.format.DecodeBase64([]byte(dmsg[1]))

	if err != nil {
		d.log.Error("Decode base64 failed")
		return nil, err
	}
	d.log.Debug("After base64 data is %s", string(raw))

	plainMsg := &Message{}
	err = d.format.DecodeJson(raw, &plainMsg)
	if err != nil {
		d.log.Error("Decode json failed")
		d.log.Error(plainMsg)
		return nil, err
	}

	d.log.Debug("%+v", plainMsg)

	arr := strings.Split(plainMsg.Rights, "|")
	plainMsg.right = arr[0]
	plainMsg.owner = arr[1]

	return plainMsg, nil
}

func (d *Deploy) reportResult(m *Message, ip string, errMsg string) error {
	return d.hc.Report(m.Loguri, m.Flowid, m.Taskruntime, ip, errMsg)
}

func InitDeploy(sc *SshContext, log *Log) (*Deploy, error) {
	d := &Deploy{}
	d.sc = sc
	d.log = log
	d.format = InitFormat()
	d.hc = InitHttpClient(d.log)

	return d, nil
}

func (d *Deploy) RunDeploy(msg []byte) error {
	plain, err := d.parse(msg)
	if err != nil {
		d.log.Error("Parse kafka message failed")
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
			defer func () { ch <- 1 }()

			sconn, err := d.sc.InitSshConn(ip, d.log)
			if err != nil {
				d.log.Error("Init ssh connection failed")
				d.reportResult(plain, ip, "Connect failed")
				return
			}

			defer sconn.SshClose()

			/* mkdir */
			_, err = sconn.SshExec("mkdir " + plain.Dir)
			if err != nil {
				/* maybe dir is exist */
				d.log.Info("Execute mkdir %s in %s failed, continue", plain.Dir, ip)
			}

			/* scp */
			err = sconn.SshScp(data, plain.File, plain.Dir, plain.right)
			if err != nil {
				d.log.Error("Scp file %s to %s failed", plain.File, ip)
				d.reportResult(plain, ip, "Scp file failed")
				return
			}
			d.log.Debug("Scp %s to %s done", plain.File, ip)

			/* Md5 check not used */
			///* md5sum */
			//res, err = sconn.SshExec("md5sum " + plain.Dir + "/" + plain.File)
			//if err != nil {
			//	d.log.Error("Execute md5sum %s in %s failed", plain.File, ip)
			//	d.reportResult(plain, ip, "Md5sum failed")
			//	return
			//}

			///* check md5 */
			//tmp := strings.Fields(string(res))
			//if !strings.EqualFold(tmp[0], plain.Localfile_md5) {
			//	d.log.Error("Check md5 failed")
			//	d.reportResult(plain, ip, "Check md5 failed")
			//	return
			//}

			/* postscript */
			if plain.Postscript != "" {
				_, err = sconn.SshExec(plain.Postscript)
				if err != nil {
					d.log.Error("Execute %s in %s failed", plain.Postscript, ip)
					d.reportResult(plain, ip, "Execute postscript failed")
					return
				}
			}

			/* report success */
			d.reportResult(plain, ip, "")
			d.log.Debug("deploy %s done", ip)
		}(host)
	}

	for i := 0; i < total; i++ {
		<-ch
	}

	d.log.Info("Job done")

	return nil
}


