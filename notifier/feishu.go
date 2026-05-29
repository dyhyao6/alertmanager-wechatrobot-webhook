package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k8stech/alertmanager-wechatrobot-webhook/model"
	"github.com/k8stech/alertmanager-wechatrobot-webhook/transformer"
)

// FeiShuNotifier sends notifications to FeiShu (Lark)
type FeiShuNotifier struct {
	name string
}

// NewFeiShuNotifier creates a new FeiShu notifier
func NewFeiShuNotifier() *FeiShuNotifier {
	return &FeiShuNotifier{name: PlatformFeiShu}
}

// Name returns the platform name
func (n *FeiShuNotifier) Name() string {
	return n.name
}

// Send sends a notification to FeiShu
func (n *FeiShuNotifier) Send(notification *model.Notification, robotKey string) error {
	markdown, message, _, err := transformer.TransformToMarkdown(*notification)
	if err != nil {
		return fmt.Errorf("[FeiShu] transform failed: %w", err)
	}

	fmt.Printf("[FeiShu] Message:\n%s\n", message)
	return n.sendMessage(markdown.Markdown.Content, robotKey)
}

// sendMessage sends the message content to FeiShu
func (n *FeiShuNotifier) sendMessage(content string, robotKey string) error {
	msg := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": content,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("[FeiShu] marshal failed: %w", err)
	}

	url := "https://open.feishu.cn/open-apis/bot/v2/hook/" + robotKey
	resp, err := n.doRequest(url, data)
	if err != nil {
		return fmt.Errorf("[FeiShu] request failed: %w", err)
	}

	fmt.Printf("[FeiShu] Response: %s\n", resp)
	return nil
}

// doRequest performs the HTTP request
func (n *FeiShuNotifier) doRequest(url string, data []byte) (string, error) {
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