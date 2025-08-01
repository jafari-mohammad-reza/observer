package conf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

func init() {

}

type Config struct {
	Port    int         `mapstructure:"port" validate:"required"`
	Kafka   KafkaConf   `mapstructure:"kafka"`
	Storage StorageConf `mapstructure:"storage"`
	Tracker TrackerConf `mapstructure:"tracker"`
}
type KafkaConf struct {
	Brokers         string   `mapstructure:"brokers" validate:"required"`
	ClientId        string   `mapstructure:"client_id" validate:"required"`
	LogChanSize     int      `mapstructure:"log_chan_size"`
	MutateChanSize  int      `mapstructure:"mutate_chan_size"`
	BootstrapTopics []string `mapstructure:"bootstrap_topics" `
	LogTopic        string   `mapstructure:"log_topic" `
	MutateTopic     string   `mapstructure:"mutate_topic" `
}
type StorageConf struct {
	Port    int    `mapstructure:"port"`
	WalPath string `mapstructure:"wal_path"`
}
type TrackerConf struct {
	Port int `mapstructure:"port"`
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.WatchConfig()

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.New("failed to unmarshal config: " + err.Error())
	}
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		var sb strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("Field '%s' failed on '%s'\n", err.Field(), err.Tag()))
		}
		return nil, errors.New(sb.String())
	}
	return &cfg, nil
}
