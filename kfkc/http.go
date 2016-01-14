package kfkc

import (
	"fmt"
	"errors"
	"net/http"
	"io/ioutil"
)

type HttpClient struct {
	log	*Log
}

var HTTP_OK = 200

var template = "%s?flow_id=%s&sid=%s&cmd=%s&status=%s&time=%s&ip=%s&content=%s"


func InitHttpClient(log *Log) *HttpClient {
	hc := &HttpClient{}
	hc.log = log
	return hc
}

func (hc *HttpClient) Getfile(uri string) ([]byte, error) {
    resp, err := http.Get(uri)
    if err != nil {
		hc.log.Error("Http get %s failed", uri)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		hc.log.Error("Read body from response failed")
		return nil, err
    }

    clen := resp.Header.Get("Content-Length")
	if len(body) != clen {
		hc.log.Error("Content length error")
		return nil, errors.New("Content length error")
	}

	reurn body, nil
}

func (hc *HttpClient) Report(uri string, fid string, sid string, time string, ip string, msg string) error {
	if len(msg) {	/* error, status is 2 */
		uri = fmt.Sprintf(template, uri, fid, sid, "cmd", "2", time, ip, msg)
	} else {		/* ok, status is 1 */
		uri = fmt.Sprintf(template, uri, fid, sid, "cmd", "1", time, ip, msg)
	}

    resp, err := http.Get(uri)
    if err != nil {
		hc.log.Error("Report get uri %s failed", uri)
        return err
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		hc.log.Error("Report read body failed")
        return err
    }

	if body.StatusCode != HTTP_OK {
		hc.log.Error("http status error: %d", body.StatusCode)
        return errors.New(body.Status)
	}

	return nil
}
