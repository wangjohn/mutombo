package main

import (
	"github.com/zenazn/goji/web"
)

type GetProxiesReq struct {
	URLParams interface{} `json:"url_params"`
	JSONData  interface{} `json:"json_data"`
	Headers   interface{} `json:"headers"`
	Blocking  bool        `json:"blocking"`
}

func GetRequest(c web.C, w http.ResponseWriter, r *http.Request) {

}
