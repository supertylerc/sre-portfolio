package internal

import "github.com/spf13/viper"

func ViperConfig() {
	viper.AutomaticEnv()
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("REDIS_HOST", "127.0.0.1")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "example")
	viper.SetDefault("REDIS_LEADER_KEY", "leader:uuid")
	viper.SetDefault("METRICS_PORT", "9834")
}
