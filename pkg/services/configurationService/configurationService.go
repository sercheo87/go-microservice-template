package configurationService

import "go-microservice-template/pkg/config"

type Executor interface {
	GetGeneralConfiguration() config.AppConfig
}

type ConfigurationService struct {
}

func NewConfigurationService() ConfigurationService {
	return ConfigurationService{}
}

func (s ConfigurationService) GetGeneralConfiguration() config.AppConfig {
	return config.Configuration
}
