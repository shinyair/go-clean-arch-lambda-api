package app

import (
	aws "github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/controller"
	bizcontroller "local.com/go-clean-lambda/internal/controller/biz"
	"local.com/go-clean-lambda/internal/logger"
	dynamodbrepo "local.com/go-clean-lambda/internal/repository/dynamodb"
	"local.com/go-clean-lambda/internal/usecase"
)

// InitControllers
//
//	@return []muxcontroller.MuxController
//	@return error
func InitDummyControllers() ([]controller.MuxController, error) {
	// init configs
	appConfig, err := NewAppConfig()
	if err != nil {
		return nil, errors.Errorf("failed to init app config. %s", err.Error())
	}
	// init logger
	logger.SetLogLevels(appConfig.LogCfg.Levels, appConfig.LogCfg.MinLevel, appConfig.LogCfg.CrNewline)
	// log configs after logger is inited
	logger.Info("app config: %s", logger.Pretty(appConfig))
	// init repo
	awsopt := awssession.Options{
		Config: aws.Config{Region: aws.String(appConfig.AwsEnvCfg.Region)},
	}
	if len(appConfig.AwsEnvCfg.Profile) > 0 {
		awsopt.Profile = appConfig.AwsEnvCfg.Profile
	}
	awssess := awssession.Must(awssession.NewSessionWithOptions(awsopt))
	dynamodbClient := awsdynamodb.New(awssess)
	dummyRepo := dynamodbrepo.NewDummyDynamodbRepo(
		appConfig.DynamodbCfg.DummyTableName,
		dynamodbClient)
	// init usecase
	dummyUsecase := usecase.NewDummyUseCase(dummyRepo)
	// init middlewares
	logMdf := controller.GetLogMiddleware()
	// init controllers
	dummyController := bizcontroller.NewDummyController(logMdf, dummyUsecase)
	return []controller.MuxController{
		dummyController,
	}, nil
}
