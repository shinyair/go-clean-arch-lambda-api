package app

import (
	"os"
	"strings"
)

type Config struct {
	Appcode     string
	Variant     string
	Stage       string
	AwsEnvCfg   *AwsEnvConfig
	LogCfg      *LogConfig
	DynamodbCfg *DynamodbConfig
}

type LogConfig struct {
	Levels    []string
	MinLevel  string
	CrNewline bool
}

type AwsEnvConfig struct {
	AccountID string
	Region    string
	Profile   string
}

type DynamodbConfig struct {
	DummyTableName string
}

func NewAppConfig() (*Config, error) {
	awsEnvConfig := &AwsEnvConfig{
		AccountID: os.Getenv("ACCOUNT_ID"),
		Region:    os.Getenv("AWS_REGION"),
		Profile:   os.Getenv("AWS_PROFILE"),
	}
	logConfig := &LogConfig{
		Levels:    strings.Split(os.Getenv("LOG_LEVELS"), ","),
		MinLevel:  os.Getenv("LOG_MIN_LEVEL"),
		CrNewline: os.Getenv("LOG_CR_NEWLINE") == "true",
	}
	dynamodbConfig := &DynamodbConfig{
		DummyTableName: os.Getenv("DUMMY_TABLE_NAME"),
	}
	appConfig := Config{
		Appcode:     os.Getenv("APPCODE"),
		Variant:     os.Getenv("VARIANT"),
		Stage:       os.Getenv("STAGE"),
		AwsEnvCfg:   awsEnvConfig,
		LogCfg:      logConfig,
		DynamodbCfg: dynamodbConfig,
	}
	return &appConfig, nil
}
