package filters

import (
	"log"
	netHttp "net/http"
)

func BaseFilter(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		log.Println("=========framework BaseFilter..=============.", r.Header, "--r.RequestURI::", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
