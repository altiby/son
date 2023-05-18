package goconf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Credentials struct {
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

// Opts for configuration.
type Opts struct {
	// Directory to look for configuration files. Default: "."
	Dir string
	// FileName of config without extension.
	FileName string
	// Env variable name which contains config filename without extension.
	EnvFileName string
	// Format of config. Available values are defined by Viper lib.
	// Default: "yaml"
	Format string
}

// NewConfig ...
func NewConfig(opts Opts, c interface{}) error {
	if opts.Dir == "" {
		opts.Dir = "."
	}

	if opts.Format == "" {
		opts.Format = "yaml"
	}

	v := viper.New()
	v.SetConfigType(opts.Format)
	v.AddConfigPath(opts.Dir)

	path, err := getFilename(opts.EnvFileName, opts.FileName)
	if err != nil {
		return err
	}

	v.SetConfigName(path)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(c); err != nil {
		return err
	}

	if err := envconfig.Process("", c); err != nil {
		return err
	}

	_, err = govalidator.ValidateStruct(c)
	return err
}

func getFilename(envvar, filepath string) (string, error) {
	if filepath != "" {
		return filepath, nil
	}

	if envvar != "" {
		if path, found := os.LookupEnv(envvar); found {
			return path, nil
		}
	}

	return "", errors.New("config not specified")
}

func Load(path, env string, cf interface{}) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	err = NewConfig(Opts{
		EnvFileName: env,
		Dir:         absPath,
	}, cf)

	return err
}

func PopulateConfig(cfg interface{}, path, envKey, config string) error {
	var (
		defaultEnvKey     = "CONFIG"
		defaultConfigPath = "configs"
		defaultConfig     = "local"
	)

	if config == "" {
		config = defaultConfig
	}

	if envKey == "" {
		envKey = defaultEnvKey
	}

	if path == "" {
		path = defaultConfigPath
	}

	if _, ok := os.LookupEnv(envKey); !ok {
		if err := os.Setenv(envKey, config); err != nil {
			return err
		}
	}

	if err := Load(path, envKey, cfg); err != nil {
		return err
	}

	return nil
}

func LoadBaseCredentials(prefix string) Credentials {
	creds := Credentials{}

	if val, ok := os.LookupEnv(fmt.Sprintf("%s_USERNAME", strings.ToUpper(prefix))); ok {
		creds.Username = val
	}

	if val, ok := os.LookupEnv(fmt.Sprintf("%s_PASSWORD", strings.ToUpper(prefix))); ok {
		creds.Password = val
	}

	return creds
}
