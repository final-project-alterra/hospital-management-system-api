package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var format string = `
{
	"time":"${time_rfc3339_nano}",
	"id":"${id}",
	"remote_ip":"${remote_ip}",` +
	`"host":"${host}",
	"method":"${method}",
	"uri":"${uri}",
	"user_agent":"${user_agent}",` +
	`"status":${status},
	"error":"${error}",
	"latency":${latency},
	"latency_human":"${latency_human}"` + `,
	"bytes_in":${bytes_in},
	"bytes_out":${bytes_out}
}
`

func Logger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		Format:           format,
		CustomTimeFormat: "2006-01-02 15:04:05",
	})
}
