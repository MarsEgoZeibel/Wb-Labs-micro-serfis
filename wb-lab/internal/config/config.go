package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB   DBConfig
	NATS NATSConfig
	HTTP HTTPConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type NATSConfig struct {
	ClusterID string
	ClientID  string
	Subject   string
}

type HTTPConfig struct {
	Port string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
		},
		NATS: NATSConfig{
			ClusterID: viper.GetString("NATS_CLUSTER_ID"),
			ClientID:  viper.GetString("NATS_CLIENT_ID"),
			Subject:   viper.GetString("NATS_SUBJECT"),
		},
		HTTP: HTTPConfig{
			Port: viper.GetString("HTTP_PORT"),
		},
	}, nil
}
