package logger

type Config struct {
	Level  string `mapstructure:"level" envconfig:"LOG_LEVEL" valid:"required"`
	Format string `mapstructure:"format" envconfig:"LOG_FORMAT" valid:"required"`
}
