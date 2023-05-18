package metrics

type Config struct {
	Port        int    `mapstructure:"port" envconfig:"METRICS_PORT" valid:"required"`
	Path        string `mapstructure:"path" valid:"required"`
	AppName     string `mapstructure:"app_name" valid:"required"`
	Environment string `mapstructure:"environment" valid:"required"`
}
