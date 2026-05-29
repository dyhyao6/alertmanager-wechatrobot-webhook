# 多平台 Webhook 技术方案

## 一、设计目标

支持企业微信、飞书、钉钉三大平台，通过不同接口路由到对应平台进行消息推送。

## 二、架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                      Alertmanager                           │
│  ────────────────────────────────────────────────────────  │
│  receivers:                                                │
│    - name: 'wecom'                                         │
│      webhook_configs:                                       │
│        - url: 'http://service:8999/webhook/wecom'           │
│    - name: 'dingtalk'                                      │
│      webhook_configs:                                       │
│        - url: 'http://service:8999/webhook/dingtalk'       │
│    - name: 'feishu'                                        │
│      webhook_configs:                                       │
│        - url: 'http://service:8999/webhook/feishu'         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Router Layer                           │
│  POST /webhook/wecom      → WeChatNotifier                │
│  POST /webhook/dingtalk   → DingTalkNotifier              │
│  POST /webhook/feishu     → FeiShuNotifier                │
│  POST /webhook/common     → 根据配置自动选择平台           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Platform Notifiers                       │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │ WeChat    │  │ DingTalk  │  │ FeiShu    │            │
│  │ Notifier │  │ Notifier │  │ Notifier │            │
│  └────────────┘  └────────────┘  └────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

## 三、目录结构

```
alertmanager-wechatrobot-webhook/
├── main.go                      # 入口，路由配置
├── config/
│   └── config.go               # 配置管理
├── model/
│   ├── alertmanager.go        # Alertmanager 数据结构
│   ├── wechat.go              # 微信消息结构
│   ├── dingtalk.go            # 钉钉消息结构
│   └── feishu.go              # 飞书消息结构
├── notifier/
│   ├── notifier.go           # 接口定义
│   ├── wechat.go             # 企业微信通知器
│   ├── dingtalk.go           # 钉钉通知器
│   └── feishu.go             # 飞书通知器
├── transformer/
│   ├── transformer.go        # 公共转换逻辑
│   ├── wechat.go             # 微信消息转换
│   ├── dingtalk.go           # 钉钉消息转换
│   └── feishu.go             # 飞书消息转换
└── router/
    └── router.go             # 路由处理
```

## 四、核心接口设计

### 4.1 Notifier 接口

```go
type Notifier interface {
    // 获取平台名称
    Name() string

    // 发送消息
    Send(notification *Notification, config *RobotConfig) error
}
```

### 4.2 Transformer 接口

```go
type Transformer interface {
    // 转换为各平台消息格式
    Transform(notification *Notification) (interface{}, error)
}
```

### 4.3 RobotConfig 配置

```go
type RobotConfig struct {
    Platform   string            // 平台类型：wecom, dingtalk, feishu
    RobotKey  string            // Webhook Key
    Extra     map[string]string // 平台特有配置
}
```

## 五、配置方案

### 5.1 启动参数

```bash
./wechat-webhook \
  -wecom.key=your-wechat-key \
  -dingtalk.key=your-dingtalk-key \
  -feishu.key=your-feishu-key \
  -addr :8999
```

### 5.2 配置文件 config.env

```bash
# 企业微信
WECOM_KEY=your-wechat-key

# 钉钉
DINGTALK_KEY=your-dingtalk-key

# 飞书
FEISHU_KEY=your-feishu-key

# 服务地址
ADDR=:8999
```

## 六、接口定义

| 接口 | 方法 | 说明 | 消息格式 |
|------|------|------|---------|
| `/webhook/wecom` | POST | 企业微信 | Markdown |
| `/webhook/dingtalk` | POST | 钉钉 | Markdown/Text |
| `/webhook/feishu` | POST | 飞书 | Markdown |
| `/webhook/common` | POST | 通用接口 | 根据请求头或配置选择平台 |

## 七、实现步骤

### Phase 1: 重构基础设施
1. 创建 `config` 包管理配置
2. 抽象 `Notifier` 接口
3. 创建 `RobotConfig` 结构
4. 实现 `router` 包处理路由

### Phase 2: 企业微信（保留现有逻辑）
1. 将现有 `notifier/wechat.go` 适配新接口
2. 将现有 `transformer` 适配新接口
3. 测试验证

### Phase 3: 钉钉支持
1. 创建 `notifier/dingtalk.go`
2. 创建 `transformer/dingtalk.go`
3. 实现钉钉 Markdown 消息格式转换
4. 测试验证

### Phase 4: 飞书支持
1. 创建 `notifier/feishu.go`
2. 创建 `transformer/feishu.go`
3. 实现飞书 Markdown 消息格式转换
4. 测试验证

## 八、各平台消息格式对比

| 平台 | Markdown 支持 | 特殊格式 |
|------|-------------|---------|
| 企业微信 | ✅ 原生支持 | `<font color='xxx'>` |
| 钉钉 | ✅ Markdown 类型 | `### 标题` |
| 飞书 | ✅ 原生支持 | `font color` |

## 九、向后兼容

保留现有 `/webhook` 接口作为 `wecom` 的别名，确保现有配置无需修改即可继续使用。

```go
router.POST("/webhook", wechatNotifier.Handle)
router.POST("/webhook/wecom", wechatNotifier.Handle)
```

## 十、风险与注意事项

1. **线程安全**：Notifier 需要考虑并发发送
2. **错误处理**：各平台返回错误码不同，需要统一错误处理
3. **消息大小**：钉钉消息限制 20KB，其他平台 4096 字符
4. **重试机制**：发送失败时需要重试
5. **敏感信息**：RobotKey 需要安全存储