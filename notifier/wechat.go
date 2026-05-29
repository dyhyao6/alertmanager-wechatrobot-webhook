package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/k8stech/alertmanager-wechatrobot-webhook/model"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/transformer"
)

const maxMessageLength = 4096

// WeChatNotifier sends notifications to WeChat
type WeChatNotifier struct {
	name string
}

// NewWeChatNotifier creates a new WeChat notifier
func NewWeChatNotifier() *WeChatNotifier {
	return &WeChatNotifier{name: PlatformWeChat}
}

// Name returns the platform name
func (n *WeChatNotifier) Name() string {
	return n.name
}

// Send sends a notification to WeChat
func (n *WeChatNotifier) Send(notification *model.Notification, robotKey string) error {
	markdown, message, _, err := transformer.TransformToMarkdown(*notification)
	if err != nil {
		return fmt.Errorf("[WeChat] transform failed: %w", err)
	}

	fmt.Printf("[WeChat] Message:\n%s\n", message)
	return n.sendMessage(markdown.Markdown.Content, robotKey)
}

// sendMessage sends the message content to WeChat
func (n *WeChatNotifier) sendMessage(content string, robotKey string) error {
	if len(content) > maxMessageLength {
		chunks := splitContent(content, maxMessageLength)
		for _, chunk := range chunks {
			if err := n.sendSingleMessage(chunk, robotKey); err != nil {
				return err
			}
		}
		return nil
	}
	return n.sendSingleMessage(content, robotKey)
}

// sendSingleMessage sends a single message to WeChat
func (n *WeChatNotifier) sendSingleMessage(content string, robotKey string) error {
	msg := &model.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &model.Markdown{
			Content: content,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("[WeChat] marshal failed: %w", err)
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + robotKey
	resp, err := n.doRequest(url, data)
	if err != nil {
		return fmt.Errorf("[WeChat] request failed: %w", err)
	}

	fmt.Printf("[WeChat] Response: %s\n", resp)
	return nil
}

// doRequest performs the HTTP request
func (n *WeChatNotifier) doRequest(url string, data []byte) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

// splitContent splits content into chunks
func splitContent(content string, maxLen int) []string {
	var chunks []string
	for len(content) > maxLen {
		splitIndex := maxLen
		if idx := strings.LastIndex(content[:maxLen], "\n"); idx != -1 {
			splitIndex = idx + 1
		}
		chunks = append(chunks, content[:splitIndex])
		content = content[splitIndex:]
	}
	chunks = append(chunks, content)
	return chunks
}