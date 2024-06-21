package filters

import (
	netHttp "net/http"
	"strings"

	"gl.king.im/king-lib/framework"
	"gl.king.im/king-lib/framework/config"
)

func CorsFilter(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		serviceConf := config.GetServiceConf()

		if len(serviceConf.Service.AllowOriginList) > 0 || serviceConf.Service.Cors != nil {
			allowedDomains := serviceConf.Service.AllowOriginList
			allowCredentials := "true"
			allowHeaders := []string{"content-type",
				strings.ToLower(framework.METADATA_KEY_AUTH_TOKEN),
				strings.ToLower(framework.METADATA_KEY_CALL_SCENARIO),
				strings.ToLower(framework.METADATA_KEY_PRISM_ACCESS_TOKEN),
			}

			if serviceConf.Service.Cors != nil {
				allowedDomains = serviceConf.Service.Cors.AllowOriginList

				if serviceConf.Service.Cors.AllowCredentials == 2 {
					allowCredentials = "false"
				}

				switch serviceConf.Service.Cors.AllowHeadersMode {
				case 2:
					allowHeaders = serviceConf.Service.Cors.AllowHeaders
				default:
					allowHeaders = append(allowHeaders, serviceConf.Service.Cors.AllowHeaders...)
				}
			}

			if origin := r.Header.Get("Origin"); origin != "" {
				for _, v := range allowedDomains {
					if origin == v {
						w.Header().Set("Access-Control-Allow-Methods", "*")
						w.Header().Set("Access-Control-Allow-Origin", origin)
						if len(allowHeaders) > 0 {
							w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
						}

						w.Header().Set("Access-Control-Allow-Credentials", allowCredentials)
						break
					}
				}

				if r.Method == netHttp.MethodOptions {
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
