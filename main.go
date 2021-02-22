package main

import (
	"fmt"
	"github.com/aasumitro/svc-telegram-notify/config"
	"github.com/aasumitro/svc-telegram-notify/src/handlers"
	"runtime"
)

func main() {
	// sets the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(runtime.NumCPU())
	// initialize and setup app configuration
	appConfig := config.InitAppConfig()
	// show current service name and version
	fmt.Println(fmt.Sprintf(
		"starting %s version %s . . . .",
		appConfig.GetAppName(),
		appConfig.GetAppVersion(),
	))
	// load server environment
	appConfig.SetupAMQPConnection()
	appConfig.SetupTelegramConnection()

	handlers.NewTelegramCommandHandler(appConfig)
}
