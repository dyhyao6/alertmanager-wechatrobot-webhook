package transformer

import (
	"bytes"
	"fmt"
	"time"

	"github.com/k8stech/alertmanager-wechatrobot-webhook/model"
)

// 新增一个函数来获取告警颜色
func getAlertColor(severity string) string {
	switch severity {
	case "critical":
		return "warning"
	case "firing":
		return "warning"
	case "resolved":
		return "info"
	default:
		return "comment"
	}
}

// 状态转译为中文
func getStatusText(status string) string {
	switch status {
	case "firing":
		return "触发中"
	case "resolved":
		return "已恢复"
	default:
		return status
	}
}

// 获取状态图标
func getStatusEmoji(status string) string {
	switch status {
	case "firing":
		return "🚨"
	case "resolved":
		return "✅"
	default:
		return "📌"
	}
}

// 获取级别显示
func getSeverityDisplay(severity string) string {
	switch severity {
	case "critical":
		return "🔴 严重"
	case "warning":
		return "🟡 警告"
	case "info":
		return "🔵 通知"
	default:
		return severity
	}
}

// 安全获取 label 值，避免字段缺失时空指针
func getLabelValue(labels map[string]string, key string) string {
	if value, ok := labels[key]; ok {
		return value
	}
	return ""
}

// 安全获取 annotation 值
func getAnnotationValue(annotations map[string]string, key string) string {
	if value, ok := annotations[key]; ok {
		return value
	}
	return ""
}

// TransformToMarkdown transform alertmanager notification to markdown message
// platform参数用于适配不同平台的格式
func TransformToMarkdown(notification model.Notification, platform string) (markdown *model.WeChatMarkdown, message string, robotURL string, err error) {

	status := notification.Status

	annotations := notification.CommonAnnotations
	robotURL = annotations["wechatRobot"]

	var buffer bytes.Buffer

	for _, alert := range notification.Alerts {
		labels := alert.Labels
		// 加载 CST 时区
		cstZone, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			fmt.Println("Error loading location:", err)
		}

		// 将 UTC 时间转换为 CST 时间
		cstTime := alert.StartsAt.In(cstZone)
		instance := getLabelValue(labels, "instance")
		severity := getLabelValue(labels, "severity")
		alertname := getLabelValue(labels, "alertname")

		// 标题
		summary := getAnnotationValue(annotations, "summary")

		// 根据平台选择不同格式
		switch platform {
		case "dingtalk":
			// 钉钉格式 - 简洁markdown
			buffer.WriteString(fmt.Sprintf("## %s %s | %s\n\n", getStatusEmoji(status), getStatusText(status), summary))
			buffer.WriteString(fmt.Sprintf("**级别** %s\n\n", getSeverityDisplay(severity)))
			buffer.WriteString(fmt.Sprintf("**类型** %s\n\n", alertname))
			buffer.WriteString(fmt.Sprintf("**主机** %s\n\n", instance))
			buffer.WriteString(fmt.Sprintf("**详情** %s\n\n", getAnnotationValue(alert.Annotations, "description")))
			buffer.WriteString(fmt.Sprintf("**时间** %s\n", cstTime.Format("2006-01-02 15:04:05")))
		default:
			// 微信/飞书格式 - 支持特殊字符
			buffer.WriteString(fmt.Sprintf("━━━━━━━━━━━━━━━━━━\n"))
			buffer.WriteString(fmt.Sprintf("%s %s | **%s**\n", getStatusEmoji(status), getStatusText(status), summary))
			buffer.WriteString(fmt.Sprintf("━━━━━━━━━━━━━━━━━━\n\n"))
			buffer.WriteString(fmt.Sprintf("📋 告警信息\n"))
			buffer.WriteString(fmt.Sprintf("├─ 🔴 级别：%s\n", getSeverityDisplay(severity)))
			buffer.WriteString(fmt.Sprintf("├─ 📛 类型：%s\n", alertname))
			buffer.WriteString(fmt.Sprintf("├─ 🖥 主机：%s\n", instance))
			buffer.WriteString(fmt.Sprintf("├─ 📝 详情：%s\n", getAnnotationValue(alert.Annotations, "description")))
			buffer.WriteString(fmt.Sprintf("└─ ⏰ 时间：%s\n", cstTime.Format("2006-01-02 15:04:05")))
		}

		message = buffer.String()
		markdown = &model.WeChatMarkdown{
			MsgType: "markdown",
			Markdown: &model.Markdown{
				Content: message,
			},
		}
	}

	return
}