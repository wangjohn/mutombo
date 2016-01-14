package main

import (
	"github.com/wangjohn/gowebutils"
	"github.com/zenazn/goji/web"
)

type GetProxiesReq struct {
	URLParams interface{} `json:"url_params"`
	JSONData  interface{} `json:"json_data"`
	Headers   interface{} `json:"headers"`
	Blocking  bool        `json:"blocking"`
}

func GetRequest(c web.C, w http.ResponseWriter, r *http.Request) {
	requestId := c.URlParams["request_id"]
	if requestId == "" {
		err = errors.New("You must specify a request_id when retrieving a request")
		gowebutils.SendError(w, err)
	}
}

func PostRequest(w http.ResponseWriter, r *http.Request) {

}
