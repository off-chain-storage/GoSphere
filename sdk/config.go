package sdk

import "github.com/spf13/viper"

type Config struct {
	Path string `mapstructure:"PM_ADDR"`
}

func loadConfig() (c *Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		log.Error("failed to read config file: ", err)
		return nil, err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Error("failed to unmarshal config: ", err)
		return nil, err
	}

	return c, nil
}
