package configuration

var appConfig *Config

func SetAppConfig(config Config) {
	if appConfig == nil || !appConfig.seal {
		appConfig = &config
		appConfig.seal = true
	}
}

func GetAppConfig() *Config {
	return appConfig
}
