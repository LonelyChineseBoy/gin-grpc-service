package handler

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ReadConfigByYamlFile(filePath string, globalVal interface{}) {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Errorf("viper read config file failed %s \n", err)
	}
	if err := viper.Unmarshal(globalVal); err != nil {
		zap.S().Errorf("unmarshal conf failed, err:%s \n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("The configuration file is modified")
		if err := viper.Unmarshal(globalVal); err != nil {
			zap.S().Errorf("unmarshal conf failed, err:%s \n", err)
		}
	})
}

func ExitProcedure(addr string, handler http.Handler) {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	zap.S().Infof("Will Listen Server At %v", addr)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			zap.S().Infof("Listen: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	zap.S().Info("Shutdown Server......")
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	err := srv.Shutdown(timeout)
	if err != nil {
		zap.S().Errorf("Server Shutdown:%v", err)
	}
	zap.S().Info("Server Exiting")
}
