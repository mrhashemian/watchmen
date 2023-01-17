package config

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const envPrefix = "watchmen"

var C *Config

type Config struct {
	Server          Server      `yaml:"server"`
	Logger          Logger      `yaml:"logger"`
	BaseAPIDatabase SQLDatabase `yaml:"base_api_database"`
	User            User        `yaml:"user"`
}

type Server struct {
	Address  string         `yaml:"address"`
	Location *time.Location `yaml:"location"`
}

type Logger struct {
	Level string `yaml:"level" validate:"required"`
	Tag   string `yaml:"tag" validate:"required"`
}

// SQLDatabase is a configuration structure
type SQLDatabase struct {
	Driver        string        `yaml:"driver"`
	Host          string        `yaml:"host"`
	Port          int           `yaml:"port"`
	DB            string        `yaml:"db"`
	User          string        `yaml:"user"`
	Password      string        `yaml:"password"`
	Location      string        `yaml:"location"`
	MaxConn       int           `yaml:"max_conn"`
	IdleConn      int           `yaml:"idle_conn"`
	Timeout       time.Duration `yaml:"timeout"`
	DialRetry     int           `yaml:"dial_retry"`
	DialTimeout   time.Duration `yaml:"dial_timeout"`
	ReadTimeout   time.Duration `yaml:"read_timeout"`
	WriteTimeout  time.Duration `yaml:"write_timeout"`
	UpdateTimeout time.Duration `yaml:"update_timeout"`
	DeleteTimeout time.Duration `yaml:"delete_timeout"`
}

// User config struct
type User struct {
	PasswordHashCost int `yaml:"password_hash_cost"`
}

type OTP struct {
	Len int           `yaml:"len"`
	TTL time.Duration `yaml:"ttl"`
}

// String representation of config struct
func (c *Config) String() string {
	s := fmt.Sprintf(
		"Log Level:\t%s\nServer Address:\t%s\nBase API DB:\t%s\n",
		c.Logger.Level,
		c.Server.Address,
		c.BaseAPIDatabase.String(),
	)

	return s
}

// String returns SQLDatabase formatted DSN
func (d *SQLDatabase) String() string {
	switch d.Driver {
	case "mysql":
		return d.mysqlDSN()
	}

	panic("SQLDatabase driver is not supported")
}

func (d *SQLDatabase) mysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_general_ci&loc=%s", d.User, d.Password, d.Host, d.Port, d.DB, url.QueryEscape(d.Location))
}

var failFunc = log.Fatalf

func Init(filename string) *Config {
	c := new(Config)
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if filename != "" {
		v.SetConfigFile(filename)
		if err := v.MergeInConfig(); err != nil {
			failFunc("opening config file [%s] failed: %v", filename, err)
		} else {
			log.Infof("config file [%s] opened successfully", filename)
		}
	}

	err := v.Unmarshal(c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			TimeLocationDecodeHook(),
		)
	})
	if err != nil {
		failFunc("failed on config unmarshal: %v", filename, err)
	}

	C = c

	return c
}

func TimeLocationDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		var timeLocation *time.Location
		if t != reflect.TypeOf(timeLocation) {
			return data, nil
		}

		return time.LoadLocation(data.(string))
	}
}
