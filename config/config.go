package config

import (
	"flag"
	"os"
)

// Config holds all configuration for the webhook server
type Config struct {
	// Server config
	Addr string

	// Platform keys
	WeChatKey   string
	DingTalkKey string
	FeiShuKey   string

	// Log config
	LogDir     string
	LogFile    string
	LogMaxSize int64
}

// Load loads configuration from command line flags
func Load() *Config {
	h := flag.Bool("h", false, "help")

	// Server config
	addr := flag.String("addr", ":8999", "listen addr")

	// Platform keys
	wechatKey := flag.String("wecom.key", "", "WeChat webhook key")
	dingtalkKey := flag.String("dingtalk.key", "", "DingTalk webhook key")
	feishuKey := flag.String("feishu.key", "", "FeiShu webhook key")

	// Log config
	logDir := flag.String("log.dir", "/var/log/wechat-webhook", "log directory")
	logFile := flag.String("log.file", "webhook.log", "log file name")
	logMaxSize := flag.Int64("log.maxsize", 100, "max log file size in MB")

	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	return &Config{
		Addr:        *addr,
		WeChatKey:   *wechatKey,
		DingTalkKey: *dingtalkKey,
		FeiShuKey:   *feishuKey,
		LogDir:      *logDir,
		LogFile:     *logFile,
		LogMaxSize:  *logMaxSize,
	}
}