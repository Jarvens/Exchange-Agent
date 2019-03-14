// date: 2019-03-12
package config

const (
	MysqlProtocol    = "websocket"
	MysqlHostDev     = "127.0.0.1"
	MysqlPortDev     = "3306"
	MysqlUserNameDev = "root"
	MysqlPasswordDev = "root"
	MysqlDatabaseDev = "eoe"
	MysqlOptionsDev  = "charset=utf8&parseTime=True"
	MysqlDSLDev      = MysqlUserNameDev + ":" + MysqlPasswordDev + "@" + MysqlProtocol + "(" + MysqlHostDev + ":" + MysqlPortDev + ")/" + MysqlDatabaseDev + "?" + MysqlOptionsDev

	MysqlHostPro     = "127.0.0.1"
	MysqlPortPro     = "3306"
	MysqlUserNamePro = "root"
	MysqlPasswordPro = "root"
	MysqlDatabasePro = "eoe"
	MysqlOptionsPro  = "charset=utf8&parseTime=True"
	MysqlDSLPro      = MysqlUserNamePro + ":" + MysqlPasswordPro + "@" + MysqlProtocol + "(" + MysqlHostPro + ":" + MysqlPortPro + ")/" + MysqlDatabasePro + "?" + MysqlOptionsPro
)

func MysqlDSL() string {
	var mysqlDSL string
	switch Environment {
	case "DEVELOPMENT":
		mysqlDSL = MysqlDSLDev
	case "PRODUCTION":
		mysqlDSL = MysqlDSLPro
	default:
		mysqlDSL = MysqlDSLDev

	}
	return mysqlDSL
}
