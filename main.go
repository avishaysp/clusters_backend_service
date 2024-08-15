package main

import (
	"log"
	"os"

	"github.com/TAULargeScaleWorkshop/RLAD/config"
	RegService "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/service"
	TestService "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/service"
	CacheService "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/service"

	"github.com/TAULargeScaleWorkshop/RLAD/utils"
	"gopkg.in/yaml.v2"
)

func main() {
	// read configuration file from command line argument
	if len(os.Args) != 2 {
		utils.Logger.Fatal("Expecting exactly one configuration file")
		os.Exit(1)
	}
	configFile := os.Args[1]
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
		os.Exit(2)
	}

	var config config.BaseConfig
	err = yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		log.Fatalf("error unmarshaling data: %v", err)
		os.Exit(3)
	}
	utils.Logger.Printf("Loading service type: %v\n", config.Type)
	switch config.Type {
	case "TestService":
		TestService.Start(configData)
	case "RegService":
		RegService.Start(configData)
	case "CacheService":
		CacheService.Start(configData)
	default:
		utils.Logger.Fatalf("Unknown configuration type: %v", config.Type)
		os.Exit(4)
	}
}
