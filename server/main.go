package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/common/utils"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/config"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/grpc"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/http"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/log"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/mq"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository"
)

var (
	sigCh = make(chan os.Signal, 1)
)

func main() {
	config.Init()
	log.InitLogger()
	utils.InitJwtSecret()
	repository.Init()
	mq.InitKafka()
	go grpc.Init(sigCh)
	go http.Init(sigCh)
	// listen terminage signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh // Block until signal is received
	log.Logger.Infof("Received signal: %v, shutting down...", sig)
}
