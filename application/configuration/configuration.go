package configuration

var appConfig *Config

func SetAppConfig(config Config) {
	appConfig = &config
}

func GetAppConfig() *Config {
	return appConfig
}
