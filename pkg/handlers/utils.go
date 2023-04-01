package handlers

import "net/http"

func connectingIP(headers http.Header) string {

	var connectingIP string
	if _, ok := headers["Cf-Connecting-Ip"]; ok {
		connectingIP = headers.Get("Cf-Connecting-Ip")
	} else {
		connectingIP = headers.Get("X-FORWARDED-FOR")
	}
	return connectingIP
}
