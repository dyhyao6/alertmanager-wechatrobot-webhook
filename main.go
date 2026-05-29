package main

import (
	"log"
	"os"

	"github.com/k8stech/alertmanager-wechatrobot-webhook/config"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/logger"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/notifier"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger, err := logger.New(cfg.LogDir, cfg.LogFile, cfg.LogMaxSize)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer appLogger.Close()

	appLogger.Info("Starting Alertmanager Webhook Server", logger.Fields{
		"addr": cfg.Addr,
	})

	// Initialize notifiers
	notifiers := map[string]notifier.Notifier{
		notifier.PlatformWeChat:   notifier.NewWeChatNotifier(),
		notifier.PlatformDingTalk: notifier.NewDingTalkNotifier(),
		notifier.PlatformFeiShu:   notifier.NewFeiShuNotifier(),
	}

	// Setup router
	r := router.New(cfg, notifiers, appLogger)
	engine := r.Setup()

	// Log endpoints
	appLogger.Info("Endpoints available", logger.Fields{
		"endpoints": []string{
			"POST /webhook - WeChat (legacy)",
			"POST /webhook/wecom - WeChat",
			"POST /webhook/dingtalk - DingTalk",
			"POST /webhook/feishu - FeiShu",
		},
	})

	appLogger.Info("Server listening", logger.Fields{
		"addr": cfg.Addr,
	})

	if err := engine.Run(cfg.Addr); err != nil {
		appLogger.Fatal("Server failed to start", logger.Fields{
			"error": err.Error(),
		})
		os.Exit(1)
	}
}