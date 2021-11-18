package util

import (
	"os"

	"github.com/spf13/viper"
)

// Config for Application.
type Config struct {
	Aksk struct {
		RegionID        string `required:"true" env:"AKSK_REGIONID"`
		AccessKeyID     string `required:"true" env:"AKSK_ACCESSKEYID"`
		AccessKeySecret string `required:"true" env:"AKSK_ACCESSKEYSECRET"`
	}

	Pool struct {
		Size int `default:"3" env:"POOL_SIZE"`
	}

	Log struct {
		Level string `default:"debug"`
	}
}

func LoadConfig(cfgFile string) (*Config, error) {
	var err error = nil
	config := &Config{}

	viper := viper.New()
	viper.SetDefault("pool.size", 3)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			CheckErr(os.Stderr, err)
		}
		// Search config in home directory with name ".arctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".arctl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		CheckErr(os.Stderr, err)
	}

	if err := viper.Unmarshal(config); err != nil {
		CheckErr(os.Stderr, err)
	}

	return config, err
}
