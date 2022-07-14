package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"local.com/go-clean-lambda/internal/app"
	"local.com/go-clean-lambda/internal/logger"

	"gopkg.in/yaml.v2"
	"local.com/go-clean-lambda/internal/controller"
)

func setLocalEnv() {
	yf, err := ioutil.ReadFile("configs/env.yml")
	if err != nil {
		log.Fatalf("faild to get yaml config file. %v ", err)
	}
	envCfgMap := make(map[string]map[string]string)
	err = yaml.Unmarshal(yf, envCfgMap)
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

func server() {
	setLocalEnv()
	controllers, err := app.InitDummyControllers()
	if err != nil {
		logger.Error("execution end. failed to init lambda.", err)
		return
	}
	r := controller.NewRouter(controllers)
	logger.Info("local(with mux) initialization done")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func main() {
	server()
}
