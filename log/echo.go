package log

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func EchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			bytesIn := req.Header.Get(echo.HeaderContentLength)
			if bytesIn == "" {
				bytesIn = "0"
			}

			fields := log.WithFields(log.Fields{
				"id":            id,
				"remote_ip":     req.RemoteAddr,
				"real_ip":       c.RealIP(),
				"host":          req.Host,
				"method":        req.Method,
				"uri":           req.RequestURI,
				"user_agent":    req.UserAgent(),
				"status":        res.Status,
				"latency":       strconv.FormatInt(int64(stop.Sub(start)), 10),
				"latency_human": stop.Sub(start).String(),
				"bytes_in":      bytesIn,
				"bytes_out":     strconv.FormatInt(res.Size, 10),
			})

			var m string

			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					if hs, ok := he.Message.(echo.Map); ok {
						m = hs["message"].(string)
					} else {
						m = he.Message.(string)
					}
				} else {
					m = err.Error()
				}
			}

			switch res.Status / 100 {
			case 1: // 1xx Informational Status Codes
				fields.Info("Informational")
			case 2: // 2xx Successful Status Code
				fields.Info("Success")
			case 3: // 3xx Redirection Status Code
				fields.Info("Redirection")
			case 4: // 4xx Client HTTPError Status Code
				fields.Info(m)
			case 5: // 5xx Server HTTPError Status Code
				fields.Error(m)
			default: // No Standard Status Code
				fields.Error("Status code not specified")
			}
			return
		}
	}
}
