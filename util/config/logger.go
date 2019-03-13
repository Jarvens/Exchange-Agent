// date: 2019-03-12
package config

const (
	AccessLogFilePath      = "log/access"
	AccessLogFileExtension = ".log"
	AccessLogMaxSize       = 5
	AccessLogMaxBackups    = 7
	AccessLogMaxAge        = 30
	ErrorLogFilePath       = "log/error"
	ErrorLogFileExtension  = ".log"
	ErrorLogMaxSize        = 10
	ErrorLogMaxBackups     = 7
	ErrorLogMaxAge         = 30

	WsAccessLogFilePath      = "log/ws_access"
	WsAccessLogFileExtension = ".log"
	WsAccessLogMaxSize       = 5
	WsAccessLogMaxBackups    = 7
	WsAccessLogMaxAge        = 30

	WsErrorLogFilePath      = "log/ws_access"
	WsErrorLogFileExtension = ".log"
	WsErrorLogMaxSize       = 5
	WsErrorLogMaxBackups    = 7
	WsErrorLogMaxAge        = 30
)
