package cofing

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig `mapstructure:"app"`
	Database Database  `mapstructure:"database"`
	SAML     SAML      `mapstructure:"saml"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type Database struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

type SAML struct {
	SP  SP  `mapstructure:"sp"`
	IDP IDP `mapstructure:"idp"`
}

type SP struct {
	RootURL  string `mapstructure:"root-url"`
	CertFile string `mapstructure:"certfile"`
	KeyFile  string `mapstructure:"keyfile"`
}

type IDP struct {
	MetadataURL string `mapstructure:"metadata-url"`
}

func LoadConfig() *Config {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v", err))
	}

	return &config
}
