package gredis


//RedisConfig 配置结构
type RedisConfig struct {
	Path         string `yaml:"path"`                             	// 服务器地址:端口
	Password     string `yaml:"password"`                 			// 数据库密码
	MaxIdleConns int    `yaml:"max-idle-conns"` 					// 空闲中的最大连接数
	MaxOpenConns int    `yaml:"max-open-conns"` 					// 打开到数据库的最大连接数
	DbNum		 int	`yaml:"db-num"`								// 使用的db编号
}
