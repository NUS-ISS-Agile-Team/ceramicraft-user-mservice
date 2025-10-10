package main

import (
	"fmt"
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %v", r)
		}
		fmt.Println("=== app exit")
	}()
	fmt.Println("Starting ceramicraft-user-mservice...")
	config.Init()
	log.InitLogger()
	fmt.Println("Logger initialized.")
	log.Logger.Info("Logger initialized.")
	utils.InitJwtSecret()
	fmt.Println("JWT secret initialized.")
	log.Logger.Info("JWT secret initialized.")
	repository.Init()
	fmt.Println("Database initialized.")
	log.Logger.Info("Database initialized.")
	mq.InitKafka()
	fmt.Println("Kafka initialized.")
	log.Logger.Info("Kafka initialized.")
	go grpc.Init(sigCh)
	go http.Init(sigCh)
	// listen terminage signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh // Block until signal is received
	log.Logger.Infof("Received signal: %v, shutting down...", sig)
	fmt.Printf("Received signal: %v shutting down...", sig)
}
