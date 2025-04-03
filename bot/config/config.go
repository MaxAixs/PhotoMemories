package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken string
	RedisConfig      string    `mapstructure:"redisConfig" yaml:"redisConfig"`
	DBConfig         DBConfig  `mapstructure:"dbConfig" yaml:"dbConfig"`
	Cmd              Commands  `mapstructure:"commands" yaml:"commands"`
	Buttons          Buttons   `mapstructure:"buttons" yaml:"buttons"`
	Msg              Messages  `mapstructure:"messages" yaml:"messages"`
	AWSConfig        AWSConfig `mapstructure:"awsConfig" yaml:"awsConfig"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string
	DBName   string `mapstructure:"db_name" yaml:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
}

type Commands struct {
	Start string `mapstructure:"start" yaml:"start"`
	Help  string `mapstructure:"help" yaml:"help"`
}

type Buttons struct {
	AddPic string `mapstructure:"addPic" yaml:"addPic"`
	DelPic string `mapstructure:"deletePic" yaml:"deletePic"`
	GetPic string `mapstructure:"getPic" yaml:"getPic"`
	MyTags string `mapstructure:"myTags" yaml:"myTags"`
}

type Messages struct {
	Response `mapstructure:"response" yaml:"response"`
	Errors   `mapstructure:"errors" yaml:"errors"`
	Success  `mapstructure:"success" yaml:"success"`
}

type Response struct {
	Start     string `mapstructure:"start" yaml:"start"`
	Help      string `mapstructure:"help" yaml:"help"`
	Default   string `mapstructure:"help" yaml:"help"`
	AddPic    string `mapstructure:"addPhoto" yaml:"addPhoto"`
	DelPic    string `mapstructure:"deletePhoto" yaml:"deletePhoto"`
	GetPic    string `mapstructure:"getPic" yaml:"getPic"`
	GetMyTags string `mapstructure:"getMyTags" yaml:"getMyTags"`
}

type Success struct {
	PicSaved   string `mapstructure:"pic_saved" yaml:"pic_saved"`
	TagSaved   string `mapstructure:"tag_saved" yaml:"tag_saved"`
	PicDeleted string `mapstructure:"pic_deleted" yaml:"pic_deleted"`
}

type AWSConfig struct {
	Bucket          string `mapstructure:"bucket" yaml:"bucket"`
	Region          string `mapstructure:"region" yaml:"region"`
	AccessKey       string `mapstructure:"accessKey"`
	AccessSecretKey string `mapstructure:"accessSecretKey"`
}

type Errors struct {
	ErrUploadPic string `mapstructure:"err_upload_pic" yaml:"err_upload_pic"`
	ErrUploadTag string `mapstructure:"err_upload_tag" yaml:"err_upload_tag"`
	ErrDelPic    string `mapstructure:"err_del_pic" yaml:"err_del_pic"`
	ErrGetPic    string `mapstructure:"err_get_pic" yaml:"err_get_pic"`
	ErrGetTags   string `mapstructure:"err_get_tags" yaml:"err_get_tags"`
}

func InitConfig() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var config Config

	if err := unmarshal(&config); err != nil {
		return nil, err
	}

	if err := fromEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setUpViper() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("internal/config")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	return nil
}

func unmarshal(cfg *Config) error {
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("unable to decode config, %v", err)
	}

	return nil
}

func fromEnv(cfg *Config) error {
	viper.AutomaticEnv()

	cfg.TelegramBotToken = viper.GetString("TELEGRAM_TOKEN")
	if cfg.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_TOKEN is not set")
	}

	cfg.DBConfig.Password = viper.GetString("DB_PASSWORD")
	if cfg.DBConfig.Port == "" {
		return fmt.Errorf("DB_PASSWORD is not set")
	}

	cfg.AWSConfig.AccessKey = viper.GetString("AWS_ACCESS_KEY_ID")
	cfg.AWSConfig.AccessSecretKey = viper.GetString("AWS_SECRET_ACCESS_KEY")
	if cfg.AWSConfig.AccessKey == "" || cfg.AWSConfig.AccessSecretKey == "" {
		return fmt.Errorf("AWS_ACCESS_KEY or AWS_SECRET_KEY is not set")
	}

	return nil
}
