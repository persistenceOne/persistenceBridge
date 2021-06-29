package configuration

var appConfig *Config

func SetAppConfiguration(config Config) {
	appConfig = &config
}

func GetAppConfiguration() *Config {
	return appConfig
}
