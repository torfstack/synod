package config

import "fmt"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func (dbCfg DBConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName)
}

type AuthConfig struct {
	DiscoveryURL string
}

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}
