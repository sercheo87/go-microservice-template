package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Configuration AppConfig

type AppConfig struct {
	MicroserviceName    string `mapstructure:"microserviceName" yaml:"microserviceName"`
	MicroserviceServer  string `mapstructure:"microserviceServer" yaml:"microserviceServer"`
	MicroservicePort    string `mapstructure:"microservicePort" yaml:"microservicePort"`
	MicroserviceVersion string `mapstructure:"microserviceVersion" yaml:"microserviceVersion"`
	Environment         string `mapstructure:"environment" yaml:"environment"`
	WorkspaceFolder     string `mapstructure:"workspaceFolder" yaml:"workspaceFolder"`
	Log                 Log    `mapstructure:"log" yaml:"log"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

func LoadConfigurationMicroservice(path string) {
	fmt.Println("Loading configuration from file [application.yml]")
	viper.SetConfigName("application")
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		panic("Error reading config file")
	}

	// Set undefined variables
	hostname, _ := os.Hostname()
	viper.SetDefault("microserviceServer", hostname)
	viper.SetDefault("microservicePathRoot", "./")

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}

	fmt.Println("Configuration loaded")
}
