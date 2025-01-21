package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	DBName   string `yaml:"dbname" validate:"required"`
}

func (dbCfg DBConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName)
}

type AuthConfig struct {
	Issuer       string `yaml:"issuer" validate:"required"`
	ClientID     string `yaml:"clientId" validate:"required"`
	ClientSecret string `yaml:"clientSecret" validate:"required"`
	RedirectURL  string `yaml:"redirectUrl" validate:"required"`
	BaseURL      string `yaml:"baseUrl" validate:"required"`
}

type Config struct {
	DB   DBConfig   `yaml:"db" validate:"required"`
	Auth AuthConfig `yaml:"auth" validate:"required"`
}

func ParseFile(path string) (*Config, error) {
	viper.AddConfigPath(".")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
