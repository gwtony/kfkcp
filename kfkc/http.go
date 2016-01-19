package kfkc

import (
	"fmt"
	"strconv"
	"errors"
	"net/http"
	"net/url"
	"io/ioutil"
)

type HttpClient struct {
	log	*Log
}

var HTTP_OK = 200

var template = "%s?flow_id=%d&sid=0&cmd=%s&status=%s&time=%s&ip=%s&content=%s"


func InitHttpClient(log *Log) *HttpClient {
	hc := &HttpClient{}
	hc.log = log
	return hc
}

func (hc *HttpClient) GetFile(uri string, file string) ([]byte, error) {
	ueuri, err := url.QueryUnescape(uri)
	if err != nil {
		hc.log.Error("Unescape uri failed")
		return nil, err
	}

    resp, err := http.Get(ueuri + "/" + file)
    if err != nil {
		hc.log.Error("Http get %s failed", ueuri + "/" + file)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		hc.log.Error("Read body from response failed")
		return nil, err
    }

    clen := resp.Header.Get("Content-Length")
	if clen != "" {
		if cl, _ := strconv.Atoi(clen); len(body) != cl {
			hc.log.Error("Content length error")
			return nil, errors.New("Content length error")
		}
	}

	return body, nil
}

func (hc *HttpClient) Report(uri string, fid int, time string, ip string, msg string) error {
	ueuri, err := url.QueryUnescape(uri)
	if err != nil {
		hc.log.Error("Unescape uri failed")
		return err
	}

	if msg != "" {	/* error, status is 2 */
		uri = fmt.Sprintf(template, ueuri, fid, "cmd", "2",
				url.QueryEscape(time), url.QueryEscape(ip), url.QueryEscape(msg))
	} else {		/* ok, status is 1 */
		uri = fmt.Sprintf(template, ueuri, fid, "cmd", "1",
				url.QueryEscape(time), url.QueryEscape(ip), url.QueryEscape(msg))
	}

	hc.log.Debug("Http get request to %s", uri)
    resp, err := http.Get(uri)
    if err != nil {
		hc.log.Error("Report get uri %s failed", uri)
        return err
    }

    defer resp.Body.Close()

	if resp.StatusCode != HTTP_OK {
		hc.log.Error("Http status error: %d", resp.StatusCode)
        return errors.New(resp.Status)
	}

	hc.log.Debug("Http status is %d", resp.StatusCode)

	return nil
}
