package router

import (
	"github.com/gin-gonic/gin"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/config"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/logger"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/model"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/notifier"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/transformer"
	"net/http"
)

// Router holds all the route handlers
type Router struct {
	engine     *gin.Engine
	cfg        *config.Config
	notifiers  map[string]notifier.Notifier
	logger     *logger.Logger
}

// New creates a new router instance
func New(cfg *config.Config, notifiers map[string]notifier.Notifier, appLogger *logger.Logger) *Router {
	return &Router{
		cfg:       cfg,
		notifiers: notifiers,
		logger:    appLogger,
	}
}

// Setup sets up all routes and returns the gin engine
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Platform-specific webhook endpoints
	webhook := engine.Group("/webhook")
	{
		// WeChat endpoint
		webhook.POST("/wecom", r.wrapHandler(notifier.PlatformWeChat))
		// Alias for backward compatibility
		webhook.POST("", r.wrapHandler(notifier.PlatformWeChat))

		// DingTalk endpoint
		webhook.POST("/dingtalk", r.wrapHandler(notifier.PlatformDingTalk))

		// FeiShu endpoint
		webhook.POST("/feishu", r.wrapHandler(notifier.PlatformFeiShu))
	}

	return engine
}

// wrapHandler creates a handler function for a specific platform
func (r *Router) wrapHandler(platform string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification model.Notification
		if err := c.BindJSON(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		n, ok := r.notifiers[platform]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "notifier not found for platform: " + platform})
			return
		}

		// Get robot key from query param or use configured key
		robotKey := c.DefaultQuery("key", r.getPlatformKey(platform))

		if err := n.Send(&notification, robotKey); err != nil {
			r.logger.Error("Failed to send notification", logger.Fields{
				"platform": platform,
				"error":    err.Error(),
			})
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取发送的消息内容并记录日志
		_, message, _, _ := transformer.TransformToMarkdown(notification, platform)
		r.logger.Info("Notification sent", logger.Fields{
			"platform": platform,
			"status":   notification.Status,
			"alerts":   len(notification.Alerts),
			"message":  message,
		})

		c.JSON(http.StatusOK, gin.H{"message": "sent to " + platform + " successfully"})
	}
}

// getPlatformKey returns the configured key for a platform
func (r *Router) getPlatformKey(platform string) string {
	switch platform {
	case notifier.PlatformWeChat:
		return r.cfg.WeChatKey
	case notifier.PlatformDingTalk:
		return r.cfg.DingTalkKey
	case notifier.PlatformFeiShu:
		return r.cfg.FeiShuKey
	default:
		return ""
	}
}