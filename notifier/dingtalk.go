package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k8stech/alertmanager-wechatrobot-webhook/model"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/transformer"
)

const maxDingTalkMessageLength = 20000 // DingTalk limit

// DingTalkNotifier sends notifications to DingTalk
type DingTalkNotifier struct {
	name string
}

// NewDingTalkNotifier creates a new DingTalk notifier
func NewDingTalkNotifier() *DingTalkNotifier {
	return &DingTalkNotifier{name: PlatformDingTalk}
}

// Name returns the platform name
func (n *DingTalkNotifier) Name() string {
	return n.name
}

// Send sends a notification to DingTalk
func (n *DingTalkNotifier) Send(notification *model.Notification, robotKey string) error {
	markdown, message, _, err := transformer.TransformToMarkdown(*notification)
	if err != nil {
		return fmt.Errorf("[DingTalk] transform failed: %w", err)
	}

	fmt.Printf("[DingTalk] Message:\n%s\n", message)
	return n.sendMessage(markdown.Markdown.Content, robotKey)
}

// sendMessage sends the message content to DingTalk
func (n *DingTalkNotifier) sendMessage(content string, robotKey string) error {
	msg := &model.DingTalkMessage{
		MsgType: "markdown",
		Markdown: &model.DingTalkMarkdown{
			Title: "告警通知",
			Text:  content,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("[DingTalk] marshal failed: %w", err)
	}

	url := "https://oapi.dingtalk.com/robot/send?access_token=" + robotKey
	resp, err := n.doRequest(url, data)
	if err != nil {
		return fmt.Errorf("[DingTalk] request failed: %w", err)
	}

	fmt.Printf("[DingTalk] Response: %s\n", resp)
	return nil
}

// doRequest performs the HTTP request
func (n *DingTalkNotifier) doRequest(url string, data []byte) (string, error) {
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