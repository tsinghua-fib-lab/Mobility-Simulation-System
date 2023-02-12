package config

// 配置获取

var (
	C Control
)

func Init(config Config) {
	C = config.Control
}
