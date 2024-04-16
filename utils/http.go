package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const MAX_RETRY_TIMES = 3

type CommonReq struct {
}

func HTTPRequest(url string, method string, req, rsp interface{}) error {
	var resp *http.Response
	var err error

	var bodyBuf *bytes.Reader
	if content, ok := req.(*string); ok {
		bodyBuf = bytes.NewReader([]byte(*content))
	} else {
		buf, err := json.Marshal(req)
		if err != nil {
			return err
		}
		bodyBuf = bytes.NewReader(buf)
	}

	for i := 1; i <= MAX_RETRY_TIMES; i++ {
		if strings.ToUpper(method) == "GET" {
			resp, err = http.Get(url)
		} else {
			resp, err = http.Post(url, "application/json", bodyBuf)
		}
		if err != nil {
			if i == MAX_RETRY_TIMES {
				return err
			}
			<-time.After(time.Second * 2)
			continue
		}
		break
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read body failed: %s\n", err)
		return err
	}

	if err = json.Unmarshal([]byte(body), rsp); err != nil {
		fmt.Printf("json unmarshr failed: %v, body: %s, url: %s", err, body, url)
		return err
	}

	return nil
}

func HTTPGet(url string, rsp interface{}) error {
	return HTTPRequest(url, "GET", new(CommonReq), rsp)
}

func HTTPPost(url string, req, rsp interface{}) error {
	return HTTPRequest(url, "POST", req, rsp)
}
