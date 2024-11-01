package mq

type RabbitConfig struct {
	Path        string `yaml:"path"`        // 服务器地址:端口
	UserName    string `yaml:"userName"`    // 用户名
	PassWord    string `yaml:"passWord"`    // 密码
	VirtualHost string `yaml:"virtualHost"` // 访问的虚拟主机
}

func (r *RabbitConfig) IsValue() bool {
	if len(r.Path) != 0 && len(r.UserName) != 0 && len(r.PassWord) != 0 && len(r.VirtualHost) != 0 {
		return true
	}
	return false
}
