package middleware

import "net/http"

//LoggerWR ...
type LoggerWR struct {
	W          http.ResponseWriter
	StatusCode int
}

//Header wrapper for original ResponseWriter Header
func (lwr *LoggerWR) Header() http.Header {
	return lwr.W.Header()
}

//Write wrapper for original ResponseWriter Write
func (lwr *LoggerWR) Write(data []byte) (int, error) {
	return lwr.W.Write(data)
}

//WriteHeader wrapper for original ResponseWriter WriteHeader
func (lwr *LoggerWR) WriteHeader(statusCode int) {
	lwr.StatusCode = statusCode
	lwr.W.WriteHeader(statusCode)
}
