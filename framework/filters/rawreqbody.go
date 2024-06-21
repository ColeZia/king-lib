package filters

import (
	"bytes"
	"io/ioutil"
	"log"
	netHttp "net/http"
	"sync"
)

var mutex sync.Mutex
var RawBody []byte

func RawReqBodyFilter(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		log.Println("RawReqBodyFilter:")
		var err error
		//将原始body暂存起来
		mutex.Lock()
		RawBody, err = ioutil.ReadAll(r.Body)
		mutex.Unlock()

		r.Body.Close()
		if err != nil {
			log.Println("ioutil.ReadAll(hReq.Body) err::", err)
		}

		//log.Println("RawBody::", string(RawBody))

		r.Body = ioutil.NopCloser(bytes.NewBuffer(RawBody))

		next.ServeHTTP(w, r)
	})
}
