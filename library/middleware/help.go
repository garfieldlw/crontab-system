package middleware

import (
	"bytes"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		rw := &ResponseWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = rw

		data, err := c.GetRawData()
		if err != nil {
			log.Warn("Request GetRawData Err", zap.Error(err))
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		c.Next()

		log.Info(c.HandlerName(), zap.Any("request", string(data[:])), zap.Any("response", rw.Body.String()))
	}
}

type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}
