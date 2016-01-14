package kfkc

import (
	"fmt"
	"net/http"
)

type HttpClient struct {
	log		*Log
}

var template = "%s?flow_id=%s&sid=%s&cmd=%s&status=%s&time=%s&ip=%s"

func InitHttpClient(log *Log) *HttpClient {
	hc := &HttpClient{}
	hc.log = log
	return hc
}


func (hc *HttpClient) Getfile() {
    resp, err := http.Get("http://10.73.31.119:8001")
    if err != nil {
        fmt.Println(err)
        return
        // handle error
    }

    defer resp.Body.Close()
    _, err = ioutil.ReadAll(resp.Body)
    //body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }

    //fmt.Println(string(body))
    fmt.Println(resp.Header.Get("Content-Length"))

}

func (hc *HttpClient) Report(uri string, fid string, sid string, time string, ip string, msg string) {
	if len(msg) {	/* error, status is 2 */
		uri = fmt.Sprintf(template, uri, fid, sid, "cmd", "2", time, ip)
	} else {		/* ok, status is 1 */
		uri = fmt.Sprintf(template, uri, fid, sid, "cmd", "1", time, ip)
	}

	//TODO: http request
    resp, err := http.Get("http://10.73.31.119:8001")
    if err != nil {
        fmt.Println(err)
        return
        // handle error
    }

    defer resp.Body.Close()
    _, err = ioutil.ReadAll(resp.Body)
    //body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }

    //fmt.Println(string(body))
    fmt.Println(resp.Header.Get("Content-Length"))

}
