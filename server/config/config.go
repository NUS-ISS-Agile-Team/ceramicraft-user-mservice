package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	Config = &Conf{}
)

type Conf struct {
	GrpcConfig  *GrpcConfig  `mapstructure:"grpc"`
	LogConfig   *LogConfig   `mapstructure:"log"`
	HttpConfig  *HttpConfig  `mapstructure:"http"`
	MySQLConfig *MySQL       `mapstructure:"mysql"`
	EmailConfig *EmailConfig `mapstructure:"email"`
	KafkaConfig *KafkaConfig `mapstructure:"kafka"`
}

type EmailConfig struct {
	SmtpHost      string `mapstructure:"smtp_host"`
	SmtpEmailFrom string `mapstructure:"smtp_email_from"`
	SmtpPass      string `mapstructure:"smtp_pass"`
}

type HttpConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

type GrpcConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	ConnectTimeout int    `mapstructure:"connect_timeout"`
	MaxPoolSize    int    `mapstructure:"max_pool_size"`
}

type MySQL struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
}

type KafkaConfig struct {
	Brokers            []string `mapstructure:"brokers"`
	UserActivatedTopic string   `mapstructure:"user_activated_topic"`
	MaxBytes           int      `mapstructure:"max_bytes"`
	Acks               int      `mapstructure:"acks"`
	Retries            int      `mapstructure:"retries"`
	BatchSize          int      `mapstructure:"batch_size"`
	BatchTimeoutMillis int      `mapstructure:"batch_timeout_millis"`
}

func Init() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/resources")
	viper.AddConfigPath(workDir)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword != "" {
		Config.MySQLConfig.Password = mysqlPassword
	} else {
		panic("MYSQL_PASSWORD environment variable is not set")
	}
	Config.EmailConfig.SmtpPass = os.Getenv("SMTP_PASSWORD")
	Config.EmailConfig.SmtpEmailFrom = os.Getenv("SMTP_EMAIL_FROM")
}
