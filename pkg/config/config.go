package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func Init(rootPath string) error {
	configFilePath, err := filepath.Abs(rootPath + "pkg/config")
	if err != nil {
		logrus.Fatal(err)
	}
	viper.AddConfigPath(configFilePath)

	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	env := os.Getenv("HOST")

	switch env {
	case "stage", "prod":
		viper.SetConfigName(env)
		return viper.MergeInConfig()
	default:
		return nil
	}
}