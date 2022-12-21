package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
	"local.com/go-clean-lambda/internal/app"
	"local.com/go-clean-lambda/internal/controller"
	"local.com/go-clean-lambda/internal/logger"
)

//go:embed .\..\..\configs\env.yml
var yf []byte

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

func server() {
	setLocalEnv()
	port := 8080
	controllers, err := app.InitDummyControllers()
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
