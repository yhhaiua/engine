package db

type MysqlConfig struct {
	Path         string `yaml:"path"`                             	// 服务器地址:端口
	Config       string `yaml:"config"`                       		// 高级配置
	Dbname       string `yaml:"db-name"`                     	  	// 数据库名
	Username     string `yaml:"username"`                 			// 数据库用户名
	Password     string `yaml:"password"`                 			// 数据库密码
	MaxIdleConns int    `yaml:"max-idle-conns"` 					// 空闲中的最大连接数
	MaxOpenConns int    `yaml:"max-open-conns"` 					// 打开到数据库的最大连接数
}


func (m *MysqlConfig) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
}
