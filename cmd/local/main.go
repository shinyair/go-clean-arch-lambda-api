package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	aws "github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"local.com/go-clean-lambda/internal/app"
	"local.com/go-clean-lambda/internal/controller"
	"local.com/go-clean-lambda/internal/controller/api/car"
	"local.com/go-clean-lambda/internal/controller/api/pet"
	"local.com/go-clean-lambda/internal/logger"
	"local.com/go-clean-lambda/internal/repository"
	"local.com/go-clean-lambda/internal/sdk/account"
	"local.com/go-clean-lambda/internal/sdk/authentication"
	"local.com/go-clean-lambda/internal/sdk/authorization"
	"local.com/go-clean-lambda/internal/usecase"
)

//nolint: all
//
//go:embed env.yml
var yf []byte

//nolint: all
//
//go:embed jwt.rsa
var privateKey []byte

//nolint: all
//
//go:embed jwt.rsa.pub
var publicKey []byte

func setLocalEnv() {
	envCfgMap := make(map[string]map[string]string)
	err := yaml.Unmarshal(yf, envCfgMap)
	if err != nil {
		log.Fatalf("failed to parse yaml config file: %v", err)
	}
	localCfgMap := envCfgMap["local"]
	for key, value := range localCfgMap {
		os.Setenv(key, value)
		// cannot filter log level now, because logger is not inited yet
		logger.Info("set mocked env: %s, %s", key, value)
	}
}

func initControllers() ([]controller.MuxController, error) {
	// init configs
	appConfig, err := app.NewAppConfig()
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
	// mock ssm
	store := map[string]string{
		appConfig.AuthCfg.PrivateKey: string(privateKey),
		appConfig.AuthCfg.PublicKey:  string(publicKey),
	}
	localSSMClient := NewLocalSSM(store)
	// init repo
	dummyRepo := repository.NewDummyDynamodbRepo(
		appConfig.DynamodbCfg.DummyTableName,
		dynamodbClient)
	// init usecase
	dummyUsecase := usecase.NewDummyUseCase(dummyRepo)
	// init sdk clients
	jwtClient := authentication.NewAuthJwtDummyClient(
		appConfig.AuthCfg.PublicKey,
		appConfig.AuthCfg.PrivateKey,
		localSSMClient,
	)
	roleClient := authorization.NewRoleDummyClient()
	userClient := account.NewUserDummmyClient()
	// init middlewares
	logMdf := controller.GetLogMiddleware()
	authMdf := controller.GetLoginAccessMiddleware(jwtClient)
	rolePingMdf := controller.GetRoleAccessMiddleware([]uint64{uint64(controller.AuthIndexAppPing)})
	// init controllers
	authController := controller.NewAuthController(logMdf, authMdf, jwtClient, roleClient, userClient)
	dummyController := controller.NewDummyController(logMdf, dummyUsecase)
	pingController := controller.NewPingController(logMdf, authMdf, rolePingMdf)
	apiPetController := pet.NewPetController(logMdf)
	apiCarController := car.NewCarController(logMdf)
	return []controller.MuxController{
		apiPetController,
		apiCarController,
		authController,
		dummyController,
		pingController,
		apiPetController,
	}, nil
}

func server() {
	setLocalEnv()
	port := 8080
	controllers, err := initControllers()
	if err != nil {
		logger.Error("execution end. failed to init lambda.", err)
		return
	}
	r := controller.NewRouter(controllers)
	logger.Info("local(with mux) initialization done on port: %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func main() {
	server()
}
