package notifier

import "github.com/k8stech/alertmanager-wechatrobot-webhook/model"

// Platform constants
const (
	PlatformWeChat   = "wecom"
	PlatformDingTalk  = "dingtalk"
	PlatformFeiShu    = "feishu"
)

// Notifier defines the interface for sending notifications to different platforms
type Notifier interface {
	// Name returns the platform name
	Name() string

	// Send sends a notification to the platform
	Send(notification *model.Notification, robotKey string) error
}

// BaseNotifier provides common functionality for all notifiers
type BaseNotifier struct {
	PlatformName string
}