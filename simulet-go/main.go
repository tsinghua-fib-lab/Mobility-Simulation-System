package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.fiblab.net/sim/simulet-go/utils/config"
	easy "git.fiblab.net/utils/logrus-easy-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	configPath = flag.String("config", "", "config file path")
	cacheDir   = flag.String("cache", "", "input cache dir path (empty means disable cache)")

	log = logrus.WithField("module", "simulet")
)

func main() {
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02T15:04:05.9999",
		LogFormat:       "[%module%] [%time%] [%lvl%] %msg%\n",
	})
	flag.Parse()

	// 获取配置
	var config config.Config
	file, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("config file load err: %v", err)
	}
	if err := yaml.UnmarshalStrict(file, &config); err != nil {
		log.Fatalf("config file load err: %v", err)
	}
	log.Printf("%+v", config)

	// 优雅退出
	// 创建监听退出chan
	signalCh := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh
		log.Info("stopping...")
		go func() {
			<-signalCh
			os.Exit(1) // 强制结束
		}()
		os.Exit(0)
	}()

	// 启动模拟器
	Init(config)
	Run()
}
