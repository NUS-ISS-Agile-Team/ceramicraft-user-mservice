package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.SugaredLogger
)

// 默认日志等级DEBUG
// 默认日志输出路径: 根目录下test.log
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, getLogLevel())
	Logger = zap.New(core, zap.AddCaller()).Sugar()
}

func getLogLevel() zapcore.Level {
	level := zapcore.Level(0)
	if config.Config.LogConfig.Level != "" {
		if err := level.UnmarshalText([]byte(config.Config.LogConfig.Level)); err != nil {
			level = zapcore.DebugLevel // fallback to DebugLevel if there's an error
		}
	} else {
		level = zapcore.DebugLevel // default to DebugLevel if Level is not set
	}
	return level
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	cwd, _ := os.Getwd()
	logPath := filepath.Join(cwd, config.Config.LogConfig.FilePath)
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create directories: %v", err))
	}
	file, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(file)
}
